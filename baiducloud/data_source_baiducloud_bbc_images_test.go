package baiducloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccBbcImagesDataSourceName          = "data.baiducloud_bbc_images.default"
	testAccBbcImagesDataSourceAttrKeyPrefix = "images.0."
)

func TestAccBaiduCloudBbcImagesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBbcImagesDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBbcImagesDataSourceName),
					resource.TestCheckResourceAttr(testAccBbcImagesDataSourceName, testAccBbcImagesDataSourceAttrKeyPrefix+"os_name", "CentOS"),
					resource.TestCheckResourceAttrSet(testAccBbcImagesDataSourceName, testAccBbcImagesDataSourceAttrKeyPrefix+"id"),
					resource.TestCheckResourceAttrSet(testAccBbcImagesDataSourceName, testAccBbcImagesDataSourceAttrKeyPrefix+"name"),
					resource.TestCheckResourceAttrSet(testAccBbcImagesDataSourceName, testAccBbcImagesDataSourceAttrKeyPrefix+"type"),
				),
			},
		},
	})
}

const testAccBbcImagesDataSourceConfig = `
data "baiducloud_bbc_images" "default" {
  image_type = "BbcSystem"
  os_name    = "CentOS"
  filter {
    name   = "id"
    values = ["m-i2aoqIlx"]
  }
}
`
