---
layout: "baiducloud"
subcategory: "CCEv2"
page_title: "BaiduCloud: baiducloud_ccev2_instance"
sidebar_current: "docs-baiducloud-resource-ccev2_instance"
description: |-
  Use this resource to bind to an instance and modify some of its attributes.
  Note that this resource will not create a real instance, it is just a way to bind to a remote instance and modify its attributes.
  If you wish to create more instances, please use baiducloud_ccev2_instance_group.
---

# baiducloud_ccev2_instance

Use this resource to bind to an instance and modify some of its attributes.
Note that this resource will not create a real instance, it is just a way to bind to a remote instance and modify its attributes.
If you wish to create more instances, please use baiducloud_ccev2_instance_group.

## Example Usage

```hcl
resource "baiducloud_ccev2_instance" "default" {
  cluster_id        = "your-cluster-id"
  instance_id       = "your-instance-id"
  spec {
    cce_instance_priority = 0
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) Cluster ID of this Instance.
* `instance_id` - (Required, ForceNew) Cluster ID of this Instance.
* `spec` - (Required) Instance Spec

The `spec` object supports the following:

* `cce_instance_priority` - (Optional) Priority of this instance.


