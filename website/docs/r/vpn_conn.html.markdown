---
layout: "baiducloud"
subcategory: "VPN"
page_title: "BaiduCloud: baiducloud_vpn_conn"
sidebar_current: "docs-baiducloud-resource-vpn_conn"
description: |-
  Provide a resource to create a VPN conn.
---

# baiducloud_vpn_conn

Provide a resource to create a VPN conn.

## Example Usage

```hcl
resource "baiducloud_vpn_conn" "default" {
  vpn_id = baiducloud_vpn_gateway.default.id
  secret_key = "ddd22@www"
  local_subnets = ["192.168.0.0/20"]
  remote_ip = "11.11.11.133"
  remote_subnets = ["192.168.100.0/24"]
  description = "just for test"
  vpn_conn_name = "vpn-conn"
  ike_config {
    ike_version = "v1"
    ike_mode = "main"
    ike_enc_alg = "aes"
    ike_auth_alg = "sha1"
    ike_pfs = "group2"
    ike_life_time = 300
  }
  ipsec_config {
    ipsec_enc_alg = "aes"
    ipsec_auth_alg = "sha1"
    ipsec_pfs = "group2"
    ipsec_life_time = 200
  }
}
```

## Argument Reference

The following arguments are supported:

* `local_subnets` - (Required) Local network cidr list.
* `remote_ip` - (Required) Public IP of the peer VPN gateway.
* `remote_subnets` - (Required) Peer network cidr list.
* `secret_key` - (Required) Shared secret key, 8 to 17 characters, English, numbers and symbols must exist at the same time, the symbols are limited to !@#$%^*()_.
* `vpn_id` - (Required) VPN id which vpn conn belong to.
* `description` - (Optional) Description of VPN conn.
* `ike_config` - (Optional) IKE config.
* `ipsec_config` - (Optional) Ipsec config details.
* `vpn_conn_name` - (Optional) Name of vpn conn.

The `ike_config` object supports the following:

* `ike_auth_alg` - (Required) IKE Authenticate Algorithm
* `ike_enc_alg` - (Required) IKE Encryption Algorithm.
* `ike_life_time` - (Required) IKE life time.
* `ike_mode` - (Required) Negotiation mode.
* `ike_pfs` - (Required) Diffie-Hellman key exchange algorithm.
* `ike_version` - (Required) Version of IKE.

The `ipsec_config` object supports the following:

* `ipsec_auth_alg` - (Required) Ipsec Authenticate Algorithm.
* `ipsec_enc_alg` - (Required) Ipsec Encryption Algorithm.
* `ipsec_life_time` - (Required) Ipsec life time.
* `ipsec_pfs` - (Required) Diffie-Hellman key exchange algorithm.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `local_ip` - Public IP of the VPN gateway.


## Import

VPN conn can be imported, e.g.

```hcl
$ terraform import baiducloud_vpn_conn.default vpn_conn_id
```

