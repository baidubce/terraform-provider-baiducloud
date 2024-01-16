package baiducloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const (
	testAccEipGroupResourceType = "baiducloud_eipgroup"
	testAccEipGroupResourceName = testAccEipGroupResourceType + "." + BaiduCloudTestResourceName
)

func TestAccBaiduCloudEipGroupWithName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccEipGroupConfigWithName("tf-test-acc-eipgroup"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccEipGroupResourceName),
					resource.TestCheckResourceAttrSet(testAccEipGroupResourceName, "bandwidth_in_mbps"),
				),
			},
		},
	})
}

func TestAccBaiduCloudEipGroupWithoutname(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccEipGroupConfigWithoutName("tf-test-acc-et-eipgroup"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccEtGatewayResourceName),
					resource.TestCheckResourceAttrSet(testAccEtGatewayResourceName, "bandwidth_in_mbps"),
				),
			},
		},
	})
}

func testAccEipGroupConfigWithName(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "baiducloud_eipgroup" "default" {
  name              = "testEIPgroup"
  eip_count         = 2
  bandwidth_in_mbps = 100
  payment_timing    = "Postpaid"
  billing_method    = "ByTraffic"
}
`, name)
}

func testAccEipGroupConfigWithoutName(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "baiducloud_eipgroup" "default" {
  eip_count         = 2
  bandwidth_in_mbps = 100
  payment_timing    = "Postpaid"
  billing_method    = "ByTraffic"
}
`, name)
}
