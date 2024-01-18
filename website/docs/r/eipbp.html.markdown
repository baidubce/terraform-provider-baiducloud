---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_eipbp"
sidebar_current: "docs-baiducloud-resource-eipbp"
description: |-
  Provide a resource to create an EIP BP.
---

# baiducloud_eipbp

Provide a resource to create an EIP BP.

## Example Usage

```hcl
resource "baiducloud_eipbp" "default" {
  name              = "testEIPbp"
  eip               = 10.23.42.12
  bandwidth_in_mbps = 100
  eip_group_id      = "xxx"
}
```

## Argument Reference

The following arguments are supported:

* `bandwidth_in_mbps` - (Required) Eip bp bandwidth(Mbps), if payment_timing is Prepaid or billing_method is ByBandWidth, support between 1 and 200, if billing_method is ByTraffic, support between 1 and 1000
* `eip_group_id` - (Required, ForceNew) eip group id of eip bp
* `eip` - (Required, ForceNew) eip of eip bp
* `auto_release_time` - (Optional) Eip bp auto release time
* `name` - (Optional) Eip bp name, length must be between 1 and 65 bytes
* `type` - (Optional) Eip bp type

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `bind_type` - Eip bp bind type
* `bp_id` - id of EIP bp
* `create_time` - Eip bp create_time
* `eips` - Eip bp eips
* `instance_bandwidth_in_mbps` - Eip bp instance bandwidth in mbps
* `instance_id` - Eip bp instance id
* `region` - Eip bp region


## Import

EIP bp can be imported, e.g.

```hcl
$ terraform import baiducloud_eipbp.default bp_id
```

