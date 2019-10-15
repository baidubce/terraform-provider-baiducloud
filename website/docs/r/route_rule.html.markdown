---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_route_rule"
sidebar_current: "docs-baiducloud-resource-route_rule"
description: |-
  Provides a resource to create a VPC routing rule.
---

# baiducloud_route_rule

Provides a resource to create a VPC routing rule.

## Example Usage

```hcl
resource "baiducloud_route_rule" "default" {
  route_table_id = "rt-as4npcsp2hve"
  source_address = "192.168.0.0/24"
  destination_address = "192.168.1.0/24"
  next_hop_id = "i-BtXnDM6y"
  next_hop_type = "custom"
  description = "baiducloud route rule created by terraform"
}
```

## Argument Reference

The following arguments are supported:

* `destination_address` - (Required, ForceNew) Destination CIDR block of the routing rule. The network segment can be 0.0.0.0/0, otherwise, the destination address cannot overlap with this VPC CIDR block(except when the destination network segment or the VPC CIDR is 0.0.0.0/0).
* `next_hop_type` - (Required, ForceNew) Type of the next hop, available values are custom、vpn and nat.
* `route_table_id` - (Required, ForceNew) ID of the routing table.
* `source_address` - (Required, ForceNew) Source CIDR block of the routing rule. The value can be all network segments 0.0.0.0/0, existing subnet segments in the VPC, or the network segment within the subnet.
* `description` - (Optional, ForceNew) Description of the routing rule.
* `next_hop_id` - (Optional, ForceNew) ID of the next hop.


