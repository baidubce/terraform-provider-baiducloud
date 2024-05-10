resource "baiducloud_abroad_cdn_domain" "default" {
  domain = "test.cdn.com"

  origin {
    backup = false
    type   = "IP"
    addr   = "1.2.3.4"
  }
  tags = {
    terraform = "terraform-test2"
  }
}

resource "baiducloud_abroad_cdn_domain_config_acl" "example" {
  domain = baiducloud_abroad_cdn_domain.default.domain
  allow_empty = false

  referer_acl {
    black_list  = ["xxx.as.com"]
  }

  ip_acl {
    black_list = ["2.3.4.5"]
  }
}