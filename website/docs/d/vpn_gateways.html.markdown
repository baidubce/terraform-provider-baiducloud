---
layout: "baiducloud"
subcategory: "VPN"
page_title: "BaiduCloud: baiducloud_vpn_gateways"
sidebar_current: "docs-baiducloud-datasource-vpn_gateways"
description: |-
  Use this data source to query VPN gateway list.
---

# baiducloud_vpn_gateways

Use this data source to query VPN gateway list.

## Example Usage

```hcl
data "baiducloud_vpn_gateways" "default" {
  vpc_id = "vpc-65cz3hu92kz2"
}

output "vpns" {
  value = "${data.baiducloud_vpn_gateways.default.vpns}"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required) ID of the VPC which vpn gateway belong to.
* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `output_file` - (Optional, ForceNew) Output file for saving result.

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `vpn_gateways` - Result of VPCs.
  * `band_width_in_mbps` - Eip bandwidth(Mbps), if payment_timing is Prepaid or billing_method is ByBandWidth, support between 1 and 200, if billing_method is ByTraffic, support between 1 and 1000
  * `create_time` - Create time of VPN gateway.
  * `description` - Description of the VPN.
  * `eip` - Eip address.
  * `expired_time` - Expired time.
  * `max_connection` - Max connection of VPN gateway.
  * `payment_timing` - Payment timing of billing, which can be Prepaid or Postpaid. The default is Postpaid.
  * `status` - Status of the VPN.
  * `vpc_id` - ID of the VPC which vpn gateway belong to.
  * `vpn_conn_num` - Number of VPN tunnels.
  * `vpn_conns` - ID List of the VPN gateway tunnels.
  * `vpn_id` - ID of the VPN gateway.
  * `vpn_name` - Name of the VPN gateway.


