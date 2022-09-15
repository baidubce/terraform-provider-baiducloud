resource "baiducloud_cdn_domain_config_acl" "example" {
  domain = "example.domain.com"

  referer_acl {
    allow_empty = false
    black_list = ["www.xxx.com", "*.abcde.com"]
  }

  ip_acl {
    white_list = ["1.2.3.4", "2.3.4.5"]
  }

  ua_acl {
    black_list = ["MQQBrowser/5.3/Mozilla/5.0", "Mozilla/5.0 (Linux; Android 7.0"]
  }

  cors {
    allow = "on"
    origin_list = ["https://www.baidu.com", "http://*.bce.com"]
  }

  access_limit {
    enabled = true
    limit   = 500
  }

  traffic_limit {
    enable           = true
    limit_start_hour = 1
    limit_end_hour   = 23
    limit_rate       = 500
  }

  request_auth {
    type = "A"
    key1 = "ABCD1234"
    key2 = "abcd5678"
    timeout = 1000
    timestamp_metric = 10
  }

}