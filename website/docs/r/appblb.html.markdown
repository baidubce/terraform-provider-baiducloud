---
layout: "baiducloud"
subcategory: "Application Load Balance (APPBLB)"
page_title: "BaiduCloud: baiducloud_appblb"
sidebar_current: "docs-baiducloud-resource-appblb"
description: |-
  Provide a resource to create an APPBLB.
---

# baiducloud_appblb

Provide a resource to create an APPBLB.

## Example Usage

```hcl
resource "baiducloud_appblb" "default" {
  name        = "testLoadBalance"
  description = "this is a test LoadBalance instance"
  vpc_id      = "vpc-gxaava4knqr1"
  subnet_id   = "sbn-m4x3f2i6c901"

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
* `address` - (Optional) LoadBalance instance's service IP, instance can be accessed through this IP
* `auto_renew_length` - (Optional) The automatic renewal time is 1-9 per month and 1-3 per year.
* `auto_renew_time_unit` - (Optional) Monthly payment or annual payment, month is month and year is year.
* `description` - (Optional) LoadBalance's description, length must be between 0 and 450 bytes, and support Chinese
* `eip` - (Optional) eip of the LoadBalance
* `enterprise_security_groups` - (Optional) enterprise security group ids of the APPBLB
* `name` - (Optional) LoadBalance instance's name, length must be between 1 and 65 bytes, and will be automatically generated if not set
* `performance_level` - (Optional, ForceNew) performance level, available values are small1, small2, medium1, medium2, large1, large2, large3
* `security_groups` - (Optional) security group ids of the APPBLB.
* `tags` - (Optional, ForceNew) Tags, do not support modify

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cidr` - Cidr of the network where the LoadBalance instance reside
* `create_time` - LoadBalance instance's create time
* `listener` - List of listeners mounted under the instance
  * `port` - Listening port
  * `type` - Listening protocol type
* `public_ip` - LoadBalance instance's public ip
* `release_time` - LoadBalance instance's auto release time
* `status` - LoadBalance instance's status, see https://cloud.baidu.com/doc/BLB/s/Pjwvxnxdm/#blbstatus for detail
* `subnet_cidr` - Cidr of the subnet which the LoadBalance instance belongs
* `subnet_name` - The subnet name to which the LoadBalance instance belongs
* `vpc_name` - The VPC name to which the LoadBalance instance belongs


## Import

APPBLB can be imported, e.g.

```hcl
$ terraform import baiducloud_appblb.default id
```

