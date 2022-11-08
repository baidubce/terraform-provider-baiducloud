---
layout: "baiducloud"
subcategory: "Data Transmission Service (DTS)"
page_title: "BaiduCloud: baiducloud_dts"
sidebar_current: "docs-baiducloud-resource-dts"
description: |-
  Provide a resource to create a DTS.
---

# baiducloud_dts

Provide a resource to create a DTS.

## Example Usage

```hcl
resource "baiducloud_dts" "default" {
    product_type         = "postpay"
	type                 = "migration"
	standard             = "Large"
	source_instance_type = "public"
	target_instance_type = "public"
	cross_region_tag     = 0

    task_name            = "taskname"
	data_type			 = ["schema","base"]
    src_connection = {
        region          = "public"
		db_type			= "mysql"
		db_user			= "baidu"
		db_pass			= "password"
		db_port			= 3306
		db_host			= "106.12.174.191"
		instance_id		= "rds-lNy3KsQQ"
		instance_type	= "public"
    }
	dst_connection = {
        region          = "public"
		db_type			= "mysql"
		db_user			= "baidu"
		db_pass			= "password"
		db_port			= 3306
		db_host			= "106.12.174.191"
		instance_id		= "rds-lNy3KsQQ"
		instance_type	= "public"
    }
    schema_mapping {
			type        = "db"
			src			= "db1"
			dst			= "db2"
			where		= ""
	}
}
```

## Argument Reference

The following arguments are supported:

* `cross_region_tag` - (Required, ForceNew) cross region tag of the task. Available value are 0, 1.
* `data_type` - (Required) Dts task data type
* `product_type` - (Required, ForceNew) product type of the task. Available value is postpay.
* `schema_mapping` - (Required) schema mapping
* `source_instance_type` - (Required, ForceNew) source instance type of the task. Available values are public, bcerds.
* `standard` - (Required, ForceNew) standard of the task. Available value is Large.
* `target_instance_type` - (Required, ForceNew) target instance type of the task. Available values are public, bcerds.
* `task_name` - (Required) Dts task name
* `type` - (Required, ForceNew) type of the task. Available values are migration, sync, subscribe.
* `dst_connection` - (Optional) Connection
* `dts_id` - (Optional, ForceNew) Dts task id
* `granularity` - (Optional) Dts granularity
* `init_position_type` - (Optional) Dts init position type
* `init_position` - (Optional) Dts init position
* `operation` - (Optional) operation of the task. Available values are precheck, getprecheck, start, pause, shutdown.
* `queue_type` - (Optional) Dts queue type
* `src_connection` - (Optional) Connection
* `sub_status` - (Optional) sub status

The `schema_mapping` object supports the following:

* `dst` - (Optional) dst
* `src` - (Optional) src
* `type` - (Optional) type
* `where` - (Optional) where

The `sub_status` object supports the following:

* `b` - b
* `i` - i
* `s` - s

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `base` - schemaInfo
  * `count` - count
  * `current` - current
  * `expect_finish_time` - expect finish time
  * `speed` - speed
* `create_time` - Dts create time
* `errmsg` - Dts error message
* `increment` - increment
* `pay_create_time` - Dts pay create time
* `pay_end_time` - Dts pay end time
* `region` - Dts region
* `running_time` - Dts task running time
* `schema` - schemaInfo
  * `count` - count
  * `current` - current
  * `expect_finish_time` - expect finish time
  * `speed` - speed
* `sdk_realtime_progress` - Dts sdk realtime progress
* `status` - Dts task status
* `sub_end_time` - Dts subDataScope end time
* `sub_start_time` - Dts subDataScope start time


## Import

DTS can be imported, e.g.

```hcl
$ terraform import baiducloud_dts.default dts
```

