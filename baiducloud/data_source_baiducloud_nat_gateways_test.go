package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
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
				Config: testAccNatGatewaysDataSourceConfigForNat(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccNatGatewaysDataSourceName),
					resource.TestCheckResourceAttr(testAccNatGatewaysDataSourceName, testAccNatGatewaysDataSourceAttrKeyPrefix+"name", testAccNatGatewayResourceAttrName),
					resource.TestCheckResourceAttr(testAccNatGatewaysDataSourceName, testAccNatGatewaysDataSourceAttrKeyPrefix+"spec", "medium"),
					resource.TestCheckResourceAttrSet(testAccNatGatewaysDataSourceName, testAccNatGatewaysDataSourceAttrKeyPrefix+"id"),
					resource.TestCheckResourceAttrSet(testAccNatGatewaysDataSourceName, testAccNatGatewaysDataSourceAttrKeyPrefix+"vpc_id"),
					resource.TestCheckResourceAttrSet(testAccNatGatewaysDataSourceName, testAccNatGatewaysDataSourceAttrKeyPrefix+"status"),
					resource.TestCheckResourceAttr(testAccNatGatewaysDataSourceName, testAccNatGatewaysDataSourceAttrKeyPrefix+"eips.#", "0"),
					resource.TestCheckResourceAttr(testAccNatGatewaysDataSourceName, testAccNatGatewaysDataSourceAttrKeyPrefix+"payment_timing", "Postpaid"),
				),
			},
			{
				Config: testAccNatGatewaysDataSourceConfigForAll(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccNatGatewaysDataSourceName),
					resource.TestCheckResourceAttr(testAccNatGatewaysDataSourceName, testAccNatGatewaysDataSourceAttrKeyPrefix+"name", testAccNatGatewayResourceAttrName),
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

func testAccNatGatewaysDataSourceConfigForNat() string {
	return fmt.Sprintf(`
resource "baiducloud_vpc" "default" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

data "baiducloud_zones" "default" {}

resource "baiducloud_subnet" "default" {
  name      = "%s"
  zone_name = data.baiducloud_zones.default.zones.0.zone_name
  cidr      = "192.168.1.0/24"
  vpc_id    = baiducloud_vpc.default.id
}

resource "baiducloud_nat_gateway" "default" {
  name   = "%s"
  vpc_id = baiducloud_vpc.default.id
  spec = "medium"
  billing = {
    payment_timing = "Postpaid"
  }
  depends_on = ["baiducloud_subnet.default"]
}

data "baiducloud_nat_gateways" "default" {
  nat_id = baiducloud_nat_gateway.default.id
}
`, BaiduCloudTestResourceAttrNamePrefix+"VPC", BaiduCloudTestResourceAttrNamePrefix+"Subnet", testAccNatGatewayResourceAttrName)
}

func testAccNatGatewaysDataSourceConfigForAll() string {
	return fmt.Sprintf(`
resource "baiducloud_vpc" "default" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

data "baiducloud_zones" "default" {}

resource "baiducloud_subnet" "default" {
  name      = "%s"
  zone_name = data.baiducloud_zones.default.zones.0.zone_name
  cidr      = "192.168.1.0/24"
  vpc_id    = baiducloud_vpc.default.id
}

resource "baiducloud_nat_gateway" "default" {
  name    = "%s"
  vpc_id  = baiducloud_vpc.default.id
  spec    = "medium"
  billing = {
    payment_timing = "Postpaid"
  }
  depends_on = ["baiducloud_subnet.default"]
}

data "baiducloud_nat_gateways" "default" {
  vpc_id = baiducloud_vpc.default.id
}
`, BaiduCloudTestResourceAttrNamePrefix+"VPC", BaiduCloudTestResourceAttrNamePrefix+"Subnet", testAccNatGatewayResourceAttrName)
}
