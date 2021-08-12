package baiducloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	ccev2 "github.com/baidubce/bce-sdk-go/services/cce/v2"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccCcev2InstanceGroupResource = "baiducloud_ccev2_instance_group"
)

func init() {
	resource.AddTestSweepers(testAccCcev2InstanceGroupResource, &resource.Sweeper{
		Name: testAccCcev2InstanceGroupResource,
		F:    testSweepCcev2InstanceGroup,
		Dependencies: []string{
			testAccSecurityGroupResourceType,
			testAccVPCResourceType,
		},
	})
}

func testSweepCcev2InstanceGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("get BaiduCloud client error: %s", err)
	}
	client := rawClient.(*connectivity.BaiduClient)

	//Get all cluster
	listClusterArgs := &ccev2.ListClustersArgs{
		KeywordType: ccev2.ClusterKeywordTypeClusterName,
		Keyword:     "",
		OrderBy:     ccev2.ClusterOrderByClusterName,
		Order:       ccev2.OrderASC,
		PageSize:    0,
		PageNum:     0,
	}
	raw, err := client.WithCCEv2Client(func(ccev2Client *ccev2.Client) (i interface{}, e error) {
		return ccev2Client.ListClusters(listClusterArgs)
	})
	if err != nil {
		return fmt.Errorf("list CCEv2 Cluster with error: %s", err)
	}

	clusterList := raw.(*ccev2.ListClustersResponse).ClusterPage.ClusterList
	for _, cluster := range clusterList {
		listIGArgs := &ccev2.ListInstanceGroupsArgs{
			ClusterID: cluster.Spec.ClusterID,
			ListOption: &ccev2.InstanceGroupListOption{
				PageSize: 0,
				PageNo:   0,
			},
		}
		raw, err := client.WithCCEv2Client(func(ccev2Client *ccev2.Client) (i interface{}, e error) {
			return ccev2Client.ListInstanceGroups(listIGArgs)
		})
		if err != nil {
			return fmt.Errorf("list CCEv2 Instance Group with error: %s", err)
		}
		instanceGroupList := raw.(*ccev2.ListInstancesByInstanceGroupIDResponse).Page.List
		for _, ig := range instanceGroupList {
			if !strings.HasPrefix(ig.Spec.InstanceGroupName, BaiduCloudTestResourceTypeName) {
				log.Printf("[INFO] Skipping CCEv2 Cluster: %s (%s)", ig.Spec.InstanceGroupName, ig.Spec.InstanceGroupID)
				continue
			}
			log.Printf("[INFO] Deleting CCE Cluster: %s (%s)", ig.Spec.InstanceGroupName, ig.Spec.InstanceGroupID)
			deleteArgs := &ccev2.DeleteInstanceGroupArgs{
				ClusterID:       ig.Spec.ClusterID,
				InstanceGroupID: ig.Spec.InstanceGroupID,
				DeleteInstances: true,
			}
			_, err := client.WithCCEv2Client(func(ccev2Client *ccev2.Client) (i interface{}, e error) {
				return ccev2Client.DeleteInstanceGroup(deleteArgs)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete CCE cluster %s (%s) with error: %v",
					ig.Spec.InstanceGroupName, ig.Spec.InstanceGroupID, err)
			}
		}
	}

	return nil
}

func TestAccBaiduCloudCCEv2InstanceGroupResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCcev2InstanceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCcev2InstanceGroupConfig(BaiduCloudTestResourceTypeNameCcev2InstanceGroup),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId("baiducloud_ccev2_instance_group.ccev2_instance_group_1"),
				),
			},
			{
				Config: testAccCcev2InstanceGroupUpdateConfig(BaiduCloudTestResourceTypeNameCcev2InstanceGroup),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId("baiducloud_ccev2_instance_group.ccev2_instance_group_1"),
				),
			},
		},
	})
}

