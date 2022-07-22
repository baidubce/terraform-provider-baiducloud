package baiducloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccZonesDataSourceName          = "data.baiducloud_zones.default"
	testAccZonesDataSourceAttrKeyPrefix = "zones.0."
)

//lintignore:AT003
func TestAccBaiduCloudZonesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccZonesDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccZonesDataSourceName),
					resource.TestCheckResourceAttrSet(testAccZonesDataSourceName, testAccZonesDataSourceAttrKeyPrefix+"zone_name"),
				),
			},
		},
	})
}

const testAccZonesDataSourceConfig = `
data "baiducloud_zones" "default" {
  name_regex = ".*e$"

  filter {
    name = "zone_name"
    values = ["cn-*"]
  }
}
`
