package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const (
	testAccRouteRulesDataSourceName          = "data.baiducloud_route_rules.default"
	testAccRouteRulesDataSourceAttrKeyPrefix = "route_rules.0."
)

func TestAccBaiduCloudRouteRulesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccRouteRulesDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccRouteRulesDataSourceName),
					resource.TestCheckResourceAttrSet(testAccRouteRulesDataSourceName, testAccRouteRulesDataSourceAttrKeyPrefix+"route_rule_id"),
					resource.TestCheckResourceAttrSet(testAccRouteRulesDataSourceName, testAccRouteRulesDataSourceAttrKeyPrefix+"route_table_id"),
					resource.TestCheckResourceAttrSet(testAccRouteRulesDataSourceName, testAccRouteRulesDataSourceAttrKeyPrefix+"next_hop_id"),
					resource.TestCheckResourceAttr(testAccRouteRulesDataSourceName, testAccRouteRulesDataSourceAttrKeyPrefix+"source_address", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(testAccRouteRulesDataSourceName, testAccRouteRulesDataSourceAttrKeyPrefix+"destination_address", "192.168.1.0/24"),
					resource.TestCheckResourceAttr(testAccRouteRulesDataSourceName, testAccRouteRulesDataSourceAttrKeyPrefix+"next_hop_type", "custom"),
					resource.TestCheckResourceAttr(testAccRouteRulesDataSourceName, testAccRouteRulesDataSourceAttrKeyPrefix+"description", "baiducloud route rule created by terraform"),
				),
			},
		},
	})
}

func testAccRouteRulesDataSourceConfig() string {
	return fmt.Sprintf(`
data "baiducloud_specs" "default" {}

data "baiducloud_images" "default" {}

data "baiducloud_zones" "default" {}

data "baiducloud_security_groups" "default" {
  vpc_id = "${baiducloud_vpc.default.id}"
}

resource "baiducloud_vpc" "default" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "baiducloud_subnet" "default" {
  name = "%s"
  zone_name = "${data.baiducloud_zones.default.zones.1.zone_name}"
  cidr = "192.168.1.0/24"
  description = "subnet created by terraform"
  vpc_id = "${baiducloud_vpc.default.id}"
}

resource "baiducloud_instance" "default" {
  name = "%s"
  image_id = "${data.baiducloud_images.default.images.0.id}"
  cpu_count = "${data.baiducloud_specs.default.specs.0.cpu_count}"
  memory_capacity_in_gb = "${data.baiducloud_specs.default.specs.0.memory_size_in_gb}"
  billing = {
    payment_timing = "Postpaid"
  }
  availability_zone = "${data.baiducloud_zones.default.zones.1.zone_name}"
  subnet_id = "${baiducloud_subnet.default.id}"
  security_groups = ["${data.baiducloud_security_groups.default.security_groups.0.id}"]
}

resource "%s" "%s" {
  route_table_id = "${baiducloud_vpc.default.route_table_id}"
  source_address = "192.168.0.0/24"
  destination_address = "192.168.1.0/24"
  next_hop_type = "custom"
  next_hop_id = "${baiducloud_instance.default.id}"
  description = "baiducloud route rule created by terraform"
}

data "baiducloud_route_rules" "default" {
  route_table_id = "${baiducloud_vpc.default.route_table_id}"
  route_rule_id = "${%s.%s.id}"
}
`, BaiduCloudTestResourceAttrNamePrefix+"VPC", BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceAttrNamePrefix+"BCC", testAccRouteRuleResourceType, BaiduCloudTestResourceName,
		testAccRouteRuleResourceType, BaiduCloudTestResourceName)
}
