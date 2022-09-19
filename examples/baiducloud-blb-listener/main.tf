provider "baiducloud" {
  # option config, you can use assume role as the operation account
  #assume_role {
  #  account_id = "your-account-id"
  #  role_name = "your-role-name"
  #}
}


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

resource "baiducloud_blb" "default" {
  depends_on  = [baiducloud_instance.default]
  name        = var.blb_name
  description = var.description
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}


# for more detailed config, please refer to https://cloud.baidu.com/doc/BLB/s/ujwvxnyux
resource "baiducloud_blb_listener" "default_UDP" {
  blb_id        = baiducloud_appblb.default.id
  listener_port = 125
  protocol      = "UDP"
  # support RoundRobin/LeastConnection/Hash
  scheduler     = "RoundRobin"

}

resource "baiducloud_blb_listener" "default_TCP" {
  blb_id              = baiducloud_appblb.default.id
  listener_port       = 124
  protocol            = "TCP"
  # support RoundRobin/LeastConnection/Hash
  scheduler           = "LeastConnection"
  tcp_session_timeout = 1000

}

resource "baiducloud_blb_listener" "default_HTTP" {
  blb_id        = baiducloud_appblb.default.id
  listener_port = 129
  protocol      = "HTTP"
  # support RoundRobin/LeastConnection
  scheduler     = "RoundRobin"

  # optional
  # keep_session  = true

  # support insert/rewrite
  # keep_session_type = "insert"

  # keep_session_timeout = 3600

  # only useful when keep_session_type is rewrite
  # keep_session_cookie_name = "aaa"

  # x_forwarded_for = false
  # server_timeout = 30
  # redirect_port = 80

}

data "baiducloud_blb_listeners" "default" {
  depends_on = [baiducloud_blb_listener.default_TCP, baiducloud_blb_listener.default_UDP, baiducloud_blb_listener.default_SSL, baiducloud_blb_listener.default_HTTP, baiducloud_blb_listener.default_HTTPS]
  blb_id     = baiducloud_blb.default.id
}
