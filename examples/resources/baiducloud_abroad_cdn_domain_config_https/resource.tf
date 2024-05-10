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

resource "baiducloud_abroad_cdn_domain_config_https" "example" {
  domain = baiducloud_abroad_cdn_domain.default.domain
  cert_id = "cert-xxxxxxxxx"
  enabled = true
  http_redirect = false
  http2_enabled = false
}