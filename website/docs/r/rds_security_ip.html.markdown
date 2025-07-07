---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_rds_security_ip"
subcategory: "Relational Database Service (RDS)"
sidebar_current: "docs-baiducloud-resource-rds_security_ip"
description: |-
  Use this resource to get information about a RDS Security Ip.
---

# baiducloud_rds_security_ip

Use this resource to get information about a RDS Security Ip.

~> **NOTE:** The terminate operation of rds instance does NOT take effect immediately，maybe takes for several minites.

## Example Usage

```hcl
resource "baiducloud_rds_security_ip" "default" {
    instance_id                    = "rds-xxxxx"
    security_ips                   = [192.168.0.8]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) ID of the instance
* `security_ips` - (Optional) securityIps

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `e_tag` - ETag of the instance.


## Import

RDS RDS Security Ip. can be imported, e.g.

```hcl
$ terraform import baiducloud_rds_security_ip.default id
```

