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
	testAccAppBLBListenerResourceType = "baiducloud_appblb_listener"
	testAccAppBLBListenerResourceName = testAccAppBLBListenerResourceType + "." + BaiduCloudTestResourceName
)

func TestAccBaiduCloudAppBLBListener_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccAppBLBListenerDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccAppBLBHTTPListenerConfigBasic(BaiduCloudTestResourceTypeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBListenerResourceName),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "listener_port", "129"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "scheduler", "RoundRobin"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "policies.#", "1"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "keep_session", "true"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "x_forwarded_for", "false"),
					resource.TestCheckResourceAttrSet(testAccAppBLBListenerResourceName, "keep_session_timeout"),
				),
			},
			{
				Config: testAccAppBLBHTTPListenerConfigBasicUpdate(BaiduCloudTestResourceTypeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBListenerResourceName),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "listener_port", "129"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "scheduler", "RoundRobin"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "policies.#", "0"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "keep_session", "true"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "x_forwarded_for", "false"),
					resource.TestCheckResourceAttrSet(testAccAppBLBListenerResourceName, "keep_session_timeout"),
				),
			},
		},
	})
}

func TestAccBaiduCloudAppBLBListener_TCPListener(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccAppBLBListenerDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccAppBLBTCPListenerConfig(BaiduCloudTestResourceTypeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBListenerResourceName),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "listener_port", "124"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "protocol", "TCP"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "scheduler", "LeastConnection"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "tcp_session_timeout", "900"),
				),
			},
			{
				Config: testAccAppBLBTCPListenerConfigUpdate(BaiduCloudTestResourceTypeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBListenerResourceName),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "listener_port", "124"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "protocol", "TCP"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "scheduler", "RoundRobin"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "policies.#", "1"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "tcp_session_timeout", "1000"),
				),
			},
		},
	})
}

func TestAccBaiduCloudAppBLBListener_UDPListener(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccAppBLBListenerDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccAppBLBUDPListenerConfig(BaiduCloudTestResourceTypeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBListenerResourceName),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "listener_port", "125"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "protocol", "UDP"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "scheduler", "LeastConnection"),
				),
			},
			{
				Config: testAccAppBLBUDPListenerConfigUpdate(BaiduCloudTestResourceTypeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBListenerResourceName),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "listener_port", "125"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "protocol", "UDP"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "scheduler", "RoundRobin"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "policies.#", "1"),
				),
			},
		},
	})
}

func TestAccBaiduCloudAppBLBListener_HTTPListener(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccAppBLBListenerDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccAppBLBHTTPListenerConfig(BaiduCloudTestResourceTypeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBListenerResourceName),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "listener_port", "126"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "scheduler", "LeastConnection"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "keep_session", "true"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "x_forwarded_for", "false"),
					resource.TestCheckResourceAttrSet(testAccAppBLBListenerResourceName, "keep_session_timeout"),
				),
			},
			{
				Config: testAccAppBLBHTTPListenerConfigUpdate(BaiduCloudTestResourceTypeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBListenerResourceName),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "listener_port", "126"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "scheduler", "RoundRobin"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "policies.#", "1"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "keep_session", "true"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "x_forwarded_for", "false"),
					resource.TestCheckResourceAttrSet(testAccAppBLBListenerResourceName, "keep_session_timeout"),
				),
			},
		},
	})
}

func TestAccBaiduCloudAppBLBListener_HTTPSListener(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccAppBLBListenerDestory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppBLBHTTPSListenerConfig(BaiduCloudTestResourceTypeNameAppblbListener),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBListenerResourceName),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "listener_port", "130"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "protocol", "HTTPS"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "scheduler", "LeastConnection"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "cert_ids.#", "1"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "keep_session", "true"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "encryption_protocols.#", "3"),
					resource.TestCheckResourceAttrSet(testAccAppBLBListenerResourceName, "keep_session_type"),
					resource.TestCheckResourceAttrSet(testAccAppBLBListenerResourceName, "x_forwarded_for"),
					resource.TestCheckResourceAttrSet(testAccAppBLBListenerResourceName, "server_timeout"),
					resource.TestCheckNoResourceAttr(testAccAppBLBListenerResourceName, "redirect_port"),
				),
			},
			{
				Config: testAccAppBLBHTTPSListenerConfigUpdate(BaiduCloudTestResourceTypeNameAppblbListener),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBListenerResourceName),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "listener_port", "130"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "protocol", "HTTPS"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "scheduler", "RoundRobin"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "cert_ids.#", "2"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "keep_session", "true"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "encryption_protocols.#", "3"),
					resource.TestCheckResourceAttrSet(testAccAppBLBListenerResourceName, "keep_session_type"),
					resource.TestCheckResourceAttrSet(testAccAppBLBListenerResourceName, "x_forwarded_for"),
					resource.TestCheckResourceAttrSet(testAccAppBLBListenerResourceName, "server_timeout"),
					resource.TestCheckNoResourceAttr(testAccAppBLBListenerResourceName, "redirect_port"),
				),
			},
		},
	})
}

func TestAccBaiduCloudAppBLBListener_SSLListener(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccAppBLBListenerDestory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppBLBSSLListenerConfig(BaiduCloudTestResourceTypeNameAppblbListener),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBListenerResourceName),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "listener_port", "131"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "protocol", "SSL"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "scheduler", "LeastConnection"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "cert_ids.#", "2"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "encryption_protocols.#", "3"),
				),
			},
			{
				Config: testAccAppBLBSSLListenerConfigUpdate(BaiduCloudTestResourceTypeNameAppblbListener),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBListenerResourceName),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "listener_port", "131"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "protocol", "SSL"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "scheduler", "RoundRobin"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "cert_ids.#", "2"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "encryption_protocols.#", "3"),
				),
			},
		},
	})
}

func testAccAppBLBListenerDestory(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	appblbService := APPBLBService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccAppBLBListenerResourceType {
			continue
		}

		blbId := rs.Primary.Attributes["blb_id"]
		protocol := rs.Primary.Attributes["protocol"]
		port, _ := strconv.Atoi(rs.Primary.Attributes["listener_port"])
		_, err := appblbService.DescribeListener(blbId, protocol, port)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}

	return nil
}

