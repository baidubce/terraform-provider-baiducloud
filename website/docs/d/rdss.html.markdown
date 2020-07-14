---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_rdss"
sidebar_current: "docs-baiducloud-datasource-rdss"
description: |-
  Use this data source to query RDS list.
---

# baiducloud_rdss

Use this data source to query RDS list.

## Example Usage

```hcl
data "baiducloud_rdss" "default" {}

output "rdss" {
 value = "${data.baiducloud_rdss.default.rdss}"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `name_regex` - (Optional, ForceNew) Regex pattern of the search name of rds instance
* `output_file` - (Optional, ForceNew) Output file of the instances search result

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `rdss` - The result of the instances list.
  * `address` - The domain used to access a instance.
  * `category` - Category of the instance. Available values are Basic、Standard(Default), only SQLServer 2012sp3 support Basic.
  * `cpu_count` - The number of CPU
  * `create_time` - Create time of the instance.
  * `engine_version` - Engine version of the instance. MySQL support 5.5、5.6、5.7, SQLServer support 2008r2、2012sp3、2016sp1, PostgreSQL support 9.4
  * `expire_time` - Expire time of the instance.
  * `instance_id` - ID of the instance.
  * `instance_name` - Name of the instance. Support for uppercase and lowercase letters, numbers, Chinese and special characters, such as "-","_","/",".", the value must start with a letter, length 1-65.
  * `instance_status` - Status of the instance.
  * `instance_type` - Type of the instance,  Available values are Master, ReadReplica, RdsProxy.
  * `memory_capacity` - Memory capacity(GB) of the instance.
  * `node_amount` - Number of proxy node.
  * `payment_timing` - RDS payment timing
  * `port` - The port used to access a instance.
  * `region` - Region of the instance.
  * `source_instance_id` - ID of the master instance
  * `source_region` - Region of the master instance
  * `subnets` - Subnets of the instance.
    * `subnet_id` - ID of the subnet.
    * `zone_name` - Zone name of the subnet.
  * `used_storage` - Memory capacity(GB) of the instance to be used.
  * `v_net_ip` - The internal ip used to access a instance.
  * `volume_capacity` - Volume capacity(GB) of the instance
  * `vpc_id` - ID of the specific VPC
  * `zone_names` - Zone name list


