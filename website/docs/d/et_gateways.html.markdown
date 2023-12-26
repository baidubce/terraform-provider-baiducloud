---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_et_gateways"
sidebar_current: "docs-baiducloud-datasource-et_gateways"
description: |-
  Use this data source to query et gateways.
---

# baiducloud_et_gateways

Use this data source to query et gateways.

## Example Usage

```hcl
data "baiducloud_et_gateways" "default" {
	vpc_id = "xxxxx"
}

output "gateways" {
 value = "${data.baiducloud_et_gateways.default.gateways}"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, ForceNew) ID of the instance
* `et_gateway_id` - (Optional, ForceNew) ID of et gateway.
* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `name` - (Optional, ForceNew) name of et gateway.
* `output_file` - (Optional, ForceNew) Query result output file path
* `status` - (Optional, ForceNew) status of et gateway.

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `gateways` - et gateway
  * `create_time` - create_time of et gateway.
  * `description` - description of the et gateway
  * `et_gateway_id` - ID of et gateway.
  * `et_id` - et id of the et gateway
  * `local_cidrs` - local cidrs of the et gateway
  * `name` - name of the et gateway
  * `speed` - speed of the et gateway (Mbps)
  * `status` - status of et gateway.
  * `vpc_id` - vpc id of the et gateway


