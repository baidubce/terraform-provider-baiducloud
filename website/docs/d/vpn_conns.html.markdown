---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_vpn_conns"
sidebar_current: "docs-baiducloud-datasource-vpn_conns"
description: |-
  Use this data source to query vpn conn list.
---

# baiducloud_vpn_conns

Use this data source to query vpn conn list.

## Example Usage

```hcl
data "baiducloud_vpn_conns" "default" {
    vpn_id = "vpn-xxxxxxx"
}

output "conns" {
  value = "${data.baiducloud_vpn_conns.default.conns}"
}
```

## Argument Reference

The following arguments are supported:

* `vpn_id` - (Required) VPN id which vpn conn belong to.
* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `output_file` - (Optional, ForceNew) Output file for saving result.

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `vpn_conns` - Result of VPN conns.
  * `created_time` - Create time of VPN conn.
  * `health_status` - Health status of the vpn conn.
  * `local_ip` - Public IP of the VPN gateway.
  * `local_subnets` - Local network cidr list.
  * `remote_ip` - Public IP of the peer VPN gateway.
  * `remote_subnets` - Peer network cidr list.
  * `secret_key` - Shared secret key, 8 to 17 characters, English, numbers and symbols must exist at the same time, the symbols are limited to !@#$%^*()_.
  * `status` - Status of the vpn conn.
  * `vpn_conn_id` - ID of the VPN conn.
  * `vpn_conn_name` - Name of vpn conn.
  * `vpn_id` - VPN id which vpn conn belong to.


