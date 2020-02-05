package baiducloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
				Config: testAccAppBLBHTTPListenerConfigBasic(),
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
				Config: testAccAppBLBHTTPListenerConfigBasicUpdate(),
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
				Config: testAccAppBLBTCPListenerConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBListenerResourceName),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "listener_port", "124"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "protocol", "TCP"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "scheduler", "LeastConnection"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "tcp_session_timeout", "900"),
				),
			},
			{
				Config: testAccAppBLBTCPListenerConfigUpdate(),
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
				Config: testAccAppBLBUDPListenerConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBListenerResourceName),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "listener_port", "125"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "protocol", "UDP"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "scheduler", "LeastConnection"),
				),
			},
			{
				Config: testAccAppBLBUDPListenerConfigUpdate(),
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
				Config: testAccAppBLBHTTPListenerConfig(),
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
				Config: testAccAppBLBHTTPListenerConfigUpdate(),
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
				Config: testAccAppBLBHTTPSListenerConfig(),
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
				Config: testAccAppBLBHTTPSListenerConfigUpdate(),
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
				Config: testAccAppBLBSSLListenerConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBListenerResourceName),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "listener_port", "131"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "protocol", "SSL"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "scheduler", "LeastConnection"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "cert_ids.#", "1"),
					resource.TestCheckResourceAttr(testAccAppBLBListenerResourceName, "encryption_protocols.#", "3"),
				),
			},
			{
				Config: testAccAppBLBSSLListenerConfigUpdate(),
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

func testAccAppBLBHTTPListenerConfigBasic() string {
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
}

resource "baiducloud_appblb_server_group" "default" {
  name         = "%s"
  description  = "acceptance test"
  blb_id       = baiducloud_appblb.default.id

  port_list {
    port = 70
    type = "HTTP"
    health_check = "HTTP"
  }
}

resource "%s" "%s" {
  blb_id        = baiducloud_appblb.default.id
  listener_port = 129
  protocol      = "HTTP"
  scheduler     = "RoundRobin"
  keep_session  = true

  policies {
    description         = "acceptance test"
    app_server_group_id = baiducloud_appblb_server_group.default.id
    backend_port        = 70
    priority            = 50

    rule_list {
      key   = "host"
      value = "baidu.com"
    }
  }
}
`, BaiduCloudTestResourceAttrNamePrefix+"BCC",
		BaiduCloudTestResourceAttrNamePrefix+"VPC",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceAttrNamePrefix+"APPBLB",
		BaiduCloudTestResourceAttrNamePrefix+"APPServerGroup",
		testAccAppBLBListenerResourceType,
		BaiduCloudTestResourceName)
}

func testAccAppBLBHTTPListenerConfigBasicUpdate() string {
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
}

resource "%s" "%s" {
  blb_id        = baiducloud_appblb.default.id
  listener_port = 129
  protocol      = "HTTP"
  scheduler     = "RoundRobin"
  keep_session  = true
}
`, BaiduCloudTestResourceAttrNamePrefix+"BCC",
		BaiduCloudTestResourceAttrNamePrefix+"VPC",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceAttrNamePrefix+"APPBLB",
		testAccAppBLBListenerResourceType,
		BaiduCloudTestResourceName)
}

func testAccAppBLBTCPListenerConfig() string {
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
}

resource "%s" "%s" {
  blb_id        = baiducloud_appblb.default.id
  listener_port = 124
  protocol      = "TCP"
  scheduler     = "LeastConnection"
}
`, BaiduCloudTestResourceAttrNamePrefix+"BCC",
		BaiduCloudTestResourceAttrNamePrefix+"VPC",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceAttrNamePrefix+"APPBLB",
		testAccAppBLBListenerResourceType,
		BaiduCloudTestResourceName)
}

func testAccAppBLBTCPListenerConfigUpdate() string {
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
}

resource "baiducloud_appblb_server_group" "default" {
  name        = "%s"
  description = "acceptance test"
  blb_id      = baiducloud_appblb.default.id

  port_list {
    port = 68
    type = "TCP"
    health_check = "TCP"
  }
}

resource "%s" "%s" {
  blb_id               = baiducloud_appblb.default.id
  listener_port        = 124
  protocol             = "TCP"
  scheduler            = "RoundRobin"
  tcp_session_timeout  = 1000

  policies {
    description         = "acceptance test"
    app_server_group_id = baiducloud_appblb_server_group.default.id
    backend_port        = 68
    priority            = 50

    rule_list {
      key   = "*"
      value = "*"
    }
  }
}
`, BaiduCloudTestResourceAttrNamePrefix+"BCC",
		BaiduCloudTestResourceAttrNamePrefix+"VPC",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceAttrNamePrefix+"APPBLB",
		BaiduCloudTestResourceAttrNamePrefix+"APPServerGroup",
		testAccAppBLBListenerResourceType,
		BaiduCloudTestResourceName)
}

func testAccAppBLBUDPListenerConfig() string {
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
}

resource "%s" "%s" {
  blb_id        = baiducloud_appblb.default.id
  listener_port = 125
  protocol      = "UDP"
  scheduler     = "LeastConnection"
}
`, BaiduCloudTestResourceAttrNamePrefix+"BCC",
		BaiduCloudTestResourceAttrNamePrefix+"VPC",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceAttrNamePrefix+"APPBLB",
		testAccAppBLBListenerResourceType,
		BaiduCloudTestResourceName)
}

func testAccAppBLBUDPListenerConfigUpdate() string {
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
}

resource "baiducloud_appblb_server_group" "default" {
  name        = "%s"
  description = "acceptance test"
  blb_id      = baiducloud_appblb.default.id

  port_list {
    port = 66
    type = "UDP"
    health_check = "UDP"
    udp_health_check_string = "baidu.com"
  }
}

resource "%s" "%s" {
  blb_id        = baiducloud_appblb.default.id
  listener_port = 125
  protocol      = "UDP"
  scheduler     = "RoundRobin"

  policies {
    description         = "acceptance test"
    app_server_group_id = baiducloud_appblb_server_group.default.id
    backend_port        = 66
    priority            = 50

    rule_list {
      key   = "*"
      value = "*"
    }
  }
}
`, BaiduCloudTestResourceAttrNamePrefix+"BCC",
		BaiduCloudTestResourceAttrNamePrefix+"VPC",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceAttrNamePrefix+"APPBLB",
		BaiduCloudTestResourceAttrNamePrefix+"APPServerGroup",
		testAccAppBLBListenerResourceType,
		BaiduCloudTestResourceName)
}

