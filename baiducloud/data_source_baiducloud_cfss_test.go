package baiducloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const (
	testAccCfsDataSourceName          = "data.baiducloud_cfss.default"
	testAccCfsDataSourceAttrKeyPrefix = "cfss.0."
)

func TestAccBaiduCloudCfssDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccCfssConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccCfsDataSourceName),
					resource.TestCheckResourceAttr(testAccCfsDataSourceName, testAccCfsDataSourceAttrKeyPrefix+"name", "terraform_test"),
				),
			},
		},
	})
}

func testAccCfssConfig() string {
	return fmt.Sprintf(`
resource "baiducloud_cfs" "default" {
  name = "terraform_test"
  zone = "zoneD"
}

data baiducloud_cfss "default" {
  filter{
    name = "fs_id"
    values = [baiducloud_cfs.default.id]
  }
}
`)
}
