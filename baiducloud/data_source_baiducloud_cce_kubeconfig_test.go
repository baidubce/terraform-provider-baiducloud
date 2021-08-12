package baiducloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

const (
	testAccCceKubeconfigDataSourceName          = "data.baiducloud_cce_kubeconfig.default"
	testAccCceKubeconfigDataSourceAttrKeyPrefix = "data"
)

//lintignore:AT003
func testAccBaiduCloudCceKubeconfigDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubeconfigDataSourceConfig(BaiduCloudTestResourceTypeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccCceKubeconfigDataSourceName),
					resource.TestCheckResourceAttrSet(testAccCceKubeconfigDataSourceName, testAccCceKubeconfigDataSourceAttrKeyPrefix),
				),
			},
		},
	})
}

func testAccKubeconfigDataSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "baiducloud_zones" "defaultA" {
  name_regex = ".*e$"
}

data "baiducloud_images" "default" {
  image_type = "System"
  name_regex = "7.5.*"
  os_name    = "CentOS"
}

resource "baiducloud_vpc" "default" {
  name        = var.name
  description = "created by terraform"
  cidr        = "192.168.0.0/16"
}

resource "baiducloud_subnet" "defaultA" {
  name        = var.name
  zone_name   = data.baiducloud_zones.defaultA.zones.0.zone_name
  cidr        = "192.168.1.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created by terraform"
}

resource "baiducloud_security_group" "defualt" {
  name   = var.name
  vpc_id = baiducloud_vpc.default.id
}

resource "baiducloud_security_group_rule" "default" {
  security_group_id = baiducloud_security_group.defualt.id
  remark            = "remark"
  protocol          = "all"
  port_range        = "1-65535"
  direction         = "ingress"
}

resource "baiducloud_security_group_rule" "default2" {
  security_group_id = baiducloud_security_group.defualt.id
  remark            = "remark"
  protocol          = "all"
  port_range        = "1-65535"
  direction         = "egress"
}

resource "baiducloud_cce_cluster" "default" {
  cluster_name        = var.name
  main_available_zone = "zoneA"
  version             = "1.13.10"
  container_net       = "172.16.0.0/16"

  worker_config {
    count = {
      "zoneA" : 1
    }

    instance_type = "10"
    cpu           = 1
    memory        = 2
    subnet_uuid = {
      "zoneA" : baiducloud_subnet.defaultA.id
    }
    security_group_id = baiducloud_security_group.defualt.id
    product_type      = "postpay"
    image_id          = data.baiducloud_images.default.images.0.id
    image_type        = "common"
  }
}

data "baiducloud_cce_kubeconfig" "default" {
    cluster_uuid = baiducloud_cce_cluster.default.id
}
`, name)
}
