---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_iam_group_policy_attachment"
sidebar_current: "docs-baiducloud-resource-iam_group_policy_attachment"
description: |-
  Provide a resource to attach an IAM Policy to IAM Group.
---

# baiducloud_iam_group_policy_attachment

Provide a resource to attach an IAM Policy to IAM Group.

## Example Usage

```hcl
resource "baiducloud_iam_group" "my-group" {
  name = "my_group_name"
  force_destroy    = true
}
resource "baiducloud_iam_policy" "my-policy" {
   name = "my_policy"
  document = <<EOF
{"accessControlList": [{"region":"bj","service":"bcc","resource":["*"],"permission":["*"],"effect":"Allow"}]}
  EOF
}
resource "baiducloud_iam_group_policy_attachment" "my-group-policy-attachment" {
  group = "${baiducloud_iam_group.my-group.name}"
  policy = "${baiducloud_iam_policy.my-policy.name}"
}
```

## Argument Reference

The following arguments are supported:

* `group` - (Required, ForceNew) Name of group.
* `policy` - (Required, ForceNew) Name of policy.
* `policy_type` - (Optional, ForceNew) Type of policy, valid values are Custom/System.


