package baiducloud

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	ccev2 "github.com/baidubce/bce-sdk-go/services/cce/v2"
	"github.com/baidubce/bce-sdk-go/services/cce/v2/types"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccCcev2ClusterResourceType        = "baiducloud_ccev2_cluster"
	testAccCCEv2ClusterFeatureResourceName = "baiducloud_ccev2_cluster.default"
)

func init() {
	resource.AddTestSweepers(testAccCcev2ClusterResourceType, &resource.Sweeper{
		Name: testAccCcev2ClusterResourceType,
		F:    testSweepCcev2Cluster,
	})
}

func testSweepCcev2Cluster(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("get BaiduCloud client error: %s", err)
	}
	client := rawClient.(*connectivity.BaiduClient)

	listArgs := &ccev2.ListClustersArgs{
		KeywordType: ccev2.ClusterKeywordTypeClusterName,
		Keyword:     "",
		OrderBy:     ccev2.ClusterOrderByClusterName,
		Order:       ccev2.OrderASC,
		PageSize:    1000,
		PageNum:     1,
	}
	raw, err := client.WithCCEv2Client(func(ccev2Client *ccev2.Client) (i interface{}, e error) {
		return ccev2Client.ListClusters(listArgs)
	})
	if err != nil {
		return fmt.Errorf("list CCEv2 Cluster with error: %s", err)
	}

	clusterList := raw.(*ccev2.ListClustersResponse).ClusterPage.ClusterList
	for _, c := range clusterList {
		if !strings.HasPrefix(c.Spec.ClusterName, BaiduCloudTestResourceTypeName) {
			log.Printf("[INFO] Skipping CCEv2 Cluster: %s (%s)", c.Spec.ClusterName, c.Spec.ClusterID)
			continue
		}

		if c.Status.ClusterPhase == types.ClusterPhaseDeleting || c.Status.ClusterPhase == types.ClusterPhaseDeleted {
			log.Printf("[INFO] Skipping CCE Cluster: %s (%s) with status %s",
				c.Spec.ClusterName, c.Spec.ClusterID, c.Status.ClusterPhase)
			continue
		}
		log.Printf("[INFO] Deleting CCE Cluster: %s (%s)", c.Spec.ClusterName, c.Spec.ClusterID)

		deleteArgs := &ccev2.DeleteClusterArgs{
			ClusterID:         c.Spec.ClusterID,
			DeleteCDSSnapshot: true,
			DeleteResource:    true,
		}
		_, err := client.WithCCEv2Client(func(ccev2Client *ccev2.Client) (i interface{}, e error) {
			return ccev2Client.DeleteCluster(deleteArgs)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete CCE cluster %s (%s) with error: %v",
				c.Spec.ClusterName, c.Spec.ClusterID, err)
		}
	}
	return nil
}

func TestAccBaiduCloudCCEv2ClusterResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCcev2ClusterDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCcev2ClusterConfig(BaiduCloudTestResourceTypeNameCcev2Cluster),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId("baiducloud_ccev2_cluster.default"),
				),
			},
		},
	})
}

func testAccCcev2ClusterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccCcev2ClusterResourceType {
			continue
		}
		raw, err := client.WithCCEv2Client(func(client *ccev2.Client) (i interface{}, e error) {
			return client.GetCluster(rs.Primary.ID)
		})

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		cluster := raw.(*ccev2.GetClusterResponse)
		if cluster.Cluster.Status.ClusterPhase == types.ClusterPhaseDeleted {
			continue
		}

		return WrapError(Error("CCE Cluster still exist"))
	}
	return nil
}

func testAccCcev2ClusterConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "cluster_pod_cidr" {
  default = "172.28.0.0/16"
}

