---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_snapshots"
sidebar_current: "docs-baiducloud-datasource-snapshots"
description: |-
  Use this data source to query Snapshot list.
---

# baiducloud_snapshots

Use this data source to query Snapshot list.

## Example Usage

```hcl
data "baiducloud_snapshots" "default" {}

output "snapshots" {
 value = "${data.baiducloud_snapshots.default.snapshots}"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `output_file` - (Optional, ForceNew) Snapshots search result output file.
* `volume_id` - (Optional) Volume ID to be attached of snapshots, if volume is system disk, volume ID is instance ID

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `snapshots` - The result of the snapshots list.
  * `create_method` - The creation method of the snapshot.
  * `create_time` - The creation time of the snapshot.
  * `description` - The description of the snapshot.
  * `id` - The ID of the snapshot.
  * `name` - The name of the snapshot.
  * `size_in_gb` - The size of the snapshot in GB.
  * `status` - The status of the snapshot.
  * `volume_id` - The volume id of the snapshot.


