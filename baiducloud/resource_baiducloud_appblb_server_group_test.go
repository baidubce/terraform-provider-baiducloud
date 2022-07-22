package baiducloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/baidubce/bce-sdk-go/services/appblb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccAppBLBServerGroupResourceType = "baiducloud_appblb_server_group"
	testAccAppBLBServerGroupResourceName = testAccAppBLBServerGroupResourceType + "." + BaiduCloudTestResourceName
)

func init() {
	resource.AddTestSweepers(testAccAppBLBServerGroupResourceType, &resource.Sweeper{
		Name: testAccAppBLBServerGroupResourceType,
		F:    testSweepAppBLBServerGroups,
		Dependencies: []string{
			testAccInstanceResourceType,
			testAccAppBLBResourceType,
			testAccVPCResourceType,
		},
	})
}

func testSweepAppBLBServerGroups(region string) error {
	log.Printf("[INFO] Skipping AppBLB Server Group,Nothing to do)")
	return nil
}

func TestAccBaiduCloudAppBLBServerGroup_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccAppBLBServerGroupDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccAppBLBServerGroupConfig(BaiduCloudTestResourceTypeNameAppblbServerGroup),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBServerGroupResourceName),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupResourceName, "name", BaiduCloudTestResourceTypeNameAppblbServerGroup),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupResourceName, "port_list.#", "1"),
					resource.TestCheckResourceAttrSet(testAccAppBLBServerGroupResourceName, "blb_id"),
					resource.TestCheckResourceAttrSet(testAccAppBLBServerGroupResourceName, "status"),
				),
			},
			{
				Config: testAccAppBLBServerGroupConfigUpdate(BaiduCloudTestResourceTypeNameAppblbServerGroup),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBServerGroupResourceName),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupResourceName, "name", BaiduCloudTestResourceTypeNameAppblbServerGroup),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupResourceName, "port_list.#", "2"),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupResourceName, "backend_server_list.#", "1"),
					resource.TestCheckResourceAttrSet(testAccAppBLBServerGroupResourceName, "blb_id"),
					resource.TestCheckResourceAttrSet(testAccAppBLBServerGroupResourceName, "status"),
				),
			},
			{
				Config: testAccAppBLBServerGroupConfigUpdate2(BaiduCloudTestResourceTypeNameAppblbServerGroup),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBServerGroupResourceName),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupResourceName, "name", BaiduCloudTestResourceTypeNameAppblbServerGroup),
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
				Config: testAccAppBLBServerGroupConfigUpdate2(BaiduCloudTestResourceTypeNameAppblbServerGroup),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBServerGroupResourceName),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupResourceName, "name", BaiduCloudTestResourceTypeNameAppblbServerGroup),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupResourceName, "port_list.#", "1"),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupResourceName, "backend_server_list.#", "1"),
					resource.TestCheckResourceAttrSet(testAccAppBLBServerGroupResourceName, "blb_id"),
					resource.TestCheckResourceAttrSet(testAccAppBLBServerGroupResourceName, "status"),
				),
			},
			{
				Config: testAccAppBLBServerGroupConfigUpdate3(BaiduCloudTestResourceTypeNameAppblbServerGroup),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBServerGroupResourceName),
					resource.TestCheckResourceAttr(testAccAppBLBServerGroupResourceName, "name", BaiduCloudTestResourceTypeNameAppblbServerGroup),
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
			Name:         BaiduCloudTestResourceTypeName,
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

func testAccAppBLBServerGroupConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "baiducloud_zones" "default" {
  name_regex = ".*e$"
}

data "baiducloud_specs" "default" {}

data "baiducloud_images" "default" {
  image_type = "System"
}

resource "baiducloud_instance" "default" {
  name                  = "${var.name}"
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
  name        = "${var.name}"
  description = "created by terraform"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = "${var.name}"
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created by terraform"
}

resource "baiducloud_security_group" "default" {
  name        = "${var.name}"
  description = "created by terraform"
  vpc_id      = baiducloud_vpc.default.id
}

resource "baiducloud_appblb" "default" {
  depends_on  = [baiducloud_instance.default]
  name        = "${var.name}"
  description = "created by terraform"
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}

resource "baiducloud_appblb_server_group" "default" {
  name        = "${var.name}"
  description = "created by terraform"
  blb_id      = baiducloud_appblb.default.id

  port_list {
    port = 66
    type = "TCP"
    health_check = "TCP"
  }
}
`, name)
}

func testAccAppBLBServerGroupConfigUpdate(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "baiducloud_zones" "default" {
  name_regex = ".*e$"
}

data "baiducloud_specs" "default" {}

data "baiducloud_images" "default" {
  image_type = "System"
}

resource "baiducloud_vpc" "default" {
  name        = "${var.name}"
  description = "created by terraform"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = "${var.name}"
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created by terraform"
}

resource "baiducloud_security_group" "default" {
  name        = "${var.name}"
  description = "created by terraform"
  vpc_id      = baiducloud_vpc.default.id
}

resource "baiducloud_instance" "default" {
  name                  = "${var.name}"
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
  name        = "${var.name}"
  description = "created by terraform"
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}

resource "baiducloud_appblb_server_group" "default" {
  name        = "${var.name}"
  description = "created by terraform"
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
`, name)
}

func testAccAppBLBServerGroupConfigUpdate2(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "baiducloud_zones" "default" {
  name_regex = ".*e$"
}

data "baiducloud_specs" "default" {}

data "baiducloud_images" "default" {
  image_type = "System"
}

resource "baiducloud_vpc" "default" {
  name        = "${var.name}"
  description = "created by terraform"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = "${var.name}"
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created by terraform"
}

resource "baiducloud_security_group" "default" {
  name        = "${var.name}"
  description = "created by terraform"
  vpc_id      = baiducloud_vpc.default.id
}

resource "baiducloud_instance" "default" {
  name                  = "${var.name}"
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
  name        = "${var.name}"
  description = "created by terraform"
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}

resource "baiducloud_appblb_server_group" "default" {
  name        = "${var.name}"
  description = "created by terraform"
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
`, name)
}

func testAccAppBLBServerGroupConfigUpdate3(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "baiducloud_zones" "default" {
  name_regex = ".*e$"
}

data "baiducloud_specs" "default" {}

data "baiducloud_images" "default" {
  image_type = "System"
}

resource "baiducloud_vpc" "default" {
  name        = "${var.name}"
  description = "created by terraform"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = "${var.name}"
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created by terraform"
}

resource "baiducloud_security_group" "default" {
  name        = "${var.name}"
  description = "created by terraform"
  vpc_id      = baiducloud_vpc.default.id
}

resource "baiducloud_instance" "default" {
  name                  = "${var.name}"
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
  name        = "${var.name}"
  description = "created by terraform"
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}

resource "baiducloud_appblb_server_group" "default" {
  name        = "${var.name}"
  description = "created by terraform"
  blb_id      = baiducloud_appblb.default.id

  port_list {
    port                    = 77
    type                    = "UDP"
    health_check            = "UDP"
    udp_health_check_string = "baidunew.com"
  }
}
`, name)
}
