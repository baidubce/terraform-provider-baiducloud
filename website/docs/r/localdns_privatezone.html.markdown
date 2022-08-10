---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_localdns_privatezone"
sidebar_current: "docs-baiducloud-resource-localdns_privatezone"
description: |-
  Use this resource to get information about a Local Dns PrivateZone.
---

# baiducloud_localdns_privatezone

Use this resource to get information about a Local Dns PrivateZone.

~> **NOTE:** The terminate operation of PrivateZone does NOT take effect immediately，maybe takes for several minites.

## Example Usage

```hcl
resource "baiducloud_localdns_privatezone" "my-server" {
  zone_name = "terrraform.com"
}
```

## Argument Reference

The following arguments are supported:

* `zone_name` - (Required, ForceNew) name of the DNS local PrivateZone

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - Creation time of the DNS local PrivateZone.
* `record_count` - record_count of the DNS local PrivateZone.
* `update_time` - update time of the DNS local PrivateZone.


## Import

Local Dns PrivateZone can be imported, e.g.

```hcl
$ terraform import baiducloud_localdns_privatezone.my-server id
```

