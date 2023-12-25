---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_scs_security_ip"
sidebar_current: "docs-baiducloud-resource-scs_security_ip"
description: |-
  Use this resource to get information about a SCS Security Ip.
---

# baiducloud_scs_security_ip

Use this resource to get information about a SCS Security Ip.

~> **NOTE:** The terminate operation of scs instance does NOT take effect immediately，maybe takes for several minites.

## Example Usage

```hcl
resource "baiducloud_scs_security_ip" "default" {
    instance_id                    = "scs-xxxxx"
    security_ips                   = [192.168.0.8]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) ID of the instance
* `security_ips` - (Optional) securityIps


## Import

SCS Security Ip. can be imported, e.g.

```hcl
$ terraform import baiducloud_scs_security_ip.default id
```

