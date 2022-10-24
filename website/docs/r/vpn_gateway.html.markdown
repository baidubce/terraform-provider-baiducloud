---
layout: "baiducloud"
subcategory: "VPN"
page_title: "BaiduCloud: baiducloud_vpn_gateway"
sidebar_current: "docs-baiducloud-resource-vpn_gateway"
description: |-
  Provide a resource to create a VPN gateway.
---

# baiducloud_vpn_gateway

Provide a resource to create a VPN gateway.

## Example Usage

```hcl
resource "baiducloud_vpn_gateway" "default" {
  vpn_name       = "test_vpn_gateway"
  vpc_id         = "vpc-65cz3hu92kz2"
  description    = "test desc"
  payment_timing = "Postpaid"
}
```

## Argument Reference

The following arguments are supported:

* `payment_timing` - (Required) Payment timing of billing, which can be Prepaid or Postpaid. The default is Postpaid.
* `vpc_id` - (Required) ID of the VPC which vpn gateway belong to.
* `vpn_name` - (Required) Name of the VPN gateway, which cannot take the value "default", the length is no more than 65 characters, and the value can be composed of numbers, characters and underscores.
* `description` - (Optional) Description of the VPN. The value is no more than 200 characters.
* `eip` - (Optional) Eip address.
* `reservation` - (Optional) Reservation of the instance.

The `reservation` object supports the following:

* `reservation_length` - (Required) The reservation length that you will pay for your resource. It is valid when payment_timing is Prepaid. Valid values: [1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36].
* `reservation_time_unit` - (Required) The reservation time unit that you will pay for your resource. It is valid when payment_timing is Prepaid. The value can only be month currently, which is also the default value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `bandwidth_in_mbps` - Eip bandwidth(Mbps), if payment_timing is Prepaid or billing_method is ByBandWidth, support between 1 and 200, if billing_method is ByTraffic, support between 1 and 1000
* `expired_time` - Expired time of the VPN gateway.
* `status` - VPN gateway status.
* `vpn_conn_num` - Number of VPN tunnels.
* `vpn_conns` - ID list of VPN tunnels.


## Import

VPN gateway can be imported, e.g.

```hcl
$ terraform import baiducloud_vpn_gateway.default vpn_gateway_id
```

