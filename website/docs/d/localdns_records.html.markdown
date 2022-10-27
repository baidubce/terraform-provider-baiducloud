---
layout: "baiducloud"
subcategory: "LOCALDNS"
page_title: "BaiduCloud: baiducloud_localdns_records"
sidebar_current: "docs-baiducloud-datasource-localdns_records"
description: |-
  Use this data source to query localdns records list.
---

# baiducloud_localdns_records

Use this data source to query localdns records list.

## Example Usage

```hcl
data "baiducloud_localdns_records" "default" {
  zone_id = "test-id"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `zone_id` - (Optional, ForceNew) id of the specific private zone to retrieve.
* `output_file` - (Optional, ForceNew) Output file for saving result.

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `records` - Result of local dns records.
    * `record_id` - ID of the local dns record.
    * `rr` - record of the host.
    * `value` - value of the record.
    * `type` - type of the record.
    * `ttl` - time to live, default is 60.
    * `priority` - MX type record priority, if other types, the value is 0.
    * `description` - description of the record.
    * `status` - status of the record. pause or enable


