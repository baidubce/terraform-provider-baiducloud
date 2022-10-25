resource "baiducloud_cdn_domain_config_origin" "example" {
  domain = "example.domain.com"

  range_switch = "on"

  origin_protocol {
    value = "*"
  }

  offline_mode = true

  client_ip {
    enabled = true
    name    = "True-Client-Ip"
  }

}