---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_dns_customlines"
sidebar_current: "docs-baiducloud-datasource-dns_customlines"
description: |-
  Use this data source to query Dns customline list.
---

# baiducloud_dns_customlines

Use this data source to query Dns customline list.

## Example Usage

```hcl
data "baiducloud_dns_customlines" "default" {
	name = "xxxx"
}

output "customlines" {
 value = "${data.baiducloud_dns_customlines.default.customlines}"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `output_file` - (Optional, ForceNew) DNS customlines search result output file

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `customlines` - customline list
  * `line_id` - Dns customline id
  * `lines` - lines of dns 
  * `name` - Dns customline name
  * `related_record_count` - Dns customline related record count
  * `related_zone_count` - Dns customline related zone count


