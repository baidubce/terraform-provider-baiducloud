package baiducloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccFlavorsDataSourceName          = "data.baiducloud_specs.default"
	testAccFlavorsDataSourceAttrKeyPrefix = "specs.0."
)

//lintignore:AT003
func TestAccBaiduCloudFlavorsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFlavorsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccFlavorsDataSourceName),
					resource.TestCheckResourceAttr(testAccFlavorsDataSourceName, testAccFlavorsDataSourceAttrKeyPrefix+"zone_name", "cn-bj-d"),
					resource.TestCheckResourceAttr(testAccFlavorsDataSourceName, testAccFlavorsDataSourceAttrKeyPrefix+"cpu_count", "1"),
					resource.TestCheckResourceAttr(testAccFlavorsDataSourceName, testAccFlavorsDataSourceAttrKeyPrefix+"memory_capacity_in_gb", "4"),
				),
			},
		},
	})
}

const testAccFlavorsDataSourceConfig = `
data "baiducloud_specs" "default" {
  zone_name = "cn-bj-d"
  filter {
    name   = "cpu_count"
    values = ["^([1])$"]
  }
  filter {
    name   = "memory_capacity_in_gb"
    values = ["^([4])$"]
  }
}
`
