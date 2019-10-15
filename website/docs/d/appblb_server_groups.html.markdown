---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_appblb_server_groups"
sidebar_current: "docs-baiducloud-datasource-appblb_server_groups"
description: |-
  Use this data source to query APPBLB Server Group list.
---

# baiducloud_appblb_server_groups

Use this data source to query APPBLB Server Group list.

## Example Usage

```hcl
data "baiducloud_appblb_server_groups" "default" {
 name = "testServerGroup"
}

output "server_groups" {
 value = "${data.baiducloud_appblb_server_groups.default.server_groups}"
}
```

## Argument Reference

The following arguments are supported:

* `blb_id` - (Required) ID of the LoadBalance instance to be queried
* `exactly_match` - (Optional) Whether the name is an exact match or not, default false
* `name` - (Optional) Name of the Server Group to be queried
* `output_file` - (Optional, ForceNew) Query result output file path

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `server_groups` - A list of Application LoadBalance Server Group
  * `backend_server_list` - Server group bound backend server list
    * `instance_id` - Backend server instance ID
    * `port_list` - Backend server open port list
      * `backend_port` - Backend open port
      * `health_check_port_type` - Health check port protocol type
      * `listener_port` - Listener port
      * `policy_id` - Port bind policy ID
      * `port_id` - Port ID
      * `port_type` - Port protocol type
      * `status` - Port status, include Alive/Dead/Unknown
    * `private_ip` - Backend server instance bind private ip
    * `weight` - Backend server instance weight in this group
  * `description` - Server Group's description
  * `name` - Server Group's name
  * `port_list` - Server Group backend port list
    * `health_check_down_retry` - Server Group health check down retry time
    * `health_check_interval_in_second` - Server Group health check interval time(second)
    * `health_check_normal_status` - Server Group health check normal http status code
    * `health_check_port` - Server Group port health check port
    * `health_check_timeout_in_second` - Server Group health check timeout(second)
    * `health_check_up_retry` - Server Group health check up retry time
    * `health_check_url_path` - Server Group health check url path
    * `health_check` - Server Group port health check protocol
    * `id` - Server Group port ID
    * `port` - Server Group port
    * `status` - Server Group port status
    * `type` - Server Group port protocol type
    * `udp_health_check_string` - Server Group udp health check string
  * `sg_id` - Server Group's ID
  * `status` - Server Group status