func testAccAppBLBHTTPListenerConfigBasic(name string) string {
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
  name                  = var.name
  image_id              = data.baiducloud_images.default.images.0.id
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  billing = {
    payment_timing = "Postpaid"
  }
}

resource "baiducloud_vpc" "default" {
  name        = var.name
  description = "created-by-terraform"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = var.name
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created-by-terraform"
}

resource "baiducloud_appblb" "default" {
  depends_on  = [baiducloud_instance.default]
  name        = var.name
  description = "created-by-terraform"
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}

resource "baiducloud_appblb_server_group" "default" {
  name         = var.name
  description  = "created-by-terraform"
  blb_id       = baiducloud_appblb.default.id

  port_list {
    port = 70
    type = "HTTP"
    health_check = "HTTP"
  }
}

resource "baiducloud_appblb_listener" "default" {
  blb_id        = baiducloud_appblb.default.id
  listener_port = 129
  protocol      = "HTTP"
  scheduler     = "RoundRobin"
  keep_session  = true

  policies {
    description         = "created-by-terraform"
    app_server_group_id = baiducloud_appblb_server_group.default.id
    backend_port        = 70
    priority            = 50

    rule_list {
      key   = "host"
      value = "baidu.com"
    }
  }
}
`, name)
}

func testAccAppBLBHTTPListenerConfigBasicUpdate(name string) string {
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
  name                  = var.name
  image_id              = data.baiducloud_images.default.images.0.id
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  billing = {
    payment_timing = "Postpaid"
  }
}

resource "baiducloud_vpc" "default" {
  name        = var.name
  description = "created-by-terraform"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = var.name
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created-by-terraform"
}

resource "baiducloud_appblb" "default" {
  depends_on  = [baiducloud_instance.default]
  name        = var.name
  description = "created-by-terraform"
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}

resource "baiducloud_appblb_listener" "default" {
  blb_id        = baiducloud_appblb.default.id
  listener_port = 129
  protocol      = "HTTP"
  scheduler     = "RoundRobin"
  keep_session  = true
}
`, name)
}

func testAccAppBLBTCPListenerConfig(name string) string {
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
  name                  = var.name
  image_id              = data.baiducloud_images.default.images.0.id
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  billing = {
    payment_timing = "Postpaid"
  }
}

resource "baiducloud_vpc" "default" {
  name        = var.name
  description = "created-by-terraform"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = var.name
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created-by-terraform"
}

resource "baiducloud_appblb" "default" {
  depends_on  = [baiducloud_instance.default]
  name        = var.name
  description = "created-by-terraform"
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}

resource "baiducloud_appblb_listener" "default" {
  blb_id        = baiducloud_appblb.default.id
  listener_port = 124
  protocol      = "TCP"
  scheduler     = "LeastConnection"
}
`, name)
}

func testAccAppBLBTCPListenerConfigUpdate(name string) string {
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
  name                  = var.name
  image_id              = data.baiducloud_images.default.images.0.id
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  billing = {
    payment_timing = "Postpaid"
  }
}

resource "baiducloud_vpc" "default" {
  name        = var.name
  description = "created-by-terraform"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = var.name
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created-by-terraform"
}

resource "baiducloud_appblb" "default" {
  depends_on  = [baiducloud_instance.default]
  name        = var.name
  description = "created-by-terraform"
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}

resource "baiducloud_appblb_server_group" "default" {
  name        = var.name
  description = "created-by-terraform"
  blb_id      = baiducloud_appblb.default.id

  port_list {
    port = 68
    type = "TCP"
    health_check = "TCP"
  }
}

resource "baiducloud_appblb_listener" "default" {
  blb_id               = baiducloud_appblb.default.id
  listener_port        = 124
  protocol             = "TCP"
  scheduler            = "RoundRobin"
  tcp_session_timeout  = 1000

  policies {
    description         = "created-by-terraform"
    app_server_group_id = baiducloud_appblb_server_group.default.id
    backend_port        = 68
    priority            = 50

    rule_list {
      key   = "*"
      value = "*"
    }
  }
}
`, name)
}

func testAccAppBLBUDPListenerConfig(name string) string {
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
  name                  = var.name
  image_id              = data.baiducloud_images.default.images.0.id
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  billing = {
    payment_timing = "Postpaid"
  }
}

resource "baiducloud_vpc" "default" {
  name        = var.name
  description = "created-by-terraform"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = var.name
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created-by-terraform"
}

resource "baiducloud_appblb" "default" {
  depends_on  = [baiducloud_instance.default]
  name        = var.name
  description = "created-by-terraform"
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}

resource "baiducloud_appblb_listener" "default" {
  blb_id        = baiducloud_appblb.default.id
  listener_port = 125
  protocol      = "UDP"
  scheduler     = "LeastConnection"
}
`, name)
}

func testAccAppBLBUDPListenerConfigUpdate(name string) string {
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
  name                  = var.name
  image_id              = data.baiducloud_images.default.images.0.id
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  billing = {
    payment_timing = "Postpaid"
  }
}

resource "baiducloud_vpc" "default" {
  name        = var.name
  description = "created-by-terraform"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = var.name
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created-by-terraform"
}

resource "baiducloud_appblb" "default" {
  depends_on  = [baiducloud_instance.default]
  name        = var.name
  description = "created-by-terraform"
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}

resource "baiducloud_appblb_server_group" "default" {
  name        = var.name
  description = "created-by-terraform"
  blb_id      = baiducloud_appblb.default.id

  port_list {
    port = 66
    type = "UDP"
    health_check = "UDP"
    udp_health_check_string = "baidu.com"
  }
}

resource "baiducloud_appblb_listener" "default" {
  blb_id        = baiducloud_appblb.default.id
  listener_port = 125
  protocol      = "UDP"
  scheduler     = "RoundRobin"

  policies {
    description         = "created-by-terraform"
    app_server_group_id = baiducloud_appblb_server_group.default.id
    backend_port        = 66
    priority            = 50

    rule_list {
      key   = "*"
      value = "*"
    }
  }
}
`, name)
}

