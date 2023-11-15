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
  description = "created by terraform"
}
```

If you want to create a `multi-path route`, you can use the following configuration:
```hcl
resource "baiducloud_route_rule" "default" {
  route_table_id = "rt-y97dkswd5hac"
  source_address = "10.0.0.0/24"
  destination_address = "10.0.0.0/16"
  next_hop_list {
      next_hop_id = "dcgw-7i066wq232"
      next_hop_type = "dcGateway"
      path_type = "ecmp"
  }
  next_hop_list {
    next_hop_id = "dcgw-5xtd6y233"
    next_hop_type = "dcGateway"
    path_type = "ecmp"
  }
  description = "created by terraform"
}
```

## Argument Reference

The following arguments are supported:

* `destination_address` - (Required, ForceNew) Destination CIDR block of the routing rule. The network segment can be 0.0.0.0/0, otherwise, the destination address cannot overlap with this VPC CIDR block(except when the destination network segment or the VPC CIDR is 0.0.0.0/0).
* `route_table_id` - (Required, ForceNew) ID of the routing table.
* `source_address` - (Required, ForceNew) Source CIDR block of the routing rule. The value can be all network segments 0.0.0.0/0, existing subnet segments in the VPC, or the network segment within the subnet.
* `description` - (Optional, ForceNew) Description of the routing rule.
* `next_hop_id` - (Optional, ForceNew) Next-hop ID, this field must be filled when creating a single path route.
* `next_hop_list` - (Optional) Create a multi-path route based on the next hop information. This field is required when creating a `multi-path` route.
* `next_hop_type` - (Optional, ForceNew) Type of the next hop, available values are custom, vpn, nat and dcGateway.This field is required when creating a `single-path` route.

The `next_hop_list` object supports the following:

* `next_hop_id` - (Required, ForceNew) Next-hop ID.
* `next_hop_type` - (Required, ForceNew) Routing type. Currently only the dedicated gateway type dcGateway is supported.
* `path_type` - (Required, ForceNew) Multi-line mode. The load balancing value is `ecmp`; the main backup mode value is `ha:active`, `ha:standby`, which represent the main and backup routes respectively.


