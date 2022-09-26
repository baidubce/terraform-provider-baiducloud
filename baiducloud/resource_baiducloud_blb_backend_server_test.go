package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccBLBServerResourceType = "baiducloud_blb_backend_server"
	testAccBLBServerResourceName = testAccBLBServerResourceType + "." + BaiduCloudTestResourceName
)

func TestAccBaiduCloudBLBBackendServer_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccBLBServerDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccBLBServerConfig(BaiduCloudTestResourceTypeNameblbServer),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBLBServerResourceName),
					resource.TestCheckResourceAttrSet(testAccBLBServerResourceName, "blb_id"),
					resource.TestCheckResourceAttr(testAccBLBServerResourceName, "backend_server_list.0.weight", "39"),
				),
			},
			{
				Config: testAccBLBServerConfigUpdate(BaiduCloudTestResourceTypeNameblbServer),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBLBServerResourceName),
					resource.TestCheckResourceAttr(testAccBLBServerResourceName, "backend_server_list.0.weight", "40"),
					resource.TestCheckResourceAttrSet(testAccBLBServerResourceName, "blb_id"),
				),
			},
		},
	})
}

func testAccBLBServerDestory(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	blbService := BLBService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccBLBServerResourceType {
			continue
		}

		raw, err := blbService.BackendServerList(rs.Primary.Attributes["blb_id"])
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		for _, sg := range raw {
			if sg["id"] == rs.Primary.ID {
				return WrapError(Error("BLB Server still exist"))
			}
		}
	}

	return nil
}

func testAccBLBServerConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "baiducloud_images" "default" {
  image_type = "System"
  name_regex = "8.4 aarch"
  os_name    = "CentOS"
}

resource "baiducloud_instance" "default1" {
  billing = {
    payment_timing = "Postpaid"
  }
  instance_spec = "bcc.gr1.c1m4"
  image_id      = data.baiducloud_images.default.images.0.id
  tags          = {
    "use"  = "xx-bcc"
  }
  availability_zone = "cn-bj-d"
  #security_groups  = [baiducloud_security_group.default.id]
}

resource "baiducloud_blb" "default2" {
  name        = "${var.name}"
  description = "created by terraform"
  vpc_id      = "${baiducloud_instance.default1.vpc_id}"
  subnet_id   = "${baiducloud_instance.default1.subnet_id}"
}

resource "baiducloud_blb_backend_server" "default" {
  blb_id       = "${baiducloud_blb.default2.id}"
  backend_server_list {
    instance_id = "${baiducloud_instance.default1.id}"
    weight      = 39
  }

}
`, name)
}

func testAccBLBServerConfigUpdate(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "baiducloud_images" "default" {
  image_type = "System"
  name_regex = "8.4 aarch"
  os_name    = "CentOS"
}

resource "baiducloud_instance" "default1" {
  billing = {
    payment_timing = "Postpaid"
  }
  instance_spec = "bcc.gr1.c1m4"
  image_id      = data.baiducloud_images.default.images.0.id
  tags          = {
    "use"  = "xx-bcc"
  }
  availability_zone = "cn-bj-d"
  #security_groups  = [baiducloud_security_group.default.id]
}

resource "baiducloud_blb" "default2" {
  name        = "${var.name}"
  description = "created by terraform"
  vpc_id      = "${baiducloud_instance.default1.vpc_id}"
  subnet_id   = "${baiducloud_instance.default1.subnet_id}"
}

resource "baiducloud_blb_backend_server" "default" {
  blb_id       = "${baiducloud_blb.default2.id}"
  backend_server_list {
    instance_id = "${baiducloud_instance.default1.id}"
    weight      = 40
  }

}
`, name)
}
