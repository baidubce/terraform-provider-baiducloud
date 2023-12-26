---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_et_gateway"
sidebar_current: "docs-baiducloud-resource-et_gateway"
description: |-
  Use this resource to get information about a ET Gateway.
---

# baiducloud_et_gateway

Use this resource to get information about a ET Gateway.

## Example Usage

```hcl
resource "baiducloud_et_gateway" "default" {
	name = "my_name"
	vpc_id = "vpc-xxx"
	speed = 200
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) name of the et gateway
* `speed` - (Required) speed of the et gateway (Mbps)
* `vpc_id` - (Required, ForceNew) vpc id of the et gateway
* `channel_id` - (Optional, ForceNew) channel id of the et gateway
* `description` - (Optional) description of the et gateway
* `et_id` - (Optional, ForceNew) et id of the et gateway
* `local_cidrs` - (Optional) local cidrs of the et gateway

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - create_time of et gateway.
* `et_gateway_id` - ID of et gateway.
* `health_check_dest_ip` - health_check_dest_ip of et gateway.
* `health_check_interval` - health_check_interval of et gateway.
* `health_check_source_ip` - health_check_source_ip of et gateway.
* `health_check_type` - health_check_type of et gateway.
* `health_threshold` - health_threshold of et gateway.
* `status` - status of et gateway.
* `unhealth_threshold` - unhealth_threshold of et gateway.


## Import

ET Gateway can be imported, e.g.

```hcl
$ terraform import baiducloud_et_gateway.default eip
```

