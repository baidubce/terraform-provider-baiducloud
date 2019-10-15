---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_security_group"
sidebar_current: "docs-baiducloud-resource-security_group"
description: |-
  Provide a resource to create a security group.
---

# baiducloud_security_group

Provide a resource to create a security group.

## Example Usage

```hcl
resource "baiducloud_security_group" "default" {
  name        = "testSecurityGroup"
  description = "default"
  tags {
    tag_key   = "testKey"
    tag_value = "testValue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) SecurityGroup name
* `description` - (Optional, ForceNew) SecurityGroup description
* `tags` - (Optional, ForceNew) Tags
* `vpc_id` - (Optional, ForceNew) SecurityGroup binded VPC id

The `tags` object supports the following:

* `tag_key` - (Required) Tag's key
* `tag_value` - (Required) Tag's value


## Import

Bcc SecurityGroup can be imported, e.g.

```hcl
$ terraform import baiducloud_security_group.default security_group_id
```

