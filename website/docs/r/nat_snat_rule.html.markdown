---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_nat_snat_rule"
sidebar_current: "docs-baiducloud-resource-nat_snat_rule"
description: |-
  Provide a resource to create a NAT SNAT rule.
---

# baiducloud_nat_snat_rule

Provide a resource to create a NAT Gateway SNAT rule.

## Example Usage

```hcl
resource "baiducloud_nat_snat_rule" "default" {
  nat_id = "nat-brkztytqzbh0"
  rule_name = "test"
  public_ips_address = ["100.88.14.90"]
  source_cidr = "192.168.1.0/24"
}
```

## Argument Reference

The following arguments are supported:

* `nat_id` - (Required) ID of NAT Gateway.
* `rule_name` - (Required) Rule name, consisting of uppercase and lowercase letters、 numbers and special characters, such as "-","_","/",".". The value must start with a letter, and the length should between 1-65.
* `public_ips_address` - (Required) Public network IPs, EIPs associated on the NAT gateway SNAT or IPs in the shared bandwidth.
* `source_cidr` - (Required) Specification of the NAT gateway, available values are small(supports up to 5 public IPs), medium(up to 10 public IPs) and large(up to 15 public IPs). Default to small.


## Import

NAT Gateway SNAT rule can be imported, e.g.

```hcl
$ terraform import baiducloud_nat_snat_rule.default nat-brkztytqzbh0:rule-wvj6b7v9cvts
```

