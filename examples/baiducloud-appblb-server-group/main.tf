provider "baiducloud" {}

data "baiducloud_specs" "default" {
  #name_regex        = "bcc.g1.tiny"
  #instance_type     = "General"
  cpu_count         = 1
  memory_size_in_gb = 4
}

data "baiducloud_zones" "default" {
  name_regex = ".*a$"
}

data "baiducloud_images" "default" {
  image_type = "System"
  name_regex = "7.5.*"
  os_name    = "CentOS"
}

resource "baiducloud_vpc" "default" {
  name        = var.vpc_name
  description = "test"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = var.subnet_name
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = var.description
}

resource "baiducloud_security_group" "default" {
  name        = var.sg_name
  description = var.description
  vpc_id      = baiducloud_vpc.default.id
}

resource "baiducloud_instance" "default" {
  name                  = var.bcc_name
  image_id              = data.baiducloud_images.default.images.0.id
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  subnet_id             = baiducloud_subnet.default.id
  security_groups       = [baiducloud_security_group.default.id]

  billing = {
    payment_timing = "Postpaid"
  }
}

resource "baiducloud_appblb" "default" {
  depends_on  = [baiducloud_instance.default]
  name        = var.appblb_name
  description = var.description
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}

resource "baiducloud_appblb_server_group" "default" {
  name        = var.servergroup_name
  description = var.description
  blb_id      = baiducloud_appblb.default.id

  backend_server_list {
    instance_id = baiducloud_instance.default.id
    weight      = 59
  }

  port_list {
    port = 66
    type = "UDP"
    # if type is UDP, health_check only support UDP
    # health_check default same as Server Group port
    health_check                    = "UDP"

    # optional
    health_check_port               = 66
    health_check_timeout_in_second  = 3
    health_check_interval_in_second = 3
    health_check_down_retry         = 3
    health_check_up_retry           = 3

    # required if health_check is UDP
    udp_health_check_string = "baidunew.com"
  }

  port_list {
    port = 77
    type = "TCP"
    # if type is TCP, health_check only support TCP
    # health_check default same as Server Group port
    health_check                    = "TCP"

    # optional
    health_check_port               = 77
    health_check_timeout_in_second  = 3
    health_check_interval_in_second = 3
    health_check_down_retry         = 3
    health_check_up_retry           = 3
  }

  port_list {
    port = 88
    type = "HTTP"
    # if type is HTTP, health_check support TCP and HTTP
    # health_check default same as Server Group port
    health_check                    = "HTTP"

    # optional
    health_check_port               = 88
    health_check_timeout_in_second  = 3
    health_check_interval_in_second = 3
    health_check_down_retry         = 3
    health_check_up_retry           = 3
    health_check_normal_status      = "http_2xx|http_3xx"
    health_check_url_path           = "/health/check"
  }

  port_list {
    port = 99
    type = "HTTP"
    # if type is HTTP, health_check support TCP and HTTP
    # health_check default same as Server Group port
    health_check                    = "TCP"

    # optional
    health_check_port               = 99
    health_check_timeout_in_second  = 3
    health_check_interval_in_second = 3
    health_check_down_retry         = 3
    health_check_up_retry           = 3
  }
}

data "baiducloud_appblb_server_groups" "default" {
  blb_id = baiducloud_appblb_server_group.default.blb_id
}