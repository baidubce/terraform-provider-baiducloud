package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccBLBResourceType = "baiducloud_blb"
	testAccBLBResourceName = testAccBLBResourceType + "." + BaiduCloudTestResourceName
)

//lintignore:AT003
func TestAccBaiduCloudBLB(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccBLBDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccBLBConfig(BaiduCloudTestResourceTypeNameblb),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBLBResourceName),
					resource.TestCheckResourceAttr(testAccBLBResourceName, "name", BaiduCloudTestResourceTypeNameblb),
					resource.TestCheckResourceAttr(testAccBLBResourceName, "cidr", "192.168.0.0/16"),
					resource.TestCheckResourceAttrSet(testAccBLBResourceName, "create_time"),
					resource.TestCheckResourceAttrSet(testAccBLBResourceName, "vpc_id"),
					resource.TestCheckResourceAttrSet(testAccBLBResourceName, "subnet_id"),
				),
			},
			{
				ResourceName:      testAccBLBResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBLBConfigUpdate(BaiduCloudTestResourceTypeNameblb),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBLBResourceName),
					resource.TestCheckResourceAttr(testAccBLBResourceName, "name", BaiduCloudTestResourceTypeNameblb+"-update"),
					resource.TestCheckResourceAttr(testAccBLBResourceName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(testAccBLBResourceName, "cidr", "192.168.0.0/16"),
					resource.TestCheckResourceAttrSet(testAccBLBResourceName, "create_time"),
					resource.TestCheckResourceAttrSet(testAccBLBResourceName, "vpc_id"),
					resource.TestCheckResourceAttrSet(testAccBLBResourceName, "subnet_id"),
				),
			},
		},
	})
}

func testAccBLBDestory(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	blbService := BLBService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccBLBResourceType {
			continue
		}

		_, _, err := blbService.GetBLBDetail(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		return WrapError(Error("BLB still exist"))
	}

	return nil
}

func testAccBLBConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_vpc" "default" {
    name = "terra-test-vpc"
    description = "created by terraform"
    cidr = "192.168.0.0/16"
    tags = {
       "product_name" = "terraform-test"
    }
}

resource "baiducloud_subnet" "default" {
  name = "terra-subnet"
  zone_name = "cn-bj-a"
  cidr = "192.168.3.0/24"
  vpc_id = "${baiducloud_vpc.default.id}"
}

resource "baiducloud_blb" "default" {
  name        = "${var.name}"
  description = "created by terraform"
  vpc_id      = "${baiducloud_vpc.default.id}"
  subnet_id   = "${baiducloud_subnet.default.id}"

  tags = {
   "product_name" = "terra-test"
  }
}
`, name)
}

func testAccBLBConfigUpdate(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_vpc" "default" {
    name = "terra-test-vpc"
    description = "created by terraform"
    cidr = "192.168.0.0/16"
    tags = {
       "product_name" = "terraform-test"
    }
}

resource "baiducloud_subnet" "default" {
  name = "terra-subnet"
  zone_name = "cn-bj-a"
  cidr = "192.168.3.0/24"
  vpc_id = "${baiducloud_vpc.default.id}"
}

resource "baiducloud_blb" "default" {
  name        = "${var.name}"
  description = "created by terraform"
  vpc_id      = "${baiducloud_vpc.default.id}"
  subnet_id   = "${baiducloud_subnet.default.id}"

  tags = {
   "product_name" = "terra-test"
  }
}
`, name+"-update")
}
