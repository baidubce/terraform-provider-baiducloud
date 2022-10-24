package baiducloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const (
	testAccVPNGatewayResourceType = "baiducloud_vpn_gateway"
	testAccVPNGatewayResourceName = testAccVPNGatewayResourceType + "." + BaiduCloudTestResourceName
)

func TestAccBaiduCloudVPNGateway(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccVPNGatewayConfig(BaiduCloudTestResourceTypeNameVPNGateway),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccVPNGatewayResourceName),
					resource.TestCheckResourceAttr(testAccVPNGatewayResourceName, "description", "test desc"),
				),
			},
		},
	})
}

func testAccVPNGatewayConfig(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_vpn_gateway" "default" {
  vpn_name       = "%s"
  vpc_id         = "vpc-65cz3hu92kz2"
  description    = "test desc"
  payment_timing = "Postpaid"
}
`, name)
}
