---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_appblb_listener"
sidebar_current: "docs-baiducloud-resource-appblb_listener"
description: |-
  Provide a resource to create an APPBLB Listener.
---

# baiducloud_appblb_listener

Provide a resource to create an APPBLB Listener.

## Example Usage

```hcl
[TCP/UDP] Listener
resource "baiducloud_appblb_listener" "default" {
  blb_id               = "lb-0d29a3f6"
  listener_port        = 124
  protocol             = "TCP"
  scheduler            = "LeastConnection"
}

[HTTP] Listener
resource "baiducloud_appblb_listener" "default" {
  blb_id        = "lb-0d29a3f6"
  listener_port = 129
  protocol      = "HTTP"
  scheduler     = "RoundRobin"
  keep_session  = true

  policies {
    description         = "acceptance test"
    app_server_group_id = "sg-11bd8054"
    backend_port        = 70
    priority            = 50

    rule_list {
      key   = "host"
      value = "baidu.com"
    }
  }
}

[HTTPS] Listener
resource "baiducloud_appblb_listener" "default" {
  blb_id               = "lb-0d29a3f6"
  listener_port        = 130
  protocol             = "HTTPS"
  scheduler            = "LeastConnection"
  keep_session         = true
  cert_ids             = ["cert-xvysj80uif1y"]
  encryption_protocols = ["sslv3", "tlsv10", "tlsv11"]
  encryption_type      = "userDefind"
}

[SSL] Listener
resource "baiducloud_appblb_listener" "default" {
  blb_id               = "lb-0d29a3f6"
  listener_port        = 131
  protocol             = "SSL"
  scheduler            = "LeastConnection"
  cert_ids             = ["cert-xvysj80uif1y"]
  encryption_protocols = ["sslv3", "tlsv10", "tlsv11"]
  encryption_type      = "userDefind"
}
```

## Argument Reference

The following arguments are supported:

* `blb_id` - (Required, ForceNew) ID of the Application LoadBalance instance
* `listener_port` - (Required, ForceNew) Listening port, range from 1-65535
* `protocol` - (Required, ForceNew) Listening protocol, support TCP/UDP/HTTP/HTTPS/SSL
* `scheduler` - (Required) Load balancing algorithm, support RoundRobin/LeastConnection/Hash, if protocol is HTTP/HTTPS, only support RoundRobin/LeastConnection
* `cert_ids` - (Optional) Listener bind certifications
* `client_cert_ids` - (Optional) Listener import cert list, only useful when dual_auth is true
* `dual_auth` - (Optional) Listener open dual authorization or not, default false
* `encryption_protocols` - (Optional) Listener encryption protocol, only useful when encryption_type is userDefind, support [sslv3, tlsv10, tlsv11, tlsv12]
* `encryption_type` - (Optional) Listener encryption option, support [compatibleIE, incompatibleIE, userDefind]
* `ie6_compatible` - (Optional) Listener support ie6 option, default true
* `keep_session_cookie_name` - (Optional) CookieName which need to covered, useful when keep_session_type is rewrite
* `keep_session_timeout` - (Optional) KeepSession Cookie timeout time(second), support in [1, 15552000], default 3600s
* `keep_session_type` - (Optional) KeepSessionType option, support insert/rewrite, default insert
* `keep_session` - (Optional) KeepSession or not
* `policies` - (Optional) Listener's policy
* `redirect_port` - (Optional) Redirect HTTP request to HTTPS Listener, HTTPS Listener port set by this parameter
* `server_timeout` - (Optional) Backend server maximum timeout time, only support in [1, 3600] second, default 30s
* `tcp_session_timeout` - (Optional) TCP Listener connection session timeout time(second), default 900, support 10-4000
* `x_forwarded_for` - (Optional) Listener xForwardedFor, determine get client real ip or not, default false

The `policies` object supports the following:

* `app_server_group_id` - (Required) Policy bind server group id
* `backend_port` - (Required) Backend port
* `priority` - (Required) Policy priority, support in [1, 32768]
* `description` - (Optional) Policy's description
* `rule_list` - (Optional) Policy rule list
* `app_server_group_name` - Policy bind server group name
* `frontend_port` - Frontend port
* `id` - Policy's id
* `port_type` - Policy bind port protocol type

The `rule_list` object supports the following:

* `key` - (Required) Rule key
* `value` - (Required) Rule value


