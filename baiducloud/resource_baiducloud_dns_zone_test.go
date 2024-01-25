package baiducloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const (
	testAccDnsZoneResourceType = "baiducloud_eipbp"
	testAccDnsZoneResourceName = testAccDnsZoneResourceType + "." + BaiduCloudTestResourceName
)

func TestAccBaiduCloudDnsZoneSimple(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccDnsZoneConfig("tf-test-acc-dns_zone"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccDnsZoneResourceName),
					resource.TestCheckResourceAttrSet(testAccDnsZoneResourceName, "name"),
				),
			},
		},
	})
}

func testAccDnsZoneConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_dns_zone" "default" {
  name              = "testDnsZone"
}
`, name)
}