func testAccAppBLBHTTPListenerConfig() string {
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
}

resource "%s" "%s" {
  blb_id        = baiducloud_appblb.default.id
  listener_port = 126
  protocol      = "HTTP"
  scheduler     = "LeastConnection"
  keep_session  = true
}
`, BaiduCloudTestResourceAttrNamePrefix+"BCC",
		BaiduCloudTestResourceAttrNamePrefix+"VPC",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceAttrNamePrefix+"APPBLB",
		testAccAppBLBListenerResourceType,
		BaiduCloudTestResourceName)
}

func testAccAppBLBHTTPListenerConfigUpdate() string {
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
}

resource "baiducloud_appblb_server_group" "default" {
  name         = "%s"
  description  = "acceptance test"
  blb_id       = baiducloud_appblb.default.id

  port_list {
    port = 67
    type = "HTTP"
    health_check = "HTTP"
  }
}

resource "%s" "%s" {
  blb_id        = baiducloud_appblb.default.id
  listener_port = 126
  protocol      = "HTTP"
  scheduler     = "RoundRobin"
  keep_session  = true

  policies {
    description         = "acceptance test"
    app_server_group_id = baiducloud_appblb_server_group.default.id
    backend_port        = 67
    priority            = 50

    rule_list {
      key   = "*"
      value = "*"
    }
  }
}
`, BaiduCloudTestResourceAttrNamePrefix+"BCC",
		BaiduCloudTestResourceAttrNamePrefix+"VPC",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceAttrNamePrefix+"APPBLB",
		BaiduCloudTestResourceAttrNamePrefix+"APPServerGroup",
		testAccAppBLBListenerResourceType,
		BaiduCloudTestResourceName)
}

func testAccAppBLBHTTPSListenerConfig() string {
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
}

resource "baiducloud_cert" "default" {
  cert_name         = "%s"
  cert_server_data  = "-----BEGIN CERTIFICATE-----\nMIIEGzCCA8CgAwIBAgIQBHVIJNCDJKsC1maaUVgqdjAKBggqhkjOPQQDAjByMQswCQYDVQQGEwJDTjElMCMGA1UEChMcVHJ1c3RBc2lhIFRlY2hub2xvZ2llcywgSW5jLjEdMBsGA1UECxMURG9tYWluIFZhbGlkYXRlZCBTU0wxHTAbBgNVBAMTFFRydXN0QXNpYSBUTFMgRUNDIENBMB4XDTE5MDkwNjAwMDAwMFoXDTIwMDkwNTEyMDAwMFowHzEdMBsGA1UEAxMUdGVzdC55aW5jaGVuZ2ZlbmcuY24wWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAR+aGvOdizh+oAWwT6829WdcZw7oBJVU1UvKQdm7dW/7SIdrMEWq6NIWaERMKkLD6gQ6Y5KFV9oDQdSocGBtBvLo4ICiTCCAoUwHwYDVR0jBBgwFoAUEoZEZiYIVCaPZTeyKU4mIeCTvtswHQYDVR0OBBYEFAichc0eFh+KdwMYjD7Pbvc8Q80IMB8GA1UdEQQYMBaCFHRlc3QueWluY2hlbmdmZW5nLmNuMA4GA1UdDwEB/wQEAwIHgDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwTAYDVR0gBEUwQzA3BglghkgBhv1sAQIwKjAoBggrBgEFBQcCARYcaHR0cHM6Ly93d3cuZGlnaWNlcnQuY29tL0NQUzAIBgZngQwBAgEwgZIGCCsGAQUFBwEBBIGFMIGCMDQGCCsGAQUFBzABhihodHRwOi8vc3RhdHVzZi5kaWdpdGFsY2VydHZhbGlkYXRpb24uY29tMEoGCCsGAQUFBzAChj5odHRwOi8vY2FjZXJ0cy5kaWdpdGFsY2VydHZhbGlkYXRpb24uY29tL1RydXN0QXNpYVRMU0VDQ0NBLmNydDAJBgNVHRMEAjAAMIIBAwYKKwYBBAHWeQIEAgSB9ASB8QDvAHUAu9nfvB+KcbWTlCOXqpJ7RzhXlQqrUugakJZkNo4e0YUAAAFtBK0O6QAABAMARjBEAiAdmHDa5NbRtLx3lc9nQ9G81RZycaqQPMj3+sazAo5vjQIgLNuFD7zperowYJAtetRR4QUi/8dORH087fWBp+Waj5MAdgCHdb/nWXz4jEOZX73zbv9WjUdWNv9KtWDBtOr/XqCDDwAAAW0ErQ9SAAAEAwBHMEUCIQDzdkB41ukE5XQGDTp8N4r+Aw/TZ/FlhPrrZryVGz9RIQIgWiuG2RHKCbh6FtJo62ml9RDYHeW/xA7c5sBBeKkSfG4wCgYIKoZIzj0EAwIDSQAwRgIhALnmf8VUwhxU0dRo2iOlfRb9uFy3hXMceU4IEvsLSwOVAiEAxsfjpOn0JyE943lhWRvjXX8FOm927cI5mbZ5F+p6dAA=\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIID4zCCAsugAwIBAgIQBz/JpHsGAhj24Khq6fw+OzANBgkqhkiG9w0BAQsFADBhMQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBDQTAeFw0xNzEyMDgxMjI4NTdaFw0yNzEyMDgxMjI4NTdaMHIxCzAJBgNVBAYTAkNOMSUwIwYDVQQKExxUcnVzdEFzaWEgVGVjaG5vbG9naWVzLCBJbmMuMR0wGwYDVQQLExREb21haW4gVmFsaWRhdGVkIFNTTDEdMBsGA1UEAxMUVHJ1c3RBc2lhIFRMUyBFQ0MgQ0EwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAASdQvDzv44jBee0APcvKOWszZsRjc4j+L6DLlYOf9tSgvfOJplfMeDNDZzOQEcJbVPD+yekJQUmObCPOrgMhqMIo4IBTzCCAUswHQYDVR0OBBYEFBKGRGYmCFQmj2U3silOJiHgk77bMB8GA1UdIwQYMBaAFAPeUDVW0Uy7ZvCj4hsbw5eyPdFVMA4GA1UdDwEB/wQEAwIBhjAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwEgYDVR0TAQH/BAgwBgEB/wIBADA0BggrBgEFBQcBAQQoMCYwJAYIKwYBBQUHMAGGGGh0dHA6Ly9vY3NwLmRpZ2ljZXJ0LmNvbTBCBgNVHR8EOzA5MDegNaAzhjFodHRwOi8vY3JsMy5kaWdpY2VydC5jb20vRGlnaUNlcnRHbG9iYWxSb290Q0EuY3JsMEwGA1UdIARFMEMwNwYJYIZIAYb9bAECMCowKAYIKwYBBQUHAgEWHGh0dHBzOi8vd3d3LmRpZ2ljZXJ0LmNvbS9DUFMwCAYGZ4EMAQIBMA0GCSqGSIb3DQEBCwUAA4IBAQBZcGGhLE09CbQD5xP93NAuNC85G1BMa1OG2Q01TWvvgp7Qt1wNfRLAnhQT5pb7kRs+E7nM4IS894ufmuL452q8gYaq5HmvOmfhXMmL6K+eICfvyqjb/tSi8iy20ULO/TZhLhPor9tle52Yx811FG4i5vqwPIUEOEJ7pXe6RPVoBiwi4rbLspQGD/vYqrj9OJV4JctoIhhGq+y/sozU6nBXHfhVSD3x+hkOOst6tyRq481IyUWQHcFtwda3gfMnaA3dsag2dtJz33RIJIUfxXmVK7w4YzHOHifn7TYk8iNrDDLtql6vS8FjiUx3kJnI6zge1C9lUHhZ/aD3RiTJrwWI\n-----END CERTIFICATE-----"
  cert_private_data = "-----BEGIN PRIVATE KEY-----\nMIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgp8yx31T7g0TyZcU4IdJS4px8p0b9FOHqx0uIMwtIjP6gCgYIKoZIzj0DAQehRANCAAR+aGvOdizh+oAWwT6829WdcZw7oBJVU1UvKQdm7dW/7SIdrMEWq6NIWaERMKkLD6gQ6Y5KFV9oDQdSocGBtBvL\n-----END PRIVATE KEY-----"
}

