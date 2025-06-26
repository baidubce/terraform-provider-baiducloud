---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_eipbps"
subcategory: "Elastic IP (EIP)"
sidebar_current: "docs-baiducloud-datasource-eipbps"
description: |-
  Use this data source to query EIP bp list.
---

# baiducloud_eipbps

Use this data source to query EIP bp list.

## Example Usage

```hcl
data "baiducloud_eipbps" "default" {}

output "eip_bps" {
 value = "${data.baiducloud_eipbps.default.eip_bps}"
}
```

## Argument Reference

The following arguments are supported:

* `bind_type` - (Optional, ForceNew) Eip bp bind type
* `bp_id` - (Optional, ForceNew) Id of Eip bp
* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `name` - (Optional, ForceNew) name of Eip bp
* `output_file` - (Optional, ForceNew) Eipbps search result output file
* `type` - (Optional, ForceNew) Eip bp type

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `eip_bps` - Eip bp list
  * `auto_release_time` - Eip bp auto release time
  * `band_width_in_mbps` - Eip bp band width in mbps
  * `create_time` - Eip bp create time
  * `eips` - Eip bp eips
  * `id` - Eip bp id
  * `instance_id` - Eip bp instance id
  * `name` - Eip bp name
  * `region` - Eip bp region
  * `type` - Eip bp type


