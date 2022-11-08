---
layout: "baiducloud"
subcategory: "Cloud File Storage (CFS)"
page_title: "BaiduCloud: baiducloud_cfss"
sidebar_current: "docs-baiducloud-datasource-cfss"
description: |-
  Use this data source to get CFS list.
---

# baiducloud_cfss

Use this data source to get CFS list.

## Example Usage

```hcl
data baiducloud_cfss "default" {

}

output "cfss" {
 value = "${data.baiducloud_cfss.default}"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `output_file` - (Optional, ForceNew) CFS search result output file

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cfss` - cfs list
  * `fs_id` - ID of the CFS.
  * `mount_target_list` - Name of the deployset.
    * `domain` - Domain of the mount target.
    * `mount_id` - ID of the mount target.
    * `subnet_id` - ID of subnet which mount target bind.
  * `name` - Name of the CFS.
  * `protocol` - CFS protocol, available value is nfs and smb, default is nfs.
  * `status` - CFS status, available value is available,updating,paused and unavailable.
  * `type` - CFS type, default is cap.
  * `vpc_id` - VPC ID.


