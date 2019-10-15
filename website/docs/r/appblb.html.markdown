---
layout: "baiducloud"
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

  tags {
    tag_key   = "tagAKey"
    tag_value = "tagAValue"
  }

  tags {
    tag_key   = "tagBKey"
    tag_value = "tagBValue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `subnet_id` - (Required, ForceNew) The subnet ID to which the LoadBalance instance belongs
* `vpc_id` - (Required, ForceNew) The VPC short ID to which the LoadBalance instance belongs
* `description` - (Optional) LoadBalance's description, length must be between 0 and 450 bytes, and support Chinese
* `name` - (Optional) LoadBalance instance's name, length must be between 1 and 65 bytes, and will be automatically generated if not set
* `tags` - (Optional, ForceNew) Tags

The `tags` object supports the following:

* `tag_key` - (Required) Tag's key
* `tag_value` - (Required) Tag's value

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `address` - LoadBalance instance's service IP, instance can be accessed through this IP
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

