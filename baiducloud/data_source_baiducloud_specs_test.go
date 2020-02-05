package baiducloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const (
	testAccSpecsDataSourceName          = "data.baiducloud_specs.default"
	testAccSpecsDataSourceAttrKeyPrefix = "specs.0."
)

//lintignore:AT003
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
					resource.TestCheckResourceAttr(testAccSpecsDataSourceName, testAccSpecsDataSourceAttrKeyPrefix+"name", "bcc.g1.tiny"),
					resource.TestCheckResourceAttr(testAccSpecsDataSourceName, testAccSpecsDataSourceAttrKeyPrefix+"instance_type", "General"),
					resource.TestCheckResourceAttr(testAccSpecsDataSourceName, testAccSpecsDataSourceAttrKeyPrefix+"cpu_count", "1"),
					resource.TestCheckResourceAttr(testAccSpecsDataSourceName, testAccSpecsDataSourceAttrKeyPrefix+"memory_size_in_gb", "2"),
					resource.TestCheckResourceAttrSet(testAccSpecsDataSourceName, testAccSpecsDataSourceAttrKeyPrefix+"local_disk_size_in_gb"),
				),
			},
		},
	})
}

const testAccSpecsDataSourceConfig = `
data "baiducloud_specs" "default" {
    name_regex        = "bcc.g1.tiny"
    instance_type     = "General"
    cpu_count         = 1
    memory_size_in_gb = 2
}
`