func testAccAppBLBHTTPListenerConfig(name string) string {
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
  name                  = var.name
  image_id              = data.baiducloud_images.default.images.0.id
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  billing = {
    payment_timing = "Postpaid"
  }
}

resource "baiducloud_vpc" "default" {
  name        = var.name
  description = "created-by-terraform"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = var.name
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created-by-terraform"
}

resource "baiducloud_appblb" "default" {
  depends_on  = [baiducloud_instance.default]
  name        = var.name
  description = "created-by-terraform"
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}

resource "baiducloud_appblb_listener" "default" {
  blb_id        = baiducloud_appblb.default.id
  listener_port = 126
  protocol      = "HTTP"
  scheduler     = "LeastConnection"
  keep_session  = true
}
`, name)
}

func testAccAppBLBHTTPListenerConfigUpdate(name string) string {
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
  name                  = var.name
  image_id              = data.baiducloud_images.default.images.0.id
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  billing = {
    payment_timing = "Postpaid"
  }
}

resource "baiducloud_vpc" "default" {
  name        = var.name
  description = "created-by-terraform"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = var.name
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created-by-terraform"
}

resource "baiducloud_appblb" "default" {
  depends_on  = [baiducloud_instance.default]
  name        = var.name
  description = "created-by-terraform"
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}

resource "baiducloud_appblb_server_group" "default" {
  name         = var.name
  description  = "created-by-terraform"
  blb_id       = baiducloud_appblb.default.id

  port_list {
    port = 67
    type = "HTTP"
    health_check = "HTTP"
  }
}

resource "baiducloud_appblb_listener" "default" {
  blb_id        = baiducloud_appblb.default.id
  listener_port = 126
  protocol      = "HTTP"
  scheduler     = "RoundRobin"
  keep_session  = true

  policies {
    description         = "created-by-terraform"
    app_server_group_id = baiducloud_appblb_server_group.default.id
    backend_port        = 67
    priority            = 50

    rule_list {
      key   = "*"
      value = "*"
    }
  }
}
`, name)
}

func testAccAppBLBHTTPSListenerConfig(name string) string {
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
  name                  = var.name
  image_id              = data.baiducloud_images.default.images.0.id
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  billing = {
    payment_timing = "Postpaid"
  }
}

resource "baiducloud_vpc" "default" {
  name        = var.name
  description = "created-by-terraform"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = var.name
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created-by-terraform"
}
resource "baiducloud_appblb" "default" {
  depends_on  = [baiducloud_instance.default]
  name        = var.name
  description = "created-by-terraform"
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}

