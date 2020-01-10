---
layout: "baiducloud"
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

* `instance_id` - (Optional) Security Group attached instance ID
* `output_file` - (Optional, ForceNew) Security Group search result output file
* `vpc_id` - (Optional) Security Group attached vpc id

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `security_groups` - Security Groups search result
  * `description` - Security Group description
  * `id` - Security Group ID
  * `name` - Security Group name
  * `tags` - Tags
  * `vpc_id` - Security Group vpc id


