package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const (
	testAccCcev2InstanceResource = "baiducloud_ccev2_instance"
)

// Resource baiducloud_ccev2_instance just gives user a way to update the spec of a specific instance in instancegroup.
// This resource will not create any real IaaS resource on the cloud.
// So there is no need to implement a sweep function.

func TestAccBaiduCloudCCEv2InstanceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccCcev2InstanceConfig(BaiduCloudTestResourceTypeNameCcev2Instance),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId("baiducloud_ccev2_instance.default"),
				),
			},
			{
				Config: testAccCcev2UpdateConfig(BaiduCloudTestResourceTypeNameCcev2Instance),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId("baiducloud_ccev2_instance.default"),
				),
			},
		},
	})
}

func testAccCcev2InstanceConfig(name string) string {
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
  name        = var.name
  vpc_id 	  = baiducloud_vpc.default.id
  description = "created by terraform"
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
data "baiducloud_images" "default" {
  image_type = "System"
  name_regex = "7.5.*"
  os_name    = "CentOS"
}
resource "baiducloud_ccev2_instance_group" "ccev2_instance_group_1" {
  spec {
    cluster_id = baiducloud_ccev2_cluster.default.id
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
data "baiducloud_ccev2_instance_group_instances" "default" {
  cluster_id = baiducloud_ccev2_cluster.default.id
  instance_group_id = baiducloud_ccev2_instance_group.ccev2_instance_group_1.id
}

resource "baiducloud_ccev2_instance" "default" {
  cluster_id        = baiducloud_ccev2_cluster.default.id
  instance_id       = data.baiducloud_ccev2_instance_group_instances.default.instance_list.0.instance_spec.0.cce_instance_id
  spec {
    cce_instance_priority = 1
  }
}
`, name)
}

func testAccCcev2UpdateConfig(name string) string {
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
data "baiducloud_images" "default" {
  image_type = "System"
  name_regex = "7.5.*"
  os_name    = "CentOS"
}
resource "baiducloud_ccev2_instance_group" "ccev2_instance_group_1" {
  spec {
    cluster_id = baiducloud_ccev2_cluster.default.id
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
data "baiducloud_ccev2_instance_group_instances" "default" {
  cluster_id = baiducloud_ccev2_cluster.default.id
  instance_group_id = baiducloud_ccev2_instance_group.ccev2_instance_group_1.id
}

resource "baiducloud_ccev2_instance" "default" {
  cluster_id        = baiducloud_ccev2_cluster.default.id
  instance_id       = data.baiducloud_ccev2_instance_group_instances.default.instance_list.0.instance_spec.0.cce_instance_id
  spec {
    cce_instance_priority = 12
  }
}
`, name)
}
