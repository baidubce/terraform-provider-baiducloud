---
layout: "baiducloud"
subcategory: "Baidu Cloud Compute (BCC)"
page_title: "BaiduCloud: baiducloud_cdss"
sidebar_current: "docs-baiducloud-datasource-cdss"
description: |-
  Use this data source to query CDS list.
---

# baiducloud_cdss

Use this data source to query CDS list.

## Example Usage

```hcl
data "baiducloud_cdss" "default" {}

output "cdss" {
 value = "${data.baiducloud_cdss.default.cdss}"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `instance_id` - (Optional, ForceNew) CDS volume bind instance ID
* `output_file` - (Optional, ForceNew) CDS volume search result output file
* `zone_name` - (Optional, ForceNew) CDS volume zone name

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cdss` - CDS volume list
  * `attachments` - CDS volume attachments
    * `device` - CDS attachment device path
    * `instance_id` - CDS attachment instance id
    * `serial` - CDS attachment serial
    * `volume_id` - CDS attachment volume id
  * `auto_snapshot_policy` - CDS volume bind auto snapshot policy info
    * `created_time` - Auto Snapshot Policy created time
    * `deleted_time` - Auto Snapshot Policy deleted time
    * `id` - Auto Snapshot Policy ID
    * `last_execute_time` - Auto Snapshot Policy last execute time
    * `name` - Auto Snapshot Policy name
    * `repeat_weekdays` - Auto Snapshot Policy repeat weekdays
    * `retention_days` - Auto Snapshot Policy retention days
    * `status` - Auto Snapshot Policy status
    * `time_points` - Auto Snapshot Policy set snapshot create time points
    * `updated_time` - Auto Snapshot Policy updated time
    * `volume_count` - Auto Snapshot Policy volume count
  * `cds_id` - CDS volume id
  * `create_time` - CDS disk create time
  * `description` - CDS description
  * `disk_size_in_gb` - CDS disk size, should in [1, 32765], when snapshot_id not set, this parameter is required.
  * `expire_time` - CDS disk expire time
  * `is_system_volume` - CDS disk is system volume or not
  * `name` - CDS disk name
  * `payment_timing` - payment method, support Prepaid or Postpaid
  * `region_id` - CDS disk region id
  * `snapshot_num` - CDS disk snapshot num
  * `source_snapshot_id` - CDS disk create source snapshot id
  * `status` - CDS volume status
  * `storage_type` - CDS dist storage type, support hp1 and std1, default hp1
  * `tags` - Tags
  * `type` - CDS disk type
  * `zone_name` - Zone name


