---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_dns_zones"
sidebar_current: "docs-baiducloud-datasource-dns_zones"
description: |-
  Use this data source to query Dns zone list.
---

# baiducloud_dns_zones

Use this data source to query Dns zone list.

## Example Usage

```hcl
data "baiducloud_dns_zones" "default" {
	name = "xxxx"
}

output "zones" {
 value = "${data.baiducloud_dns_zones.default.zones}"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `name` - (Optional, ForceNew) name of DNS ZONE
* `output_file` - (Optional, ForceNew) DNS Zones search result output file

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `zones` - zone list
  * `create_time` - Dns zone create_time
  * `expire_time` - Dns zone expire_time
  * `product_version` - Dns zone product_version
  * `status` - Dns zone status
  * `zone_id` - Dns zone id


