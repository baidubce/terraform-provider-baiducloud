---
layout: "baiducloud"
subcategory: "Baidu Load Balance (BLB)"
page_title: "BaiduCloud: baiducloud_blb_listeners"
sidebar_current: "docs-baiducloud-datasource-blb_listeners"
description: |-
  Use this data source to query BLB Listener list.
---

# baiducloud_blb_listeners

Use this data source to query BLB Listener list.

## Example Usage

```hcl
data "baiducloud_blb_listeners" "default" {
 blb_id = "lb-0d2xxxx6"
}

output "listeners" {
 value = "${data.baiducloud_blb_listeners.default.listeners}"
}
```

## Argument Reference

The following arguments are supported:

* `blb_id` - (Required, ForceNew) ID of the LoadBalance instance
* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `listener_port` - (Optional, ForceNew) The port of the Listener to be queried
* `output_file` - (Optional, ForceNew) Query result output file path
* `protocol` - (Optional, ForceNew) Protocol of the Listener to be queried

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `listeners` - A list of LoadBalance Listener
  * `applied_ciphers` - applied ciphers
  * `backend_port` - backend port, range from 1-65535
  * `cert_ids` - Listener bind certifications
  * `client_cert_ids` - Listener import cert list, only useful when dual_auth is true
  * `dual_auth` - Listener open dual authorization or not, default false
  * `encryption_protocols` - Listener encryption protocol
  * `encryption_type` - Listener encryption option
  * `get_blb_ip` - get blb ip or not
  * `health_check_interval` - health check interval
  * `health_check_normal_status` - health check normal status
  * `health_check_port` - health check port
  * `health_check_string` - health check string, This parameter is mandatory when the listening protocol is UDP
  * `health_check_timeout_in_second` - health check timeout in second
  * `health_check_type` - health check type
  * `health_check_uri` - health check uri
  * `healthy_threshold` - healthy threshold
  * `ie6_compatible` - Listener support ie6 option, default true
  * `keep_session_cookie_name` - Listener keepSeesionCookieName
  * `keep_session_duration` - keep session duration
  * `keep_session_timeout` - Listener keepSessionTimeout value
  * `keep_session_type` - Listener keepSessionType option
  * `keep_session` - Listener keepSession or not
  * `listener_port` - Listener bind port
  * `protocol` - Listener protocol
  * `redirect_port` - Listener redirect request to https listener port
  * `scheduler` - Load balancing algorithm
  * `server_timeout` - Backend server maximum timeout time, only support in [1, 3600] second, default 30s
  * `tcp_session_timeout` - TCP Listener connetion session timeout time
  * `udp_session_timeout` - UDP Listener connection session timeout time(second), default 900, support 10-4000
  * `unhealthy_threshold` - unhealthy threshold
  * `x_forwarded_for` - Listener xForwardedFor, determine get client real ip or not, default false


