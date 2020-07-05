package baiducloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const (
	testAccCceVersionDataSourceName          = "data.baiducloud_cce_versions.default"
	testAccCceVersionDataSourceAttrKeyPrefix = "versions.#"
)

//lintignore:AT003
func TestAccBaiduCloudCceVersionsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVersionsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccCceVersionDataSourceName),
					resource.TestCheckResourceAttrSet(testAccCceVersionDataSourceName, testAccCceVersionDataSourceAttrKeyPrefix),
				),
			},
		},
	})
}

const testAccVersionsDataSourceConfig = `
data "baiducloud_cce_versions" "default" {
    version_regex        = ".*13.*"
}
`
