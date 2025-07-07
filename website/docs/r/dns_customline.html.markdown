---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_dns_customline"
subcategory: "DNS"
sidebar_current: "docs-baiducloud-resource-dns_customline"
description: |-
  Provide a resource to create an Dns customline.
---

# baiducloud_dns_customline

Provide a resource to create an Dns customline.

## Example Usage

```hcl
resource "baiducloud_dns_customline" "default" {
 name              = "testDnscustomline"
}
```

## Argument Reference

The following arguments are supported:

* `lines` - (Required) lines of dns 
* `name` - (Required, ForceNew) Dns customline name

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `line_id` - Dns customline id
* `related_record_count` - Dns customline related record count
* `related_zone_count` - Dns customline related zone count


