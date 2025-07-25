---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_rds_readonly_instance"
subcategory: "Relational Database Service (RDS)"
sidebar_current: "docs-baiducloud-resource-rds_readonly_instance"
description: |-
  Use this resource to get information about a RDS readonly instance.
---

# baiducloud_rds_readonly_instance

Use this resource to get information about a RDS readonly instance.

~> **NOTE:** The terminate operation of rds readonly instance does NOT take effect immediately，maybe takes for several minites.

## Example Usage

```hcl
resource "baiducloud_rds_readonly_instance" "default" {
    billing = {
        payment_timing        = "Postpaid"
    }
    source_instance_id        = baiducloud_rds_instance.default.instance_id
    cpu_count                 = 1
    memory_capacity           = 1
    volume_capacity           = 5
}
```

## Argument Reference

The following arguments are supported:

* `billing` - (Required) Billing information of the Rds.
* `cpu_count` - (Required) The number of CPU
* `memory_capacity` - (Required) Memory capacity(GB) of the instance.
* `source_instance_id` - (Required, ForceNew) ID of the master instance
* `subnets` - (Required) Subnets of the instance.
* `volume_capacity` - (Required) Volume capacity(GB) of the instance
* `vpc_id` - (Required, ForceNew) ID of the specific VPC
* `auto_renew_time_length` - (Optional, ForceNew) The time length of automatic renewal, empty indicating no automatic renewal. It is valid when payment_timing is Prepaid, and the value should be 1-9 when the auto_renew_time_unit is month and 1-3 when the auto_renew_time_unit is year.
* `auto_renew_time_unit` - (Optional, ForceNew) Time unit of automatic renewal, the value can be month or year. The default value is month. It is valid only when the payment_timing is Prepaid.
* `category` - (Optional, ForceNew) Category of the instance. Available values are Basic、Standard(Default), only SQLServer 2012sp3 support Basic.
* `instance_name` - (Optional) Name of the instance. Support for uppercase and lowercase letters, numbers, Chinese and special characters, such as "-","_","/",".", the value must start with a letter, length 1-65.
* `reservation` - (Optional) Reservation of the Rds.
* `tags` - (Optional, ForceNew) Tags, do not support modify

The `billing` object supports the following:

* `payment_timing` - (Required) Payment timing of billing, which can be Prepaid or Postpaid. The default is Postpaid.

The `reservation` object supports the following:

* `reservation_length` - (Required) The reservation length that you will pay for your resource. It is valid when payment_timing is Prepaid. Valid values: [1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36].
* `reservation_time_unit` - (Required) The reservation time unit that you will pay for your resource. It is valid when payment_timing is Prepaid. The value can only be month currently, which is also the default value.

The `subnets` object supports the following:

* `subnet_id` - (Optional, ForceNew) ID of the subnet.
* `zone_name` - (Optional, ForceNew) Zone name of the subnet.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `address` - The domain used to access a instance.
* `create_time` - Create time of the instance.
* `expire_time` - Expire time of the instance.
* `instance_id` - ID of the instance.
* `instance_status` - Status of the instance.
* `instance_type` - Type of the instance,  Available values are Master, ReadReplica, RdsProxy.
* `node_amount` - Number of proxy node.
* `port` - The port used to access a instance.
* `region` - Region of the instance.
* `used_storage` - Memory capacity(GB) of the instance to be used.
* `v_net_ip` - The internal ip used to access a instance.
* `zone_names` - Zone name list


## Import

RDS readonly instance can be imported, e.g.

```hcl
$ terraform import baiducloud_rds_readonly_instance.default id
```

