---
layout: "baiducloud"
subcategory: "Virtual private Cloud (VPC)"
page_title: "BaiduCloud: baiducloud_peer_conn_acceptor"
sidebar_current: "docs-baiducloud-resource-peer_conn_acceptor"
description: |-
  Provide a resource to create a Peer Conn Acceptor.
---

# baiducloud_peer_conn_acceptor

Provide a resource to create a Peer Conn Acceptor.

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

resource "baiducloud_peer_conn_acceptor" "default" {
  peer_conn_id = "${baiducloud_peer_conn.default.id}"
  auto_accept = true
  dns_sync = true
}
```

## Argument Reference

The following arguments are supported:

* `peer_conn_id` - (Required, ForceNew) ID of the peer connection.
* `auto_accept` - (Optional) Whether to accept the peer connection request, default to false.
* `auto_reject` - (Optional) Whether to reject the peer connection request, default to false.
* `dns_sync` - (Optional) Whether to open the switch of dns synchronization.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `bandwidth_in_mbps` - Bandwidth(Mbps) of the peer connection.
* `billing` - Billing information of the peer connection.
  * `payment_timing` - Payment timing of billing, which can be Prepaid or Postpaid. The default is Postpaid.
* `created_time` - Created time of the peer connection.
* `description` - Description of the peer connection.
* `dns_status` - DNS status of the peer connection.
* `expired_time` - Expired time of the peer connection, which will be empty when the payment_timing is Postpaid.
* `local_if_id` - Local interface ID of the peer connection.
* `local_if_name` - Local interface name of the peer connection.
* `local_vpc_id` - Local VPC ID of the peer connection.
* `peer_account_id` - Peer account ID of the peer VPC, which is required only when creating a peer connection across accounts.
* `peer_if_name` - Peer interface name of the peer connection, which is allowed to be set only when the peer connection within this account.
* `peer_region` - Peer region of the peer connection.
* `peer_vpc_id` - Peer VPC ID of the peer connection.
* `role` - Role of the peer connection, which can be initiator or acceptor.
* `status` - Status of the peer connection.


## Import

Peer Conn Acceptor instance can be imported, e.g.

```hcl
$ terraform import baiducloud_peer_conn_acceptor.default peer_conn_id
```

