package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const (
	testAccAppBLBServerGroupsDataSourceName          = "data.baiducloud_appblb_server_groups.default"
	testAccAppBLBServerGroupsDataSourceAttrKeyPrefix = "server_groups.0."
)

//lintignore:AT003
func TestAccBaiduCloudAppBLBServerGroupsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAppBLBServerGroupDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBServerGroupsDataSourceName),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupsDataSourceName, testAccAppBLBServerGroupsDataSourceAttrKeyPrefix+"name", BaiduCloudTestResourceAttrNamePrefix+"ServerGroup"),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupsDataSourceName, testAccAppBLBServerGroupsDataSourceAttrKeyPrefix+"port_list.#", "1"),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupsDataSourceName, testAccAppBLBServerGroupsDataSourceAttrKeyPrefix+"port_list.0.type", "TCP"),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupsDataSourceName, testAccAppBLBServerGroupsDataSourceAttrKeyPrefix+"port_list.0.port", "66"),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupsDataSourceName, testAccAppBLBServerGroupsDataSourceAttrKeyPrefix+"backend_server_list.#", "1"),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupsDataSourceName, testAccAppBLBServerGroupsDataSourceAttrKeyPrefix+"backend_server_list.0.weight", "50"),
					resource.TestCheckResourceAttrSet(testAccAppBLBServerGroupsDataSourceName, testAccAppBLBServerGroupsDataSourceAttrKeyPrefix+"sg_id"),
					resource.TestCheckResourceAttrSet(testAccAppBLBServerGroupsDataSourceName, testAccAppBLBServerGroupsDataSourceAttrKeyPrefix+"status"),
					resource.TestCheckResourceAttrSet(testAccAppBLBServerGroupsDataSourceName, testAccAppBLBServerGroupsDataSourceAttrKeyPrefix+"backend_server_list.0.instance_id"),
				),
			},
		},
	})
}

func testAccAppBLBServerGroupDataSourceConfig() string {
	return fmt.Sprintf(`
data "baiducloud_zones" "default" {}

data "baiducloud_specs" "default" {}

data "baiducloud_images" "default" {
  image_type = "System"
}

resource "baiducloud_vpc" "default" {
  name        = "%s"
  description = "test"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = "%s"
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "test description"
}

resource "baiducloud_security_group" "default" {
  name        = "%s"
  description = "Baidu acceptance test"
  vpc_id      = baiducloud_vpc.default.id
}

resource "baiducloud_instance" "default" {
  name                  = "%s"
  image_id              = data.baiducloud_images.default.images.0.id
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  subnet_id             = baiducloud_subnet.default.id
  security_groups       = [baiducloud_security_group.default.id]

  billing = {
    payment_timing = "Postpaid"
  }
}

resource "baiducloud_appblb" "default" {
  depends_on  = [baiducloud_instance.default]
  name        = "%s"
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}

resource "baiducloud_appblb_server_group" "default" {
  name        = "%s"
  description = "acceptance test"
  blb_id      = baiducloud_appblb.default.id

  backend_server_list {
    instance_id = baiducloud_instance.default.id
    weight      = 50
  }

  port_list {
    port = 66
    type = "TCP"
    health_check = "TCP"
  }
}

data "baiducloud_appblb_server_groups" "default" {
  blb_id = baiducloud_appblb.default.id
  name   = baiducloud_appblb_server_group.default.name

  filter {
    name = "name"
    values = ["test-BaiduAcc*"]
  }
}
`, BaiduCloudTestResourceAttrNamePrefix+"VPC",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceAttrNamePrefix+"SecurityGroup",
		BaiduCloudTestResourceAttrNamePrefix+"Instance",
		BaiduCloudTestResourceAttrNamePrefix+"APPBLB",
		BaiduCloudTestResourceAttrNamePrefix+"ServerGroup")
}
