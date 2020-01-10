---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_eip"
sidebar_current: "docs-baiducloud-resource-eip"
description: |-
  Provide a resource to create an EIP.
---

# baiducloud_eip

Provide a resource to create an EIP.

## Example Usage

```hcl
resource "baiducloud_eip" "default" {
  name              = "testEIP"
  bandwidth_in_mbps = 100
  payment_timing    = "Postpaid"
  billing_method    = "ByTraffic"
}
```

## Argument Reference

The following arguments are supported:

* `bandwidth_in_mbps` - (Required) Eip bandwidth(Mbps), if payment_timing is Prepaid or billing_method is ByBandWidth, support between 1 and 200, if billing_method is ByTraffic, support between 1 and 1000
* `billing_method` - (Required, ForceNew) Eip billing method, support ByTraffic or ByBandwidth
* `payment_timing` - (Required, ForceNew) Eip payment timing, support Prepaid and Postpaid
* `name` - (Optional, ForceNew) Eip name, length must be between 1 and 65 bytes
* `reservation_length` - (Optional) Eip Prepaid billing reservation length, only useful when payment_timing is Prepaid
* `reservation_time_unit` - (Optional) Eip Prepaid billing reservation time unit, only useful when payment_timing is Prepaid
* `tags` - (Optional, ForceNew) Tags, do not support modify

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - Eip create time
* `eip_instance_type` - Eip instance type
* `eip` - Eip address
* `expire_time` - Eip expire time
* `share_group_id` - Eip share group id
* `status` - Eip status


## Import

EIP can be imported, e.g.

```hcl
$ terraform import baiducloud_eip.default eip
```

