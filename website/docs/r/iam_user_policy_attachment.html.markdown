---
layout: "baiducloud"
subcategory: "Identity and Access Management (IAM)"
page_title: "BaiduCloud: baiducloud_iam_user_policy_attachment"
sidebar_current: "docs-baiducloud-resource-iam_user_policy_attachment"
description: |-
  Provide a resource to attach an IAM Policy to IAM User.
---

# baiducloud_iam_user_policy_attachment

Provide a resource to attach an IAM Policy to IAM User.

## Example Usage

```hcl
resource "baiducloud_iam_user" "my-user" {
  name = "my_user_name"
  force_destroy    = true
}
resource "baiducloud_iam_policy" "my-policy" {
  name = "my_policy"
  document = <<EOF
{"accessControlList": [{"region":"bj","service":"bcc","resource":["*"],"permission":["*"],"effect":"Allow"}]}
  EOF
}
resource "baiducloud_iam_user_policy_attachment" "my-user-policy-attachment" {
  user = "${baiducloud_iam_user.my-user.name}"
  policy = "${baiducloud_iam_policy.my-policy.name}"
}
```

## Argument Reference

The following arguments are supported:

* `policy` - (Required, ForceNew) Name of policy.
* `user` - (Required, ForceNew) Name of user.
* `policy_type` - (Optional, ForceNew) Type of policy, valid values are Custom/System.


