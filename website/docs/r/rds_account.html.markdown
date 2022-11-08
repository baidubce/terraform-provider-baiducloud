---
layout: "baiducloud"
subcategory: "Relational Database Service (RDS)"
page_title: "BaiduCloud: baiducloud_rds_account"
sidebar_current: "docs-baiducloud-resource-rds_account"
description: |-
  Use this resource to get information about a RDS Account.
---

# baiducloud_rds_account

Use this resource to get information about a RDS Account.

## Example Usage

```hcl
resource "baiducloud_rds_account" "default" {
}
```

## Argument Reference

The following arguments are supported:

* `account_name` - (Required, ForceNew) Account name.
* `instance_id` - (Required, ForceNew) ID of the rds instance.
* `password` - (Required, ForceNew) Operation password.
* `account_type` - (Optional, ForceNew) Type of the Account, Available values are Common、Super. The default is Common
* `desc` - (Optional, ForceNew) description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `status` - Status of the Account.


## Import

RDS Account can be imported, e.g.

```hcl
$ terraform import baiducloud_rds_account.default id
```

