resource "baiducloud_appblb_ip_group" "example" {
  blb_id = "blb-xxxxx"
  name   = "ip-group-example"
  desc   = "example ip group"

  backend_policy_list {
    type                = "HTTP"
    enable_health_check = true
    health_check        = "HTTP"
    health_check_port   = 80
    health_check_timeout_in_second  = 3
    health_check_interval_in_second = 3
    health_check_down_retry         = 3
    health_check_up_retry           = 3
    health_check_normal_status = "http_2xx|http_3xx"
    health_check_url_path      = "/"
  }

  member_list {
    ip     = "192.168.0.10"
    port   = 80
    weight = 80
  }
}
