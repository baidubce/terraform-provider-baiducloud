package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccNatGatewaysDataSourceName          = "data.baiducloud_nat_gateways.default"
	testAccNatGatewaysDataSourceAttrKeyPrefix = "nat_gateways.0."
)

//lintignore:AT003
func TestAccBaiduCloudNatGatewaysDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccNatGatewaysDataSourceConfigForNat(BaiduCloudTestResourceTypeNameNatGateway),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccNatGatewaysDataSourceName),
					resource.TestCheckResourceAttr(testAccNatGatewaysDataSourceName, testAccNatGatewaysDataSourceAttrKeyPrefix+"name", BaiduCloudTestResourceTypeNameNatGateway),
					resource.TestCheckResourceAttr(testAccNatGatewaysDataSourceName, testAccNatGatewaysDataSourceAttrKeyPrefix+"spec", "medium"),
					resource.TestCheckResourceAttrSet(testAccNatGatewaysDataSourceName, testAccNatGatewaysDataSourceAttrKeyPrefix+"id"),
					resource.TestCheckResourceAttrSet(testAccNatGatewaysDataSourceName, testAccNatGatewaysDataSourceAttrKeyPrefix+"vpc_id"),
					resource.TestCheckResourceAttrSet(testAccNatGatewaysDataSourceName, testAccNatGatewaysDataSourceAttrKeyPrefix+"status"),
					resource.TestCheckResourceAttr(testAccNatGatewaysDataSourceName, testAccNatGatewaysDataSourceAttrKeyPrefix+"eips.#", "0"),
					resource.TestCheckResourceAttr(testAccNatGatewaysDataSourceName, testAccNatGatewaysDataSourceAttrKeyPrefix+"payment_timing", "Postpaid"),
				),
			},
			{
				Config: testAccNatGatewaysDataSourceConfigForAll(BaiduCloudTestResourceTypeNameNatGateway),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccNatGatewaysDataSourceName),
					resource.TestCheckResourceAttr(testAccNatGatewaysDataSourceName, testAccNatGatewaysDataSourceAttrKeyPrefix+"name", BaiduCloudTestResourceTypeNameNatGateway),
					resource.TestCheckResourceAttr(testAccNatGatewaysDataSourceName, testAccNatGatewaysDataSourceAttrKeyPrefix+"spec", "medium"),
					resource.TestCheckResourceAttrSet(testAccNatGatewaysDataSourceName, testAccNatGatewaysDataSourceAttrKeyPrefix+"id"),
					resource.TestCheckResourceAttrSet(testAccNatGatewaysDataSourceName, testAccNatGatewaysDataSourceAttrKeyPrefix+"vpc_id"),
					resource.TestCheckResourceAttrSet(testAccNatGatewaysDataSourceName, testAccNatGatewaysDataSourceAttrKeyPrefix+"status"),
					resource.TestCheckResourceAttr(testAccNatGatewaysDataSourceName, testAccNatGatewaysDataSourceAttrKeyPrefix+"eips.#", "0"),
					resource.TestCheckResourceAttr(testAccNatGatewaysDataSourceName, testAccNatGatewaysDataSourceAttrKeyPrefix+"payment_timing", "Postpaid"),
				),
			},
		},
	})
}

func testAccNatGatewaysDataSourceConfigForNat(name string) string {
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

data "baiducloud_nat_gateways" "default" {
  nat_id = baiducloud_nat_gateway.default.id

  filter {
    name = "name"
    values = ["tf-test-acc*"]
  }
}
`, name)
}

func testAccNatGatewaysDataSourceConfigForAll(name string) string {
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
  name    = var.name
  vpc_id  = baiducloud_vpc.default.id
  spec    = "medium"
  billing = {
    payment_timing = "Postpaid"
  }
  depends_on = ["baiducloud_subnet.default"]
}

data "baiducloud_nat_gateways" "default" {
  vpc_id = baiducloud_vpc.default.id

  filter {
    name = "spec"
    values = ["medium"]
  }
}
`, name)
}
