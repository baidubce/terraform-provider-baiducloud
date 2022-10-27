---
layout: "baiducloud"
subcategory: "BCC"
page_title: "BaiduCloud: baiducloud_security_group_rules"
sidebar_current: "docs-baiducloud-datasource-security_group_rules"
description: |-
  Use this data source to query Security Group list.
---

# baiducloud_security_group_rules

Use this data source to query Security Group list.

## Example Usage

```hcl
data "baiducloud_security_group_rules" "default" {}

output "security_group_rules" {
 value = "${data.baiducloud_security_group_rules.default.rules}"
}
```

## Argument Reference

The following arguments are supported:

* `security_group_id` - (Required) Security Group ID
* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `instance_id` - (Optional) Security Group attached instance ID
* `output_file` - (Optional, ForceNew) Security Group search result output file
* `vpc_id` - (Optional) Security Group attached vpc id

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `rules` - Security Group rules
  * `dest_group_id` - SecurityGroup rule's destination group id
  * `dest_ip` - SecurityGroup rule's destination ip
  * `direction` - SecurityGroup rule's direction
  * `ether_type` - SecurityGroup rule's ether type
  * `port_range` - SecurityGroup rule's port range
  * `protocol` - SecurityGroup rule's protocol
  * `remark` - SecurityGroup rule's remark
  * `security_group_id` - SecurityGroup rule's security group id
  * `source_group_id` - SecurityGroup rule's source group id
  * `source_ip` - SecurityGroup rule's source ip


