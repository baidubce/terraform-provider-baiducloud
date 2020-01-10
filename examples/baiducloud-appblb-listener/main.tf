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

resource "baiducloud_cert" "default" {
  cert_name         = var.cert_name
  cert_server_data  = file("${path.module}/cert.crt")
  cert_private_data = file("${path.module}/cert.key")
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

  port_list {
    port         = 70
    type         = "HTTP"
    health_check = "HTTP"
  }
  port_list {
    port         = 68
    type         = "TCP"
    health_check = "TCP"
  }
  port_list {
    port                    = 66
    type                    = "UDP"
    health_check            = "UDP"
    udp_health_check_string = "baidu.com"
  }
}

resource "baiducloud_appblb_listener" "default_UDP" {
  blb_id        = baiducloud_appblb.default.id
  listener_port = 125
  protocol      = "UDP"
  scheduler     = "RoundRobin"

  policies {
    description         = "acceptance test"
    app_server_group_id = baiducloud_appblb_server_group.default.id
    # UDP listener should has backend_port which server group has UDP port_list
    backend_port = 66
    priority     = 50

    # UDP listener only support *:*
    rule_list {
      key   = "*"
      value = "*"
    }
  }
}

resource "baiducloud_appblb_listener" "default_TCP" {
  blb_id              = baiducloud_appblb.default.id
  listener_port       = 124
  protocol            = "TCP"
  scheduler           = "LeastConnection"
  tcp_session_timeout = 1000

  policies {
    description         = "acceptance test"
    app_server_group_id = baiducloud_appblb_server_group.default.id
    # TCP listener should has backend_port which server group has TCP port_list
    backend_port = 68
    priority     = 50

    # TCP listener only support *:*
    rule_list {
      key   = "*"
      value = "*"
    }
  }
}

resource "baiducloud_appblb_listener" "default_HTTP" {
  blb_id        = baiducloud_appblb.default.id
  listener_port = 129
  protocol      = "HTTP"
  scheduler     = "RoundRobin"
  keep_session  = true

  policies {
    description         = "acceptance test"
    app_server_group_id = baiducloud_appblb_server_group.default.id
    # HTTP listener should has backend_port which server group has HTTP port_list
    backend_port = 70
    priority     = 50

    rule_list {
      key   = "host"
      value = "baidu.com"
    }
  }
}

resource "baiducloud_appblb_listener" "default_HTTPS" {
  blb_id               = baiducloud_appblb.default.id
  listener_port        = 130
  protocol             = "HTTPS"
  scheduler            = "LeastConnection"
  keep_session         = true
  cert_ids             = [baiducloud_cert.default.id]
  encryption_protocols = ["sslv3", "tlsv10", "tlsv11"]
  encryption_type      = "userDefind"

  policies {
    description         = "acceptance test"
    app_server_group_id = baiducloud_appblb_server_group.default.id
    # HTTPS listener should has backend_port which server group has HTTP port_list
    backend_port = 70
    priority     = 50

    rule_list {
      key   = "host"
      value = "baidu.com"
    }
  }
}

resource "baiducloud_appblb_listener" "default_SSL" {
  blb_id               = baiducloud_appblb.default.id
  listener_port        = 131
  protocol             = "SSL"
  scheduler            = "RoundRobin"
  cert_ids             = [baiducloud_cert.default.id]
  encryption_protocols = ["tlsv10", "tlsv11", "tlsv12"]
  encryption_type      = "userDefind"

  policies {
    description         = "acceptance test"
    app_server_group_id = baiducloud_appblb_server_group.default.id
    # SSL listener should has backend_port which server group has TCP port_list
    backend_port = 68
    priority     = 50

    # SSL listener only support *:*
    rule_list {
      key   = "*"
      value = "*"
    }
  }
}

data "baiducloud_appblb_listeners" "default" {
  depends_on = [baiducloud_appblb_listener.default_TCP, baiducloud_appblb_listener.default_UDP, baiducloud_appblb_listener.default_SSL, baiducloud_appblb_listener.default_HTTP, baiducloud_appblb_listener.default_HTTPS]
  blb_id     = baiducloud_appblb.default.id
}
