---
layout: "baiducloud"
subcategory: "Virtual private Cloud (VPC)"
page_title: "BaiduCloud: baiducloud_peer_conns"
sidebar_current: "docs-baiducloud-datasource-peer_conns"
description: |-
  Use this data source to query Peer Conn list.
---

# baiducloud_peer_conns

Use this data source to query Peer Conn list.

## Example Usage

```hcl
data "baiducloud_peer_conns" "default" {
  vpc_id = "vpc-y4p102r3mz6m"
}

output "peer_conns" {
  value = "${data.baiducloud_peer_conns.default.peer_conns}"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `output_file` - (Optional, ForceNew) Output file for saving result.
* `peer_conn_id` - (Optional) ID of the peer connection to retrieve.
* `vpc_id` - (Optional) VPC ID where the peer connections located.

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `peer_conns` - The list of the peer connections.
  * `bandwidth_in_mbps` - Bandwidth(Mbps) of the peer connection.
  * `created_time` - Created time of the peer connection.
  * `description` - Description of the peer connection.
  * `dns_status` - DNS status of the peer connection.
  * `expired_time` - Expired time of the peer connection.
  * `local_if_id` - Local interface ID of the peer connection.
  * `local_if_name` - Local interface name of the peer connection.
  * `local_region` - Local region of the peer connection.
  * `local_vpc_id` - Local VPC ID of the peer connection.
  * `payment_timing` - Payment timing of the peer connection.
  * `peer_account_id` - Peer account ID of the peer connection.
  * `peer_conn_id` - ID of the peer connection.
  * `peer_region` - Peer region of the peer connection.
  * `peer_vpc_id` - Peer VPC ID of the peer connection.
  * `role` - Role of the peer connection.
  * `status` - Status of the peer connection.