func testAccCcev2InstanceGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccCcev2InstanceGroupResource {
			continue
		}
		args := &ccev2.GetInstanceGroupArgs{
			ClusterID:       rs.Primary.Attributes["baiducloud_ccev2_instance_group.ccev2_instance_group_1.spec.0.cluster_id"],
			InstanceGroupID: rs.Primary.ID,
		}
		log.Println("CheckDestroy获取的Cluster ID是" + args.ClusterID)
		_, err := client.WithCCEv2Client(func(client *ccev2.Client) (i interface{}, e error) {
			return client.GetInstanceGroup(args)
		})
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}
	return nil
}

func testAccCcev2InstanceGroupConfig(name string) string {
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
  name        = var.name
  description = "created by terraform"
  cidr        = "192.168.0.0/16"
}
data "baiducloud_zones" "default" {
  name_regex = ".*e$"
}
resource "baiducloud_subnet" "default" {
  name        = var.name
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.1.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created by terraform"
}
resource "baiducloud_security_group" "default" {
  name   = var.name
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
resource "baiducloud_ccev2_cluster" "default_managed" {
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
data "baiducloud_images" "default" {
  image_type = "System"
  name_regex = "7.5.*"
  os_name    = "CentOS"
}
resource "baiducloud_ccev2_instance_group" "ccev2_instance_group_1" {
  spec {
    cluster_id = baiducloud_ccev2_cluster.default_managed.id
    replicas = 1
    instance_group_name = var.name
    instance_template {
      cce_instance_id = ""
      instance_name = var.name
      cluster_role = "node"
      existed = false
      machine_type = "BCC"
      instance_type = "N3"
      vpc_config {
        vpc_id = baiducloud_vpc.default.id
        vpc_subnet_id = baiducloud_subnet.default.id
        security_group_id = baiducloud_security_group.default.id
        available_zone = "zoneA"
      }
      instance_resource {
        cpu = 4
        mem = 8
        root_disk_size = 40
        local_disk_size = 0
      }
      image_id = data.baiducloud_images.default.images.0.id
      instance_os {
        image_type = "System"
      }
      need_eip = false
      admin_password = "test123!YT"
      ssh_key_id = ""
      instance_charging_type = "Postpaid"
      runtime_type = "docker"
    }
  }
}
`, name)
}

func testAccCcev2InstanceGroupUpdateConfig(name string) string {
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
  name        = var.name
  description = "created by terraform"
  cidr        = "192.168.0.0/16"
}
data "baiducloud_zones" "default" {
  name_regex = ".*e$"
}
resource "baiducloud_subnet" "default" {
  name        = var.name
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.1.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created by terraform"
}
resource "baiducloud_security_group" "default" {
  name   = var.name
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
data "baiducloud_images" "default" {
  image_type = "System"
  name_regex = "7.5.*"
  os_name    = "CentOS"
}
resource "baiducloud_ccev2_instance_group" "ccev2_instance_group_1" {
  spec {
    cluster_id = baiducloud_ccev2_cluster.default_managed.id
    replicas = 0
    instance_group_name = var.name
    instance_template {
      cce_instance_id = ""
      instance_name = var.name
      cluster_role = "node"
      existed = false
      machine_type = "BCC"
      instance_type = "N3"
      vpc_config {
        vpc_id = baiducloud_vpc.default.id
        vpc_subnet_id = baiducloud_subnet.default.id
        security_group_id = baiducloud_security_group.default.id
        available_zone = "zoneA"
      }
      instance_resource {
        cpu = 4
        mem = 8
        root_disk_size = 40
        local_disk_size = 0
      }
      image_id = data.baiducloud_images.default.images.0.id
      instance_os {
        image_type = "System"
      }
      need_eip = false
      admin_password = "test123!YT"
      ssh_key_id = ""
      instance_charging_type = "Postpaid"
      runtime_type = "docker"
    }
  }
}
`, name)
}
