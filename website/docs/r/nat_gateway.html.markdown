---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_nat_gateway"
subcategory: "Virtual private Cloud (VPC)"
sidebar_current: "docs-baiducloud-resource-nat_gateway"
description: |-
Provide a resource to create a NAT Gateway.
---

# baiducloud_nat_gateway

Provide a resource to create a NAT Gateway.

## Example Usage

```hcl
resource "baiducloud_eip" "eip1" {
  bandwidth_in_mbps = 1
  payment_timing    = "Postpaid"
  billing_method = "ByBandwidth"
}
resource "baiducloud_eip" "eip2" {
  bandwidth_in_mbps = 1
  payment_timing    = "Postpaid"
  billing_method = "ByBandwidth"
}

resource "baiducloud_nat_gateway" "default" {
  cu_num  = 1
  vpc_id  = "vpc-xxxxxx"
  name    = "test"
  snat_eips = [baiducloud_eip.eip1.eip]
  dnat_eips = [baiducloud_eip.eip2.eip]
  billing = {
    payment_timing = "Postpaid"
  }
}
```

## Argument Reference

The following arguments are supported:

* `billing` - (Required) Billing information of the NAT gateway.
* `name` - (Required) Name of the NAT gateway, consisting of uppercase and lowercase letters、numbers and special characters, such as "-","_","/",".". The value must start with a letter, and the length should between 1-65.
* `vpc_id` - (Required, ForceNew) VPC ID of the NAT gateway.
* `cu_num` - (Optional) cu num.
* `snat_eips` - (Optional) One public network EIP associated with the NAT gateway SNATs or one or more EIPs in the shared bandwidth.
* `dnat_eips` - (Optional) One public network EIP associated with the NAT gateway DNATs or one or more EIPs in the shared bandwidth.
* `spec` - (Optional, ForceNew) Specification of the NAT gateway, available values are small(supports up to 5 public IPs), medium(up to 10 public IPs) and large(up to 15 public IPs).

The `billing` object supports the following:

* `payment_timing` - (Required, ForceNew) Payment timing of the billing, which can be Prepaid or Postpaid. The default is Postpaid.
* `reservation` - (Optional) Reservation of the NAT gateway.

The `reservation` object supports the following:

* `reservation_length` - (Optional, ForceNew) Reservation length that you will pay for your resource. It is valid when payment_timing is Prepaid. Valid values: [1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36].
* `reservation_time_unit` - (Optional) Reservation time unit that you will pay for your resource. It is valid when payment_timing is Prepaid. The value can only be month currently, which is also the default value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `expired_time` - Expired time of the NAT gateway, which will be empty when the payment_timing is Postpaid.
* `status` - Status of the NAT gateway.


## Import

NAT Gateway instance can be imported, e.g.

```hcl
$ terraform import baiducloud_nat_gateway.default nat_gateway_id
```

