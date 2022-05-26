package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const (
	testAccNatSnatRulesDataSourceName          = "data.baiducloud_nat_snat_rules.default"
	testAccNatSnatRulesDataSourceAttrKeyPrefix = "nat_snat_rules.0."
)

func TestAccBaiduCloudNatSnatRulesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccNatSnatRulesDataSourceConfigForNat(BaiduCloudTestResourceTypeNameNatSnatRule),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccNatSnatRulesDataSourceName),
					resource.TestCheckResourceAttr(testAccNatSnatRulesDataSourceName, testAccNatSnatRulesDataSourceAttrKeyPrefix+"name", BaiduCloudTestResourceTypeNameNatGateway),
					resource.TestCheckResourceAttr(testAccNatSnatRulesDataSourceName, testAccNatSnatRulesDataSourceAttrKeyPrefix+"spec", "medium"),
					resource.TestCheckResourceAttrSet(testAccNatSnatRulesDataSourceName, testAccNatSnatRulesDataSourceAttrKeyPrefix+"id"),
					resource.TestCheckResourceAttrSet(testAccNatSnatRulesDataSourceName, testAccNatSnatRulesDataSourceAttrKeyPrefix+"vpc_id"),
					resource.TestCheckResourceAttrSet(testAccNatSnatRulesDataSourceName, testAccNatSnatRulesDataSourceAttrKeyPrefix+"status"),
					resource.TestCheckResourceAttr(testAccNatSnatRulesDataSourceName, testAccNatSnatRulesDataSourceAttrKeyPrefix+"eips.#", "0"),
					resource.TestCheckResourceAttr(testAccNatSnatRulesDataSourceName, testAccNatSnatRulesDataSourceAttrKeyPrefix+"payment_timing", "Postpaid"),
				),
			},
		},
	})
}

func testAccNatSnatRulesDataSourceConfigForNat(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_vpc" "default" {
  name = var.name
  cidr = "192.168.0.0/16"
}

data "baiducloud_zones" "default" {
  name_regex = ".*e$"
}

resource "baiducloud_subnet" "default" {
  name      = var.name
  zone_name = data.baiducloud_zones.default.zones.0.zone_name
  cidr      = "192.168.1.0/24"
  vpc_id    = baiducloud_vpc.default.id
}

resource "baiducloud_nat_gateway" "default" {
  name   = var.name
  vpc_id = baiducloud_vpc.default.id
  spec = "medium"
  billing = {
    payment_timing = "Postpaid"
  }
  depends_on = ["baiducloud_subnet.default"]
}

data "baiducloud_nat_snat_rules" "default" {
  nat_id = baiducloud_nat_gateway.default.id
}
`, name)
}
