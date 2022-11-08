---
layout: "baiducloud"
subcategory: "Cloud File Storage (CFS)"
page_title: "BaiduCloud: baiducloud_cfs_mount_targets"
sidebar_current: "docs-baiducloud-datasource-cfs_mount_targets"
description: |-
  Use this data source to get CFS mount target list.
---

# baiducloud_cfs_mount_targets

Use this data source to get CFS mount target list.

## Example Usage

```hcl
data "baiducloud_cfs_mount_targets" "default" {
  fs_id = "cfs-xxxxxxxxxxx"
}

output "cfss" {
 value = "${baiducloud_cfs_mount_targets.default}"
}
```

## Argument Reference

The following arguments are supported:

* `fs_id` - (Required, ForceNew) CFS ID which you want query
* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `output_file` - (Optional, ForceNew) CFS search result output file

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `mount_targets` - Mount targets info list.
  * `domain` - Domain of mount target.
  * `mount_id` - ID of the mount target
  * `subnet_id` - Subnet ID which mount target belong to.