resource "%s" "%s" {
  blb_id               = baiducloud_appblb.default.id
  listener_port        = 130
  protocol             = "HTTPS"
  scheduler            = "LeastConnection"
  keep_session         = true
  cert_ids             = [baiducloud_cert.default.id]
  encryption_protocols = ["sslv3", "tlsv10", "tlsv11"]
  encryption_type      = "userDefind"
}
`, BaiduCloudTestResourceAttrNamePrefix+"BCC",
		BaiduCloudTestResourceAttrNamePrefix+"VPC",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceAttrNamePrefix+"APPBLB",
		BaiduCloudTestResourceAttrNamePrefix+"Cert",
		testAccAppBLBListenerResourceType,
		BaiduCloudTestResourceName)
}
func testAccAppBLBHTTPSListenerConfigUpdate() string {
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
}

resource "baiducloud_cert" "default" {
  cert_name         = "%s"
  cert_server_data  = "-----BEGIN CERTIFICATE-----\nMIIEGzCCA8CgAwIBAgIQBHVIJNCDJKsC1maaUVgqdjAKBggqhkjOPQQDAjByMQswCQYDVQQGEwJDTjElMCMGA1UEChMcVHJ1c3RBc2lhIFRlY2hub2xvZ2llcywgSW5jLjEdMBsGA1UECxMURG9tYWluIFZhbGlkYXRlZCBTU0wxHTAbBgNVBAMTFFRydXN0QXNpYSBUTFMgRUNDIENBMB4XDTE5MDkwNjAwMDAwMFoXDTIwMDkwNTEyMDAwMFowHzEdMBsGA1UEAxMUdGVzdC55aW5jaGVuZ2ZlbmcuY24wWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAR+aGvOdizh+oAWwT6829WdcZw7oBJVU1UvKQdm7dW/7SIdrMEWq6NIWaERMKkLD6gQ6Y5KFV9oDQdSocGBtBvLo4ICiTCCAoUwHwYDVR0jBBgwFoAUEoZEZiYIVCaPZTeyKU4mIeCTvtswHQYDVR0OBBYEFAichc0eFh+KdwMYjD7Pbvc8Q80IMB8GA1UdEQQYMBaCFHRlc3QueWluY2hlbmdmZW5nLmNuMA4GA1UdDwEB/wQEAwIHgDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwTAYDVR0gBEUwQzA3BglghkgBhv1sAQIwKjAoBggrBgEFBQcCARYcaHR0cHM6Ly93d3cuZGlnaWNlcnQuY29tL0NQUzAIBgZngQwBAgEwgZIGCCsGAQUFBwEBBIGFMIGCMDQGCCsGAQUFBzABhihodHRwOi8vc3RhdHVzZi5kaWdpdGFsY2VydHZhbGlkYXRpb24uY29tMEoGCCsGAQUFBzAChj5odHRwOi8vY2FjZXJ0cy5kaWdpdGFsY2VydHZhbGlkYXRpb24uY29tL1RydXN0QXNpYVRMU0VDQ0NBLmNydDAJBgNVHRMEAjAAMIIBAwYKKwYBBAHWeQIEAgSB9ASB8QDvAHUAu9nfvB+KcbWTlCOXqpJ7RzhXlQqrUugakJZkNo4e0YUAAAFtBK0O6QAABAMARjBEAiAdmHDa5NbRtLx3lc9nQ9G81RZycaqQPMj3+sazAo5vjQIgLNuFD7zperowYJAtetRR4QUi/8dORH087fWBp+Waj5MAdgCHdb/nWXz4jEOZX73zbv9WjUdWNv9KtWDBtOr/XqCDDwAAAW0ErQ9SAAAEAwBHMEUCIQDzdkB41ukE5XQGDTp8N4r+Aw/TZ/FlhPrrZryVGz9RIQIgWiuG2RHKCbh6FtJo62ml9RDYHeW/xA7c5sBBeKkSfG4wCgYIKoZIzj0EAwIDSQAwRgIhALnmf8VUwhxU0dRo2iOlfRb9uFy3hXMceU4IEvsLSwOVAiEAxsfjpOn0JyE943lhWRvjXX8FOm927cI5mbZ5F+p6dAA=\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIID4zCCAsugAwIBAgIQBz/JpHsGAhj24Khq6fw+OzANBgkqhkiG9w0BAQsFADBhMQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBDQTAeFw0xNzEyMDgxMjI4NTdaFw0yNzEyMDgxMjI4NTdaMHIxCzAJBgNVBAYTAkNOMSUwIwYDVQQKExxUcnVzdEFzaWEgVGVjaG5vbG9naWVzLCBJbmMuMR0wGwYDVQQLExREb21haW4gVmFsaWRhdGVkIFNTTDEdMBsGA1UEAxMUVHJ1c3RBc2lhIFRMUyBFQ0MgQ0EwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAASdQvDzv44jBee0APcvKOWszZsRjc4j+L6DLlYOf9tSgvfOJplfMeDNDZzOQEcJbVPD+yekJQUmObCPOrgMhqMIo4IBTzCCAUswHQYDVR0OBBYEFBKGRGYmCFQmj2U3silOJiHgk77bMB8GA1UdIwQYMBaAFAPeUDVW0Uy7ZvCj4hsbw5eyPdFVMA4GA1UdDwEB/wQEAwIBhjAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwEgYDVR0TAQH/BAgwBgEB/wIBADA0BggrBgEFBQcBAQQoMCYwJAYIKwYBBQUHMAGGGGh0dHA6Ly9vY3NwLmRpZ2ljZXJ0LmNvbTBCBgNVHR8EOzA5MDegNaAzhjFodHRwOi8vY3JsMy5kaWdpY2VydC5jb20vRGlnaUNlcnRHbG9iYWxSb290Q0EuY3JsMEwGA1UdIARFMEMwNwYJYIZIAYb9bAECMCowKAYIKwYBBQUHAgEWHGh0dHBzOi8vd3d3LmRpZ2ljZXJ0LmNvbS9DUFMwCAYGZ4EMAQIBMA0GCSqGSIb3DQEBCwUAA4IBAQBZcGGhLE09CbQD5xP93NAuNC85G1BMa1OG2Q01TWvvgp7Qt1wNfRLAnhQT5pb7kRs+E7nM4IS894ufmuL452q8gYaq5HmvOmfhXMmL6K+eICfvyqjb/tSi8iy20ULO/TZhLhPor9tle52Yx811FG4i5vqwPIUEOEJ7pXe6RPVoBiwi4rbLspQGD/vYqrj9OJV4JctoIhhGq+y/sozU6nBXHfhVSD3x+hkOOst6tyRq481IyUWQHcFtwda3gfMnaA3dsag2dtJz33RIJIUfxXmVK7w4YzHOHifn7TYk8iNrDDLtql6vS8FjiUx3kJnI6zge1C9lUHhZ/aD3RiTJrwWI\n-----END CERTIFICATE-----"
  cert_private_data = "-----BEGIN PRIVATE KEY-----\nMIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgp8yx31T7g0TyZcU4IdJS4px8p0b9FOHqx0uIMwtIjP6gCgYIKoZIzj0DAQehRANCAAR+aGvOdizh+oAWwT6829WdcZw7oBJVU1UvKQdm7dW/7SIdrMEWq6NIWaERMKkLD6gQ6Y5KFV9oDQdSocGBtBvL\n-----END PRIVATE KEY-----"
}

resource "baiducloud_cert" "default2" {
  cert_name         = "%s"
  cert_server_data  = "-----BEGIN CERTIFICATE-----\nMIIEHjCCA8SgAwIBAgIQD7e2kCM5IFr1AhZhtHco3DAKBggqhkjOPQQDAjByMQswCQYDVQQGEwJDTjElMCMGA1UEChMcVHJ1c3RBc2lhIFRlY2hub2xvZ2llcywgSW5jLjEdMBsGA1UECxMURG9tYWluIFZhbGlkYXRlZCBTU0wxHTAbBgNVBAMTFFRydXN0QXNpYSBUTFMgRUNDIENBMB4XDTE5MDkwNjAwMDAwMFoXDTIwMDkwNTEyMDAwMFowIDEeMBwGA1UEAxMVdGVzdDIueWluY2hlbmdmZW5nLmNuMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE5aIFysizmk3WriZXuYXzgcqcF7ORRPFIxQXvYTDGuuR9ybqBkT3zCt7n7YUW3z9AN4ux1Yxj2VnGM79YpPszGqOCAowwggKIMB8GA1UdIwQYMBaAFBKGRGYmCFQmj2U3silOJiHgk77bMB0GA1UdDgQWBBSoycYcJp+vvxdIWaM9QS4IchsYKDAgBgNVHREEGTAXghV0ZXN0Mi55aW5jaGVuZ2ZlbmcuY24wDgYDVR0PAQH/BAQDAgeAMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjBMBgNVHSAERTBDMDcGCWCGSAGG/WwBAjAqMCgGCCsGAQUFBwIBFhxodHRwczovL3d3dy5kaWdpY2VydC5jb20vQ1BTMAgGBmeBDAECATCBkgYIKwYBBQUHAQEEgYUwgYIwNAYIKwYBBQUHMAGGKGh0dHA6Ly9zdGF0dXNmLmRpZ2l0YWxjZXJ0dmFsaWRhdGlvbi5jb20wSgYIKwYBBQUHMAKGPmh0dHA6Ly9jYWNlcnRzLmRpZ2l0YWxjZXJ0dmFsaWRhdGlvbi5jb20vVHJ1c3RBc2lhVExTRUNDQ0EuY3J0MAkGA1UdEwQCMAAwggEFBgorBgEEAdZ5AgQCBIH2BIHzAPEAdgDuS723dc5guuFCaR+r4Z5mow9+X7By2IMAxHuJeqj9ywAAAW0FFJZQAAAEAwBHMEUCIDq3C14Mq4CaueNUWVIBKI3HGphyj4JqRKVvfGP4qBR4AiEAsgc3/WUucxBeK/+2vQJmFgE+kUwAa3ZGgoq4fmKsxlcAdwCHdb/nWXz4jEOZX73zbv9WjUdWNv9KtWDBtOr/XqCDDwAAAW0FFJa9AAAEAwBIMEYCIQDoRpKHe+ljJ6JmJoMzK3IE+f3AfLrN5f07D9eRIwqBNQIhAMYw+Sn8HZ53sxE5ttkJGetSu4mUf1bqrXG7CoSo5rjFMAoGCCqGSM49BAMCA0gAMEUCIQDjzWnH6V/OHVQvPZuaNXD6P/U4rdoUvhLnqoFkrRZxYAIgU7qPXUAOdwAWy0LuINOz0OmoXc5angeJAqK67hULNI4=\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIID4zCCAsugAwIBAgIQBz/JpHsGAhj24Khq6fw+OzANBgkqhkiG9w0BAQsFADBhMQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBDQTAeFw0xNzEyMDgxMjI4NTdaFw0yNzEyMDgxMjI4NTdaMHIxCzAJBgNVBAYTAkNOMSUwIwYDVQQKExxUcnVzdEFzaWEgVGVjaG5vbG9naWVzLCBJbmMuMR0wGwYDVQQLExREb21haW4gVmFsaWRhdGVkIFNTTDEdMBsGA1UEAxMUVHJ1c3RBc2lhIFRMUyBFQ0MgQ0EwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAASdQvDzv44jBee0APcvKOWszZsRjc4j+L6DLlYOf9tSgvfOJplfMeDNDZzOQEcJbVPD+yekJQUmObCPOrgMhqMIo4IBTzCCAUswHQYDVR0OBBYEFBKGRGYmCFQmj2U3silOJiHgk77bMB8GA1UdIwQYMBaAFAPeUDVW0Uy7ZvCj4hsbw5eyPdFVMA4GA1UdDwEB/wQEAwIBhjAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwEgYDVR0TAQH/BAgwBgEB/wIBADA0BggrBgEFBQcBAQQoMCYwJAYIKwYBBQUHMAGGGGh0dHA6Ly9vY3NwLmRpZ2ljZXJ0LmNvbTBCBgNVHR8EOzA5MDegNaAzhjFodHRwOi8vY3JsMy5kaWdpY2VydC5jb20vRGlnaUNlcnRHbG9iYWxSb290Q0EuY3JsMEwGA1UdIARFMEMwNwYJYIZIAYb9bAECMCowKAYIKwYBBQUHAgEWHGh0dHBzOi8vd3d3LmRpZ2ljZXJ0LmNvbS9DUFMwCAYGZ4EMAQIBMA0GCSqGSIb3DQEBCwUAA4IBAQBZcGGhLE09CbQD5xP93NAuNC85G1BMa1OG2Q01TWvvgp7Qt1wNfRLAnhQT5pb7kRs+E7nM4IS894ufmuL452q8gYaq5HmvOmfhXMmL6K+eICfvyqjb/tSi8iy20ULO/TZhLhPor9tle52Yx811FG4i5vqwPIUEOEJ7pXe6RPVoBiwi4rbLspQGD/vYqrj9OJV4JctoIhhGq+y/sozU6nBXHfhVSD3x+hkOOst6tyRq481IyUWQHcFtwda3gfMnaA3dsag2dtJz33RIJIUfxXmVK7w4YzHOHifn7TYk8iNrDDLtql6vS8FjiUx3kJnI6zge1C9lUHhZ/aD3RiTJrwWI\n-----END CERTIFICATE-----"
  cert_private_data = "-----BEGIN PRIVATE KEY-----\nMIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQg4vsAo5xhUZD92opgs+dSIDFHgFjikrZylNHvSSIyJjegCgYIKoZIzj0DAQehRANCAATlogXKyLOaTdauJle5hfOBypwXs5FE8UjFBe9hMMa65H3JuoGRPfMK3ufthRbfP0A3i7HVjGPZWcYzv1ik+zMa\n-----END PRIVATE KEY-----"
}

resource "%s" "%s" {
  blb_id               = baiducloud_appblb.default.id
  listener_port        = 130
  protocol             = "HTTPS"
  scheduler            = "RoundRobin"
  keep_session         = true
  cert_ids             = [baiducloud_cert.default.id, baiducloud_cert.default2.id]
  encryption_protocols = ["tlsv10", "tlsv11", "tlsv12"]
  encryption_type      = "userDefind"
}
`, BaiduCloudTestResourceAttrNamePrefix+"BCC",
		BaiduCloudTestResourceAttrNamePrefix+"VPC",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceAttrNamePrefix+"APPBLB",
		BaiduCloudTestResourceAttrNamePrefix+"Cert",
		BaiduCloudTestResourceAttrNamePrefix+"Cert2",
		testAccAppBLBListenerResourceType,
		BaiduCloudTestResourceName)
}

