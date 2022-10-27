---
layout: "baiducloud"
subcategory: "BCC"
page_title: "BaiduCloud: baiducloud_security_groups"
sidebar_current: "docs-baiducloud-datasource-security_groups"
description: |-
  Use this data source to query Security Group list.
---

# baiducloud_security_groups

Use this data source to query Security Group list.

## Example Usage

```hcl
data "baiducloud_security_groups" "default" {}

output "security_groups" {
 value = "${data.baiducloud_security_groups.default.security_groups}"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `instance_id` - (Optional) Security Group attached instance ID
* `output_file` - (Optional, ForceNew) Security Group search result output file
* `vpc_id` - (Optional) Security Group attached vpc id

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `security_groups` - Security Groups search result
  * `description` - Security Group description
  * `id` - Security Group ID
  * `name` - Security Group name
  * `tags` - Tags
  * `vpc_id` - Security Group vpc id


