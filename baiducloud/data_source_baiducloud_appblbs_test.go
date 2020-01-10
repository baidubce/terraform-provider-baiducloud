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
					resource.TestCheckResourceAttr(testAccAppBLBsDataSourceName, testAccAppBLBsDataSourceAttrKeyPrefix+"tags.testKey", "testValue"),
				),
			},
		},
	})
}

func testAccAppBLBDataSourceConfig() string {
	return fmt.Sprintf(`
data "baiducloud_specs" "default" {}

data "baiducloud_zones" "default" {}

data "baiducloud_images" "default" {
  image_type = "System"
}

resource "baiducloud_instance" "default" {
  name                  = "%s"
  image_id              = data.baiducloud_images.default.images.0.id
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  billing = {
    payment_timing = "Postpaid"
  }
}

resource "baiducloud_vpc" "default" {
  name        = "%s"
  description = "test"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = "%s"
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "test description"
}

resource "baiducloud_appblb" "default" {
  depends_on  = [baiducloud_instance.default]
  name        = "%s"
  description = ""
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id

  tags = {
    "testKey" = "testValue"
  }
}

data "baiducloud_appblbs" "default" {
  blb_id  = baiducloud_appblb.default.id
  name    = baiducloud_appblb.default.name
  address = baiducloud_appblb.default.address
}
`, BaiduCloudTestResourceAttrNamePrefix+"BCC",
		BaiduCloudTestResourceAttrNamePrefix+"VPC",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceName+"APPBLB")
}
