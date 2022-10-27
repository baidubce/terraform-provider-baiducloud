---
layout: "baiducloud"
subcategory: "IAM"
page_title: "BaiduCloud: baiducloud_iam_policy"
sidebar_current: "docs-baiducloud-resource-iam_policy"
description: |-
  Provide a resource to manage an IAM Policy.
---

# baiducloud_iam_policy

Provide a resource to manage an IAM Policy.

## Example Usage

```hcl
resource "baiducloud_iam_policy" "my-policy" {
  name = "my_policy"
  description = "my description"
  document = <<EOF
{"accessControlList": [{"region":"bj","service":"bcc","resource":["*"],"permission":["*"],"effect":"Allow"}]}
  EOF
}
```

## Argument Reference

The following arguments are supported:

* `document` - (Required, ForceNew) Json serialized ACL string.
* `name` - (Required, ForceNew) Name of policy.
* `description` - (Optional, ForceNew) Description of the policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `unique_id` - Unique ID of policy.


