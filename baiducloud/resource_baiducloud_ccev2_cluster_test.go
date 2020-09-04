package baiducloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	ccev2 "github.com/baidubce/bce-sdk-go/services/cce/v2"
	"github.com/baidubce/bce-sdk-go/services/cce/v2/types"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccCcev2ClusterResource = "baiducloud_ccev2_cluster"
)

func init() {
	resource.AddTestSweepers(testAccCcev2ClusterResource, &resource.Sweeper{
		Name: testAccCcev2ClusterResource,
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
		if !strings.HasPrefix(c.Spec.ClusterName, BaiduCloudTestResourceAttrNamePrefix) {
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
				Config: testAccCcev2ClusterConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId("baiducloud_ccev2_cluster.default_managed"),
				),
			},
		},
	})
}

func testAccCcev2ClusterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccCcev2ClusterResource {
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

func testAccCcev2ClusterConfig() string {
	return fmt.Sprintf(`
variable "cluster_name" {
  default = "%s"
}
variable "cluster_pod_cidr" {
  default = "172.28.0.0/16"
}

variable "cluster_ip_service_cidr" {
  default = "172.31.0.0/16"
}
resource "baiducloud_vpc" "default" {
  name        = "%s"
  description = "test-BaiduAcc_test-vpc-tf-auto"
  cidr        = "192.168.0.0/16"
}
data "baiducloud_zones" "defaultA" {
  name_regex = ".*a$"
}
resource "baiducloud_subnet" "defaultA" {
  name        = "%s"
  zone_name   = data.baiducloud_zones.defaultA.zones.0.zone_name
  cidr        = "192.168.1.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "test-subnet-tf-auto"
}
resource "baiducloud_security_group" "default" {
  name   = "%s"
  vpc_id = baiducloud_vpc.default.id
}
resource "baiducloud_security_group_rule" "default" {
  security_group_id = baiducloud_security_group.default.id
  remark            = "remark"
  protocol          = "all"
  port_range        = "1-65535"
  direction         = "ingress"
}
resource "baiducloud_security_group_rule" "default2" {
  security_group_id = baiducloud_security_group.default.id
  remark            = "remark"
  protocol          = "all"
  port_range        = "1-65535"
  direction         = "egress"
}
resource "baiducloud_ccev2_cluster" "default_managed" {
  cluster_spec  {
    cluster_name = var.cluster_name
    k8s_version = "1.16.8"
    runtime_type = "docker"
    vpc_id = baiducloud_vpc.default.id

    master_config {
      master_type = "managed"
      cluster_ha = 1
      exposed_public = false
      cluster_blb_vpc_subnet_id = baiducloud_subnet.defaultA.id
      managed_cluster_master_option {
        master_vpc_subnet_zone = "zoneA"
      }
    }
    container_network_config  {
      mode = "kubenet"
      lb_service_vpc_subnet_id = baiducloud_subnet.defaultA.id
      cluster_pod_cidr = var.cluster_pod_cidr
      cluster_ip_service_cidr = var.cluster_ip_service_cidr
    }
    cluster_delete_option {
      delete_resource = true
      delete_cds_snapshot = true
    }
  }
}
`,
		BaiduCloudTestResourceAttrNamePrefix+"_ccev2_cluster_1",
		BaiduCloudTestResourceAttrNamePrefix+"_test-vpc-tf-auto",
		BaiduCloudTestResourceAttrNamePrefix+"_test-subnet-tf-auto",
		BaiduCloudTestResourceAttrNamePrefix+"_test-security-group-tf-auto")
}
