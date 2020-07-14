---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_scs"
sidebar_current: "docs-baiducloud-resource-scs"
description: |-
  Use this resource to get information about a SCS.
---

# baiducloud_scs

Use this resource to get information about a SCS.

~> **NOTE:** The terminate operation of scs does NOT take effect immediately，maybe takes for several minites.

## Example Usage

```hcl
resource "baiducloud_scs" "default" {
	billing = {
		payment_timing = "Postpaid"
	}
	instance_name = "terraform-redis"
	purchase_count = 1
	port = 6379
	engine_version = "3.2"
	node_type = "cache.n1.micro"
	architecture_type = "master_slave"
	replication_num = 1
	shard_num = 1
}
```

## Argument Reference

The following arguments are supported:

* `billing` - (Required) Billing information of the Scs.
* `instance_name` - (Required) Name of the instance. Support for uppercase and lowercase letters, numbers, Chinese and special characters, such as "-","_","/",".", the value must start with a letter, length 1-65.
* `node_type` - (Required) Type of the instance. Available values are cache.n1.micro, cache.n1.small, cache.n1.medium...cache.n1hs3.4xlarge.
* `cluster_type` - (Optional, ForceNew) Type of the instance,  Available values are cluster, master_slave.
* `engine_version` - (Optional, ForceNew) Engine version of the instance. Available values are 3.2, 4.0.
* `port` - (Optional, ForceNew) The port used to access a instance.
* `proxy_num` - (Optional, ForceNew) The number of instance proxy.
* `purchase_count` - (Optional) Count of the instance to buy
* `replication_num` - (Optional, ForceNew) The number of instance copies.
* `shard_num` - (Optional) The number of instance shard. IF cluster_type is cluster, support 2/4/6/8/12/16/24/32/48/64/96/128, if cluster_type is master_slave, support 1.
* `subnets` - (Optional) Subnets of the instance.
* `vpc_id` - (Optional, ForceNew) ID of the specific VPC

The `billing` object supports the following:

* `payment_timing` - (Required) Payment timing of billing, which can be Prepaid or Postpaid. The default is Postpaid.
* `reservation` - (Optional) Reservation of the Scs.

The `reservation` object supports the following:

* `reservation_length` - (Required) The reservation length that you will pay for your resource. It is valid when payment_timing is Prepaid. Valid values: [1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36].
* `reservation_time_unit` - (Required) The reservation time unit that you will pay for your resource. It is valid when payment_timing is Prepaid. The value can only be month currently, which is also the default value.

The `subnets` object supports the following:

* `subnet_id` - (Optional, ForceNew) ID of the subnet.
* `zone_name` - (Optional, ForceNew) Zone name of the subnet.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `auto_renew_time_length` - The time length of automatic renewal. It is valid when payment_timing is Prepaid, and the value should be 1-9 when the auto_renew_time_unit is month and 1-3 when the auto_renew_time_unit is year. Default to 1.
* `auto_renew_time_unit` - Time unit of automatic renewal, the value can be month or year. The default value is empty, indicating no automatic renewal. It is valid only when the payment_timing is Prepaid.
* `auto_renew` - Whether to automatically renew.
* `capacity` - Memory capacity(GB) of the instance.
* `create_time` - Create time of the instance.
* `domain` - Domain of the instance.
* `engine` - Engine of the instance. Available values are redis, memcache.
* `expire_time` - Expire time of the instance.
* `instance_id` - ID of the instance.
* `instance_status` - Status of the instance.
* `payment_timing` - SCS payment timing
* `tags` - Tags
* `used_capacity` - Memory capacity(GB) of the instance to be used.
* `v_net_ip` - The internal ip used to access a instance.
* `zone_names` - Zone name list


## Import

SCS can be imported, e.g.

```hcl
$ terraform import baiducloud_scs.default id
```

