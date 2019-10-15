---
layout: "baiducloud"
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
* `tags` - (Optional, ForceNew) Tags

The `tags` object supports the following:

* `tag_key` - (Required) Tag's key
* `tag_value` - (Required) Tag's value

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `route_table_id` - Route table ID created by default on VPC creation.
* `secondary_cidrs` - Secondary cidr list of the VPC. They will not be repeated.


## Import

VPC instance can be imported, e.g.

```hcl
$ terraform import baiducloud_vpc.default vpc_id
```

