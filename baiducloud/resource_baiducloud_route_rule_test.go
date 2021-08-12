package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccRouteRuleResourceType = "baiducloud_route_rule"
	testAccRouteRuleResourceName = testAccRouteRuleResourceType + "." + BaiduCloudTestResourceName
)

//lintignore:AT003
func TestAccBaiduCloudRouteRule(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccRouteRuleDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccRouteRuleConfig(BaiduCloudTestResourceTypeNameRouteRule),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccRouteRuleResourceName),
					resource.TestCheckResourceAttr(testAccRouteRuleResourceName, "source_address", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(testAccRouteRuleResourceName, "destination_address", "192.168.1.0/24"),
					resource.TestCheckResourceAttr(testAccRouteRuleResourceName, "next_hop_type", "custom"),
					resource.TestCheckResourceAttr(testAccRouteRuleResourceName, "description", "created by terraform"),
					resource.TestCheckResourceAttrSet(testAccRouteRuleResourceName, "route_table_id"),
					resource.TestCheckResourceAttrSet(testAccRouteRuleResourceName, "next_hop_id"),
				),
			},
			{
				Config: testAccRouteRuleConfigUpdate(BaiduCloudTestResourceTypeNameRouteRule),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccRouteRuleResourceName),
					resource.TestCheckResourceAttr(testAccRouteRuleResourceName, "source_address", "192.168.2.0/24"),
					resource.TestCheckResourceAttr(testAccRouteRuleResourceName, "destination_address", "192.168.3.0/24"),
					resource.TestCheckResourceAttr(testAccRouteRuleResourceName, "next_hop_type", "custom"),
					resource.TestCheckResourceAttr(testAccRouteRuleResourceName, "description", "created by terraform"),
					resource.TestCheckResourceAttrSet(testAccRouteRuleResourceName, "route_table_id"),
					resource.TestCheckResourceAttrSet(testAccRouteRuleResourceName, "next_hop_id"),
				),
			},
		},
	})
}

func testAccRouteRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	vpcService := &VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRouteRuleResourceType {
			continue
		}

		routeTableID := rs.Primary.Attributes["route_table_id"]
		result, err := vpcService.GetRouteTableDetail(routeTableID, "")
		if err != nil {
			if NotFoundError(err) || IsExceptedErrors(err, []string{"BadRequest"}) {
				continue
			}
			return WrapError(err)
		}

		for _, rule := range result.RouteRules {
			if rule.RouteRuleId == rs.Primary.ID {
				return WrapError(Error("Route Rule still exist"))
			}
		}
	}

	return nil
}

func testAccRouteRuleConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "baiducloud_specs" "default" {}

data "baiducloud_images" "default" {}

data "baiducloud_zones" "default" {
  name_regex = ".*e$"
}

data "baiducloud_security_groups" "default" {
  vpc_id = baiducloud_vpc.default.id
}

resource "baiducloud_vpc" "default" {
  name = var.name
  cidr = "192.168.0.0/16"
}

resource "baiducloud_subnet" "default" {
  name        = var.name
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.1.0/24"
  description = "created by terraform"
  vpc_id      = baiducloud_vpc.default.id
}

resource "baiducloud_instance" "default" {
  name                  = var.name
  image_id              = data.baiducloud_images.default.images.0.id
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  billing = {
    payment_timing = "Postpaid"
  }
  availability_zone = data.baiducloud_zones.default.zones.0.zone_name
  subnet_id         = baiducloud_subnet.default.id
  security_groups   = [data.baiducloud_security_groups.default.security_groups.0.id]
}

resource "baiducloud_route_rule" "default" {
  route_table_id      = baiducloud_vpc.default.route_table_id
  source_address      = "192.168.0.0/24"
  destination_address = "192.168.1.0/24"
  next_hop_type       = "custom"
  next_hop_id         = baiducloud_instance.default.id
  description         = "created by terraform"
}
`, name)
}

func testAccRouteRuleConfigUpdate(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "baiducloud_specs" "default" {}

data "baiducloud_images" "default" {}

data "baiducloud_zones" "default" {
  name_regex = ".*e$"
}

data "baiducloud_security_groups" "default" {
  vpc_id = baiducloud_vpc.default.id
}

resource "baiducloud_vpc" "default" {
  name = var.name
  cidr = "192.168.0.0/16"
}

resource "baiducloud_subnet" "default" {
  name        = var.name
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.1.0/24"
  description = "created by terraform"
  vpc_id      = baiducloud_vpc.default.id
}

resource "baiducloud_instance" "default" {
  name                  = var.name
  image_id              = data.baiducloud_images.default.images.0.id
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  billing = {
    payment_timing = "Postpaid"
  }
  subnet_id       = baiducloud_subnet.default.id
  security_groups = [data.baiducloud_security_groups.default.security_groups.0.id]
}

resource "baiducloud_route_rule" "default" {
  route_table_id      = baiducloud_vpc.default.route_table_id
  source_address      = "192.168.2.0/24"
  destination_address = "192.168.3.0/24"
  next_hop_type       = "custom"
  next_hop_id         = baiducloud_instance.default.id
  description         = "created by terraform"
}
`, name)
}
