---
layout: "baiducloud"
subcategory: "Cloud File Storage (CFS)"
page_title: "BaiduCloud: baiducloud_cfs_mount_target"
sidebar_current: "docs-baiducloud-resource-cfs_mount_target"
description: |-
  Use this resource to create a CFS mount target.
---

# baiducloud_cfs_mount_target

Use this resource to create a CFS mount target.

~> **NOTE:** 
For multiple subnet IPs in the same VPC, different mount targets can be added respectively; A subnet IP under a VPC can only add one mount target;All subnets under the VPC can create mount targets, and you can switch the availability zone where the virtual machine/container is located to obtain the best access performance.
## Example Usage

```hcl
resource "baiducloud_cfs_mount_target" "default" {
  fs_id = "cfs-xxxxxx"
  subnet_id = "sbn-xxxxxxx"
  vpc_id = "vpc-xxxxxxx"
}
```
## Argument Reference

The following arguments are supported:

* `fs_id` - (Required, ForceNew) CFS id which mount target belong to.
* `subnet_id` - (Required, ForceNew) Subnet ID which mount target belong to.
* `vpc_id` - (Required, ForceNew) VPC ID which mount target belong to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `domain` - Domain of mount target.


## Import

CFS mount target can be imported, e.g.

```hcl
$ terraform import baiducloud_cfs_mount_target.default id
```

