---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_acl"
sidebar_current: "docs-baiducloud-resource-acl"
description: |-
  Provide a resource to create an ACL Rule.
---

# baiducloud_acl

Provide a resource to create an ACL Rule.

## Example Usage

```hcl
resource "baiducloud_acl" "default" {
  subnet_id = "sbn-86c3v6pnt8b4"
  protocol = "tcp"
  source_ip_address = "192.168.0.0/24"
  destination_ip_address = "192.168.1.0/24"
  source_port = "8888"
  destination_port = "9999"
  position = 20
  direction = "ingress"
  action = "allow"
}
```

## Argument Reference

The following arguments are supported:

* `action` - (Required) Action of the acl. Valid values are allow and deny.
* `destination_ip_address` - (Required) Destination ip address of the acl.
* `destination_port` - (Required) Destination port of the acl.
* `direction` - (Required, ForceNew) Direction of the acl. Valid values are ingress and egress, respectively indicating the inbound of the rule and the outbound rule.
* `position` - (Required) Position of the acl, representing the priority of the acl rule. The value should be 1-5000 and cannot be duplicated with existing entries. The smaller the value, the higher the priority, and the rule matching order is to match the priority from high to low.
* `protocol` - (Required) Protocol of the acl, available values are all, tcp, udp and icmp.
* `source_ip_address` - (Required) Source ip address of the acl.
* `source_port` - (Required) Source port of the acl.
* `subnet_id` - (Required, ForceNew) Subnet ID of the acl.
* `description` - (Optional) Description of the acl.


