---
layout: "baiducloud"
subcategory: "BCC"
page_title: "BaiduCloud: baiducloud_auto_snapshot_policy"
sidebar_current: "docs-baiducloud-resource-auto_snapshot_policy"
description: |-
  Provide a resource to create an AutoSnapshotPolicy.
---

# baiducloud_auto_snapshot_policy

Provide a resource to create an AutoSnapshotPolicy.

## Example Usage

```hcl
resource "baiducloud_auto_snapshot_policy" "my-asp" {
  name            = "${var.name}"
  time_points     = [0, 22]
  repeat_weekdays = [0, 3]
  retention_days  = -1
  volume_ids      = ["v-Trb3rQXa"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the automatic snapshot policy, which supports uppercase and lowercase letters, numbers, Chinese and special characters, such as "-","_","/",",".", and the value must start with a letter, length 1-65.
* `repeat_weekdays` - (Required) Repeat time of the automatic snapshot policy, supporting in range of [0, 6]
* `retention_days` - (Required) Number of days to retain the automatic snapshot, and -1 means permanently retained.
* `time_points` - (Required) Time point of generate snapshot in a day, the minimum unit is hour, supporting in range of [0, 23]
* `volume_ids` - (Optional) Volume id list to be attached of the automatic snapshot policy, these CDS volumes must be in-used.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `created_time` - Creation time of the automatic snapshot policy.
* `deleted_time` - Deletion time of the automatic snapshot policy.
* `last_execute_time` - Last execution time of the automatic snapshot policy.
* `status` - Status of the automatic snapshot policy.
* `updated_time` - Update time of the automatic snapshot policy.
* `volume_count` - The count of volumes associated with the snapshot.


## Import

AutoSnapshotPolicy can be imported, e.g.

```hcl
$ terraform import baiducloud_auto_snapshot_policy.default id
```