func testAccAppBLBSSLListenerConfig() string {
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
}

resource "baiducloud_cert" "default" {
  cert_name         = "%s"
  cert_server_data  = "-----BEGIN CERTIFICATE-----\nMIIEGzCCA8CgAwIBAgIQBHVIJNCDJKsC1maaUVgqdjAKBggqhkjOPQQDAjByMQswCQYDVQQGEwJDTjElMCMGA1UEChMcVHJ1c3RBc2lhIFRlY2hub2xvZ2llcywgSW5jLjEdMBsGA1UECxMURG9tYWluIFZhbGlkYXRlZCBTU0wxHTAbBgNVBAMTFFRydXN0QXNpYSBUTFMgRUNDIENBMB4XDTE5MDkwNjAwMDAwMFoXDTIwMDkwNTEyMDAwMFowHzEdMBsGA1UEAxMUdGVzdC55aW5jaGVuZ2ZlbmcuY24wWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAR+aGvOdizh+oAWwT6829WdcZw7oBJVU1UvKQdm7dW/7SIdrMEWq6NIWaERMKkLD6gQ6Y5KFV9oDQdSocGBtBvLo4ICiTCCAoUwHwYDVR0jBBgwFoAUEoZEZiYIVCaPZTeyKU4mIeCTvtswHQYDVR0OBBYEFAichc0eFh+KdwMYjD7Pbvc8Q80IMB8GA1UdEQQYMBaCFHRlc3QueWluY2hlbmdmZW5nLmNuMA4GA1UdDwEB/wQEAwIHgDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwTAYDVR0gBEUwQzA3BglghkgBhv1sAQIwKjAoBggrBgEFBQcCARYcaHR0cHM6Ly93d3cuZGlnaWNlcnQuY29tL0NQUzAIBgZngQwBAgEwgZIGCCsGAQUFBwEBBIGFMIGCMDQGCCsGAQUFBzABhihodHRwOi8vc3RhdHVzZi5kaWdpdGFsY2VydHZhbGlkYXRpb24uY29tMEoGCCsGAQUFBzAChj5odHRwOi8vY2FjZXJ0cy5kaWdpdGFsY2VydHZhbGlkYXRpb24uY29tL1RydXN0QXNpYVRMU0VDQ0NBLmNydDAJBgNVHRMEAjAAMIIBAwYKKwYBBAHWeQIEAgSB9ASB8QDvAHUAu9nfvB+KcbWTlCOXqpJ7RzhXlQqrUugakJZkNo4e0YUAAAFtBK0O6QAABAMARjBEAiAdmHDa5NbRtLx3lc9nQ9G81RZycaqQPMj3+sazAo5vjQIgLNuFD7zperowYJAtetRR4QUi/8dORH087fWBp+Waj5MAdgCHdb/nWXz4jEOZX73zbv9WjUdWNv9KtWDBtOr/XqCDDwAAAW0ErQ9SAAAEAwBHMEUCIQDzdkB41ukE5XQGDTp8N4r+Aw/TZ/FlhPrrZryVGz9RIQIgWiuG2RHKCbh6FtJo62ml9RDYHeW/xA7c5sBBeKkSfG4wCgYIKoZIzj0EAwIDSQAwRgIhALnmf8VUwhxU0dRo2iOlfRb9uFy3hXMceU4IEvsLSwOVAiEAxsfjpOn0JyE943lhWRvjXX8FOm927cI5mbZ5F+p6dAA=\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIID4zCCAsugAwIBAgIQBz/JpHsGAhj24Khq6fw+OzANBgkqhkiG9w0BAQsFADBhMQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBDQTAeFw0xNzEyMDgxMjI4NTdaFw0yNzEyMDgxMjI4NTdaMHIxCzAJBgNVBAYTAkNOMSUwIwYDVQQKExxUcnVzdEFzaWEgVGVjaG5vbG9naWVzLCBJbmMuMR0wGwYDVQQLExREb21haW4gVmFsaWRhdGVkIFNTTDEdMBsGA1UEAxMUVHJ1c3RBc2lhIFRMUyBFQ0MgQ0EwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAASdQvDzv44jBee0APcvKOWszZsRjc4j+L6DLlYOf9tSgvfOJplfMeDNDZzOQEcJbVPD+yekJQUmObCPOrgMhqMIo4IBTzCCAUswHQYDVR0OBBYEFBKGRGYmCFQmj2U3silOJiHgk77bMB8GA1UdIwQYMBaAFAPeUDVW0Uy7ZvCj4hsbw5eyPdFVMA4GA1UdDwEB/wQEAwIBhjAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwEgYDVR0TAQH/BAgwBgEB/wIBADA0BggrBgEFBQcBAQQoMCYwJAYIKwYBBQUHMAGGGGh0dHA6Ly9vY3NwLmRpZ2ljZXJ0LmNvbTBCBgNVHR8EOzA5MDegNaAzhjFodHRwOi8vY3JsMy5kaWdpY2VydC5jb20vRGlnaUNlcnRHbG9iYWxSb290Q0EuY3JsMEwGA1UdIARFMEMwNwYJYIZIAYb9bAECMCowKAYIKwYBBQUHAgEWHGh0dHBzOi8vd3d3LmRpZ2ljZXJ0LmNvbS9DUFMwCAYGZ4EMAQIBMA0GCSqGSIb3DQEBCwUAA4IBAQBZcGGhLE09CbQD5xP93NAuNC85G1BMa1OG2Q01TWvvgp7Qt1wNfRLAnhQT5pb7kRs+E7nM4IS894ufmuL452q8gYaq5HmvOmfhXMmL6K+eICfvyqjb/tSi8iy20ULO/TZhLhPor9tle52Yx811FG4i5vqwPIUEOEJ7pXe6RPVoBiwi4rbLspQGD/vYqrj9OJV4JctoIhhGq+y/sozU6nBXHfhVSD3x+hkOOst6tyRq481IyUWQHcFtwda3gfMnaA3dsag2dtJz33RIJIUfxXmVK7w4YzHOHifn7TYk8iNrDDLtql6vS8FjiUx3kJnI6zge1C9lUHhZ/aD3RiTJrwWI\n-----END CERTIFICATE-----"
  cert_private_data = "-----BEGIN PRIVATE KEY-----\nMIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgp8yx31T7g0TyZcU4IdJS4px8p0b9FOHqx0uIMwtIjP6gCgYIKoZIzj0DAQehRANCAAR+aGvOdizh+oAWwT6829WdcZw7oBJVU1UvKQdm7dW/7SIdrMEWq6NIWaERMKkLD6gQ6Y5KFV9oDQdSocGBtBvL\n-----END PRIVATE KEY-----"
}

