resource "baiducloud_cdn_domain" "example" {
  domain = "example.domain.com"

  origin {
    backup = false
    host   = "example1r.domain.com"
    peer   = "https://1.2.3.4:443"
  }
  origin {
    backup = false
    host   = "example2r.domain.com"
    peer   = "http://2.3.4.5:80"
  }
  origin {
    backup = true
    peer   = "http://3.4.5.6:80"
  }

  default_host = "example3.domain.com"
  form         = "image"

}