---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_subnet"
sidebar_current: "docs-baiducloud-resource-subnet"
description: |-
  Provide a resource to create a VPC subnet.
---

# baiducloud_subnet

Provide a resource to create a VPC subnet.

## Example Usage

```hcl
resource "baiducloud_subnet" "default" {
  name = "my-subnet"
  zone_name = "cn-bj-a"
  cidr = "192.168.3.0/24"
  vpc_id = "${baiducloud_vpc.default.id}"
}

resource "baiducloud_vpc" "default" {
  name = "my-vpc"
  cidr = "192.168.0.0/16"
}
```

## Argument Reference

The following arguments are supported:

* `cidr` - (Required, ForceNew) CIDR block of the subnet.
* `name` - (Required) Name of the subnet, which cannot take the value "default", the length is no more than 65 characters, and the value can be composed of numbers, characters and underscores.
* `vpc_id` - (Required, ForceNew) ID of the VPC.
* `zone_name` - (Required, ForceNew) The availability zone name within which the subnet should be created.
* `description` - (Optional) Description of the subnet, and the value must be no more than 200 characters.
* `subnet_type` - (Optional, ForceNew) Type of the subnet, valid values are BCC, BCC_NAT and BBC. Default to BCC.
* `tags` - (Optional, ForceNew) Tags, do not support modify


## Import

VPC subnet instance can be imported, e.g.

```hcl
$ terraform import baiducloud_subnet.default subnet_id
```

