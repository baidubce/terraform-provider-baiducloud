package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccBLBsDataSourceName          = "data.baiducloud_blbs.default"
	testAccBLBsDataSourceAttrKeyPrefix = "blbs.0."
)

//lintignore:AT003
func TestAccBaiduCloudBLBsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBLBDataSourceConfig(BaiduCloudTestResourceTypeNameblb),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBLBsDataSourceName),
					resource.TestCheckResourceAttrSet(testAccBLBsDataSourceName, testAccBLBsDataSourceAttrKeyPrefix+"blb_id"),
					resource.TestCheckResourceAttrSet(testAccBLBsDataSourceName, testAccBLBsDataSourceAttrKeyPrefix+"vpc_id"),
					resource.TestCheckResourceAttrSet(testAccBLBsDataSourceName, testAccBLBsDataSourceAttrKeyPrefix+"subnet_id"),
					resource.TestCheckResourceAttr(testAccBLBsDataSourceName, testAccBLBsDataSourceAttrKeyPrefix+"vpc_name", BaiduCloudTestResourceTypeNameblb),
				),
			},
		},
	})
}

func testAccBLBDataSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "baiducloud_vpc" "default" {
    name        = "${var.name}"
    description = "baiducloud vpc created by terraform"
    cidr = "192.168.0.0/16"
    tags = {
       "product_name" = "terraform-test"
    }
}

resource "baiducloud_subnet" "default" {
  name        = "${var.name}"
  zone_name = "cn-bj-a"
  cidr = "192.168.3.0/24"
  vpc_id = "${baiducloud_vpc.default.id}"
}

resource "baiducloud_blb" "default" {
  name        = "${var.name}"
  description = "this is a test LoadBalance instance"
  vpc_id      = "${baiducloud_vpc.default.id}"
  subnet_id   = "${baiducloud_subnet.default.id}"

  tags = {
   "product_name" = "terra-test"
  }
}

data "baiducloud_blbs" "default" {
  blb_id  = baiducloud_blb.default.id

}
`, name)
}
