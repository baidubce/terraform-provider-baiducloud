---
layout: "baiducloud"
subcategory: "Cloud Container Engine v2 (CCEv2)"
page_title: "BaiduCloud: baiducloud_ccev2_instance_group_attachment"
sidebar_current: "docs-baiducloud-resource-ccev2_instance_group_attachment"
description: |-
  Use this resource to attach instances to a CCE InstanceGroup.
---

# baiducloud_ccev2_instance_group_attachment

Use this resource to attach instances to a CCE InstanceGroup.

~> **NOTE:** After creation, instances may take several minutes to reach the `running` state.
Destroying this resource **does not** remove instances from the instance group.

## Example Usage

```hcl
resource "baiducloud_ccev2_instance_group_attachment" "example" {
  cluster_id = "cce-example"
  instance_group_id = "cce-ig-example"
  existed_instances = ["i-example"]

  existed_instances_config {
    rebuild = true
    image_id = "m-example"
    admin_password = "pass@word"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) The ID of the CCE cluster.
* `instance_group_id` - (Required, ForceNew) The ID of the instance group.
* `existed_instances_config` - (Optional) Configuration for adding instances from outside the cluster. Required with `existed_instances`.
* `existed_instances_in_cluster` - (Optional) IDs of instances already in the cluster to be added to the instance group.
* `existed_instances` - (Optional) IDs of instances outside the cluster to be added. Requires `existed_instances_config`.

The `existed_instances_config` object supports the following:

* `admin_password` - (Optional) Admin password for login.
* `image_id` - (Optional) Image ID used for rebuild.
* `machine_type` - (Optional) Machine type. Valid values: `BCC`, `BBC`, `EBC`, `HPAS`. Defaults to `BCC`.
* `rebuild` - (Optional) Whether to reinstall the operating system. This will reinstall the OS on the selected instances, clearing all data on the system disk (irrecoverable). Data on cloud disks will not be affected. Only 'true' is supported currently.
* `ssh_key_id` - (Optional) Key pair ID for login.
* `use_instance_group_config_with_disk_info` - (Optional) Whether to apply the instance group’s disk mount configuration. Defaults to `false`.
* `use_instance_group_config` - (Optional) Whether to apply the instance group’s configuration. Only 'true' is supported currently.