variable "cluster_ip_service_cidr" {
  default = "172.31.0.0/16"
}
resource "baiducloud_vpc" "default" {
  name        = "${var.name}"
  description = "created by terraform"
  cidr        = "192.168.0.0/16"
}
data "baiducloud_zones" "defaultA" {
  name_regex = ".*e$"
}
resource "baiducloud_subnet" "default" {
  name        = "${var.name}"
  zone_name   = data.baiducloud_zones.defaultA.zones.0.zone_name
  cidr        = "192.168.1.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created by terraform"
}
resource "baiducloud_security_group" "default" {
  name   = "${var.name}"
  vpc_id = baiducloud_vpc.default.id
}
resource "baiducloud_security_group_rule" "ingress" {
  security_group_id = baiducloud_security_group.default.id
  remark            = "remark"
  protocol          = "all"
  port_range        = "1-65535"
  direction         = "ingress"
}
resource "baiducloud_security_group_rule" "egress" {
  security_group_id = baiducloud_security_group.default.id
  remark            = "remark"
  protocol          = "all"
  port_range        = "1-65535"
  direction         = "egress"
}
resource "baiducloud_ccev2_cluster" "default" {
  cluster_spec  {
    cluster_name = var.name
    k8s_version = "1.16.8"
    runtime_type = "docker"
    vpc_id = baiducloud_vpc.default.id

    master_config {
      master_type = "managed"
      cluster_ha = 1
      exposed_public = false
      cluster_blb_vpc_subnet_id = baiducloud_subnet.default.id
      managed_cluster_master_option {
        master_vpc_subnet_zone = "zoneA"
      }
    }
    container_network_config  {
      mode = "kubenet"
      lb_service_vpc_subnet_id = baiducloud_subnet.default.id
      cluster_pod_cidr = var.cluster_pod_cidr
      cluster_ip_service_cidr = var.cluster_ip_service_cidr
    }
    cluster_delete_option {
      delete_resource = true
      delete_cds_snapshot = true
    }
  }
}
`, name)
}

func testAccCCEv2ClusterFeaturePreCheck(t *testing.T) {
	testAccPreCheck(t)

	requiredEnvKeys := []string{
		"BAIDUCLOUD_TEST_CCEV2_VPC_ID",
		"BAIDUCLOUD_TEST_CCEV2_SUBNET_ID",
		"BAIDUCLOUD_TEST_CCEV2_KMS_KEY_ID",
	}
	for _, key := range requiredEnvKeys {
		if os.Getenv(key) == "" {
			t.Fatalf("%s must be set for CCEv2 cluster feature acceptance tests", key)
		}
	}
}

func TestAccBaiduCloudCCEv2ClusterResourceKMSAndSAN(t *testing.T) {
	clusterName := fmt.Sprintf("%s-kms-san-%d", BaiduCloudTestResourceTypeNameCcev2Cluster, time.Now().Unix())
	kmsKeyID := os.Getenv("BAIDUCLOUD_TEST_CCEV2_KMS_KEY_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccCCEv2ClusterFeaturePreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCcev2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEv2ClusterFeatureConfig(clusterName, []string{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccCCEv2ClusterFeatureResourceName),
					resource.TestCheckResourceAttr(testAccCCEv2ClusterFeatureResourceName, "kms_encryption.#", "1"),
					resource.TestCheckResourceAttr(testAccCCEv2ClusterFeatureResourceName, "kms_encryption.0.enabled", "true"),
					resource.TestCheckResourceAttr(testAccCCEv2ClusterFeatureResourceName, "kms_encryption.0.kms_key_id", kmsKeyID),
					resource.TestCheckResourceAttr(testAccCCEv2ClusterFeatureResourceName, "api_server_cert_san.#", "0"),
				),
			},
			{
				Config: testAccCCEv2ClusterFeatureConfig(clusterName, []string{"k8s.sdk-test.internal"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccCCEv2ClusterFeatureResourceName),
					resource.TestCheckResourceAttr(testAccCCEv2ClusterFeatureResourceName, "kms_encryption.#", "1"),
					resource.TestCheckResourceAttr(testAccCCEv2ClusterFeatureResourceName, "kms_encryption.0.enabled", "true"),
					resource.TestCheckResourceAttr(testAccCCEv2ClusterFeatureResourceName, "kms_encryption.0.kms_key_id", kmsKeyID),
					resource.TestCheckResourceAttr(testAccCCEv2ClusterFeatureResourceName, "api_server_cert_san.#", "1"),
					resource.TestCheckResourceAttr(testAccCCEv2ClusterFeatureResourceName, "api_server_cert_san.0", "k8s.sdk-test.internal"),
				),
			},
		},
	})
}

func TestAccBaiduCloudCCEv2ClusterImportUpdateSANAndDelete(t *testing.T) {
	existingClusterID := os.Getenv("BAIDUCLOUD_TEST_CCEV2_EXISTING_CLUSTER_ID")
	if existingClusterID == "" {
		t.Skip("BAIDUCLOUD_TEST_CCEV2_EXISTING_CLUSTER_ID is not set")
	}

	clusterName := BaiduCloudTestResourceTypeNameCcev2Cluster + "-kms-san"
	kmsKeyID := os.Getenv("BAIDUCLOUD_TEST_CCEV2_KMS_KEY_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccCCEv2ClusterFeaturePreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCcev2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config:            testAccCCEv2ClusterFeatureConfig(clusterName, []string{}),
				ResourceName:      testAccCCEv2ClusterFeatureResourceName,
				ImportState:       true,
				ImportStateId:     existingClusterID,
				ImportStateVerify: false,
			},
			{
				Config: testAccCCEv2ClusterFeatureConfig(clusterName, []string{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccCCEv2ClusterFeatureResourceName),
					resource.TestCheckResourceAttr(testAccCCEv2ClusterFeatureResourceName, "kms_encryption.#", "1"),
					resource.TestCheckResourceAttr(testAccCCEv2ClusterFeatureResourceName, "kms_encryption.0.enabled", "true"),
					resource.TestCheckResourceAttr(testAccCCEv2ClusterFeatureResourceName, "kms_encryption.0.kms_key_id", kmsKeyID),
					resource.TestCheckResourceAttr(testAccCCEv2ClusterFeatureResourceName, "api_server_cert_san.#", "0"),
				),
			},
			{
				Config: testAccCCEv2ClusterFeatureConfig(clusterName, []string{"k8s.sdk-test.internal"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccCCEv2ClusterFeatureResourceName),
					resource.TestCheckResourceAttr(testAccCCEv2ClusterFeatureResourceName, "api_server_cert_san.#", "1"),
					resource.TestCheckResourceAttr(testAccCCEv2ClusterFeatureResourceName, "api_server_cert_san.0", "k8s.sdk-test.internal"),
				),
			},
		},
	})
}

func testAccCCEv2ClusterFeatureConfig(name string, apiServerCertSAN []string) string {
	vpcID := os.Getenv("BAIDUCLOUD_TEST_CCEV2_VPC_ID")
	subnetID := os.Getenv("BAIDUCLOUD_TEST_CCEV2_SUBNET_ID")
	kmsKeyID := os.Getenv("BAIDUCLOUD_TEST_CCEV2_KMS_KEY_ID")

	return fmt.Sprintf(`
