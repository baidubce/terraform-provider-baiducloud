package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const (
	testAccCcev2ClusterInstancesSourceName = "data.baiducloud_ccev2_cluster_instances.default"
)

func TestAccBaiduCloudCCEv2ClusterInstancesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccCcev2ClusterNodesDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccCcev2ClusterInstancesSourceName),
					resource.TestCheckResourceAttrSet(testAccCcev2ClusterInstancesSourceName, "total_count"),
				),
			},
		},
	})
}

func testAccCcev2ClusterNodesDataSourceConfig() string {
	return fmt.Sprintf(`
variable "vpc_cidr" {
  default = "192.168.0.0/16"
}
variable "container_cidr" {
  default = "172.28.0.0/16"
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
    cluster_name = "%s"
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
data "baiducloud_ccev2_cluster_instances" "default" {
  cluster_id = baiducloud_ccev2_cluster.default_managed.id
  keyword_type = "instanceName"
  keyword = "t"
  order_by = "instanceName"
  order = "ASC"
  page_no = 0
  page_size = 0
}
`,
		BaiduCloudTestResourceAttrNamePrefix+"_test-vpc-tf-auto",
		BaiduCloudTestResourceAttrNamePrefix+"_test-subnet-tf-auto",
		BaiduCloudTestResourceAttrNamePrefix+"_test-security-group-tf-auto",
		BaiduCloudTestResourceAttrNamePrefix+"_ccev2_cluster_1")
}
