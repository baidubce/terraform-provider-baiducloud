---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_iam_user"
sidebar_current: "docs-baiducloud-resource-iam_user"
description: |-
  Provide a resource to manage an IAM user.
---

# baiducloud_iam_user

Provide a resource to manage an IAM user.

## Example Usage

```hcl
resource "baiducloud_iam_user" "my-user" {
  name = "my_user_name"
  description = "user description"
  force_destroy    = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of user.
* `description` - (Optional) Description of the user.
* `force_destroy` - (Optional) Delete user and its related access keys, group memberships and policy attachments.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `unique_id` - Unique ID of user.


