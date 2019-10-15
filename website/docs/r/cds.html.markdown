---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_cds"
sidebar_current: "docs-baiducloud-resource-cds"
description: |-
  Provide a resource to create a CDS.
---

# baiducloud_cds

Provide a resource to create a CDS.

## Example Usage

```hcl
resource "baiducloud_cds" "default" {
  name                    = "terraformCreate"
  description             = "terraform create cds"
  payment_timing          = "Postpaid"
  auto_snapshot_policy_id = "asp-xyYk0XFC"
  snapshot_id             = "s-WTGlKBR1"
}
```

## Argument Reference

The following arguments are supported:

* `payment_timing` - (Required) payment method, support Prepaid or Postpaid
* `auto_snapshot` - (Optional) Delete relate auto snapshot when release this cds volume
* `description` - (Optional) CDS volume description
* `disk_size_in_gb` - (Optional) CDS disk size, support between 5 and 32765, if snapshot_id not set, this parameter is required.
* `manual_snapshot` - (Optional) Delete relate snapshot when release this cds volume
* `name` - (Optional) CDS volume name
* `reservation_length` - (Optional) Prepaid reservation length, support [1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36], only useful when payment_timing is Prepaid
* `reservation_time_unit` - (Optional) Prepaid reservation time unit, only support Month now
* `snapshot_id` - (Optional, ForceNew) Snapshot id, support create cds use snapshot, when set this parameter, cds_disk_size is ignored
* `storage_type` - (Optional) CDS dist storage type, support hp1, std1, cloud_hp1 and hdd, default hp1, see https://cloud.baidu.com/doc/BCC/s/6jwvyo0q2/#storagetype for detail
* `zone_name` - (Optional) Zone name

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `auto_snapshot_policy_id` - CDS bind Auto Snapshot policy id
* `create_time` - CDS volume create time
* `expire_time` - CDS volume expire time
* `status` - CDS volume status
* `type` - CDS volume type


## Import

CDS can be imported, e.g.

```hcl
$ terraform import baiducloud_cds.default id
```

