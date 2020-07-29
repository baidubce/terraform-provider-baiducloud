---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_dtss"
sidebar_current: "docs-baiducloud-datasource-dtss"
description: |-
  Use this data source to query DTS list.
---

# baiducloud_dtss

Use this data source to query DTS list.

## Example Usage

```hcl
data "baiducloud_dtss" "default" {}

output "dtss" {
 value = "${data.baiducloud_dtss.default.dtss}"
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Required) type of the task. Available values are migration, sync, subscribe.
* `dts_name` - (Optional, ForceNew) Name of the Dts to be queried
* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `output_file` - (Optional, ForceNew) Dtss search result output file

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `dtss` - Dts list
  * `base` - schemaInfo
    * `count` - count
    * `current` - current
    * `expect_finish_time` - expect finish time
    * `speed` - speed
  * `create_time` - Dts create time
  * `cross_region_tag` - Dts cross region tag
  * `data_type` - Dts task data type
  * `dst_connection` - Connection
  * `dts_id` - Dts task id
  * `errmsg` - Dts errmsg
  * `granularity` - Dts granularity
  * `increment` - increment
  * `pay_create_time` - Dts pay create time
  * `pay_end_time` - Dts pay end time
  * `product_type` - Dts product type
  * `region` - Dts region
  * `running_time` - Dts task running time
  * `schema_mapping` - schema mapping
    * `dst` - dst
    * `src` - src
    * `type` - type
    * `where` - where
  * `schema` - schemaInfo
    * `count` - count
    * `current` - current
    * `expect_finish_time` - expect finish time
    * `speed` - speed
  * `sdk_realtime_progress` - Dts sdk realtime progress
  * `source_instance_type` - Dts source instance type
  * `src_connection` - Connection
  * `standard` - Dts standard
  * `status` - Dts task status
  * `sub_end_time` - Dts subDataScope end time
  * `sub_start_time` - Dts subDataScope start time
  * `sub_status` - sub status
    * `b` - b
    * `i` - i
    * `s` - s
  * `target_instance_type` - Dts target instance type
  * `task_name` - Dts task name


