package baiducloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccNatSnatRuleResourceType = "baiducloud_nat_snat_rule"
	testAccNatSnatRuleResourceName = testAccNatSnatRuleResourceType + "." + BaiduCloudTestResourceName
)

//lintignore:AT003
func TestAccBaiduCloudNatSnatRule(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccNatSnatRuleDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccNatSnatRuleConfig(BaiduCloudTestResourceTypeNameNatSnatRule),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccNatSnatRuleResourceName),
					resource.TestCheckResourceAttr(testAccNatSnatRuleResourceName, "rule_name", BaiduCloudTestResourceTypeNameNatSnatRule),
					resource.TestCheckResourceAttr(testAccNatSnatRuleResourceName, "public_ips_address.#", "1"),
					resource.TestCheckResourceAttr(testAccNatSnatRuleResourceName, "source_cidr", "192.168.1.0/24"),
					resource.TestCheckResourceAttrSet(testAccNatSnatRuleResourceName, "nat_id"),
				),
			},
			{
				Config: testAccNatSnatRuleConfigUpdate(BaiduCloudTestResourceTypeNameNatSnatRule),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccNatSnatRuleResourceName),
					resource.TestCheckResourceAttr(testAccNatSnatRuleResourceName, "rule_name", BaiduCloudTestResourceTypeNameNatSnatRule+"-update"),
					resource.TestCheckResourceAttr(testAccNatSnatRuleResourceName, "public_ips_address.#", "1"),
					resource.TestCheckResourceAttr(testAccNatSnatRuleResourceName, "source_cidr", "192.168.1.0/24"),
					resource.TestCheckResourceAttrSet(testAccNatSnatRuleResourceName, "nat_id"),
				),
			},
		},
	})
}

func testAccNatSnatRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	vpcService := &VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccNatSnatRuleResourceType {
			continue
		}

		natId := rs.Primary.Attributes["nat_id"]
		ruleId := rs.Primary.Attributes["rule_id"]
		snatRules, err := vpcService.ListAllNatSnatRulesWithNatID(natId)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		for _, p := range snatRules {
			if p.RuleId == ruleId {
				return WrapError(Error("NatSnatRule still exist"))
			}
		}

	}

	return nil
}

func testAccNatSnatRuleConfig(name string) string {
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
  spec   = "medium"
  billing = {
    payment_timing = "Postpaid"
  }
  depends_on = [baiducloud_subnet.default]
}

resource "baiducloud_nat_snat_rule" "default" {
  nat_id = baiducloud_nat_gateway.default.id
  rule_name = var.name
  public_ips_address = ["100.88.14.90"]
  source_cidr = "192.168.1.0/24"
}
`, name)
}

func testAccNatSnatRuleConfigUpdate(name string) string {
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

resource "baiducloud_eip" "default" {
  name              = var.name
  bandwidth_in_mbps = 1
  payment_timing    = "Postpaid"
  billing_method    = "ByTraffic"
}

resource "baiducloud_nat_gateway" "default" {
  name   = var.name
  vpc_id = baiducloud_vpc.default.id
  spec   = "medium"
  billing = {
    payment_timing = "Postpaid"
  }
  depends_on = ["baiducloud_subnet.default"]
}

resource "baiducloud_eip_association" "default" {
  eip           = baiducloud_eip.default.id
  instance_type = "NAT"
  instance_id   = baiducloud_nat_gateway.default.id
}

resource "baiducloud_nat_snat_rule" "default" {
  nat_id = baiducloud_nat_gateway.default.id
  rule_name = "%s"
  public_ips_address = ["100.88.14.90"]
  source_cidr = "192.168.1.0/24"
}
`, name, name+"-update")
}
