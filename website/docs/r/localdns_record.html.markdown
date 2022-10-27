---
layout: "baiducloud"
subcategory: "LOCALDNS"
page_title: "BaiduCloud: baiducloud_localdns_record"
sidebar_current: "docs-baiducloud-resource-localdns_record"
description: |-
Provide a resource to create a local dns record.
---

# baiducloud_localdns_record

Provide a resource to create a local dns record.

## Example Usage

```hcl
resource "baiducloud_localdns_record" "default" {
  zone_id     = "test-id"
  rr          = "www"
  value       = "1.1.1.1"
  type        = "A"
  ttl         = "3000"
  priority    = 0
  description = "terraform_test"
  status      = "enable"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - id of the private zone
* `rr` - record of the host.
* `value` - value of the record.
* `type` - type of the record.
* `ttl` - time to live, default is 60.
* `priority` - MX type record priority, if other types, the value is 0.
* `description` - description of the record.
* `status` - status of the record. pause or enable