resource "baiducloud_cert" "default" {
  cert_name         = var.name
  cert_server_data  = "-----BEGIN CERTIFICATE-----\n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\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIIErjCCA5agAwIBAgIQBYAmfwbylVM0jhwYWl7uLjANBgkqhkiG9w0BAQsFADBhMQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBDQTAeFw0xNzEyMDgxMjI4MjZaFw0yNzEyMDgxMjI4MjZaMHIxCzAJBgNVBAYTAkNOMSUwIwYDVQQKExxUcnVzdEFzaWEgVGVjaG5vbG9naWVzLCBJbmMuMR0wGwYDVQQLExREb21haW4gVmFsaWRhdGVkIFNTTDEdMBsGA1UEAxMUVHJ1c3RBc2lhIFRMUyBSU0EgQ0EwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCgWa9X+ph+wAm8Yh1Fk1MjKbQ5QwBOOKVaZR/OfCh+F6f93u7vZHGcUU/lvVGgUQnbzJhR1UV2epJae+m7cxnXIKdD0/VS9btAgwJszGFvwoqXeaCqFoP71wPmXjjUwLT70+qvX4hdyYfOJcjeTz5QKtg8zQwxaK9x4JT9CoOmoVdVhEBAiD3DwR5fFgOHDwwGxdJWVBvktnoAzjdTLXDdbSVC5jZ0u8oq9BiTDv7jAlsB5F8aZgvSZDOQeFrwaOTbKWSEInEhnchKZTD1dz6aBlk1xGEI5PZWAnVAba/ofH33ktymaTDsE6xRDnW97pDkimCRak6CEbfe3dXw6OV5AgMBAAGjggFPMIIBSzAdBgNVHQ4EFgQUf9OZ86BHDjEAVlYijrfMnt3KAYowHwYDVR0jBBgwFoAUA95QNVbRTLtm8KPiGxvDl7I90VUwDgYDVR0PAQH/BAQDAgGGMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjASBgNVHRMBAf8ECDAGAQH/AgEAMDQGCCsGAQUFBwEBBCgwJjAkBggrBgEFBQcwAYYYaHR0cDovL29jc3AuZGlnaWNlcnQuY29tMEIGA1UdHwQ7MDkwN6A1oDOGMWh0dHA6Ly9jcmwzLmRpZ2ljZXJ0LmNvbS9EaWdpQ2VydEdsb2JhbFJvb3RDQS5jcmwwTAYDVR0gBEUwQzA3BglghkgBhv1sAQIwKjAoBggrBgEFBQcCARYcaHR0cHM6Ly93d3cuZGlnaWNlcnQuY29tL0NQUzAIBgZngQwBAgEwDQYJKoZIhvcNAQELBQADggEBAK3dVOj5dlv4MzK2i233lDYvyJ3slFY2X2HKTYGte8nbK6i5/fsDImMYihAkp6VaNY/en8WZ5qcrQPVLuJrJDSXT04NnMeZOQDUoj/NHAmdfCBB/h1bZ5OGK6Sf1h5Yx/5wR4f3TUoPgGlnU7EuPISLNdMRiDrXntcImDAiRvkh5GJuH4YCVE6XEntqaNIgGkRwxKSgnU3Id3iuFbW9FUQ9Qqtb1GX91AJ7i4153TikGgYCdwYkBURD8gSVe8OAco6IfZOYt/TEwii1Ivi1CqnuUlWpsF1LdQNIdfbW3TSe0BhQa7ifbVIfvPWHYOu3rkg1ZeMo6XRU9B4n5VyJYRmE=\n-----END CERTIFICATE-----"
  cert_private_data = "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAkZCmcbCZyzdYQ379jndWgcKB9Q0/ITUYE2sUlh39HwyAJYP9DP2fStqF9pW+KVYJDPlhdFbYTrUqp+dO9001Hj3bL2JKMwx6zKmKtjUecpSkFtOYKAF+e6odhcS5HUTq0maydjvLSp1kPKd0Q6n4GW14+gkjoUvMVWK59FmjfO47yqETTWK/zA6JAp4ZMXuI6ujFlsnoxY4nFD0+HKiPNkYaFwDr0AUEwdAaDG/fp7/AP2xnQLdGXVn/gflulflUmA9tqwdF2ul9UmsY8iZcbPUbO0EtC78LOGPCgC5SsYVioUCzGO1HCbnDz3jvshAsZ/EMR07KBzqretZwEnuTSQIDAQABAoIBAAzBl4cfWfLljY4TVbFY7ZNJ0i1Wilbkz2XQPJ8aegFGYqp8TROI3EnpKX6I89UCgvYzRSI2rsEC/lMgIZrpa1i+70jRPRMJKm+/VyENjvatO6NRH/ni26HcWrb2HN90Qnx1XyPzrHvZnBxL876EPseCVkIvGoNliulb+/4Y/DXpNthA28UOB9RafPsEoDNinrTqlZf0gNLxm1LOgcj/NEqsDwuwzwfCky9GAhQgZpwic2IAEwKoCbfeRNNraVgG+IdCC8Nn3/uMcy9Zft3fV7xNE6HdfkW1SKnEvN+sFxKhH7ad0FNtaE+kSAcxTWXOg/xErvUBIcDrZv23BgN4JVMCgYEAwiNb00eRuBcPTHAaEb9JqrFRtUlqLnFJe1ang1QRfn+FrlTnijGACTjEFpzaXavaGNKi+To8OZjSTL2OW6ewEwSA9siPXUkq3ldPj5uPIhr80Jn1Ox/K5+X5ZBkQg8Iw9GIY6P6Kgf/prihVIbGZVNa0U/8H/1RvQIBxvA21dfMCgYEAv/L8iGiSwcgqMv0NTzfiW4fA9L7yLE04mfs9QI1V/uHPX5ufb/Y3LCS1RSuOdjrCdD2Ru7OKMi1v7mwPg1+NJBZjLIlCw/oVCJZabd8KGXZUNSH+PNuQAbIGdotEpO+LPgVgwi4ovrx6oJYEED/1FFjfU2bBFfuZtrDBWz2yNNMCgYEAvdoKQJHq5RZX9a5jMBvbFLwXZawH1Kcg7ycM5hdejFB1EMkjLTe/OEV1LY/y1EvtGv1SN1xF7SWP81AkWWmhfNeYrr3vxZB6Bbloqs27qeSue+kzssAik6mIu+TvC4rqiPMt3RyfowX7Jj93EV42zoqxCruKvJ17tp5lmzvkyxUCgYBRN60mwqimGd3RKUWCaXD7rZs1c73ghOQYMzgdoi/q4vztxVlW9GUv5nBUzjM/T2mL6alKNJOa26LqzQpbWgjMZjScWY/IgH553bRxnNgXIfxLZxC+C2EJdpxJeHAZIcpW+cuRHhrbacCxRgh+H7HBZEFKdsXoWUcXB/8obhiDRQKBgCwOE+1hfrV7/gFaMBWSML1n+LVV2ns80jCDtkhN9yF+9iJTjMwW4wuvFx8t8o2XICOwJPog4IvXFJLVZeed/zhgqe4qImHRW0aMYGyGEpgkLtHIFFFCxGd57Df/qEbUL55LU53rlCv2QKVBBs/6XDkiVRBk8izT7ihF2U8qb6t4\n-----END RSA PRIVATE KEY-----"
}

