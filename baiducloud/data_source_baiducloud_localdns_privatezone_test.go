package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccLocalDnsPrivatezonesDataSourceName          = "data.baiducloud_localdns_privatezones.default"
	testAccLocalDnsPrivatezonesDataSourceAttrKeyPrefix = "zones.0."
)

//lintignore:AT003
func TestAccBaiduCloudLocalDnsPrivatezonesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLocalDnsPrivatezonesDataSourceConfig(BaiduCloudTestResourceTypeNameLocalDnsPrivatezone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccLocalDnsPrivatezonesDataSourceName),
					// check attr
					resource.TestCheckResourceAttrSet(testAccLocalDnsPrivatezonesDataSourceName, testAccLocalDnsPrivatezonesDataSourceAttrKeyPrefix+"zone_name"),
					resource.TestCheckResourceAttrSet(testAccLocalDnsPrivatezonesDataSourceName, testAccLocalDnsPrivatezonesDataSourceAttrKeyPrefix+"zone_id"),
					// check attr value
					resource.TestCheckResourceAttr(testAccLocalDnsPrivatezonesDataSourceName, testAccLocalDnsPrivatezonesDataSourceAttrKeyPrefix+"record_count", "3"),
				),
			},
		},
	})
}

func testAccLocalDnsPrivatezonesDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_localdns_privatezone" "default" {
  zone_name         = "%s.com"
}

data "baiducloud_localdns_privatezones" "default" {

}
`, name)
}
