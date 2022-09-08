package baiducloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const (
	testAccBbcFlavorsDataSourceName          = "data.baiducloud_bbc_flavors.default"
	testAccBbcFlavorsDataSourceAttrKeyPrefix = "flavors.0."
)

func TestAccBaiduCloudBbcFlavorsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBbcFlavorsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBbcFlavorsDataSourceName),
					resource.TestCheckResourceAttr(testAccBbcFlavorsDataSourceName, testAccBbcFlavorsDataSourceAttrKeyPrefix+"cpu_count", "96"),
					resource.TestCheckResourceAttr(testAccBbcFlavorsDataSourceName, testAccBbcFlavorsDataSourceAttrKeyPrefix+"flavor_id", "BBC-I4-HC04S"),
					resource.TestCheckResourceAttrSet(testAccBbcFlavorsDataSourceName, testAccBbcFlavorsDataSourceAttrKeyPrefix+"cpu_type"),
					resource.TestCheckResourceAttrSet(testAccBbcFlavorsDataSourceName, testAccBbcFlavorsDataSourceAttrKeyPrefix+"disk"),
				),
			},
		},
	})
}

const testAccBbcFlavorsDataSourceConfig = `
data "baiducloud_bbc_flavors" "default" {
  filter {
    name   = "flavor_id"
    values = ["BBC-I4-HC04S"]
  }
}
`
