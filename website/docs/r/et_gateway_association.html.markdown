---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_et_gateway_association"
subcategory: "Virtual private Cloud (VPC)"
sidebar_current: "docs-baiducloud-resource-et_gateway_association"
description: |-
  Provide a resource to manage an et gateway association.
---

# baiducloud_et_gateway_association

Provide a resource to manage an et gateway association.

## Example Usage

```hcl
resource "baiducloud_et_gateway_association" "default" {
  et_gateway_id = "xxx"
  et_id = "xxx"
  channel_id = "xxx"
  local_cidrs = ["192.168.0.0/20"]
}
```

## Argument Reference

The following arguments are supported:

* `et_gateway_id` - (Required, ForceNew) ID of et gateway.
* `channel_id` - (Optional, ForceNew) channel id of the et gateway
* `et_id` - (Optional, ForceNew) et id of the et gateway
* `local_cidrs` - (Optional) local cidrs of the et gateway

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - create_time of et gateway.
* `health_check_dest_ip` - health_check_dest_ip of et gateway.
* `health_check_interval` - health_check_interval of et gateway.
* `health_check_source_ip` - health_check_source_ip of et gateway.
* `health_check_type` - health_check_type of et gateway.
* `health_threshold` - health_threshold of et gateway.
* `name` - name of et gateway.
* `speed` - speed of et gateway.
* `status` - status of et gateway.
* `unhealth_threshold` - unhealth_threshold of et gateway.
* `vpc_id` - vpc id of et gateway.


## Import

ET Gateway Association can be imported, e.g.

```hcl
$ terraform import baiducloud_et_gateway_association.default et_gateway_id
```

