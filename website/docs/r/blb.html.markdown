---
layout: "baiducloud"
subcategory: "BLB"
page_title: "BaiduCloud: baiducloud_blb"
sidebar_current: "docs-baiducloud-resource-blb"
description: |-
  Provide a resource to create an BLB.
---

# baiducloud_blb

Provide a resource to create an BLB.

## Example Usage

```hcl
resource "baiducloud_blb" "default" {
  name        = "testLoadBalance"
  description = "this is a test LoadBalance instance"
  vpc_id      = "vpc-xxxx"
  subnet_id   = "sbn-xxxx"

  tags = {
    "tagAKey" = "tagAValue"
    "tagBKey" = "tagBValue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `subnet_id` - (Required, ForceNew) The subnet ID to which the LoadBalance instance belongs
* `vpc_id` - (Required, ForceNew) The VPC short ID to which the LoadBalance instance belongs
* `description` - (Optional) LoadBalance's description, length must be between 0 and 450 bytes, and support Chinese
* `name` - (Optional) LoadBalance instance's name, length must be between 1 and 65 bytes, and will be automatically generated if not set
* `tags` - (Optional, ForceNew) Tags, do not support modify

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `address` - LoadBalance instance's service IP, instance can be accessed through this IP
* `cidr` - Cidr of the network where the LoadBalance instance reside
* `create_time` - LoadBalance instance's create time
* `listener` - List of listeners mounted under the instance
  * `port` - Listening port
  * `type` - Listening protocol type
* `public_ip` - LoadBalance instance's public ip
* `status` - LoadBalance instance's status, see https://cloud.baidu.com/doc/BLB/s/Pjwvxnxdm/#blbstatus for detail
* `vpc_name` - The VPC name to which the LoadBalance instance belongs


## Import

BLB can be imported, e.g.

```hcl
$ terraform import baiducloud_blb.default id
```

