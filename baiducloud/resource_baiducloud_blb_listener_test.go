package baiducloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccBLBListenerResourceType = "baiducloud_blb_listener"
	testAccBLBListenerResourceName = testAccBLBListenerResourceType + "." + BaiduCloudTestResourceName
)

func TestAccBaiduCloudBLBListener_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccBLBListenerDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccBLBHTTPListenerConfigBasic(BaiduCloudTestResourceTypeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBLBListenerResourceName),
					resource.TestCheckResourceAttr(testAccBLBListenerResourceName, "listener_port", "129"),
					resource.TestCheckResourceAttr(testAccBLBListenerResourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(testAccBLBListenerResourceName, "scheduler", "RoundRobin"),
				),
			},
			{
				Config: testAccBLBHTTPListenerConfigBasicUpdate(BaiduCloudTestResourceTypeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBLBListenerResourceName),
					resource.TestCheckResourceAttr(testAccBLBListenerResourceName, "listener_port", "129"),
					resource.TestCheckResourceAttr(testAccBLBListenerResourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(testAccBLBListenerResourceName, "scheduler", "RoundRobin"),
				),
			},
		},
	})
}

func TestAccBaiduCloudBLBListener_TCPListener(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccBLBListenerDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccBLBTCPListenerConfig(BaiduCloudTestResourceTypeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBLBListenerResourceName),
					resource.TestCheckResourceAttr(testAccBLBListenerResourceName, "listener_port", "124"),
					resource.TestCheckResourceAttr(testAccBLBListenerResourceName, "protocol", "TCP"),
					resource.TestCheckResourceAttr(testAccBLBListenerResourceName, "scheduler", "LeastConnection"),
				),
			},
			{
				Config: testAccBLBTCPListenerConfigUpdate(BaiduCloudTestResourceTypeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBLBListenerResourceName),
					resource.TestCheckResourceAttr(testAccBLBListenerResourceName, "listener_port", "124"),
					resource.TestCheckResourceAttr(testAccBLBListenerResourceName, "protocol", "TCP"),
					resource.TestCheckResourceAttr(testAccBLBListenerResourceName, "scheduler", "RoundRobin"),
				),
			},
		},
	})
}

func TestAccBaiduCloudBLBListener_UDPListener(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccBLBListenerDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccBLBUDPListenerConfig(BaiduCloudTestResourceTypeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBLBListenerResourceName),
					resource.TestCheckResourceAttr(testAccBLBListenerResourceName, "listener_port", "125"),
					resource.TestCheckResourceAttr(testAccBLBListenerResourceName, "protocol", "UDP"),
					resource.TestCheckResourceAttr(testAccBLBListenerResourceName, "scheduler", "LeastConnection"),
				),
			},
			{
				Config: testAccBLBUDPListenerConfigUpdate(BaiduCloudTestResourceTypeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBLBListenerResourceName),
					resource.TestCheckResourceAttr(testAccBLBListenerResourceName, "listener_port", "125"),
					resource.TestCheckResourceAttr(testAccBLBListenerResourceName, "protocol", "UDP"),
					resource.TestCheckResourceAttr(testAccBLBListenerResourceName, "scheduler", "RoundRobin"),
				),
			},
		},
	})
}

func TestAccBaiduCloudBLBListener_HTTPListener(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccBLBListenerDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccBLBHTTPListenerConfig(BaiduCloudTestResourceTypeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBLBListenerResourceName),
					resource.TestCheckResourceAttr(testAccBLBListenerResourceName, "listener_port", "126"),
					resource.TestCheckResourceAttr(testAccBLBListenerResourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(testAccBLBListenerResourceName, "scheduler", "LeastConnection"),
				),
			},
			{
				Config: testAccBLBHTTPListenerConfigUpdate(BaiduCloudTestResourceTypeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBLBListenerResourceName),
					resource.TestCheckResourceAttr(testAccBLBListenerResourceName, "listener_port", "126"),
					resource.TestCheckResourceAttr(testAccBLBListenerResourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(testAccBLBListenerResourceName, "scheduler", "RoundRobin"),
				),
			},
		},
	})
}

func testAccBLBListenerDestory(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	blbService := BLBService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccBLBListenerResourceType {
			continue
		}

		blbId := rs.Primary.Attributes["blb_id"]
		protocol := rs.Primary.Attributes["protocol"]
		port, _ := strconv.Atoi(rs.Primary.Attributes["listener_port"])
		_, err := blbService.DescribeListener(blbId, protocol, port)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}

	return nil
}

