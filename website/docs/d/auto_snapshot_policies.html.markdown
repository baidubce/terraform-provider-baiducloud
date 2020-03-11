---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_auto_snapshot_policies"
sidebar_current: "docs-baiducloud-datasource-auto_snapshot_policies"
description: |-
  Use this data source to query Auto Snapshot Policy list.
---

# baiducloud_auto_snapshot_policies

Use this data source to query Auto Snapshot Policy list.

## Example Usage

```hcl
data "baiducloud_auto_snapshot_policies" "default" {}

output "auto_snapshot_policiess" {
 value = "${data.baiducloud_auto_snapshot_policies.default.auto_snapshot_policies}"
}
```

## Argument Reference

The following arguments are supported:

* `asp_name` - (Optional) Name of the automatic snapshot policy.
* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `output_file` - (Optional, ForceNew) Automatic snapshot policies search result output file.
* `volume_name` - (Optional) Name of the volume.

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `auto_snapshot_policies` - The automatic snapshot policies search result list.
  * `created_time` - The creation time of the automatic snapshot policy.
  * `deleted_time` - The deletion time of the automatic snapshot policy.
  * `id` - The ID of the automatic snapshot policy.
  * `last_execute_time` - The last execution time of the automatic snapshot policy.
  * `name` - The name of the automatic snapshot policy.
  * `repeat_weekdays` - The repeat weekdays of the automatic snapshot policy.
  * `retention_days` - The retention days of the automatic snapshot policy.
  * `status` - The status of the automatic snapshot policy.
  * `time_points` - The time points of the automatic snapshot policy.
  * `updated_time` - The updation time of the automatic snapshot policy.
  * `volume_count` - The volume count of the automatic snapshot policy.


