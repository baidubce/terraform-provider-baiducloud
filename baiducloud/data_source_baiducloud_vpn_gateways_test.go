package baiducloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccVPNGatewaysDataSourceName          = "data.baiducloud_vpn_gateways.default"
	testAccVPNGatewaysDataSourceAttrKeyPrefix = "vpn_gateways.0."
)

//lintignore:AT003
func TestAccBaiduCloudVPNGatewaysDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNGatewaysDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccVPNGatewaysDataSourceName),
					resource.TestCheckResourceAttrSet(testAccVPNGatewaysDataSourceName, testAccVPNGatewaysDataSourceAttrKeyPrefix+"vpc_id"),
				),
			},
		},
	})
}

const testAccVPNGatewaysDataSourceConfig = `
data "baiducloud_vpn_gateways" "default" {
  vpc_id = "vpc-65cz3hu92kz2"
}
`
