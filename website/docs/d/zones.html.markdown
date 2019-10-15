---
layout: "baiducloud"
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

* `output_file` - (Optional, ForceNew) Output file for saving result.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `zones` - Useful zone list
  * `zone_name` - Useful zone name


