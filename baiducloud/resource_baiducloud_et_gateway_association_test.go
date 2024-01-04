package baiducloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const (
	testAccEtGatewayAssociationResourceType = "baiducloud_et_gateway_association"

	testAccEtGatewayAssociationResourceName = testAccEtGatewayAssociationResourceType + "." + BaiduCloudTestResourceName
)

func TestAccBaiduCloudEtGatewayAssociation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccEtGatewayAssociationConfig("tf-test-acc-et-gateway-association"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccEtGatewayAssociationResourceName),
					resource.TestCheckResourceAttrSet(testAccEtGatewayAssociationResourceName, "vpc_id"),
				),
			},
		},
	})
}

func TestAccBaiduCloudEtGatewayAssociationMultiIp(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccEtGatewayAssociationMultiIpConfig("tf-test-acc-et-gateway-association"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccEtGatewayAssociationResourceName),
					resource.TestCheckResourceAttrSet(testAccEtGatewayAssociationResourceName, "vpc_id"),
				),
			},
		},
	})
}

func TestAccBaiduCloudEtGatewayAssociationNilIp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccEtGatewayAssociationNilIpConfig("tf-test-acc-et-gateway-association"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccEtGatewayAssociationResourceName),
					resource.TestCheckResourceAttrSet(testAccEtGatewayAssociationResourceName, "vpc_id"),
				),
			},
		},
	})
}

func testAccEtGatewayAssociationConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_et_gateway" "default" {
  name = "my_name"
  vpc_id = "vpc-ud3hmp5ziuvm"
  speed = 200
  description = "description"
}

resource "baiducloud_et_gateway_association" "default" {
  et_gateway_id = baiducloud_et_gateway.default.et_gateway_id
  et_id = "xxx"
  channel_id = "xxx"
  local_cidrs = ["192.168.0.0/20"]
}
`, name)
}

func testAccEtGatewayAssociationMultiIpConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_et_gateway" "default" {
  name = "my_name"
  vpc_id = "vpc-xxx"
  speed = 200
  description = "description"
}

resource "baiducloud_et_gateway_association" "default" {
  et_gateway_id = baiducloud_et_gateway.default.et_gateway_id
  et_id = "xxx"
  channel_id = "xxx"
  local_cidrs = ["192.168.3.5","192.168.3.6","192.168.3.7"]
}
`, name)
}

func testAccEtGatewayAssociationNilIpConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_et_gateway" "default" {
  name = "my_name"
  vpc_id = "vpc-xxx"
  speed = 200
  description = "description"
}

resource "baiducloud_et_gateway_association" "default" {
  et_gateway_id = baiducloud_et_gateway.default.et_gateway_id
  et_id = "xxx"
  channel_id = "xxx"
}

`, name)
}
