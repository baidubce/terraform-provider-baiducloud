---
layout: "baiducloud"
subcategory: "VPC"
page_title: "BaiduCloud: baiducloud_peer_conn"
sidebar_current: "docs-baiducloud-resource-peer_conn"
description: |-
  Provide a resource to create a Peer Conn.
---

# baiducloud_peer_conn

Provide a resource to create a Peer Conn.

## Example Usage

```hcl
resource "baiducloud_peer_conn" "default" {
  bandwidth_in_mbps = 10
  local_vpc_id = "vpc-y4p102r3mz6m"
  peer_vpc_id = "vpc-4njbqurm0uag"
  peer_region = "bj"
  billing = {
    payment_timing = "Postpaid"
  }
}
```

## Argument Reference

The following arguments are supported:

* `bandwidth_in_mbps` - (Required) Bandwidth(Mbps) of the peer connection.
* `billing` - (Required) Billing information of the peer connection.
* `local_vpc_id` - (Required, ForceNew) Local VPC ID of the peer connection.
* `peer_region` - (Required, ForceNew) Peer region of the peer connection.
* `peer_vpc_id` - (Required, ForceNew) Peer VPC ID of the peer connection.
* `description` - (Optional) Description of the peer connection.
* `dns_sync` - (Optional) Whether to open the switch of dns synchronization.
* `local_if_name` - (Optional) Local interface name of the peer connection.
* `peer_account_id` - (Optional, ForceNew) Peer account ID of the peer VPC, which is required only when creating a peer connection across accounts.
* `peer_if_name` - (Optional) Peer interface name of the peer connection, which is allowed to be set only when the peer connection within this account.

The `billing` object supports the following:

* `payment_timing` - (Required, ForceNew) Payment timing of the billing, which can be Prepaid or Postpaid. The default is Postpaid.
* `reservation` - (Optional) Reservation of the peer connection.

The `reservation` object supports the following:

* `reservation_length` - (Optional, ForceNew) Reservation length that you will pay for your resource. It is valid when payment_timing is Prepaid. Valid values: [1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36].
* `reservation_time_unit` - (Optional) Reservation time unit that you will pay for your resource. It is valid when payment_timing is Prepaid. The value can only be month currently, which is also the default value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `created_time` - Created time of the peer connection.
* `dns_status` - DNS status of the peer connection.
* `expired_time` - Expired time of the peer connection, which will be empty when the payment_timing is Postpaid.
* `local_if_id` - Local interface ID of the peer connection.
* `role` - Role of the peer connection, which can be initiator or acceptor.
* `status` - Status of the peer connection.


## Import

Peer Conn instance can be imported, e.g.

```hcl
$ terraform import baiducloud_peer_conn.default peer_conn_id
```

