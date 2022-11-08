---
layout: "baiducloud"
subcategory: "Cloud File Storage (CFS)"
page_title: "BaiduCloud: baiducloud_cfs"
sidebar_current: "docs-baiducloud-resource-cfs"
description: |-
 Use this resource to create a CFS.
---

# baiducloud_cfs

Use this resource to create a CFS.

## Example Usage

```hcl
resource "baiducloud_cfs" "default" {
  name = "terraform_test"
  zone = "zoneD"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) cfs name, length must be between 1 and 64 bytes
* `protocol` - (Optional, ForceNew) CFS protocol, available value is nfs and smb, default is nfs
* `type` - (Optional, ForceNew) CFS type, default is cap
* `zone` - (Optional, ForceNew) cfs zone

## Attributes Reference

In addition to all arguments above, the following attributes are exported :

* `status` - CFS status, available value is available, updating, paused and unavailable
* `vpc_id` - VPC ID


## Import

CFS can be imported, e.g.

```hcl
$ terraform import baiducloud_cfs.default cfs_id
```

