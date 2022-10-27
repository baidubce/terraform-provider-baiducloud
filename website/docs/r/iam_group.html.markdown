---
layout: "baiducloud"
subcategory: "IAM"
page_title: "BaiduCloud: baiducloud_iam_group"
sidebar_current: "docs-baiducloud-resource-iam_group"
description: |-
  Provide a resource to manage an IAM group.
---

# baiducloud_iam_group

Provide a resource to manage an IAM group.

## Example Usage

```hcl
resource "baiducloud_iam_group" "my-group" {
  name = "my_group_name"
  description = "group description"
  force_destroy    = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of group.
* `description` - (Optional) Description of the group.
* `force_destroy` - (Optional) Delete group and its related user memberships and policy attachments.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `unique_id` - Unique ID of group.


