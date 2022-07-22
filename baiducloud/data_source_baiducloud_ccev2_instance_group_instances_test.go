package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccCcev2InstanceGroupInstancesSourceName = "data.baiducloud_ccev2_instance_group_instances.default"
)

func TestAccBaiduCloudCCEv2InstanceGroupInstancesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccCcev2InstanceGroupInstancesDataSourceConfig(BaiduCloudTestResourceTypeNameCcev2InstanceGroup),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccCcev2InstanceGroupInstancesSourceName),
					resource.TestCheckResourceAttrSet(testAccCcev2InstanceGroupInstancesSourceName, "total_count"),
				),
			},
		},
	})
}

func testAccCcev2InstanceGroupInstancesDataSourceConfig(name string) string {
	return fmt.Sprintf(
		`
variable "name" {
  default = "%s"
}
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
resource "baiducloud_ccev2_instance_group" "default" {
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
      deploy_custom_config {
        pre_user_script  = "ls"
        post_user_script = "date"
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
data "baiducloud_ccev2_instance_group_instances" "default" {
  cluster_id = baiducloud_ccev2_cluster.default_managed.id
  instance_group_id = baiducloud_ccev2_instance_group.default.id
  page_no = 0
  page_size = 0
}
`, name)

}
