---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_scs_security_ips"
sidebar_current: "docs-baiducloud-datasource-scs_security_ips"
description: |-
  Use this data source to query SCS security ips.
---

# baiducloud_scs_security_ips

Use this data source to query SCS security ips.

## Example Usage

```hcl
data "baiducloud_scs_security_ips" "default" {
	instance_id = "scs-xxxxx"
}

output "security_ips" {
 value = "${data.baiducloud_scs.default.security_ips}"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) ID of the instance

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `security_ips` - security_ips
  * `ip` - securityIp


