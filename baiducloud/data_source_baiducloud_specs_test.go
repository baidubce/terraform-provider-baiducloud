package baiducloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const (
	testAccSpecsDataSourceName          = "data.baiducloud_specs.default"
	testAccSpecsDataSourceAttrKeyPrefix = "specs.0."
)

func TestAccBaiduCloudSpecsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSpecsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccSpecsDataSourceName),
					resource.TestCheckResourceAttrSet(testAccSpecsDataSourceName, testAccSpecsDataSourceAttrKeyPrefix+"name"),
					resource.TestCheckResourceAttrSet(testAccSpecsDataSourceName, testAccSpecsDataSourceAttrKeyPrefix+"instance_type"),
					resource.TestCheckResourceAttrSet(testAccSpecsDataSourceName, testAccSpecsDataSourceAttrKeyPrefix+"cpu_count"),
					resource.TestCheckResourceAttrSet(testAccSpecsDataSourceName, testAccSpecsDataSourceAttrKeyPrefix+"memory_size_in_gb"),
					resource.TestCheckResourceAttrSet(testAccSpecsDataSourceName, testAccSpecsDataSourceAttrKeyPrefix+"local_disk_size_in_gb"),
				),
			},
		},
	})
}

const testAccSpecsDataSourceConfig = `
data "baiducloud_specs" "default" {}
`
