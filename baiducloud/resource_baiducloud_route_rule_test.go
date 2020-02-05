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
				Config: testAccRouteRuleConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccRouteRuleResourceName),
					resource.TestCheckResourceAttr(testAccRouteRuleResourceName, "source_address", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(testAccRouteRuleResourceName, "destination_address", "192.168.1.0/24"),
					resource.TestCheckResourceAttr(testAccRouteRuleResourceName, "next_hop_type", "custom"),
					resource.TestCheckResourceAttr(testAccRouteRuleResourceName, "description", "baiducloud route rule created by terraform"),
					resource.TestCheckResourceAttrSet(testAccRouteRuleResourceName, "route_table_id"),
					resource.TestCheckResourceAttrSet(testAccRouteRuleResourceName, "next_hop_id"),
				),
			},
			{
				Config: testAccRouteRuleConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccRouteRuleResourceName),
					resource.TestCheckResourceAttr(testAccRouteRuleResourceName, "source_address", "192.168.2.0/24"),
					resource.TestCheckResourceAttr(testAccRouteRuleResourceName, "destination_address", "192.168.3.0/24"),
					resource.TestCheckResourceAttr(testAccRouteRuleResourceName, "next_hop_type", "custom"),
					resource.TestCheckResourceAttr(testAccRouteRuleResourceName, "description", "test route rule update"),
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

func testAccRouteRuleConfig() string {
	return fmt.Sprintf(`
data "baiducloud_specs" "default" {}

data "baiducloud_images" "default" {}

data "baiducloud_zones" "default" {}

data "baiducloud_security_groups" "default" {
  vpc_id = baiducloud_vpc.default.id
}

resource "baiducloud_vpc" "default" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "baiducloud_subnet" "default" {
  name        = "%s"
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.1.0/24"
  description = "subnet created by terraform"
  vpc_id      = baiducloud_vpc.default.id
}

resource "baiducloud_instance" "default" {
  name                  = "%s"
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

resource "%s" "%s" {
  route_table_id      = baiducloud_vpc.default.route_table_id
  source_address      = "192.168.0.0/24"
  destination_address = "192.168.1.0/24"
  next_hop_type       = "custom"
  next_hop_id         = baiducloud_instance.default.id
  description         = "baiducloud route rule created by terraform"
}
`, BaiduCloudTestResourceAttrNamePrefix+"VPC", BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceAttrNamePrefix+"BCC", testAccRouteRuleResourceType, BaiduCloudTestResourceName)
}

func testAccRouteRuleConfigUpdate() string {
	return fmt.Sprintf(`
data "baiducloud_specs" "default" {}

data "baiducloud_images" "default" {}

data "baiducloud_zones" "default" {}

data "baiducloud_security_groups" "default" {
  vpc_id = baiducloud_vpc.default.id
}

resource "baiducloud_vpc" "default" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "baiducloud_subnet" "default" {
  name        = "%s"
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.1.0/24"
  description = "subnet created by terraform"
  vpc_id      = baiducloud_vpc.default.id
}

resource "baiducloud_instance" "default" {
  name                  = "%s"
  image_id              = data.baiducloud_images.default.images.0.id
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  billing = {
    payment_timing = "Postpaid"
  }
  subnet_id       = baiducloud_subnet.default.id
  security_groups = [data.baiducloud_security_groups.default.security_groups.0.id]
}

resource "%s" "%s" {
  route_table_id      = baiducloud_vpc.default.route_table_id
  source_address      = "192.168.2.0/24"
  destination_address = "192.168.3.0/24"
  next_hop_type       = "custom"
  next_hop_id         = baiducloud_instance.default.id
  description         = "test route rule update"
}
`, BaiduCloudTestResourceAttrNamePrefix+"VPC", BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceAttrNamePrefix+"BCC", testAccRouteRuleResourceType, BaiduCloudTestResourceName)
}
