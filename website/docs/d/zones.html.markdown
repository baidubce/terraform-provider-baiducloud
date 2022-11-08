---
layout: "baiducloud"
subcategory: "Baidu Cloud Compute (BCC)"
page_title: "BaiduCloud: baiducloud_zones"
sidebar_current: "docs-baiducloud-datasource-zones"
description: |-
  Use this data source to query zone list.
---

# baiducloud_zones

Use this data source to query zone list.

## Example Usage

```hcl
data "baiducloud_zones" "default" {}

output "zone" {
  value = "${data.baiducloud_zones.default.zones}"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `name_regex` - (Optional, ForceNew) Regex pattern of the search zone name
* `output_file` - (Optional, ForceNew) Output file for saving result.

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `zones` - Useful zone list
  * `zone_name` - Useful zone name


