---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_eipgroup"
subcategory: "Elastic IP (EIP)"
sidebar_current: "docs-baiducloud-resource-eipgroup"
description: |-
  Provide a resource to create an EIP GROUP.
---

# baiducloud_eipgroup

Provide a resource to create an EIP GROUP.

## Example Usage

```hcl
resource "baiducloud_eipgroup" "default" {
  name              = "testEIPgroup"
  eip_count         = 2
  bandwidth_in_mbps = 100
  payment_timing    = "Postpaid"
  billing_method    = "ByTraffic"
}
```

## Argument Reference

The following arguments are supported:

* `bandwidth_in_mbps` - (Required) Eip group bandwidth(Mbps), if payment_timing is Prepaid or billing_method is ByBandWidth, support between 1 and 200, if billing_method is ByTraffic, support between 1 and 1000
* `billing_method` - (Required, ForceNew) Eip group billing method, support ByTraffic or ByBandwidth
* `eip_count` - (Required) count of eip group
* `payment_timing` - (Required, ForceNew) Eip group payment timing, support Prepaid and Postpaid
* `name` - (Optional) Eip group name, length must be between 1 and 65 bytes
* `reservation_length` - (Optional, Sensitive) Eip group Prepaid billing reservation length, only useful when payment_timing is Prepaid
* `reservation_time_unit` - (Optional, Sensitive) Eip group Prepaid billing reservation time unit, only useful when payment_timing is Prepaid
* `tags` - (Optional, ForceNew) Tags, do not support modify

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `bw_bandwidth_in_mbps` - Eip group status
* `bw_short_id` - Eip group status
* `continuous` - Eip group continuous
* `create_time` - Eip group create time
* `default_domestic_bandwidth` - Eip group status
* `domestic_bw_bandwidth_in_mbps` - Eip group status
* `domestic_bw_short_id` - Eip group status
* `eips` - Eip list
  * `bandwidth_in_mbps` - Eip bandwidth(Mbps)
  * `billing_method` - Eip billing method
  * `create_time` - Eip create time
  * `eip_instance_type` - Eip instance type
  * `eip` - Eip address
  * `expire_time` - Eip expire time
  * `name` - Eip name
  * `payment_timing` - Eip payment timing
  * `share_group_id` - Eip share group id
  * `status` - Eip status
  * `tags` - Tags
* `expire_time` - Eip group expire time
* `group_id` - id of EIP group
* `idc` - idc of Eip group
* `region` - region of eip group
* `route_type` - Eip Group routeType
* `status` - Eip group status


## Import

EIP group can be imported, e.g.

```hcl
$ terraform import baiducloud_eipgroup.default group_id
```

