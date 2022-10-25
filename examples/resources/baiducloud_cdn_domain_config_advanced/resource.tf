resource "baiducloud_cdn_domain_config_advanced" "example" {
  domain = "example.domain.com"

  ipv6_dispatch {
    enable = true
  }

  http_header {
    type     = "response"
    header   = "Cache-Control"
    value    = "allowFull"
    action   = "add"
    describe = "Specifies the caching mechanism."
  }
  http_header {
    type     = "origin"
    header   = "Cache-Control"
    action   = "remove"
    describe = "Specifies the caching mechanism."
  }

  media_drag {
    mp4 {
      file_suffix = ["mp4"]
      start_arg_name = "abcd"
      end_arg_name = "1234"
      drag_mode = "second"
    }
    flv {
      file_suffix = ["flv"]
      drag_mode = "byteAV"
    }
  }

  seo_switch {
    directly_origin = "ON"
  }

  file_trim = true

  compress {
    allow = true
    type = "all"
  }

  quic = true
}