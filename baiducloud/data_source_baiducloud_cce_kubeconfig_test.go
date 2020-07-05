package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const (
	testAccCceKubeconfigDataSourceName          = "data.baiducloud_cce_kubeconfig.default"
	testAccCceKubeconfigDataSourceAttrKeyPrefix = "data"
)

//lintignore:AT003
func TestAccBaiduCloudCceKubeconfigDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubeconfigDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccCceKubeconfigDataSourceName),
					resource.TestCheckResourceAttrSet(testAccCceKubeconfigDataSourceName, testAccCceKubeconfigDataSourceAttrKeyPrefix),
				),
			},
		},
	})
}

func testAccKubeconfigDataSourceConfig() string {
	return fmt.Sprintf(`
data "baiducloud_zones" "defaultA" {
  name_regex = ".*a$"
}

data "baiducloud_images" "default" {
  image_type = "System"
  name_regex = "7.5.*"
  os_name    = "CentOS"
}

resource "baiducloud_vpc" "default" {
  name        = "%s"
  description = var.description
  cidr        = "192.168.0.0/16"
}

resource "baiducloud_subnet" "defaultA" {
  name        = "%s"
  zone_name   = data.baiducloud_zones.defaultA.zones.0.zone_name
  cidr        = "192.168.1.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "terraform create"
}

resource "baiducloud_security_group" "defualt" {
  name   = "%s"
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

resource "baiducloud_cce_cluster" "default_independent" {
  cluster_name        = "%s"
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
    cluster_uuid = baiducloud_cce_cluster.default_managed.id
}
`, BaiduCloudTestResourceAttrNamePrefix+"VPC", BaiduCloudTestResourceAttrNamePrefix+"SubnetA",
		BaiduCloudTestResourceAttrNamePrefix+"SG", BaiduCloudTestResourceAttrNamePrefix+"CCE")
}
