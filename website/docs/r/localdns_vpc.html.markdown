---
layout: "baiducloud"
subcategory: "LOCALDNS"
page_title: "BaiduCloud: baiducloud_localdns_vpc"
sidebar_current: "docs-baiducloud-resource-localdns_vpc"
description: |-
  Use this resource to get information about a Local Dns VPC.
---

# baiducloud_localdns_vpc

Use this resource to get information about a Local Dns VPC.

~> **NOTE:** The terminate operation of vpc does NOT take effect immediately，maybe takes for several minites.

## Example Usage

```hcl
resource "baiducloud_localdns_vpc" "default" {
   zone_id = "zone-test-id"
   vpc_ids = ["vpc-test-id"]
   region = "bj"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required, ForceNew) region of the DNS  vpc
* `vpc_ids` - (Required, ForceNew) vpc_ids  of the DNS  vpc.
* `zone_id` - (Required, ForceNew) zone_id of the DNS privatezone 

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `bind_vpcs` - privatezone bind vpcs
  * `vpc_id` - bind vpc id
  * `vpc_name` - name of vpc
  * `vpc_region` - region of vpc


## Import

Local Dns vpc can be imported, e.g.

```hcl
$ terraform import baiducloud_localdns_vpc.my-server id
```

