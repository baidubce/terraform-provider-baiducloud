---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_route_rules"
sidebar_current: "docs-baiducloud-datasource-route_rules"
description: |-
  Use this data source to query route rule list.
---

# baiducloud_route_rules

Use this data source to query route rule list.

## Example Usage

```hcl
data "baiducloud_route_rules" "default" {
 vpc_id = "vpc-y4p102r3mz6m"
}

output "route_rules" {
 value = "${data.baiducloud_route_rules.default.route_rules}"
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional, ForceNew) Output file for saving result.
* `route_rule_id` - (Optional) ID of the routing rule to be retrieved.
* `route_table_id` - (Optional) Routing table ID for the routing rules to retrieve.
* `vpc_id` - (Optional) VPC ID for the routing rules.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `route_rules` - Result of the routing rules.
  * `description` - Description of the routing rule.
  * `destination_address` - Destination address of the routing rule.
  * `next_hop_id` - Next hop ID of the routing rule.
  * `next_hop_type` - Next hop type of the routing rule.
  * `route_rule_id` - ID of the routing rule.
  * `route_table_id` - Routing table ID of the routing rule.
  * `source_address` - Source address of the routing rule.


