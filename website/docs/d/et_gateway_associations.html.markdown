---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_et_gateway_associations"
subcategory: "Virtual private Cloud (VPC)"
sidebar_current: "docs-baiducloud-datasource-et_gateway_associations"
description: |-
  Use this data source to query et gateway associations.
---

# baiducloud_et_gateway_associations

Use this data source to query et gateway associations.

## Example Usage

```hcl
data "baiducloud_et_gateway_associations" "default" {
	et_gateway_id = "xxxxx"
}

output "gateway" {
 value = "${data.baiducloud_et_gateway_associations.default.gateway_associations}"
}
```

## Argument Reference

The following arguments are supported:

* `et_gateway_id` - (Optional, ForceNew) ID of et gateway.
* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `output_file` - (Optional, ForceNew) Query result output file path

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `gateway_associations` - et gateway associations
  * `create_time` - create_time of et gateway.
  * `health_check_dest_ip` - health_check_dest_ip of et gateway.
  * `health_check_interval` - health_check_interval of et gateway.
  * `health_check_source_ip` - health_check_source_ip of et gateway.
  * `health_check_type` - health_check_type of et gateway.
  * `health_threshold` - health_threshold of et gateway.
  * `local_cidrs` - local cidrs of the et gateway
  * `name` - name of et gateway.
  * `speed` - speed of the et gateway (Mbps)
  * `status` - status of et gateway.
  * `unhealth_threshold` - unhealth_threshold of et gateway.