resource "baiducloud_ccev2_cluster" "default" {
  cluster_spec {
    cluster_name = %q
    k8s_version  = "1.31.1"
    vpc_id       = %q

    master_config {
      master_type               = "managedPro"
      exposed_public            = false
      cluster_blb_vpc_subnet_id = %q

      managed_cluster_master_option {
        master_flavor = "l50"
      }
    }

    container_network_config {
      mode                     = "vpc-secondary-ip-veth"
      ip_version               = "ipv4"
      node_port_range_min      = 30000
      node_port_range_max      = 32767
      cluster_ip_service_cidr  = "172.16.0.0/16"
      kube_proxy_mode          = "ipvs"
      lb_service_vpc_subnet_id = %q
      network_policy_type      = "none"
      enable_node_local_dns    = false
      net_device_driver        = "veth-pair"
      enable_rdma              = true

      ebpf_config {
        enabled = false
      }

      eni_vpc_subnet_ids {
        zone_and_id = {
          zoneD = %q
        }
      }
    }

    k8s_custom_config {
      enable_hostname                = false
      etcd_data_path                 = "/home/cce/etcd"
      enable_cloud_node_controller   = true
      enable_lb_service_controller   = true
      disable_kubelet_read_only_port = false
    }
  }

  create_options {
    skip_network_check = false
  }

  kms_encryption {
    enabled    = true
    kms_key_id = %q
  }

  api_server_cert_san = %s
}
`, name, vpcID, subnetID, subnetID, subnetID, kmsKeyID, terraformStringList(apiServerCertSAN))
}

func terraformStringList(values []string) string {
	if len(values) == 0 {
		return "[]"
	}

	quoted := make([]string, 0, len(values))
	for _, value := range values {
		quoted = append(quoted, fmt.Sprintf("%q", value))
	}
	return fmt.Sprintf("[%s]", strings.Join(quoted, ", "))
}
