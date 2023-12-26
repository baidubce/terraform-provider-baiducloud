package baiducloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const (
	testAccEtGatewayResourceType = "baiducloud_et_gateway"
	testAccEtGatewayResourceName = testAccEtGatewayResourceType + "." + BaiduCloudTestResourceName
)

func TestAccBaiduCloudEtGateway(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccEtGatewayConfig("tf-test-acc-et-gateway"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccEtGatewayResourceName),
					resource.TestCheckResourceAttrSet(testAccEtGatewayResourceName, "vpc_id"),
				),
			},
		},
	})
}

func TestAccBaiduCloudEtGatewayNilEt(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccEtGatewayNilEtAndChannelConfig("tf-test-acc-et-gateway"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccEtGatewayResourceName),
					resource.TestCheckResourceAttrSet(testAccEtGatewayResourceName, "vpc_id"),
				),
			},
		},
	})
}

func TestAccBaiduCloudEtGatewayMultiIp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccEtGatewayMultiIpConfig("tf-test-acc-et-gateway"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccEtGatewayResourceName),
					resource.TestCheckResourceAttrSet(testAccEtGatewayResourceName, "vpc_id"),
				),
			},
		},
	})
}

func TestAccBaiduCloudEtGatewayNilIp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccEtGatewayNilIpConfig("tf-test-acc-et-gateway"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccEtGatewayResourceName),
					resource.TestCheckResourceAttrSet(testAccEtGatewayResourceName, "vpc_id"),
				),
			},
		},
	})
}

func testAccEtGatewayConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "baiducloud_et_gateway" "default" {
  name = "my_name"
  vpc_id = "vpc-ud3hmp5ziuvm"
  speed = 200
  description = "description"
  et_id = "xxx"
  channel_id = "xxx"
  local_cidrs = ["192.168.0.0/20"]
}
`, name)
}

func testAccEtGatewayNilEtAndChannelConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "baiducloud_et_gateway" "default" {
  name = "my_name"
  vpc_id = "vpc-xxx"
  speed = 200
  description = "description"
  local_cidrs = ["192.168.0.0/20"]
}
`, name)
}

func testAccEtGatewayMultiIpConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "baiducloud_et_gateway" "default" {
  name = "my_name"
  vpc_id = "vpc-xxx"
  speed = 200
  description = "description"
  et_id = "xxx"
  channel_id = "xxx"
  local_cidrs = ["192.168.3.5","192.168.3.6","192.168.3.7"]
}
`, name)
}

func testAccEtGatewayNilIpConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "baiducloud_et_gateway" "default" {
  name = "my_name"
  vpc_id = "vpc-xxx"
  speed = 200
  description = "description"
  et_id = "xxx"
  channel_id = "xxx"
}
`, name)
}
