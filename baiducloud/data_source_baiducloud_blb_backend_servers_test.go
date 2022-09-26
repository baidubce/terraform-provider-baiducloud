package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccBLBServersDataSourceName = "data.baiducloud_blb_backend_servers.default"
)

//lintignore:AT003
func TestAccBaiduCloudBLBServersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBLBServerDataSourceConfig(BaiduCloudTestResourceTypeNameblbServer),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBLBServersDataSourceName),
					resource.TestCheckResourceAttrSet(testAccBLBServersDataSourceName, "blb_id"),
				),
			},
		},
	})
}

func testAccBLBServerDataSourceConfig(name string) string {
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
    "use"  = "zq-bcc"
  }
  availability_zone = "cn-bj-d"
  #security_groups  = [baiducloud_security_group.default.id]
}

resource "baiducloud_blb" "default2" {
  name        = "terratestLoadBalance"
  description = "this is a test LoadBalance instance"
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
data "baiducloud_blb_backend_servers" "default" {
    blb_id = baiducloud_blb.default2.id
}

`, name)
}
