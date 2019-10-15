---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_eip_association"
sidebar_current: "docs-baiducloud-resource-eip_association"
description: |-
  Provide a resource to create an EIP association, bind an EIP with instance.
---

# baiducloud_eip_association

Provide a resource to create an EIP association, bind an EIP with instance.

## Example Usage

```hcl
resource "baiducloud_eip_association" "default" {
  eip           = "1.1.1.1"
  instance_type = "BCC"
  instance_id   = "i-7xc9Q6KR"
}
```

## Argument Reference

The following arguments are supported:

* `eip` - (Required, ForceNew) EIP which need to associate with instance
* `instance_id` - (Required, ForceNew) Instance ID which need to associate with EIP
* `instance_type` - (Required, ForceNew) Instance type which need to associate with EIP, support BCC/BLB/NAT/VPN


## Import

EIP association can be imported, e.g.

```hcl
$ terraform import baiducloud_eip_association.default eip
```

