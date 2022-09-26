---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_blb_backend_servers"
sidebar_current: "docs-baiducloud-datasource-blb_backend_servers"
description: |-
  Use this data source to query BLB Backend Server list.
---

# baiducloud_blb_backend_servers

Use this data source to query BLB Backend Server list.

## Example Usage

```hcl
data "baiducloud_blb_backend_servers" "default" {
 blb_id = "xxxx"
}

output "server_groups" {
 value = "${data.baiducloud_blb_backend_servers.default.backend_server_list}"
}
```

## Argument Reference

The following arguments are supported:

* `blb_id` - (Required, ForceNew) ID of the LoadBalance instance to be queried
* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `output_file` - (Optional, ForceNew) Query result output file path

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `backend_server_list` - backend server list
  * `instance_id` - Backend server instance ID
  * `private_ip` - Backend server instance bind private ip
  * `weight` - Backend server instance weight in this group, range from 0-100


