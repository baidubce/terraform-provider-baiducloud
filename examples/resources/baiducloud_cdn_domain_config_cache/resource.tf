resource "baiducloud_cdn_domain_config_cache" "example" {
  domain = "example.domain.com"

  cache_ttl {
    type   = "suffix"
    value  = ".jpg"
    ttl    = 36000
    weight = 30
  }
  cache_ttl {
    type   = "path"
    value  = "/path/to/my/file"
    ttl    = 1800
    weight = 5
  }

  cache_url_args {
    cache_full_url = false
    cache_url_args = ["test1", "test2", "test3"]
  }

  error_page {
    code = 403
    url  = "403.html"
  }
  error_page {
    code = 404
    url  = "404.html"
  }

  cache_share {
    enabled = true
    domain  = "example2.domain.com"
  }

  mobile_access {
    distinguish_client = true
  }

}