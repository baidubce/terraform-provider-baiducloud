package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccNatGatewayResourceType = "baiducloud_nat_gateway"
	testAccNatGatewayResourceName = testAccNatGatewayResourceType + "." + BaiduCloudTestResourceName
)

//lintignore:AT003
func TestAccBaiduCloudNatGateway(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccNatGatewayDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccNatGatewayConfig(BaiduCloudTestResourceTypeNameNatGateway),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccNatGatewayResourceName),
					resource.TestCheckResourceAttr(testAccNatGatewayResourceName, "name", BaiduCloudTestResourceTypeNameNatGateway),
					resource.TestCheckResourceAttr(testAccNatGatewayResourceName, "spec", "medium"),
					resource.TestCheckResourceAttrSet(testAccNatGatewayResourceName, "id"),
					resource.TestCheckResourceAttrSet(testAccNatGatewayResourceName, "vpc_id"),
					resource.TestCheckResourceAttrSet(testAccNatGatewayResourceName, "status"),
				),
			},
			{
				ResourceName:      testAccNatGatewayResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNatGatewayConfigUpdate(BaiduCloudTestResourceTypeNameNatGateway),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccNatGatewayResourceName),
					resource.TestCheckResourceAttr(testAccNatGatewayResourceName, "name", BaiduCloudTestResourceTypeNameNatGateway+"-update"),
					resource.TestCheckResourceAttr(testAccNatGatewayResourceName, "spec", "medium"),
					resource.TestCheckResourceAttrSet(testAccNatGatewayResourceName, "id"),
					resource.TestCheckResourceAttrSet(testAccNatGatewayResourceName, "vpc_id"),
					resource.TestCheckResourceAttrSet(testAccNatGatewayResourceName, "status"),
				),
			},
		},
	})
}

func testAccNatGatewayDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	vpcService := &VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccNatGatewayResourceType {
			continue
		}

		_, err := vpcService.GetNatGatewayDetail(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		return WrapError(Error("NatGateway still exist"))
	}

	return nil
}

func testAccNatGatewayConfig(name string) string {
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
`, name)
}

func testAccNatGatewayConfigUpdate(name string) string {
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
  name   = "%s"
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
`, name, name+"-update")
}