func testAccBLBHTTPListenerConfigBasic(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_vpc" "default" {
  name        = var.name
  description = "created-by-terraform"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = var.name
  zone_name   = "cn-bj-a"
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created-by-terraform"
}

resource "baiducloud_blb" "default" {
  name        = var.name
  description = "created-by-terraform"
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}

resource "baiducloud_blb_listener" "default" {
  blb_id        = baiducloud_blb.default.id
  listener_port = 129
  backend_port  = 129
  protocol      = "HTTP"
  scheduler     = "RoundRobin"

}
`, name)
}

func testAccBLBHTTPListenerConfigBasicUpdate(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_vpc" "default" {
  name        = var.name
  description = "created-by-terraform"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = var.name
  zone_name   = "cn-bj-a"
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created-by-terraform"
}

resource "baiducloud_blb" "default" {
  name        = var.name
  description = "created-by-terraform"
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}

resource "baiducloud_blb_listener" "default" {
  blb_id        = baiducloud_blb.default.id
  listener_port = 129
  backend_port  = 129
  protocol      = "HTTP"
  scheduler     = "RoundRobin"
}
`, name)
}

func testAccBLBTCPListenerConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_vpc" "default" {
  name        = var.name
  description = "created-by-terraform"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = var.name
  zone_name   = "cn-bj-a"
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created-by-terraform"
}

resource "baiducloud_blb" "default" {
  name        = var.name
  description = "created-by-terraform"
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}

resource "baiducloud_blb_listener" "default" {
  blb_id        = baiducloud_blb.default.id
  listener_port = 124
  backend_port  = 124
  protocol      = "TCP"
  scheduler     = "LeastConnection"
}
`, name)
}

func testAccBLBTCPListenerConfigUpdate(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_vpc" "default" {
  name        = var.name
  description = "created-by-terraform"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = var.name
  zone_name   = "cn-bj-a"
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created-by-terraform"
}

resource "baiducloud_blb" "default" {
  name        = var.name
  description = "created-by-terraform"
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}


resource "baiducloud_blb_listener" "default" {
  blb_id               = baiducloud_blb.default.id
  listener_port        = 124
  backend_port  = 124
  protocol             = "TCP"
  scheduler            = "RoundRobin"

}
`, name)
}

func testAccBLBUDPListenerConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_vpc" "default" {
  name        = var.name
  description = "created-by-terraform"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = var.name
  zone_name   = "cn-bj-a"
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created-by-terraform"
}

resource "baiducloud_blb" "default" {
  name        = var.name
  description = "created-by-terraform"
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}

resource "baiducloud_blb_listener" "default" {
  blb_id        = baiducloud_blb.default.id
  listener_port = 125  
  backend_port  = 125
  protocol      = "UDP"
  health_check_string  = "healthy"
  scheduler     = "LeastConnection"
}
`, name)
}

func testAccBLBUDPListenerConfigUpdate(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_vpc" "default" {
  name        = var.name
  description = "created-by-terraform"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = var.name
  zone_name   = "cn-bj-a"
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created-by-terraform"
}

resource "baiducloud_blb" "default" {
  name        = var.name
  description = "created-by-terraform"
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}

resource "baiducloud_blb_listener" "default" {
  blb_id        = baiducloud_blb.default.id
  listener_port = 125
  protocol      = "UDP"
  scheduler     = "RoundRobin"

}
`, name)
}

func testAccBLBHTTPListenerConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "baiducloud_vpc" "default" {
  name        = var.name
  description = "created-by-terraform"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = var.name
  zone_name   = "cn-bj-a"
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created-by-terraform"
}

resource "baiducloud_blb" "default" {
  depends_on  = [baiducloud_instance.default]
  name        = var.name
  description = "created-by-terraform"
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}

resource "baiducloud_blb_listener" "default" {
  blb_id        = baiducloud_blb.default.id
  listener_port = 126
  backend_port  = 126
  protocol      = "HTTP"
  scheduler     = "LeastConnection"
}
`, name)
}

func testAccBLBHTTPListenerConfigUpdate(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_vpc" "default" {
  name        = var.name
  description = "created-by-terraform"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = var.name
  zone_name   = "cn-bj-a"
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created-by-terraform"
}

resource "baiducloud_blb" "default" {
  name        = var.name
  description = "created-by-terraform"
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}


resource "baiducloud_blb_listener" "default" {
  blb_id        = baiducloud_blb.default.id
  listener_port = 126
  protocol      = "HTTP"
  scheduler     = "RoundRobin"

}
`, name)
}
