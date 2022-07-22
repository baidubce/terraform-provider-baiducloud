package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccAppBLBsDataSourceName          = "data.baiducloud_appblbs.default"
	testAccAppBLBsDataSourceAttrKeyPrefix = "appblbs.0."
)

//lintignore:AT003
func TestAccBaiduCloudAppBLBsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAppBLBDataSourceConfig(BaiduCloudTestResourceTypeNameAppblb),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBsDataSourceName),
					resource.TestCheckResourceAttrSet(testAccAppBLBsDataSourceName, testAccAppBLBsDataSourceAttrKeyPrefix+"blb_id"),
					resource.TestCheckResourceAttrSet(testAccAppBLBsDataSourceName, testAccAppBLBsDataSourceAttrKeyPrefix+"vpc_id"),
					resource.TestCheckResourceAttrSet(testAccAppBLBsDataSourceName, testAccAppBLBsDataSourceAttrKeyPrefix+"subnet_id"),
					resource.TestCheckResourceAttr(testAccAppBLBsDataSourceName, testAccAppBLBsDataSourceAttrKeyPrefix+"vpc_name", BaiduCloudTestResourceTypeNameAppblb),
					resource.TestCheckResourceAttr(testAccAppBLBsDataSourceName, testAccAppBLBsDataSourceAttrKeyPrefix+"subnet_name", BaiduCloudTestResourceTypeNameAppblb),
					resource.TestCheckResourceAttr(testAccAppBLBsDataSourceName, testAccAppBLBsDataSourceAttrKeyPrefix+"tags.testKey", "testValue"),
				),
			},
		},
	})
}

func testAccAppBLBDataSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "baiducloud_specs" "default" {}

data "baiducloud_zones" "default" {
  name_regex = ".*e$"
}

data "baiducloud_images" "default" {
  image_type = "System"
}

resource "baiducloud_instance" "default" {
  name                  = "${var.name}"
  image_id              = data.baiducloud_images.default.images.0.id
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  billing = {
    payment_timing = "Postpaid"
  }
}

resource "baiducloud_vpc" "default" {
  name        = "${var.name}"
  description = "created by terraform"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = "${var.name}"
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created by terraform"
}

resource "baiducloud_appblb" "default" {
  depends_on  = [baiducloud_instance.default]
  name        = "${var.name}"
  description = "created by terraform"
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

  filter {
    name = "vpc_id"
    values = [baiducloud_vpc.default.id]
  }
}
`, name)
}
