package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccDnscustomlinesDataSourceName          = "data.baiducloud_dns_customlines.default"
	testAccDnscustomlinesDataSourceAttrKeyPrefix = "customlines.0."
)

//lintignore:AT003
func TestAccBaiduCloudDnscustomlinesDataSourceSimple(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnscustomlineByNameDataSourceSimpleConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccDnscustomlinesDataSourceName),
					resource.TestCheckResourceAttrSet(testAccDnscustomlinesDataSourceAttrKeyPrefix, testAccDnscustomlinesDataSourceAttrKeyPrefix+"name"),
				),
			},
		},
	})
}

func TestAccBaiduCloudDnscustomlinesDataSourceFull(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnscustomlineDataSourceFullConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccDnscustomlinesDataSourceName),
					resource.TestCheckResourceAttrSet(testAccDnscustomlinesDataSourceName, testAccDnscustomlinesDataSourceAttrKeyPrefix+"name"),
				),
			},
		},
	})
}

func testAccDnscustomlineByNameDataSourceSimpleConfig() string {
	return fmt.Sprintf(`
data "baiducloud_dns_customlines" "default" { 
}
`)
}

func testAccDnscustomlineDataSourceFullConfig() string {
	return fmt.Sprintf(`
resource "baiducloud_dns_customline" "default" {
  name              = "testname"
  lines             = ["zhejiang.ct"]
}

data "baiducloud_dns_customlines" "default" {
}
`)
}
