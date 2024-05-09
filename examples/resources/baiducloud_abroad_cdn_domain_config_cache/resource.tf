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

resource "baiducloud_abroad_cdn_domain_config_cache" "example" {
  domain = baiducloud_abroad_cdn_domain.default.domain
  cache_ttl {
    type   = "suffix"
    value  = ".png"
    ttl    = 36000
    weight = 30
  }
  cache_ttl {
    type   = "path"
    value  = "/to/my/file"
    ttl    = 1800
    weight = 5
  }
}