resource "baiducloud_appblb_listener" "default" {
  blb_id               = baiducloud_appblb.default.id
  listener_port        = 130
  protocol             = "HTTPS"
  scheduler            = "LeastConnection"
  keep_session         = true
  cert_ids             = [baiducloud_cert.default.id]
  encryption_protocols = ["sslv3", "tlsv10", "tlsv11"]
  encryption_type      = "userDefind"
}
`, name)
}
func testAccAppBLBHTTPSListenerConfigUpdate(name string) string {
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
  name                  = var.name
  image_id              = data.baiducloud_images.default.images.0.id
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  billing = {
    payment_timing = "Postpaid"
  }
}

resource "baiducloud_vpc" "default" {
  name        = "var.name
  description = "created-by-terraform"
  cidr        = "192.168.0.0/24"
}
resource "baiducloud_subnet" "default" {
  name        = var.name
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created-by-terraform"
}
resource "baiducloud_appblb" "default" {
  depends_on  = [baiducloud_instance.default]
  name        = var.name
  description = "created-by-terraform"
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}

resource "baiducloud_cert" "default" {
  cert_name         = var.name
  cert_server_data  = "-----BEGIN CERTIFICATE-----\n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\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIIErjCCA5agAwIBAgIQBYAmfwbylVM0jhwYWl7uLjANBgkqhkiG9w0BAQsFADBhMQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBDQTAeFw0xNzEyMDgxMjI4MjZaFw0yNzEyMDgxMjI4MjZaMHIxCzAJBgNVBAYTAkNOMSUwIwYDVQQKExxUcnVzdEFzaWEgVGVjaG5vbG9naWVzLCBJbmMuMR0wGwYDVQQLExREb21haW4gVmFsaWRhdGVkIFNTTDEdMBsGA1UEAxMUVHJ1c3RBc2lhIFRMUyBSU0EgQ0EwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCgWa9X+ph+wAm8Yh1Fk1MjKbQ5QwBOOKVaZR/OfCh+F6f93u7vZHGcUU/lvVGgUQnbzJhR1UV2epJae+m7cxnXIKdD0/VS9btAgwJszGFvwoqXeaCqFoP71wPmXjjUwLT70+qvX4hdyYfOJcjeTz5QKtg8zQwxaK9x4JT9CoOmoVdVhEBAiD3DwR5fFgOHDwwGxdJWVBvktnoAzjdTLXDdbSVC5jZ0u8oq9BiTDv7jAlsB5F8aZgvSZDOQeFrwaOTbKWSEInEhnchKZTD1dz6aBlk1xGEI5PZWAnVAba/ofH33ktymaTDsE6xRDnW97pDkimCRak6CEbfe3dXw6OV5AgMBAAGjggFPMIIBSzAdBgNVHQ4EFgQUf9OZ86BHDjEAVlYijrfMnt3KAYowHwYDVR0jBBgwFoAUA95QNVbRTLtm8KPiGxvDl7I90VUwDgYDVR0PAQH/BAQDAgGGMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjASBgNVHRMBAf8ECDAGAQH/AgEAMDQGCCsGAQUFBwEBBCgwJjAkBggrBgEFBQcwAYYYaHR0cDovL29jc3AuZGlnaWNlcnQuY29tMEIGA1UdHwQ7MDkwN6A1oDOGMWh0dHA6Ly9jcmwzLmRpZ2ljZXJ0LmNvbS9EaWdpQ2VydEdsb2JhbFJvb3RDQS5jcmwwTAYDVR0gBEUwQzA3BglghkgBhv1sAQIwKjAoBggrBgEFBQcCARYcaHR0cHM6Ly93d3cuZGlnaWNlcnQuY29tL0NQUzAIBgZngQwBAgEwDQYJKoZIhvcNAQELBQADggEBAK3dVOj5dlv4MzK2i233lDYvyJ3slFY2X2HKTYGte8nbK6i5/fsDImMYihAkp6VaNY/en8WZ5qcrQPVLuJrJDSXT04NnMeZOQDUoj/NHAmdfCBB/h1bZ5OGK6Sf1h5Yx/5wR4f3TUoPgGlnU7EuPISLNdMRiDrXntcImDAiRvkh5GJuH4YCVE6XEntqaNIgGkRwxKSgnU3Id3iuFbW9FUQ9Qqtb1GX91AJ7i4153TikGgYCdwYkBURD8gSVe8OAco6IfZOYt/TEwii1Ivi1CqnuUlWpsF1LdQNIdfbW3TSe0BhQa7ifbVIfvPWHYOu3rkg1ZeMo6XRU9B4n5VyJYRmE=\n-----END CERTIFICATE-----"
  cert_private_data = "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAkZCmcbCZyzdYQ379jndWgcKB9Q0/ITUYE2sUlh39HwyAJYP9DP2fStqF9pW+KVYJDPlhdFbYTrUqp+dO9001Hj3bL2JKMwx6zKmKtjUecpSkFtOYKAF+e6odhcS5HUTq0maydjvLSp1kPKd0Q6n4GW14+gkjoUvMVWK59FmjfO47yqETTWK/zA6JAp4ZMXuI6ujFlsnoxY4nFD0+HKiPNkYaFwDr0AUEwdAaDG/fp7/AP2xnQLdGXVn/gflulflUmA9tqwdF2ul9UmsY8iZcbPUbO0EtC78LOGPCgC5SsYVioUCzGO1HCbnDz3jvshAsZ/EMR07KBzqretZwEnuTSQIDAQABAoIBAAzBl4cfWfLljY4TVbFY7ZNJ0i1Wilbkz2XQPJ8aegFGYqp8TROI3EnpKX6I89UCgvYzRSI2rsEC/lMgIZrpa1i+70jRPRMJKm+/VyENjvatO6NRH/ni26HcWrb2HN90Qnx1XyPzrHvZnBxL876EPseCVkIvGoNliulb+/4Y/DXpNthA28UOB9RafPsEoDNinrTqlZf0gNLxm1LOgcj/NEqsDwuwzwfCky9GAhQgZpwic2IAEwKoCbfeRNNraVgG+IdCC8Nn3/uMcy9Zft3fV7xNE6HdfkW1SKnEvN+sFxKhH7ad0FNtaE+kSAcxTWXOg/xErvUBIcDrZv23BgN4JVMCgYEAwiNb00eRuBcPTHAaEb9JqrFRtUlqLnFJe1ang1QRfn+FrlTnijGACTjEFpzaXavaGNKi+To8OZjSTL2OW6ewEwSA9siPXUkq3ldPj5uPIhr80Jn1Ox/K5+X5ZBkQg8Iw9GIY6P6Kgf/prihVIbGZVNa0U/8H/1RvQIBxvA21dfMCgYEAv/L8iGiSwcgqMv0NTzfiW4fA9L7yLE04mfs9QI1V/uHPX5ufb/Y3LCS1RSuOdjrCdD2Ru7OKMi1v7mwPg1+NJBZjLIlCw/oVCJZabd8KGXZUNSH+PNuQAbIGdotEpO+LPgVgwi4ovrx6oJYEED/1FFjfU2bBFfuZtrDBWz2yNNMCgYEAvdoKQJHq5RZX9a5jMBvbFLwXZawH1Kcg7ycM5hdejFB1EMkjLTe/OEV1LY/y1EvtGv1SN1xF7SWP81AkWWmhfNeYrr3vxZB6Bbloqs27qeSue+kzssAik6mIu+TvC4rqiPMt3RyfowX7Jj93EV42zoqxCruKvJ17tp5lmzvkyxUCgYBRN60mwqimGd3RKUWCaXD7rZs1c73ghOQYMzgdoi/q4vztxVlW9GUv5nBUzjM/T2mL6alKNJOa26LqzQpbWgjMZjScWY/IgH553bRxnNgXIfxLZxC+C2EJdpxJeHAZIcpW+cuRHhrbacCxRgh+H7HBZEFKdsXoWUcXB/8obhiDRQKBgCwOE+1hfrV7/gFaMBWSML1n+LVV2ns80jCDtkhN9yF+9iJTjMwW4wuvFx8t8o2XICOwJPog4IvXFJLVZeed/zhgqe4qImHRW0aMYGyGEpgkLtHIFFFCxGd57Df/qEbUL55LU53rlCv2QKVBBs/6XDkiVRBk8izT7ihF2U8qb6t4\n-----END RSA PRIVATE KEY-----"
}

resource "baiducloud_appblb_listener" "default" {
  blb_id               = baiducloud_appblb.default.id
  listener_port        = 130
  protocol             = "HTTPS"
  scheduler            = "RoundRobin"
  keep_session         = true
  cert_ids             = [baiducloud_cert.default.id]
  encryption_protocols = ["tlsv10", "tlsv11", "tlsv12"]
  encryption_type      = "userDefind"
}
`, name)
}

