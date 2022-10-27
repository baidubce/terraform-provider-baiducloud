---
layout: "baiducloud"
subcategory: "CCE"
page_title: "BaiduCloud: baiducloud_cce_container_net"
sidebar_current: "docs-baiducloud-datasource-cce_container_net"
description: |-
  Use this data source to get cce container network.
---

# baiducloud_cce_container_net

Use this data source to get cce container network.

## Example Usage

```hcl
data "baiducloud_cce_container_net" "default" {
	vpc_id   = "vpc-t6d16myuuqyu"
	vpc_cidr = "192.168.0.0/20"
}

output "net" {
  value = "${data.baiducloud_cce_container_net.default.container_net}"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_cidr` - (Required, ForceNew) CCE used vpc cidr
* `vpc_id` - (Required, ForceNew) CCE used vpc id
* `size` - (Optional, ForceNew) CCE used max container count

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `capacity` - container net capacity
* `container_net` - container net


