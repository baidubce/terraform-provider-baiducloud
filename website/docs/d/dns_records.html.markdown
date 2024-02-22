---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_dns_records"
sidebar_current: "docs-baiducloud-datasource-dns_records"
description: |-
  Use this data source to query Dns record list.
---

# baiducloud_dns_records

Use this data source to query Dns record list.

## Example Usage

```hcl
data "baiducloud_dns_records" "default" {
	zone_name = "xxxx"
}

output "records" {
 value = "${data.baiducloud_dns_records.default.records}"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `output_file` - (Optional, ForceNew) DNS records search result output file

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `records` - record list
  * `rr` - Dns record rr
  * `type` - Dns record type
  * `value` - Dns record value
  * `ttl` - Dns record ttl
  * `line` - Dns record line
  * `description` - Dns record description
  * `priority` - Dns record priority
  * `record_id` - Dns record id
  * `status` - Dns record status
