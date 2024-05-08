resource "baiducloud_abroad_cdn_domain" "default" {
  domain = "test.cdn.cloud"

  origin {
    backup = false
    type   = "IP"
    addr   = "1.2.3.4"
  }
  tags = {
    terraform = "terraform-test2"
  }
}