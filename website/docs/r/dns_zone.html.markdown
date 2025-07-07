---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_dns_zone"
subcategory: "DNS"
sidebar_current: "docs-baiducloud-resource-dns_zone"
description: |-
  Provide a resource to create an Dns zone.
---

# baiducloud_dns_zone

Provide a resource to create an Dns zone.

## Example Usage

```hcl
resource "baiducloud_dns_zone" "default" {
  name              = "testDnsZone"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) Dns zone name
* `tags` - (Optional, ForceNew) Tags, do not support modify

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - Dns zone create_time
* `expire_time` - Dns zone expire_time
* `product_version` - Dns zone product_version
* `status` - Dns zone status
* `zone_id` - Dns zone id


