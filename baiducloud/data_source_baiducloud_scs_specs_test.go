package baiducloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccScsSpecsDataSourceName          = "data.baiducloud_scs_specs.default"
	testAccScsSpecsDataSourceAttrKeyPrefix = "specs.0."
)

func TestAccBaiduCloudScsSpecsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScsSpecsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccScsSpecsDataSourceName),
					resource.TestCheckResourceAttr(testAccScsSpecsDataSourceName, testAccScsSpecsDataSourceAttrKeyPrefix+"node_capacity", "1"),
					resource.TestCheckResourceAttr(testAccScsSpecsDataSourceName, testAccScsSpecsDataSourceAttrKeyPrefix+"node_type", "cache.n1.micro"),
				),
			},
		},
	})
}

const testAccScsSpecsDataSourceConfig = `
data "baiducloud_scs_specs" "default" {
    cluster_type     = "cluster"
    node_capacity  = 1
}
`