resource "%s" "%s" {
  blb_id               = baiducloud_appblb.default.id
  listener_port        = 131
  protocol             = "SSL"
  scheduler            = "LeastConnection"
  cert_ids             = [baiducloud_cert.default.id]
  encryption_protocols = ["sslv3", "tlsv10", "tlsv11"]
  encryption_type      = "userDefind"
}
`, BaiduCloudTestResourceAttrNamePrefix+"BCC",
		BaiduCloudTestResourceAttrNamePrefix+"VPC",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceAttrNamePrefix+"APPBLB",
		BaiduCloudTestResourceAttrNamePrefix+"Cert",
		testAccAppBLBListenerResourceType,
		BaiduCloudTestResourceName)
}
func testAccAppBLBSSLListenerConfigUpdate() string {
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
}

resource "baiducloud_cert" "default" {
  cert_name         = "%s"
  cert_server_data  = "-----BEGIN CERTIFICATE-----\nMIIEGzCCA8CgAwIBAgIQBHVIJNCDJKsC1maaUVgqdjAKBggqhkjOPQQDAjByMQswCQYDVQQGEwJDTjElMCMGA1UEChMcVHJ1c3RBc2lhIFRlY2hub2xvZ2llcywgSW5jLjEdMBsGA1UECxMURG9tYWluIFZhbGlkYXRlZCBTU0wxHTAbBgNVBAMTFFRydXN0QXNpYSBUTFMgRUNDIENBMB4XDTE5MDkwNjAwMDAwMFoXDTIwMDkwNTEyMDAwMFowHzEdMBsGA1UEAxMUdGVzdC55aW5jaGVuZ2ZlbmcuY24wWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAR+aGvOdizh+oAWwT6829WdcZw7oBJVU1UvKQdm7dW/7SIdrMEWq6NIWaERMKkLD6gQ6Y5KFV9oDQdSocGBtBvLo4ICiTCCAoUwHwYDVR0jBBgwFoAUEoZEZiYIVCaPZTeyKU4mIeCTvtswHQYDVR0OBBYEFAichc0eFh+KdwMYjD7Pbvc8Q80IMB8GA1UdEQQYMBaCFHRlc3QueWluY2hlbmdmZW5nLmNuMA4GA1UdDwEB/wQEAwIHgDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwTAYDVR0gBEUwQzA3BglghkgBhv1sAQIwKjAoBggrBgEFBQcCARYcaHR0cHM6Ly93d3cuZGlnaWNlcnQuY29tL0NQUzAIBgZngQwBAgEwgZIGCCsGAQUFBwEBBIGFMIGCMDQGCCsGAQUFBzABhihodHRwOi8vc3RhdHVzZi5kaWdpdGFsY2VydHZhbGlkYXRpb24uY29tMEoGCCsGAQUFBzAChj5odHRwOi8vY2FjZXJ0cy5kaWdpdGFsY2VydHZhbGlkYXRpb24uY29tL1RydXN0QXNpYVRMU0VDQ0NBLmNydDAJBgNVHRMEAjAAMIIBAwYKKwYBBAHWeQIEAgSB9ASB8QDvAHUAu9nfvB+KcbWTlCOXqpJ7RzhXlQqrUugakJZkNo4e0YUAAAFtBK0O6QAABAMARjBEAiAdmHDa5NbRtLx3lc9nQ9G81RZycaqQPMj3+sazAo5vjQIgLNuFD7zperowYJAtetRR4QUi/8dORH087fWBp+Waj5MAdgCHdb/nWXz4jEOZX73zbv9WjUdWNv9KtWDBtOr/XqCDDwAAAW0ErQ9SAAAEAwBHMEUCIQDzdkB41ukE5XQGDTp8N4r+Aw/TZ/FlhPrrZryVGz9RIQIgWiuG2RHKCbh6FtJo62ml9RDYHeW/xA7c5sBBeKkSfG4wCgYIKoZIzj0EAwIDSQAwRgIhALnmf8VUwhxU0dRo2iOlfRb9uFy3hXMceU4IEvsLSwOVAiEAxsfjpOn0JyE943lhWRvjXX8FOm927cI5mbZ5F+p6dAA=\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIID4zCCAsugAwIBAgIQBz/JpHsGAhj24Khq6fw+OzANBgkqhkiG9w0BAQsFADBhMQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBDQTAeFw0xNzEyMDgxMjI4NTdaFw0yNzEyMDgxMjI4NTdaMHIxCzAJBgNVBAYTAkNOMSUwIwYDVQQKExxUcnVzdEFzaWEgVGVjaG5vbG9naWVzLCBJbmMuMR0wGwYDVQQLExREb21haW4gVmFsaWRhdGVkIFNTTDEdMBsGA1UEAxMUVHJ1c3RBc2lhIFRMUyBFQ0MgQ0EwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAASdQvDzv44jBee0APcvKOWszZsRjc4j+L6DLlYOf9tSgvfOJplfMeDNDZzOQEcJbVPD+yekJQUmObCPOrgMhqMIo4IBTzCCAUswHQYDVR0OBBYEFBKGRGYmCFQmj2U3silOJiHgk77bMB8GA1UdIwQYMBaAFAPeUDVW0Uy7ZvCj4hsbw5eyPdFVMA4GA1UdDwEB/wQEAwIBhjAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwEgYDVR0TAQH/BAgwBgEB/wIBADA0BggrBgEFBQcBAQQoMCYwJAYIKwYBBQUHMAGGGGh0dHA6Ly9vY3NwLmRpZ2ljZXJ0LmNvbTBCBgNVHR8EOzA5MDegNaAzhjFodHRwOi8vY3JsMy5kaWdpY2VydC5jb20vRGlnaUNlcnRHbG9iYWxSb290Q0EuY3JsMEwGA1UdIARFMEMwNwYJYIZIAYb9bAECMCowKAYIKwYBBQUHAgEWHGh0dHBzOi8vd3d3LmRpZ2ljZXJ0LmNvbS9DUFMwCAYGZ4EMAQIBMA0GCSqGSIb3DQEBCwUAA4IBAQBZcGGhLE09CbQD5xP93NAuNC85G1BMa1OG2Q01TWvvgp7Qt1wNfRLAnhQT5pb7kRs+E7nM4IS894ufmuL452q8gYaq5HmvOmfhXMmL6K+eICfvyqjb/tSi8iy20ULO/TZhLhPor9tle52Yx811FG4i5vqwPIUEOEJ7pXe6RPVoBiwi4rbLspQGD/vYqrj9OJV4JctoIhhGq+y/sozU6nBXHfhVSD3x+hkOOst6tyRq481IyUWQHcFtwda3gfMnaA3dsag2dtJz33RIJIUfxXmVK7w4YzHOHifn7TYk8iNrDDLtql6vS8FjiUx3kJnI6zge1C9lUHhZ/aD3RiTJrwWI\n-----END CERTIFICATE-----"
  cert_private_data = "-----BEGIN PRIVATE KEY-----\nMIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgp8yx31T7g0TyZcU4IdJS4px8p0b9FOHqx0uIMwtIjP6gCgYIKoZIzj0DAQehRANCAAR+aGvOdizh+oAWwT6829WdcZw7oBJVU1UvKQdm7dW/7SIdrMEWq6NIWaERMKkLD6gQ6Y5KFV9oDQdSocGBtBvL\n-----END PRIVATE KEY-----"
}

resource "baiducloud_cert" "default2" {
  cert_name         = "%s"
  cert_server_data  = "-----BEGIN CERTIFICATE-----\nMIIEHjCCA8SgAwIBAgIQD7e2kCM5IFr1AhZhtHco3DAKBggqhkjOPQQDAjByMQswCQYDVQQGEwJDTjElMCMGA1UEChMcVHJ1c3RBc2lhIFRlY2hub2xvZ2llcywgSW5jLjEdMBsGA1UECxMURG9tYWluIFZhbGlkYXRlZCBTU0wxHTAbBgNVBAMTFFRydXN0QXNpYSBUTFMgRUNDIENBMB4XDTE5MDkwNjAwMDAwMFoXDTIwMDkwNTEyMDAwMFowIDEeMBwGA1UEAxMVdGVzdDIueWluY2hlbmdmZW5nLmNuMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE5aIFysizmk3WriZXuYXzgcqcF7ORRPFIxQXvYTDGuuR9ybqBkT3zCt7n7YUW3z9AN4ux1Yxj2VnGM79YpPszGqOCAowwggKIMB8GA1UdIwQYMBaAFBKGRGYmCFQmj2U3silOJiHgk77bMB0GA1UdDgQWBBSoycYcJp+vvxdIWaM9QS4IchsYKDAgBgNVHREEGTAXghV0ZXN0Mi55aW5jaGVuZ2ZlbmcuY24wDgYDVR0PAQH/BAQDAgeAMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjBMBgNVHSAERTBDMDcGCWCGSAGG/WwBAjAqMCgGCCsGAQUFBwIBFhxodHRwczovL3d3dy5kaWdpY2VydC5jb20vQ1BTMAgGBmeBDAECATCBkgYIKwYBBQUHAQEEgYUwgYIwNAYIKwYBBQUHMAGGKGh0dHA6Ly9zdGF0dXNmLmRpZ2l0YWxjZXJ0dmFsaWRhdGlvbi5jb20wSgYIKwYBBQUHMAKGPmh0dHA6Ly9jYWNlcnRzLmRpZ2l0YWxjZXJ0dmFsaWRhdGlvbi5jb20vVHJ1c3RBc2lhVExTRUNDQ0EuY3J0MAkGA1UdEwQCMAAwggEFBgorBgEEAdZ5AgQCBIH2BIHzAPEAdgDuS723dc5guuFCaR+r4Z5mow9+X7By2IMAxHuJeqj9ywAAAW0FFJZQAAAEAwBHMEUCIDq3C14Mq4CaueNUWVIBKI3HGphyj4JqRKVvfGP4qBR4AiEAsgc3/WUucxBeK/+2vQJmFgE+kUwAa3ZGgoq4fmKsxlcAdwCHdb/nWXz4jEOZX73zbv9WjUdWNv9KtWDBtOr/XqCDDwAAAW0FFJa9AAAEAwBIMEYCIQDoRpKHe+ljJ6JmJoMzK3IE+f3AfLrN5f07D9eRIwqBNQIhAMYw+Sn8HZ53sxE5ttkJGetSu4mUf1bqrXG7CoSo5rjFMAoGCCqGSM49BAMCA0gAMEUCIQDjzWnH6V/OHVQvPZuaNXD6P/U4rdoUvhLnqoFkrRZxYAIgU7qPXUAOdwAWy0LuINOz0OmoXc5angeJAqK67hULNI4=\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIID4zCCAsugAwIBAgIQBz/JpHsGAhj24Khq6fw+OzANBgkqhkiG9w0BAQsFADBhMQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBDQTAeFw0xNzEyMDgxMjI4NTdaFw0yNzEyMDgxMjI4NTdaMHIxCzAJBgNVBAYTAkNOMSUwIwYDVQQKExxUcnVzdEFzaWEgVGVjaG5vbG9naWVzLCBJbmMuMR0wGwYDVQQLExREb21haW4gVmFsaWRhdGVkIFNTTDEdMBsGA1UEAxMUVHJ1c3RBc2lhIFRMUyBFQ0MgQ0EwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAASdQvDzv44jBee0APcvKOWszZsRjc4j+L6DLlYOf9tSgvfOJplfMeDNDZzOQEcJbVPD+yekJQUmObCPOrgMhqMIo4IBTzCCAUswHQYDVR0OBBYEFBKGRGYmCFQmj2U3silOJiHgk77bMB8GA1UdIwQYMBaAFAPeUDVW0Uy7ZvCj4hsbw5eyPdFVMA4GA1UdDwEB/wQEAwIBhjAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwEgYDVR0TAQH/BAgwBgEB/wIBADA0BggrBgEFBQcBAQQoMCYwJAYIKwYBBQUHMAGGGGh0dHA6Ly9vY3NwLmRpZ2ljZXJ0LmNvbTBCBgNVHR8EOzA5MDegNaAzhjFodHRwOi8vY3JsMy5kaWdpY2VydC5jb20vRGlnaUNlcnRHbG9iYWxSb290Q0EuY3JsMEwGA1UdIARFMEMwNwYJYIZIAYb9bAECMCowKAYIKwYBBQUHAgEWHGh0dHBzOi8vd3d3LmRpZ2ljZXJ0LmNvbS9DUFMwCAYGZ4EMAQIBMA0GCSqGSIb3DQEBCwUAA4IBAQBZcGGhLE09CbQD5xP93NAuNC85G1BMa1OG2Q01TWvvgp7Qt1wNfRLAnhQT5pb7kRs+E7nM4IS894ufmuL452q8gYaq5HmvOmfhXMmL6K+eICfvyqjb/tSi8iy20ULO/TZhLhPor9tle52Yx811FG4i5vqwPIUEOEJ7pXe6RPVoBiwi4rbLspQGD/vYqrj9OJV4JctoIhhGq+y/sozU6nBXHfhVSD3x+hkOOst6tyRq481IyUWQHcFtwda3gfMnaA3dsag2dtJz33RIJIUfxXmVK7w4YzHOHifn7TYk8iNrDDLtql6vS8FjiUx3kJnI6zge1C9lUHhZ/aD3RiTJrwWI\n-----END CERTIFICATE-----"
  cert_private_data = "-----BEGIN PRIVATE KEY-----\nMIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQg4vsAo5xhUZD92opgs+dSIDFHgFjikrZylNHvSSIyJjegCgYIKoZIzj0DAQehRANCAATlogXKyLOaTdauJle5hfOBypwXs5FE8UjFBe9hMMa65H3JuoGRPfMK3ufthRbfP0A3i7HVjGPZWcYzv1ik+zMa\n-----END PRIVATE KEY-----"
}

resource "%s" "%s" {
  blb_id               = baiducloud_appblb.default.id
  listener_port        = 131
  protocol             = "SSL"
  scheduler            = "RoundRobin"
  cert_ids             = [baiducloud_cert.default.id, baiducloud_cert.default2.id]
  encryption_protocols = ["tlsv10", "tlsv11", "tlsv12"]
  encryption_type      = "userDefind"
}
`, BaiduCloudTestResourceAttrNamePrefix+"BCC",
		BaiduCloudTestResourceAttrNamePrefix+"VPC",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceAttrNamePrefix+"APPBLB",
		BaiduCloudTestResourceAttrNamePrefix+"Cert",
		BaiduCloudTestResourceAttrNamePrefix+"Cert2",
		testAccAppBLBListenerResourceType,
		BaiduCloudTestResourceName)
}
