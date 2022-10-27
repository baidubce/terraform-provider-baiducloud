---
layout: "baiducloud"
subcategory: "BLB"
page_title: "BaiduCloud: baiducloud_blb_securitygroup"
sidebar_current: "docs-baiducloud-resource-blb_securitygroup"
description: |-
  Use this resource to get information about a Blb SecurityGroup.
---

# baiducloud_blb_securitygroup

Use this resource to create a Blb SecurityGroup.

~> **NOTE:** The terminate operation of SecurityGroup does NOT take effect immediately，maybe takes for several minites.

## Example Usage

```hcl
resource "baiducloud_blb_securitygroup" "my-server" {
 blb_id = "xxxx"
 security_group_ids = ["xxxxxx"]
}
```

## Argument Reference

The following arguments are supported:

* `blb_id` - (Required, ForceNew) id of the blb
* `security_group_ids` - (Required, ForceNew) ids of the security.

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


## Import

Blb SecurityGroup can be imported, e.g.

```hcl
$ terraform import baiducloud_blb_securitygroup.my-server id
```

