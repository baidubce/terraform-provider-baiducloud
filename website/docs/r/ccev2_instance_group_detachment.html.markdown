---
layout: "baiducloud"
subcategory: "Cloud Container Engine v2 (CCEv2)"
page_title: "BaiduCloud: baiducloud_ccev2_instance_group_detachment"
sidebar_current: "docs-baiducloud-resource-ccev2_instance_group_detachment"
description: |-
  Use this resource to remove instances from a CCE InstanceGroup.
---

# baiducloud_ccev2_instance_group_detachment

Use this resource to remove instances from a CCE InstanceGroup.

~> **NOTE:** After creation, it may take several minutes for the instances to be fully removed from the instance group.

## Example Usage

```hcl
resource "baiducloud_ccev2_instance_group_detachment" "example" {
  cluster_id = "cce-example"
  instance_group_id = "cce-ig-example"
  instances_to_be_removed = ["cce-example-node"]
  clean_policy = "Delete"
  delete_option {
    move_out = false
    delete_resource = true
    delete_cds_snapshot = true
    drain_node = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `clean_policy` - (Required) Whether to remove instances from the CCE cluster. `Remain` retains the instances in the cluster, `Delete` removes the instances from the cluster.
* `cluster_id` - (Required, ForceNew) The ID of the CCE cluster.
* `instance_group_id` - (Required, ForceNew) The ID of the instance group.
* `instances_to_be_removed` - (Required) IDs of node to be removed. Note this refers to the node ID within the cluster, not the actual instance ID.
* `delete_option` - (Optional) Node deletion options.Required when `clean_policy` is set to `Delete`.

The `delete_option` object supports the following:

* `delete_cds_snapshot` - (Optional) Whether to delete associated CDS snapshots when removing the node. Defaults to `false`.
* `delete_resource` - (Optional) Whether to release related resources when removing the node. Defaults to `false`.
* `drain_node` - (Optional) Whether to perform node draining before removal. Defaults to `false`.
* `move_out` - (Optional) Whether to release the instance when removing the node. `true` keeps the instance, `false` releases it. Defaults to `true`.


