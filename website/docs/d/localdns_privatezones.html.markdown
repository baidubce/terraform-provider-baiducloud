---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_localdns_privatezones"
sidebar_current: "docs-baiducloud-datasource-localdns_privatezones"
description: |-
  Use this data source to query localdns privatezones.
---

# baiducloud_localdns_privatezones

Use this data source to query localdns privatezones.

## Example Usage

```hcl
data "baiducloud_localdns_privatezones" "default" {}

output "privatezones" {
   value = "${data.baiducloud_localdns_privatezones.default.zones}"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `output_file` - (Optional, ForceNew) local dns privatezones search result output file

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `zones` - privatezone info
  * `zone_id` - zone id
  * `zone_name` - name of privatezone
  * `record_count` - record_count of the local DNS PrivateZone

