---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_appblb_listeners"
sidebar_current: "docs-baiducloud-datasource-appblb_listeners"
description: |-
  Use this data source to query APPBLB Listener list.
---

# baiducloud_appblb_listeners

Use this data source to query APPBLB Listener list.

## Example Usage

```hcl
data "baiducloud_appblb_listeners" "default" {
 blb_id = "lb-0d29a3f6"
}

output "listeners" {
 value = "${data.baiducloud_appblb_listeners.default.listeners}"
}
```

## Argument Reference

The following arguments are supported:

* `blb_id` - (Required, ForceNew) ID of the Application LoadBalance instance
* `listener_port` - (Optional, ForceNew) The port of the Listener to be queried
* `output_file` - (Optional, ForceNew) Query result output file path
* `protocol` - (Optional, ForceNew) Protocol of the Listener to be queried

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `listeners` - A list of Application LoadBalance Listener
  * `cert_ids` - Listener bind certifications
  * `client_cert_ids` - Listener import cert list, only useful when dual_auth is true
  * `dual_auth` - Listener open dual authorization or not, default false
  * `encryption_protocols` - Listener encryption protocol
  * `encryption_type` - Listener encryption option
  * `ie6_compatible` - Listener support ie6 option, default true
  * `keep_session_cookie_name` - Listener keepSeesionCookieName
  * `keep_session_timeout` - Listener keepSessionTimeout value
  * `keep_session_type` - Listener keepSessionType option
  * `keep_session` - Listener keepSession or not
  * `listener_port` - Listener bind port
  * `policys` - Listener's policy
    * `app_server_group_id` - Policy bind server group ID
    * `app_server_group_name` - Policy bind server group name
    * `backend_port` - Backend port
    * `description` - Policy's description
    * `frontend_port` - Frontend port
    * `id` - Policy's ID
    * `port_type` - Policy bind port protocol
    * `priority` - Policy priority
    * `rule_list` - Policy rule list
      * `key` - Rule key
      * `value` - Rule value
  * `protocol` - Listener protocol
  * `redirect_port` - Listener redirect request to https listener port
  * `scheduler` - Load balancing algorithm
  * `server_timeout` - Backend server maximum timeout time, only support in [1, 3600] second, default 30s
  * `tcp_session_timeout` - TCP Listener connetion session timeout time
  * `x_forwarded_for` - Listener xForwardedFor, determine get client real ip or not, default false


