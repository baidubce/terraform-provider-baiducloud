---
layout: "baiducloud"
subcategory: "Virtual private Cloud (VPC)"
page_title: "BaiduCloud: baiducloud_vpc"
sidebar_current: "docs-baiducloud-resource-vpc"
description: |-
  Provide a resource to create a VPC.
---

# baiducloud_vpc

Provide a resource to create a VPC.

## Example Usage

```hcl
resource "baiducloud_vpc" "default" {
    name = "my-vpc"
    description = "baiducloud vpc created by terraform"
	cidr = "192.168.0.0/24"
}
```

## Argument Reference

The following arguments are supported:

* `cidr` - (Required, ForceNew) CIDR block for the VPC.
* `name` - (Required) Name of the VPC, which cannot take the value "default", the length is no more than 65 characters, and the value can be composed of numbers, characters and underscores.
* `description` - (Optional) Description of the VPC. The value is no more than 200 characters.
* `enable_ipv6` - (Optional) Whether to enable ipv6. Default is false.
* `enable_relay` - (Optional) Whether to enable route relay to allow the route table to forward traffic not originated from this VPC. When disabled, only traffic originated from this VPC will be forwarded. Default is false.
* `secondary_cidrs` - (Optional) Secondary cidr list of the VPC. They will not be repeated. replacement update.
* `tags` - (Optional, ForceNew) Tags, do not support modify

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `route_table_id` - Route table ID created by default on VPC creation.


## Import

VPC instance can be imported, e.g.

```hcl
$ terraform import baiducloud_vpc.default vpc_id
```

