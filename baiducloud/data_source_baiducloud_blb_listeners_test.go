package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccBLBListenersDataSourceName          = "data.baiducloud_blb_listeners.default"
	testAccBLBListenersDataSourceAttrKeyPrefix = "listeners.0."
)

//lintignore:AT003
func TestAccBaiduCloudBLBListenersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBLBListenersDataSourceConfig(BaiduCloudTestResourceTypeNameblbListener),
				Check: resource.ComposeTestCheckFunc(
					// TCP Listener
					testAccCheckBaiduCloudDataSourceId(testAccBLBListenersDataSourceName+"_TCP"),
					resource.TestCheckResourceAttr(testAccBLBListenersDataSourceName+"_TCP", testAccBLBListenersDataSourceAttrKeyPrefix+"listener_port", "125"),
					resource.TestCheckResourceAttr(testAccBLBListenersDataSourceName+"_TCP", testAccBLBListenersDataSourceAttrKeyPrefix+"protocol", "TCP"),
					resource.TestCheckResourceAttr(testAccBLBListenersDataSourceName+"_TCP", testAccBLBListenersDataSourceAttrKeyPrefix+"scheduler", "LeastConnection"),

					// UDP Listener
					testAccCheckBaiduCloudDataSourceId(testAccBLBListenersDataSourceName+"_UDP"),
					resource.TestCheckResourceAttr(testAccBLBListenersDataSourceName+"_UDP", testAccBLBListenersDataSourceAttrKeyPrefix+"listener_port", "126"),
					resource.TestCheckResourceAttr(testAccBLBListenersDataSourceName+"_UDP", testAccBLBListenersDataSourceAttrKeyPrefix+"protocol", "UDP"),
					resource.TestCheckResourceAttr(testAccBLBListenersDataSourceName+"_UDP", testAccBLBListenersDataSourceAttrKeyPrefix+"scheduler", "LeastConnection"),

					// HTTP Listener
					testAccCheckBaiduCloudDataSourceId(testAccBLBListenersDataSourceName+"_HTTP"),
					resource.TestCheckResourceAttr(testAccBLBListenersDataSourceName+"_HTTP", testAccBLBListenersDataSourceAttrKeyPrefix+"listener_port", "127"),
					resource.TestCheckResourceAttr(testAccBLBListenersDataSourceName+"_HTTP", testAccBLBListenersDataSourceAttrKeyPrefix+"protocol", "HTTP"),
					resource.TestCheckResourceAttr(testAccBLBListenersDataSourceName+"_HTTP", testAccBLBListenersDataSourceAttrKeyPrefix+"scheduler", "LeastConnection"),
				),
			},
		},
	})
}

func testAccBLBListenersDataSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_vpc" "default" {
  name        = "${var.name}"
  description = "created by terraform"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = "${var.name}"
  zone_name   = "cn-bj-a"
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "test description"
}

resource "baiducloud_blb" "default" {
  name        = "${var.name}"
  description = "created by terraform"
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}

resource "baiducloud_blb_listener" "default_TCP" {
  blb_id        = baiducloud_blb.default.id
  listener_port = 125
  backend_port  = 125
  protocol      = "TCP"
  scheduler     = "LeastConnection"
}

resource "baiducloud_blb_listener" "default_UDP" {
  blb_id         = baiducloud_blb.default.id
  listener_port  = 126
  backend_port   = 126
  protocol       = "UDP"
  health_check_string  = "healthy"
  scheduler      = "LeastConnection"
}

resource "baiducloud_blb_listener" "default_HTTP" {
  blb_id        = baiducloud_blb.default.id
  listener_port = 127
  backend_port  = 127
  protocol      = "HTTP"
  scheduler     = "LeastConnection"
}


data "baiducloud_blb_listeners" "default_TCP" {
  blb_id        = baiducloud_blb.default.id
  protocol      = baiducloud_blb_listener.default_TCP.protocol
  listener_port = baiducloud_blb_listener.default_TCP.listener_port

  filter {
    name = "protocol"
    values = ["TCP"]
  }
}

data "baiducloud_blb_listeners" "default_UDP" {
  blb_id        = baiducloud_blb.default.id
  protocol      = baiducloud_blb_listener.default_UDP.protocol
  listener_port = baiducloud_blb_listener.default_UDP.listener_port

  filter {
    name = "protocol"
    values = ["UDP"]
  }
}

data "baiducloud_blb_listeners" "default_HTTP" {
  blb_id        = baiducloud_blb.default.id
  protocol      = baiducloud_blb_listener.default_HTTP.protocol
  listener_port = baiducloud_blb_listener.default_HTTP.listener_port

  filter {
    name = "protocol"
    values = ["HTTP"]
  }
}

`, name)
}
