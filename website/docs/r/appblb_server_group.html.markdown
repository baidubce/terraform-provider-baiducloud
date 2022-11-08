---
layout: "baiducloud"
subcategory: "Application Load Balance (APPBLB)"
page_title: "BaiduCloud: baiducloud_appblb_server_group"
sidebar_current: "docs-baiducloud-resource-appblb_server_group"
description: |-
  Provide a resource to create an APPBLB Server Group.
---

# baiducloud_appblb_server_group

Provide a resource to create an APPBLB Server Group.

## Example Usage

```hcl
resource "baiducloud_appblb_server_group" "default" {
  name        = "testServerGroup"
  description = "this is a test Server Group"
  blb_id      = "lb-0d29a3f6"

  backend_server_list {
    instance_id = "i-VRKxC21a"
    weight = 50
  }

  port_list {
    port = 66
    type = "TCP"
  }
}
```

## Argument Reference

The following arguments are supported:

* `blb_id` - (Required, ForceNew) ID of the Application LoadBalance instance
* `backend_server_list` - (Optional) Server group bound backend server list
* `description` - (Optional) Server Group's description, length must be between 0 and 450 bytes, and support Chinese
* `name` - (Optional) Name of the Server Group, length must be between 1 and 65 bytes, and will be automatically generated if not set
* `port_list` - (Optional) Server Group backend port list

The `backend_server_list` object supports the following:

* `instance_id` - (Required) Backend server instance ID
* `weight` - (Required) Backend server instance weight in this group, range from 0-100
* `port_list` - Backend server open port list
  * `backend_port` - Backend open port
  * `health_check_port_type` - Health check port protocol type
  * `listener_port` - Listener port
  * `policy_id` - Port bind policy id
  * `port_id` - Port id
  * `port_type` - Port protocol type
  * `status` - Port status, include Alive/Dead/Unknown
* `private_ip` - Backend server instance bind private ip

The `port_list` object supports the following:

* `backend_port` - Backend open port
* `health_check_port_type` - Health check port protocol type
* `listener_port` - Listener port
* `policy_id` - Port bind policy id
* `port_id` - Port id
* `port_type` - Port protocol type
* `status` - Port status, include Alive/Dead/Unknown

The `port_list` object supports the following:

* `health_check` - (Required) Server Group port health check protocol, support TCP/UDP/HTTP, default same as port protocol type
* `port` - (Required) App Server Group port, range from 1-65535
* `type` - (Required) Server Group port protocol type, support TCP/UDP/HTTP
* `health_check_down_retry` - (Optional) Server Group health check down retry time, support in [2, 5], default 3
* `health_check_interval_in_second` - (Optional) Server Group health check interval time(second), support in [1, 10], default 3
* `health_check_normal_status` - (Optional) Server Group health check normal http status code, only useful when health_check is HTTP
* `health_check_port` - (Optional) Server Group port health check port, default same as Server Group port
* `health_check_timeout_in_second` - (Optional) Server Group health check timeout(second), support in [1, 60], default 3
* `health_check_up_retry` - (Optional) Server Group health check up retry time, support in [2, 5], default 3
* `health_check_url_path` - (Optional) Server Group health check url path, only useful when health_check is HTTP
* `udp_health_check_string` - (Optional) Server Group udp health check string, if type is UDP, this parameter is required
* `id` - Server Group port id
* `status` - Server Group port status

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `status` - Server Group's status, see https://cloud.baidu.com/doc/BLB/s/Pjwvxnxdm/#blbstatus for detail


