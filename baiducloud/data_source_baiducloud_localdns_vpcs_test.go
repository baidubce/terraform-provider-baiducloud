package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccLocalDnsVpcsDataSourceName          = "data.baiducloud_localdns_vpcs.default"
	testAccLocalDnsVpcsDataSourceAttrKeyPrefix = "bind_vpcs.0."
)

//lintignore:AT003
func TestAccBaiduCloudLocalDnsVpcsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLocalDnsVpcsDataSourceConfig(BaiduCloudTestResourceTypeNameLocalDnsVPC),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccLocalDnsVpcsDataSourceName),
					// check attr
					resource.TestCheckResourceAttrSet(testAccLocalDnsVpcsDataSourceName, testAccLocalDnsVpcsDataSourceAttrKeyPrefix+"vpc_id"),
					resource.TestCheckResourceAttrSet(testAccLocalDnsVpcsDataSourceName, testAccLocalDnsVpcsDataSourceAttrKeyPrefix+"vpc_name"),
					// check attr value
					resource.TestCheckResourceAttr(testAccLocalDnsVpcsDataSourceName, testAccLocalDnsVpcsDataSourceAttrKeyPrefix+"vpc_region", "bj"),
				),
			},
		},
	})
}

func testAccLocalDnsVpcsDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_localdns_privatezone" "default" {
    zone_name = "%s.com"
}

data "baiducloud_localdns_vpcs" "default" {
	zone_id = "zone-1mytixsfqpku"
}
`, name)
}
