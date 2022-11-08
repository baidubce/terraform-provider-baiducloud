---
layout: "baiducloud"
subcategory: "Virtual private Cloud (VPC)"
page_title: "BaiduCloud: baiducloud_nat_gateway"
sidebar_current: "docs-baiducloud-resource-nat_gateway"
description: |-
  Provide a resource to create a NAT Gateway.
---

# baiducloud_nat_gateway

Provide a resource to create a NAT Gateway.

## Example Usage

```hcl
resource "baiducloud_nat_gateway" "default" {
  name = "terraform-nat-gateway"
  vpc_id = "vpc-ggm7drdgyvha"
  spec = "medium"
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
* `spec` - (Optional, ForceNew) Specification of the NAT gateway, available values are small(supports up to 5 public IPs), medium(up to 10 public IPs) and large(up to 15 public IPs). Default to small.
* `cu_num` - (Optional) Number of NAT gateway CU, max is 100.

The `billing` object supports the following:

* `payment_timing` - (Required, ForceNew) Payment timing of the billing, which can be Prepaid or Postpaid. The default is Postpaid.
* `reservation` - (Optional) Reservation of the NAT gateway.

The `reservation` object supports the following:

* `reservation_length` - (Optional, ForceNew) Reservation length that you will pay for your resource. It is valid when payment_timing is Prepaid. Valid values: [1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36].
* `reservation_time_unit` - (Optional) Reservation time unit that you will pay for your resource. It is valid when payment_timing is Prepaid. The value can only be month currently, which is also the default value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `eips` - One public network EIP associated with the NAT gateway or one or more EIPs in the shared bandwidth.
* `expired_time` - Expired time of the NAT gateway, which will be empty when the payment_timing is Postpaid.
* `status` - Status of the NAT gateway.


## Import

NAT Gateway instance can be imported, e.g.

```hcl
$ terraform import baiducloud_nat_gateway.default nat_gateway_id
```