func testAccAppBLBSSLListenerConfig(name string) string {
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
 name                  = var.name
 image_id              = data.baiducloud_images.default.images.0.id
 availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
 cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
 memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
 billing = {
   payment_timing = "Postpaid"
 }
}

resource "baiducloud_vpc" "default" {
 name        = var.name
 description = "created-by-terraform"
 cidr        = "192.168.0.0/24"
}
resource "baiducloud_subnet" "default" {
 name        = var.name
 zone_name   = data.baiducloud_zones.default.zones.0.zone_name
 cidr        = "192.168.0.0/24"
 vpc_id      = baiducloud_vpc.default.id
 description = "created-by-terraform"
}
resource "baiducloud_appblb" "default" {
 depends_on  = [baiducloud_instance.default]
 name        = var.name
 description = "created-by-terraform"
 vpc_id      = baiducloud_vpc.default.id
 subnet_id   = baiducloud_subnet.default.id
}

resource "baiducloud_cert" "default" {
  cert_name         = var.name
  cert_server_data  = "-----BEGIN CERTIFICATE-----\n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\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIIErjCCA5agAwIBAgIQBYAmfwbylVM0jhwYWl7uLjANBgkqhkiG9w0BAQsFADBhMQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBDQTAeFw0xNzEyMDgxMjI4MjZaFw0yNzEyMDgxMjI4MjZaMHIxCzAJBgNVBAYTAkNOMSUwIwYDVQQKExxUcnVzdEFzaWEgVGVjaG5vbG9naWVzLCBJbmMuMR0wGwYDVQQLExREb21haW4gVmFsaWRhdGVkIFNTTDEdMBsGA1UEAxMUVHJ1c3RBc2lhIFRMUyBSU0EgQ0EwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCgWa9X+ph+wAm8Yh1Fk1MjKbQ5QwBOOKVaZR/OfCh+F6f93u7vZHGcUU/lvVGgUQnbzJhR1UV2epJae+m7cxnXIKdD0/VS9btAgwJszGFvwoqXeaCqFoP71wPmXjjUwLT70+qvX4hdyYfOJcjeTz5QKtg8zQwxaK9x4JT9CoOmoVdVhEBAiD3DwR5fFgOHDwwGxdJWVBvktnoAzjdTLXDdbSVC5jZ0u8oq9BiTDv7jAlsB5F8aZgvSZDOQeFrwaOTbKWSEInEhnchKZTD1dz6aBlk1xGEI5PZWAnVAba/ofH33ktymaTDsE6xRDnW97pDkimCRak6CEbfe3dXw6OV5AgMBAAGjggFPMIIBSzAdBgNVHQ4EFgQUf9OZ86BHDjEAVlYijrfMnt3KAYowHwYDVR0jBBgwFoAUA95QNVbRTLtm8KPiGxvDl7I90VUwDgYDVR0PAQH/BAQDAgGGMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjASBgNVHRMBAf8ECDAGAQH/AgEAMDQGCCsGAQUFBwEBBCgwJjAkBggrBgEFBQcwAYYYaHR0cDovL29jc3AuZGlnaWNlcnQuY29tMEIGA1UdHwQ7MDkwN6A1oDOGMWh0dHA6Ly9jcmwzLmRpZ2ljZXJ0LmNvbS9EaWdpQ2VydEdsb2JhbFJvb3RDQS5jcmwwTAYDVR0gBEUwQzA3BglghkgBhv1sAQIwKjAoBggrBgEFBQcCARYcaHR0cHM6Ly93d3cuZGlnaWNlcnQuY29tL0NQUzAIBgZngQwBAgEwDQYJKoZIhvcNAQELBQADggEBAK3dVOj5dlv4MzK2i233lDYvyJ3slFY2X2HKTYGte8nbK6i5/fsDImMYihAkp6VaNY/en8WZ5qcrQPVLuJrJDSXT04NnMeZOQDUoj/NHAmdfCBB/h1bZ5OGK6Sf1h5Yx/5wR4f3TUoPgGlnU7EuPISLNdMRiDrXntcImDAiRvkh5GJuH4YCVE6XEntqaNIgGkRwxKSgnU3Id3iuFbW9FUQ9Qqtb1GX91AJ7i4153TikGgYCdwYkBURD8gSVe8OAco6IfZOYt/TEwii1Ivi1CqnuUlWpsF1LdQNIdfbW3TSe0BhQa7ifbVIfvPWHYOu3rkg1ZeMo6XRU9B4n5VyJYRmE=\n-----END CERTIFICATE-----"
  cert_private_data = "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAkZCmcbCZyzdYQ379jndWgcKB9Q0/ITUYE2sUlh39HwyAJYP9DP2fStqF9pW+KVYJDPlhdFbYTrUqp+dO9001Hj3bL2JKMwx6zKmKtjUecpSkFtOYKAF+e6odhcS5HUTq0maydjvLSp1kPKd0Q6n4GW14+gkjoUvMVWK59FmjfO47yqETTWK/zA6JAp4ZMXuI6ujFlsnoxY4nFD0+HKiPNkYaFwDr0AUEwdAaDG/fp7/AP2xnQLdGXVn/gflulflUmA9tqwdF2ul9UmsY8iZcbPUbO0EtC78LOGPCgC5SsYVioUCzGO1HCbnDz3jvshAsZ/EMR07KBzqretZwEnuTSQIDAQABAoIBAAzBl4cfWfLljY4TVbFY7ZNJ0i1Wilbkz2XQPJ8aegFGYqp8TROI3EnpKX6I89UCgvYzRSI2rsEC/lMgIZrpa1i+70jRPRMJKm+/VyENjvatO6NRH/ni26HcWrb2HN90Qnx1XyPzrHvZnBxL876EPseCVkIvGoNliulb+/4Y/DXpNthA28UOB9RafPsEoDNinrTqlZf0gNLxm1LOgcj/NEqsDwuwzwfCky9GAhQgZpwic2IAEwKoCbfeRNNraVgG+IdCC8Nn3/uMcy9Zft3fV7xNE6HdfkW1SKnEvN+sFxKhH7ad0FNtaE+kSAcxTWXOg/xErvUBIcDrZv23BgN4JVMCgYEAwiNb00eRuBcPTHAaEb9JqrFRtUlqLnFJe1ang1QRfn+FrlTnijGACTjEFpzaXavaGNKi+To8OZjSTL2OW6ewEwSA9siPXUkq3ldPj5uPIhr80Jn1Ox/K5+X5ZBkQg8Iw9GIY6P6Kgf/prihVIbGZVNa0U/8H/1RvQIBxvA21dfMCgYEAv/L8iGiSwcgqMv0NTzfiW4fA9L7yLE04mfs9QI1V/uHPX5ufb/Y3LCS1RSuOdjrCdD2Ru7OKMi1v7mwPg1+NJBZjLIlCw/oVCJZabd8KGXZUNSH+PNuQAbIGdotEpO+LPgVgwi4ovrx6oJYEED/1FFjfU2bBFfuZtrDBWz2yNNMCgYEAvdoKQJHq5RZX9a5jMBvbFLwXZawH1Kcg7ycM5hdejFB1EMkjLTe/OEV1LY/y1EvtGv1SN1xF7SWP81AkWWmhfNeYrr3vxZB6Bbloqs27qeSue+kzssAik6mIu+TvC4rqiPMt3RyfowX7Jj93EV42zoqxCruKvJ17tp5lmzvkyxUCgYBRN60mwqimGd3RKUWCaXD7rZs1c73ghOQYMzgdoi/q4vztxVlW9GUv5nBUzjM/T2mL6alKNJOa26LqzQpbWgjMZjScWY/IgH553bRxnNgXIfxLZxC+C2EJdpxJeHAZIcpW+cuRHhrbacCxRgh+H7HBZEFKdsXoWUcXB/8obhiDRQKBgCwOE+1hfrV7/gFaMBWSML1n+LVV2ns80jCDtkhN9yF+9iJTjMwW4wuvFx8t8o2XICOwJPog4IvXFJLVZeed/zhgqe4qImHRW0aMYGyGEpgkLtHIFFFCxGd57Df/qEbUL55LU53rlCv2QKVBBs/6XDkiVRBk8izT7ihF2U8qb6t4\n-----END RSA PRIVATE KEY-----"
}

