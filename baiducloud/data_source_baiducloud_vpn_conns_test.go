package baiducloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccVPNConnsDataSourceName          = "data.baiducloud_vpn_conns.default"
	testAccVPNConnsDataSourceAttrKeyPrefix = "vpn_conns.0."
)

//lintignore:AT003
func TestAccBaiduCloudVPNConnsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNConnsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccVPNConnsDataSourceName),
					resource.TestCheckResourceAttrSet(testAccVPNConnsDataSourceName, testAccVPNConnsDataSourceAttrKeyPrefix+"vpc_id"),
				),
			},
		},
	})
}

const testAccVPNConnsDataSourceConfig = `
data "baiducloud_vpn_conns" "default" {
  vpn_id = "vpn-b2gmbd51pk57"
}
`
