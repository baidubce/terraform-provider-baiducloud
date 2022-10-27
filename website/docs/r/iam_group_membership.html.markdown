---
layout: "baiducloud"
subcategory: "IAM"
page_title: "BaiduCloud: baiducloud_iam_group_membership"
sidebar_current: "docs-baiducloud-resource-iam_group_membership"
description: |-
  Provide a resource to manage IAM Group membership for IAM Users.
---

# baiducloud_iam_group_membership

Provide a resource to manage IAM Group membership for IAM Users.

## Example Usage

```hcl
resource "baiducloud_iam_group" "my-group" {
  name = "my_group_name"
  force_destroy = true
}
resource "baiducloud_iam_user" "my-user" {
  name = "my_user_name"
  force_destroy = true
}
resource "baiducloud_iam_group_membership" "my-group-membership" {
  group = "${baiducloud_iam_group.my-group.name}"
  users = ["${baiducloud_iam_user.my-user.name}"]
}
```

## Argument Reference

The following arguments are supported:

* `group` - (Required, ForceNew) Name of group.
* `users` - (Required) Names of users to add into group.


