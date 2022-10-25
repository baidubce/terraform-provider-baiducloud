data "baiducloud_cdn_domain_certificate" "example" {
  domain = "example.domain.com"
}

resource "baiducloud_cdn_domain_config_https" "example" {
  domain = "example.domain.com"

  https {
    enabled             = true
    cert_id             = "${data.baiducloud_cdn_domain_certificate.example.certificate.0.cert_id}"
    http_redirect       = true
    http_redirect_code  = 301
    https_redirect      = false
    http2_enabled       = true
    verify_client       = true
    ssl_protocols       = ["TLSv1.1", "TLSv1.2", "TLSv1.3"]
  }

  ocsp = true
}