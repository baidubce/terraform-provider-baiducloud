package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const (
	testAccAppBLBsDataSourceName          = "data.baiducloud_appblbs.default"
	testAccAppBLBsDataSourceAttrKeyPrefix = "appblbs.0."
)

func TestAccBaiduCloudAppBLBsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAppBLBDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBsDataSourceName),
					resource.TestCheckResourceAttrSet(testAccAppBLBsDataSourceName, testAccAppBLBsDataSourceAttrKeyPrefix+"blb_id"),
					resource.TestCheckResourceAttrSet(testAccAppBLBsDataSourceName, testAccAppBLBsDataSourceAttrKeyPrefix+"vpc_id"),
					resource.TestCheckResourceAttrSet(testAccAppBLBsDataSourceName, testAccAppBLBsDataSourceAttrKeyPrefix+"subnet_id"),
					resource.TestCheckResourceAttr(testAccAppBLBsDataSourceName, testAccAppBLBsDataSourceAttrKeyPrefix+"vpc_name", BaiduCloudTestResourceAttrNamePrefix+"VPC"),
					resource.TestCheckResourceAttr(testAccAppBLBsDataSourceName, testAccAppBLBsDataSourceAttrKeyPrefix+"subnet_name", BaiduCloudTestResourceAttrNamePrefix+"Subnet"),
					resource.TestCheckResourceAttr(testAccAppBLBsDataSourceName, testAccAppBLBsDataSourceAttrKeyPrefix+"tags.0.tag_key", "testKey"),
					resource.TestCheckResourceAttr(testAccAppBLBsDataSourceName, testAccAppBLBsDataSourceAttrKeyPrefix+"tags.0.tag_value", "testValue"),
				),
			},
		},
	})
}

func testAccAppBLBDataSourceConfig() string {
	return fmt.Sprintf(`
data "baiducloud_zones" "default" {}

resource "baiducloud_vpc" "default" {
  name        = "%s"
  description = "test"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = "%s"
  zone_name   = "${data.baiducloud_zones.default.zones.0.zone_name}"
  cidr        = "192.168.0.0/24"
  vpc_id      = "${baiducloud_vpc.default.id}"
  description = "test description"
}

resource "baiducloud_appblb" "default" {
  name        = "%s"
  description = ""
  vpc_id      = "${baiducloud_vpc.default.id}"
  subnet_id   = "${baiducloud_subnet.default.id}"

  tags {
    tag_key   = "testKey"
    tag_value = "testValue"
  }
}

data "baiducloud_appblbs" "default" {
  blb_id  = "${baiducloud_appblb.default.id}"
  name    = "${baiducloud_appblb.default.name}"
  address = "${baiducloud_appblb.default.address}"
}
`, BaiduCloudTestResourceAttrNamePrefix+"VPC",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceName+"APPBLB")
}
