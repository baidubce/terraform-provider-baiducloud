---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_nat_snat_rules"
sidebar_current: "docs-baiducloud-datasource-nat_snat_rules"
description: |-
  Use this data source to query NAT Gateway SNAT rule list.
---

# baiducloud_nat_snat_rules

Use this data source to query NAT Gateway SNAT rule list.

## Example Usage

```hcl
data "baiducloud_nat_snat_rules" "default" {
 nat_id = "nat-brkztytqzbh0"
}

output "nat_snat_rules" {
 value = "${data.baiducloud_nat_snat_rules.default.nat_snat_rules}"
}
```

## Argument Reference

The following arguments are supported:

* `nat_id` - (Required) ID of the NAT gateway to retrieve.
* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `output_file` - (Optional, ForceNew) Output file for saving result.

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `nat_snat_rules` - The list of NAT Gateway SNAT rules.
  * `rule_id` - ID of the NAT Gateway SNAT rule.
  * `rule_name` - Name of the NAT Gateway SNAT rule.
  * `public_ips_address` - Public network IPs, EIPs associated on the NAT gateway SNAT or IPs in the shared bandwidth.
  * `source_cidr` - Intranet IP/segment.
  * `status` - Status of the NAT Gateway SNAT rule.


