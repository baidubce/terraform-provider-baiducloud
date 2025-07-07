---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_rds_security_ips"
subcategory: "Relational Database Service (RDS)"
sidebar_current: "docs-baiducloud-datasource-rds_security_ips"
description: |-
  Use this data source to query RDS security ips.
---

# baiducloud_rds_security_ips

Use this data source to query RDS security ips.

## Example Usage

```hcl
data "baiducloud_rds_security_ips" "default" {
	instance_id = "rds-xxxxx"
}

output "security_ips" {
 value = "${data.baiducloud_rdss.default.security_ips}"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) ID of the instance

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `security_ips` - security_ips
  * `ip` - securityIp


