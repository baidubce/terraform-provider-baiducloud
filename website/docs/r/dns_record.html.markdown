---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_dns_record"
subcategory: "DNS"
sidebar_current: "docs-baiducloud-resource-dns_record"
description: |-
  Provide a resource to create an Dns record.
---

# baiducloud_dns_record

Provide a resource to create an Dns record.

## Example Usage

```hcl
resource "baiducloud_dns_record" "default" {
  zone_name              = "testZoneName"
  rr                     = "test"
  type                   = "test"
  value                  = "test"
}
```

## Argument Reference

The following arguments are supported:

* `rr` - (Required) Dns record rr
* `type` - (Required) Dns record type
* `value` - (Required) Dns record value
* `zone_name` - (Required, ForceNew) Dns record zone name
* `description` - (Optional) Dns record description
* `line` - (Optional) Dns record line
* `priority` - (Optional) Dns record priority
* `record_action` - (Optional) Dns record action
* `ttl` - (Optional) Dns record ttl

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `record_id` - Dns record id
* `status` - Dns record status


