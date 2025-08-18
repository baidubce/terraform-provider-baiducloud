---
layout: "baiducloud"
subcategory: "Baidu Cloud Compute (BCC)"
page_title: "BaiduCloud: baiducloud_security_group_rule"
sidebar_current: "docs-baiducloud-resource-security_group_rule"
description: |-
  Provide a resource to create a security group rule.
---

# baiducloud_security_group_rule

Provide a resource to create a security group rule.

~> **NOTE:** When a rule with the same combination of protocol, direction, port range, and either source/destination IP or source/destination group ID already exists in the security group, 
creating such a rule will update the existing one instead of adding a new entry.
This behavior is specific to the Terraform provider and ensures resource uniqueness. 
It may differ from how rules appear or behave when managed directly in the cloud provider’s console.

## Example Usage

```hcl
resource "baiducloud_security_group" "default" {
  name = "my-sg"
  description = "default"
}

resource "baiducloud_security_group_rule" "default" {
  security_group_id = "${baiducloud_security_group.default.id}"
  remark            = "remark"
  protocol          = "udp"
  port_range        = "1-65523"
  direction         = "ingress"
}
```

## Argument Reference

The following arguments are supported:

* `direction` - (Required, ForceNew) SecurityGroup rule's direction, support ingress/egress
* `security_group_id` - (Required, ForceNew) SecurityGroup rule's security group id
* `dest_group_id` - (Optional, ForceNew) SecurityGroup rule's destination group id, dest_group_id and dest_ip can not set in the same time
* `dest_ip` - (Optional, ForceNew) SecurityGroup rule's destination ip, dest_group_id and dest_ip can not set in the same time
* `ether_type` - (Optional, ForceNew) SecurityGroup rule's ether type, support IPv4/IPv6
* `port_range` - (Optional, ForceNew) SecurityGroup rule's port range, you can set single port like 80, or set a port range, like 1-65535, default 1-65535. If protocol is all, only support 1-65535
* `protocol` - (Optional, ForceNew) SecurityGroup rule's protocol, support tcp/udp/icmp/all, default all
* `remark` - (Optional, ForceNew) SecurityGroup rule's remark
* `source_group_id` - (Optional, ForceNew) SecurityGroup rule's source group id, source_group_id and source_ip can not set in the same time
* `source_ip` - (Optional, ForceNew) SecurityGroup rule's source ip, source_group_id and source_ip can not set in the same time


