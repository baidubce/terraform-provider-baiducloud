---
layout: "baiducloud"
subcategory: "Baidu Load Balance (BLB)"
page_title: "BaiduCloud: baiducloud_blb_backend_server"
sidebar_current: "docs-baiducloud-resource-blb_backend_server"
description: |-
  Provide a resource to create an BLB Backend Server.
---

# baiducloud_blb_backend_server

Provide a resource to create an BLB Backend Server.

## Example Usage

```hcl
resource "baiducloud_blb_backend_server" "default" {
  blb_id      = "lb-0d29xxx6"

  backend_server_list {
    instance_id = "i-VRxxxx1a"
    weight = 50
  }
}
```

## Argument Reference

The following arguments are supported:

* `backend_server_list` - (Required) Server group bound backend server list
* `blb_id` - (Required, ForceNew) ID of the lication LoadBalance instance

The `backend_server_list` object supports the following:

* `instance_id` - (Required, ForceNew) Backend server instance ID
* `weight` - (Required) Backend server instance weight in this group, range from 0-100
* `private_ip` - Backend server instance bind private ip


