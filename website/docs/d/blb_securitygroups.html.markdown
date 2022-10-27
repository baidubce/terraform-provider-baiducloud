---
layout: "baiducloud"
subcategory: "BLB"
page_title: "BaiduCloud: baiducloud_blb_securitygroups"
sidebar_current: "docs-baiducloud-datasource-blb_securitygroups"
description: |-
  Use this data source to query blb SecurityGroups.
---

# baiducloud_blb_securitygroups

Use this data source to query blb SecurityGroups.

## Example Usage

```hcl
data "baiducloud_blb_securitygroups" "default" {
   blb_id = "lb-0d29axxx6"
}

output "security_groups" {
   value = "${data.baiducloud_blb_securitygroups.default.bind_security_groups}"
}
```

## Argument Reference

The following arguments are supported:

* `blb_id` - (Required, ForceNew) id of the blb
* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `output_file` - (Optional, ForceNew) blb securitygroup search result output file

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `bind_security_groups` - blb bind security_groups
  * `security_group_desc` - desc of security group
  * `security_group_id` - bind security id
  * `security_group_name` - name of security group
  * `security_group_rules` - rules of security groups
    * `dest_group_id` - dest group id
    * `dest_ip` - dest ip
    * `direction` - direction
    * `ethertype` - ethertype
    * `port_range` - portRange
    * `protocol` - protocol
    * `security_group_rule_id` - id of security group rule
  * `vpc_name` - name of vpc