resource "baiducloud_appblb_listener" "default" {
 blb_id               = baiducloud_appblb.default.id
 listener_port        = 131
 protocol             = "SSL"
 scheduler            = "LeastConnection"
 cert_ids             = [baiducloud_cert.default.id]
 encryption_protocols = ["sslv3", "tlsv10", "tlsv11"]
 encryption_type      = "userDefind"
}
`, name)
}
func testAccAppBLBSSLListenerConfigUpdate(name string) string {
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
 name                  = var.name
 image_id              = data.baiducloud_images.default.images.0.id
 availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
 cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
 memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
 billing = {
   payment_timing = "Postpaid"
 }
}

resource "baiducloud_vpc" "default" {
 name        = var.name
 description = "created-by-terraform"
 cidr        = "192.168.0.0/24"
}
resource "baiducloud_subnet" "default" {
 name        = var.name
 zone_name   = data.baiducloud_zones.default.zones.0.zone_name
 cidr        = "192.168.0.0/24"
 vpc_id      = baiducloud_vpc.default.id
 description = "created-by-terraform"
}
resource "baiducloud_appblb" "default" {
 depends_on  = [baiducloud_instance.default]
 name        = var.name
 description = "created-by-terraform"
 vpc_id      = baiducloud_vpc.default.id
 subnet_id   = baiducloud_subnet.default.id
}

resource "baiducloud_cert" "default" {
  cert_name         = var.name
  cert_server_data  = "-----BEGIN CERTIFICATE-----\n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\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIIErjCCA5agAwIBAgIQBYAmfwbylVM0jhwYWl7uLjANBgkqhkiG9w0BAQsFADBhMQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBDQTAeFw0xNzEyMDgxMjI4MjZaFw0yNzEyMDgxMjI4MjZaMHIxCzAJBgNVBAYTAkNOMSUwIwYDVQQKExxUcnVzdEFzaWEgVGVjaG5vbG9naWVzLCBJbmMuMR0wGwYDVQQLExREb21haW4gVmFsaWRhdGVkIFNTTDEdMBsGA1UEAxMUVHJ1c3RBc2lhIFRMUyBSU0EgQ0EwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCgWa9X+ph+wAm8Yh1Fk1MjKbQ5QwBOOKVaZR/OfCh+F6f93u7vZHGcUU/lvVGgUQnbzJhR1UV2epJae+m7cxnXIKdD0/VS9btAgwJszGFvwoqXeaCqFoP71wPmXjjUwLT70+qvX4hdyYfOJcjeTz5QKtg8zQwxaK9x4JT9CoOmoVdVhEBAiD3DwR5fFgOHDwwGxdJWVBvktnoAzjdTLXDdbSVC5jZ0u8oq9BiTDv7jAlsB5F8aZgvSZDOQeFrwaOTbKWSEInEhnchKZTD1dz6aBlk1xGEI5PZWAnVAba/ofH33ktymaTDsE6xRDnW97pDkimCRak6CEbfe3dXw6OV5AgMBAAGjggFPMIIBSzAdBgNVHQ4EFgQUf9OZ86BHDjEAVlYijrfMnt3KAYowHwYDVR0jBBgwFoAUA95QNVbRTLtm8KPiGxvDl7I90VUwDgYDVR0PAQH/BAQDAgGGMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjASBgNVHRMBAf8ECDAGAQH/AgEAMDQGCCsGAQUFBwEBBCgwJjAkBggrBgEFBQcwAYYYaHR0cDovL29jc3AuZGlnaWNlcnQuY29tMEIGA1UdHwQ7MDkwN6A1oDOGMWh0dHA6Ly9jcmwzLmRpZ2ljZXJ0LmNvbS9EaWdpQ2VydEdsb2JhbFJvb3RDQS5jcmwwTAYDVR0gBEUwQzA3BglghkgBhv1sAQIwKjAoBggrBgEFBQcCARYcaHR0cHM6Ly93d3cuZGlnaWNlcnQuY29tL0NQUzAIBgZngQwBAgEwDQYJKoZIhvcNAQELBQADggEBAK3dVOj5dlv4MzK2i233lDYvyJ3slFY2X2HKTYGte8nbK6i5/fsDImMYihAkp6VaNY/en8WZ5qcrQPVLuJrJDSXT04NnMeZOQDUoj/NHAmdfCBB/h1bZ5OGK6Sf1h5Yx/5wR4f3TUoPgGlnU7EuPISLNdMRiDrXntcImDAiRvkh5GJuH4YCVE6XEntqaNIgGkRwxKSgnU3Id3iuFbW9FUQ9Qqtb1GX91AJ7i4153TikGgYCdwYkBURD8gSVe8OAco6IfZOYt/TEwii1Ivi1CqnuUlWpsF1LdQNIdfbW3TSe0BhQa7ifbVIfvPWHYOu3rkg1ZeMo6XRU9B4n5VyJYRmE=\n-----END CERTIFICATE-----"
  cert_private_data = "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAkZCmcbCZyzdYQ379jndWgcKB9Q0/ITUYE2sUlh39HwyAJYP9DP2fStqF9pW+KVYJDPlhdFbYTrUqp+dO9001Hj3bL2JKMwx6zKmKtjUecpSkFtOYKAF+e6odhcS5HUTq0maydjvLSp1kPKd0Q6n4GW14+gkjoUvMVWK59FmjfO47yqETTWK/zA6JAp4ZMXuI6ujFlsnoxY4nFD0+HKiPNkYaFwDr0AUEwdAaDG/fp7/AP2xnQLdGXVn/gflulflUmA9tqwdF2ul9UmsY8iZcbPUbO0EtC78LOGPCgC5SsYVioUCzGO1HCbnDz3jvshAsZ/EMR07KBzqretZwEnuTSQIDAQABAoIBAAzBl4cfWfLljY4TVbFY7ZNJ0i1Wilbkz2XQPJ8aegFGYqp8TROI3EnpKX6I89UCgvYzRSI2rsEC/lMgIZrpa1i+70jRPRMJKm+/VyENjvatO6NRH/ni26HcWrb2HN90Qnx1XyPzrHvZnBxL876EPseCVkIvGoNliulb+/4Y/DXpNthA28UOB9RafPsEoDNinrTqlZf0gNLxm1LOgcj/NEqsDwuwzwfCky9GAhQgZpwic2IAEwKoCbfeRNNraVgG+IdCC8Nn3/uMcy9Zft3fV7xNE6HdfkW1SKnEvN+sFxKhH7ad0FNtaE+kSAcxTWXOg/xErvUBIcDrZv23BgN4JVMCgYEAwiNb00eRuBcPTHAaEb9JqrFRtUlqLnFJe1ang1QRfn+FrlTnijGACTjEFpzaXavaGNKi+To8OZjSTL2OW6ewEwSA9siPXUkq3ldPj5uPIhr80Jn1Ox/K5+X5ZBkQg8Iw9GIY6P6Kgf/prihVIbGZVNa0U/8H/1RvQIBxvA21dfMCgYEAv/L8iGiSwcgqMv0NTzfiW4fA9L7yLE04mfs9QI1V/uHPX5ufb/Y3LCS1RSuOdjrCdD2Ru7OKMi1v7mwPg1+NJBZjLIlCw/oVCJZabd8KGXZUNSH+PNuQAbIGdotEpO+LPgVgwi4ovrx6oJYEED/1FFjfU2bBFfuZtrDBWz2yNNMCgYEAvdoKQJHq5RZX9a5jMBvbFLwXZawH1Kcg7ycM5hdejFB1EMkjLTe/OEV1LY/y1EvtGv1SN1xF7SWP81AkWWmhfNeYrr3vxZB6Bbloqs27qeSue+kzssAik6mIu+TvC4rqiPMt3RyfowX7Jj93EV42zoqxCruKvJ17tp5lmzvkyxUCgYBRN60mwqimGd3RKUWCaXD7rZs1c73ghOQYMzgdoi/q4vztxVlW9GUv5nBUzjM/T2mL6alKNJOa26LqzQpbWgjMZjScWY/IgH553bRxnNgXIfxLZxC+C2EJdpxJeHAZIcpW+cuRHhrbacCxRgh+H7HBZEFKdsXoWUcXB/8obhiDRQKBgCwOE+1hfrV7/gFaMBWSML1n+LVV2ns80jCDtkhN9yF+9iJTjMwW4wuvFx8t8o2XICOwJPog4IvXFJLVZeed/zhgqe4qImHRW0aMYGyGEpgkLtHIFFFCxGd57Df/qEbUL55LU53rlCv2QKVBBs/6XDkiVRBk8izT7ihF2U8qb6t4\n-----END RSA PRIVATE KEY-----"
}

resource "baiducloud_appblb_listener" "default" {
 blb_id               = baiducloud_appblb.default.id
 listener_port        = 131
 protocol             = "SSL"
 scheduler            = "RoundRobin"
 cert_ids             = [baiducloud_cert.default.id]
 encryption_protocols = ["tlsv10", "tlsv11", "tlsv12"]
 encryption_type      = "userDefind"
}
`, name)
}
