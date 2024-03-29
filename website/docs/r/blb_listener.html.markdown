---
layout: "baiducloud"
subcategory: "Baidu Load Balance (BLB)"
page_title: "BaiduCloud: baiducloud_blb_listener"
sidebar_current: "docs-baiducloud-resource-blb_listener"
description: |-
  Provide a resource to create an BLB Listener.
---

# baiducloud_blb_listener

Provide a resource to create an BLB Listener.

## Example Usage

```hcl
[TCP/UDP] Listener
resource "baiducloud_blb_listener" "default" {
  blb_id               = "lb-0d29a3f6"
  listener_port        = 124
  protocol             = "TCP"
  scheduler            = "LeastConnection"
}

[HTTP] Listener
resource "baiducloud_blb_listener" "default" {
  blb_id        = "lb-0d29a3f6"
  listener_port = 129
  protocol      = "HTTP"
  scheduler     = "RoundRobin"

}

[HTTPS] Listener
resource "baiducloud_blb_listener" "default" {
  blb_id               = "lb-0d29a3f6"
  listener_port        = 130
  protocol             = "HTTPS"
  scheduler            = "LeastConnection"
  keep_session         = true
  cert_ids             = ["cert-xvysj8xxx"]
  encryption_protocols = ["sslv3", "tlsv10", "tlsv11"]
  encryption_type      = "userDefind"
}

[SSL] Listener
resource "baiducloud_blb_listener" "default" {
  blb_id               = "lb-0d29a3f6"
  listener_port        = 131
  protocol             = "SSL"
  scheduler            = "LeastConnection"
  cert_ids             = ["cert-xvysjxxxx"]
  encryption_protocols = ["sslv3", "tlsv10", "tlsv11"]
  encryption_type      = "userDefind"
}
```

## Argument Reference

The following arguments are supported:

* `backend_port` - (Required, ForceNew) backend port, range from 1-65535
* `blb_id` - (Required, ForceNew) ID of the Application LoadBalance instance
* `listener_port` - (Required, ForceNew) Listening port, range from 1-65535
* `protocol` - (Required, ForceNew) Listening protocol, support TCP/UDP/HTTP/HTTPS/SSL
* `scheduler` - (Required, ForceNew) Load balancing algorithm, support RoundRobin/LeastConnection/Hash, if protocol is HTTP/HTTPS, only support RoundRobin/LeastConnection
* `applied_ciphers` - (Optional) applied ciphers
* `cert_ids` - (Optional) Listener bind certifications
* `client_cert_ids` - (Optional) Listener import cert list, only useful when dual_auth is true
* `dual_auth` - (Optional) Listener open dual authorization or not, default false
* `encryption_protocols` - (Optional) Listener encryption protocol, only useful when encryption_type is userDefind, support [sslv3, tlsv10, tlsv11, tlsv12]
* `encryption_type` - (Optional) Listener encryption option, support [compatibleIE, incompatibleIE, userDefind]
* `get_blb_ip` - (Optional) get blb ip or not
* `health_check_interval` - (Optional) health check interval
* `health_check_normal_status` - (Optional) health check normal status
* `health_check_port` - (Optional) health check port
* `health_check_string` - (Optional) health check string, This parameter is mandatory when the listening protocol is UDP
* `health_check_timeout_in_second` - (Optional) health check timeout in second
* `health_check_type` - (Optional) health check type
* `health_check_uri` - (Optional) health check uri
* `healthy_threshold` - (Optional) healthy threshold
* `ie6_compatible` - (Optional) Listener support ie6 option, default true
* `keep_session_cookie_name` - (Optional) CookieName which need to covered, useful when keep_session_type is rewrite
* `keep_session_duration` - (Optional) keep session duration
* `keep_session_type` - (Optional) KeepSessionType option, support insert/rewrite, default insert
* `keep_session` - (Optional) KeepSession or not
* `redirect_port` - (Optional) Redirect HTTP request to HTTPS Listener, HTTPS Listener port set by this parameter
* `server_timeout` - (Optional) Backend server maximum timeout time, only support in [1, 3600] second, default 30s
* `tcp_session_timeout` - (Optional) TCP Listener connection session timeout time(second), default 900, support 10-4000
* `udp_session_timeout` - (Optional) UDP Listener connection session timeout time(second), default 900, support 10-4000
* `unhealthy_threshold` - (Optional) unhealthy threshold
* `x_forwarded_for` - (Optional) Listener xForwardedFor, determine get client real ip or not, default false


