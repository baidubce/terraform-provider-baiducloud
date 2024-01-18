package baiducloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const (
	testAccEipBpResourceType = "baiducloud_eipbp"
	testAccEipBpResourceName = testAccEipBpResourceType + "." + BaiduCloudTestResourceName
)

func TestAccBaiduCloudEipbpWithNameSimple(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccEipBpConfigWithNameSimple("tf-test-acc-eipbp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccEipBpResourceName),
					resource.TestCheckResourceAttrSet(testAccEipBpResourceName, "name"),
				),
			},
		},
	})
}

func TestAccBaiduCloudEipBpWithoutnameSimple(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccEipbpConfigWithoutNameSimple("tf-test-acc-et-eipbp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccEtGatewayResourceName),
					resource.TestCheckResourceAttrSet(testAccEtGatewayResourceName, "bandwidth_in_mbps"),
				),
			},
		},
	})
}

func TestAccBaiduCloudEipbpWithNameFull(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccEipBpConfigWithName("tf-test-acc-eipbp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccEipBpResourceName),
					resource.TestCheckResourceAttrSet(testAccEipBpResourceName, "name"),
				),
			},
		},
	})
}

func TestAccBaiduCloudEipBpWithoutnameFull(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccEipbpConfigWithoutName("tf-test-acc-et-eipbp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccEtGatewayResourceName),
					resource.TestCheckResourceAttrSet(testAccEtGatewayResourceName, "bandwidth_in_mbps"),
				),
			},
		},
	})
}

func testAccEipBpConfigWithName(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_eip" "default" {
  name              = "%s"
  bandwidth_in_mbps = 2
  payment_timing    = "Postpaid"
  billing_method    = "ByTraffic"

  tags = {
    "testKey" = "testValue"
  }
}

resource "baiducloud_eipgroup" "default" {
  name              = "testEIPgroup"
  eip_count         = 2
  bandwidth_in_mbps = 100
  payment_timing    = "Postpaid"
  billing_method    = "ByTraffic"
}

resource "baiducloud_eipbp" "default" {
  name              = "testEIPbp"
  eip               = baiducloud_eip.default.eip
  bandwidth_in_mbps = 100
  eip_group_id      = baiducloud_eipgroup.default.group_id
}
`, name)
}

func testAccEipbpConfigWithoutName(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_eip" "default" {
  name              = "%s"
  bandwidth_in_mbps = 2
  payment_timing    = "Postpaid"
  billing_method    = "ByTraffic"

  tags = {
    "testKey" = "testValue"
  }
}

resource "baiducloud_eipgroup" "default" {
  name              = "testEIPgroup"
  eip_count         = 2
  bandwidth_in_mbps = 100
  payment_timing    = "Postpaid"
  billing_method    = "ByTraffic"
}

resource "baiducloud_eipbp" "default" {
  eip               = baiducloud_eip.default.eip
  bandwidth_in_mbps = 100
  eip_group_id      = baiducloud_eipgroup.default.group_id
}
`, name)
}

func testAccEipBpConfigWithNameSimple(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_eipbp" "default" {
  name              = "testEIPbp"
  eip               = "10.23.42.12""
  bandwidth_in_mbps = 100
  eip_group_id      = "eg-1MxtUX7c""
}
`, name)
}

func testAccEipbpConfigWithoutNameSimple(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_eipbp" "default" {
  eip               = "10.23.42.12"
  bandwidth_in_mbps = 100
  eip_group_id      = "eg-1MxtUX7c"
}
`, name)
}
