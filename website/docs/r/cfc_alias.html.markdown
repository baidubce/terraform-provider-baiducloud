---
layout: "baiducloud"
subcategory: "CFC"
page_title: "BaiduCloud: baiducloud_cfc_alias"
sidebar_current: "docs-baiducloud-resource-cfc_alias"
description: |-
  Provide a resource to create a CFC Function Alias.
---

# baiducloud_cfc_alias

Provide a resource to create a CFC Function Alias.

## Example Usage

```hcl
resource "baiducloud_cfc_alias" "default" {
  function_name    = "terraform-cfc"
  function_version = "$LATEST"
  alias_name       = "terraformAlias"
  description      = "terraform create alias"
}
```

```

## Argument Reference

The following arguments are supported:

* `alias_name` - (Required, ForceNew) CFC Function alias name
* `function_name` - (Required, ForceNew) CFC Function name
* `function_version` - (Required) CFC Function version this alias binded
* `description` - (Optional) CFC Function alias description

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `alias_arn` - CFC Function alias arn
* `alias_brn` - CFC Function alias brn
* `create_time` - CFC Function alias create time
* `uid` - CFC Function uid
* `update_time` - CFC Function alias update time


