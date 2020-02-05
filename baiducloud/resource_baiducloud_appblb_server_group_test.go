package baiducloud

import (
	"fmt"
	"testing"

	"github.com/baidubce/bce-sdk-go/services/appblb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccAppBLBServerGroupResourceType     = "baiducloud_appblb_server_group"
	testAccAppBLBServerGroupResourceName     = testAccAppBLBServerGroupResourceType + "." + BaiduCloudTestResourceName
	testAccAppBLBServerGroupResourceAttrName = BaiduCloudTestResourceAttrNamePrefix + "APPBLBServerGroup"
)

func TestAccBaiduCloudAppBLBServerGroup_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccAppBLBServerGroupDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccAppBLBServerGroupConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBServerGroupResourceName),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupResourceName, "name", testAccAppBLBServerGroupResourceAttrName),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupResourceName, "port_list.#", "1"),
					resource.TestCheckResourceAttrSet(testAccAppBLBServerGroupResourceName, "blb_id"),
					resource.TestCheckResourceAttrSet(testAccAppBLBServerGroupResourceName, "status"),
				),
			},
			{
				Config: testAccAppBLBServerGroupConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBServerGroupResourceName),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupResourceName, "name", testAccAppBLBServerGroupResourceAttrName),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupResourceName, "port_list.#", "2"),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupResourceName, "backend_server_list.#", "1"),
					resource.TestCheckResourceAttrSet(testAccAppBLBServerGroupResourceName, "blb_id"),
					resource.TestCheckResourceAttrSet(testAccAppBLBServerGroupResourceName, "status"),
				),
			},
			{
				Config: testAccAppBLBServerGroupConfigUpdate2(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBServerGroupResourceName),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupResourceName, "name", testAccAppBLBServerGroupResourceAttrName),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupResourceName, "port_list.#", "1"),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupResourceName, "backend_server_list.#", "1"),
					resource.TestCheckResourceAttrSet(testAccAppBLBServerGroupResourceName, "blb_id"),
					resource.TestCheckResourceAttrSet(testAccAppBLBServerGroupResourceName, "status"),
				),
			},
		},
	})
}

func TestAccBaiduCloudAppBLBServerGroup_Rs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccAppBLBServerGroupDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccAppBLBServerGroupConfigUpdate2(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBServerGroupResourceName),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupResourceName, "name", testAccAppBLBServerGroupResourceAttrName),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupResourceName, "port_list.#", "1"),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupResourceName, "backend_server_list.#", "1"),
					resource.TestCheckResourceAttrSet(testAccAppBLBServerGroupResourceName, "blb_id"),
					resource.TestCheckResourceAttrSet(testAccAppBLBServerGroupResourceName, "status"),
				),
			},
			{
				Config: testAccAppBLBServerGroupConfigUpdate3(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBServerGroupResourceName),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupResourceName, "name", testAccAppBLBServerGroupResourceAttrName),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupResourceName, "port_list.#", "1"),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupResourceName, "backend_server_list.#", "0"),
					resource.TestCheckResourceAttrSet(testAccAppBLBServerGroupResourceName, "blb_id"),
					resource.TestCheckResourceAttrSet(testAccAppBLBServerGroupResourceName, "status"),
				),
			},
		},
	})
}

func testAccAppBLBServerGroupDestory(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	appblbService := APPBLBService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccAppBLBServerGroupResourceType {
			continue
		}

		listArgs := &appblb.DescribeAppServerGroupArgs{
			Name:         testAccAppBLBServerGroupResourceAttrName,
			ExactlyMatch: true,
		}

		raw, err := appblbService.ListAllServerGroups(rs.Primary.Attributes["blb_id"], listArgs)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		for _, sg := range raw {
			if sg["id"] == rs.Primary.ID {
				return WrapError(Error("APPBLB ServerGroup still exist"))
			}
		}
	}

	return nil
}

func testAccAppBLBServerGroupConfig() string {
	return fmt.Sprintf(`
data "baiducloud_zones" "default" {}

data "baiducloud_specs" "default" {}

data "baiducloud_images" "default" {
  image_type = "System"
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

resource "baiducloud_appblb" "default" {
  depends_on  = [baiducloud_instance.default]
  name        = "%s"
  description = ""
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}

resource "%s" "%s" {
  name        = "%s"
  description = "acceptance test"
  blb_id      = baiducloud_appblb.default.id

  port_list {
    port = 66
    type = "TCP"
    health_check = "TCP"
  }
}
`, BaiduCloudTestResourceAttrNamePrefix+"Instance",
		BaiduCloudTestResourceAttrNamePrefix+"VPC",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceAttrNamePrefix+"SecurityGroup",
		BaiduCloudTestResourceAttrNamePrefix+"APPBLB",
		testAccAppBLBServerGroupResourceType,
		BaiduCloudTestResourceName,
		testAccAppBLBServerGroupResourceAttrName)
}

func testAccAppBLBServerGroupConfigUpdate() string {
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
  description = ""
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}

resource "%s" "%s" {
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

  port_list {
    port         = 77
    type         = "UDP"
    health_check = "UDP"
    udp_health_check_string = "baidu.com"
  }
}
`, BaiduCloudTestResourceAttrNamePrefix+"VPC",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceAttrNamePrefix+"SecurityGroup",
		BaiduCloudTestResourceAttrNamePrefix+"Instance",
		BaiduCloudTestResourceAttrNamePrefix+"APPBLB",
		testAccAppBLBServerGroupResourceType,
		BaiduCloudTestResourceName,
		testAccAppBLBServerGroupResourceAttrName)
}

func testAccAppBLBServerGroupConfigUpdate2() string {
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
  description = ""
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}

resource "%s" "%s" {
  name        = "%s"
  description = "acceptance test"
  blb_id      = baiducloud_appblb.default.id

  backend_server_list {
    instance_id = baiducloud_instance.default.id
    weight      = 60
  }

  port_list {
    port                    = 77
    type                    = "UDP"
    health_check            = "UDP"
    udp_health_check_string = "baidunew.com"
  }
}
`, BaiduCloudTestResourceAttrNamePrefix+"VPC",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceAttrNamePrefix+"SecurityGroup",
		BaiduCloudTestResourceAttrNamePrefix+"Instance",
		BaiduCloudTestResourceAttrNamePrefix+"APPBLB",
		testAccAppBLBServerGroupResourceType,
		BaiduCloudTestResourceName,
		testAccAppBLBServerGroupResourceAttrName)
}

func testAccAppBLBServerGroupConfigUpdate3() string {
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
  description = ""
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}

resource "%s" "%s" {
  name        = "%s"
  description = "acceptance test"
  blb_id      = baiducloud_appblb.default.id

  port_list {
    port                    = 77
    type                    = "UDP"
    health_check            = "UDP"
    udp_health_check_string = "baidunew.com"
  }
}
`, BaiduCloudTestResourceAttrNamePrefix+"VPC",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceAttrNamePrefix+"SecurityGroup",
		BaiduCloudTestResourceAttrNamePrefix+"Instance",
		BaiduCloudTestResourceAttrNamePrefix+"APPBLB",
		testAccAppBLBServerGroupResourceType,
		BaiduCloudTestResourceName,
		testAccAppBLBServerGroupResourceAttrName)
}
