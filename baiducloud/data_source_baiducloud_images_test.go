package baiducloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccImagesDataSourceName          = "data.baiducloud_images.default"
	testAccImagesDataSourceAttrKeyPrefix = "images.0."
)

//lintignore:AT003
func TestAccBaiduCloudImagesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccImagesDataSourceName),
					resource.TestCheckResourceAttr(testAccImagesDataSourceName, testAccImagesDataSourceAttrKeyPrefix+"os_name", "CentOS"),
					resource.TestCheckResourceAttrSet(testAccImagesDataSourceName, testAccImagesDataSourceAttrKeyPrefix+"id"),
					resource.TestCheckResourceAttrSet(testAccImagesDataSourceName, testAccImagesDataSourceAttrKeyPrefix+"name"),
					resource.TestCheckResourceAttrSet(testAccImagesDataSourceName, testAccImagesDataSourceAttrKeyPrefix+"type"),
				),
			},
		},
	})
}

const testAccImagesDataSourceConfig = `
data "baiducloud_images" "default" {
  image_type = "System"
  name_regex = "7.5.*"
  os_name    = "CentOs"

  filter {
    name = "name"
    values = ["7.5.*"]
  }
}
`
