resource "baiducloud_cdn_domain" "example" {
  domain = "example.domain.com"

  origin {
    addr              = "1.2.3.4"
    type              = "IP"
    backup            = false
    host              = "example1r.domain.com"
    weight            = 20
    isp               = "un"
    upstream_protocol = "https"
  }
  origin {
    addr   = "2.3.4.5"
    type   = "IP"
    backup = false
    host   = "example2r.domain.com"
    weight = 20
  }
  origin {
    addr   = "3.4.5.6"
    type   = "IP"
    backup = true
    weight = 20
  }

  form = "image"
}
