---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_nat_gateways"
sidebar_current: "docs-baiducloud-datasource-nat_gateways"
description: |-
  Use this data source to query NAT gateway list.
---

# baiducloud_nat_gateways

Use this data source to query NAT gateway list.

## Example Usage

```hcl
data "baiducloud_nat_gateways" "default" {
 vpc_id = "vpc-y4p102r3mz6m"
}

output "nat_gateways" {
 value = "${data.baiducloud_nat_gateways.default.nat_gateways}"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `ip` - (Optional) Specify the EIP binded by the NAT gateway to retrieve.
* `name` - (Optional) Name of the NAT gateway.
* `nat_id` - (Optional) ID of the NAT gateway to retrieve.
* `output_file` - (Optional, ForceNew) Output file for saving result.
* `vpc_id` - (Optional) VPC ID where the NAT gateways located.

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `nat_gateways` - The list of NAT gateways.
  * `eips` - EIP list of the NAT gateway.
  * `expired_time` - Expired time of the NAT gateway.
  * `id` - ID of the NAT gateway.
  * `name` - Name of the NAT gateway.
  * `payment_timing` - Payment timing of the NAT gateway.
  * `spec` - Spec of the NAT gateway.
  * `status` - Status of the NAT gateway.
  * `vpc_id` - VPC ID of the NAT gateway.


