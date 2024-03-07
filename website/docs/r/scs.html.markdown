---
layout: "baiducloud"
subcategory: "Simple Cache Service for Redis (SCS)"
page_title: "BaiduCloud: baiducloud_scs"
sidebar_current: "docs-baiducloud-resource-scs"
description: |-
  Use this resource to get information about a SCS.
---

# baiducloud_scs

Use this resource to get information about a SCS.

More information about SCS can be found in the [Developer Guide](https://cloud.baidu.com/doc/SCS/index.html).

~> **NOTE:** The terminate operation of scs does NOT take effect immediately，maybe takes for several minites.

## Example Usage

### Memcache
~> **NOTE:** Memcache currently does NOT support specifying `node_type`, set to `cache.n1.micro` directly.
```terraform
resource "baiducloud_scs" "default" {
	payment_timing = "Postpaid"
	instance_name = "terraform-memcache"
	engine = "memcache"
	port = 11211
	node_type = "cache.n1.micro"
	cluster_type = "default"
	shard_num = 2
}
```

### Redis
```terraform
resource "baiducloud_scs" "default" {
	payment_timing = "Postpaid"
	instance_name = "terraform-redis"
	port = 6379
	engine_version = "3.2"
	node_type = "cache.n1.micro"
	cluster_type = "master_slave"
	replication_num = 1
	shard_num = 1
}
```

### PegaDb
```terraform
resource "baiducloud_scs" "default" {
	payment_timing = "Prepaid"
	reservation_length = 2
	reservation_time_unit = "month"
	instance_name = "terraform-pegadb"
	purchase_count = 1
	engine = "PegaDB"
	node_type = "pega.g4s1.micro"
	cluster_type = "cluster"
	store_type = 3
	disk_flavor = 60
	port = 6379
	replication_num = 2
	shard_num = 1
	proxy_num = 2
	vpc_id = "vpc-ne32rahkaceu"
	subnets {
		subnet_id = "sbn-vhnqd71mivjq"
		zone_name = "cn-bj-d"
	}
	replication_info {
		availability_zone = "cn-bj-d"
		is_master         = 1
		subnet_id         = "sbn-vhnqd71mivjq"
	}
	replication_info {
		availability_zone = "cn-bj-d"
		is_master         = 0
		subnet_id         = "sbn-vhnqd71mivjq"
	}
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required) Name of the instance. Support for uppercase and lowercase letters, numbers, Chinese and special characters, such as `-`, `_`, `/`, `.`. Must start with a letter, length 1-65.
* `node_type` - (Required) Node type of the instance. e.g. `cache.n1.micro`. To learn about supported node type, see documentation on [Supported Node Types](https://cloud.baidu.com/doc/SCS/s/1jwvxtsh0#%E5%AE%9E%E4%BE%8B%E8%A7%84%E6%A0%BC)
* `backup_days` - (Optional) Identifies which days of the week the backup cycle is performed: Mon (Monday) Tue (Tuesday) Wed (Wednesday) Thu (Thursday) Fri (Friday) Sat (Saturday) Sun (Sunday) comma separated, the values are as follows: Sun,Mon,Tue,Wed,Thu,Fri,Sta. Note: Automatic backup is only supported if the number of slave nodes is greater than 1
* `backup_time` - (Optional) Identifies when to perform backup in a day, UTC time (+8 is Beijing time) value such as: 01:05:00
* `billing` - (Optional) **Deprecated**. Use `payment_timing`, `reservation_length`, `reservation_time_unit` instead. Billing information of the Scs.
* `client_auth` - (Optional, Sensitive) Access password of the instance. Should be 8-16 characters, and contains at least two types of letters, numbers and symbols. Allowed symbols include `$ ^ * ( ) _ + - =`.
* `cluster_type` - (Optional, ForceNew) Type of the instance. If `engine` is `memcache`, must be `default`. Valid values for other engine type: `cluster`, `master_slave`.  Defaults to `master_slave`.
* `disk_flavor` - (Optional) Storage size(GB) when use PegaDB. Must be between `50` and `160`
* `disk_type` - (Optional) Disk type of the instance. Valid values: `cloud_hp1`, `enhanced_ssd_pl1`.
* `enable_read_only` - (Optional) Whether the copies are read only. Valid values: `1`(enabled), `2`(disabled). Defaults to `2`.
* `engine_version` - (Optional) Engine version of the instance. Must be set when `engine` is `redis`. Valid values: `3.2`, `4.0`, `5.0`, `6.0`.
* `engine` - (Optional) Engine of the instance. Valid values: `memcache`, `redis`, `PegaDB`. Defaults to `redis`.
* `expire_day` - (Optional) Backup file expiration time, value such as: 3
* `payment_timing` - (Optional) Payment timing of billing, Valid values: `Prepaid`, `Postpaid`.
* `port` - (Optional, ForceNew) Port number used to access the instance. Must be between `1025` and `65534`. Defaults to `6379`.
* `proxy_num` - (Optional, ForceNew) The number of instance proxy. If `cluster_type` is `cluster`, set to the value of `shard_num` (if `shard_num` equals `1`, set to `2`). If `cluster_type` is `master_slave`, set to `0`. Defaults to `0`.
* `purchase_count` - (Optional) Count of the instance to buy. Must be between `1` and `10`. Defaults to `1`.
* `replication_info` - (Optional) Replica info of the instance. Adding and removing replicas at same time in one operation is not supported.
* `replication_num` - (Optional, ForceNew) The number of instance replicas. If `cluster_type` is `cluster`, must be between `2` and `5`. If `cluster_type` is `master_slave`, must be between `1` and `5`. Defaults to `2`.
* `replication_resize_type` - (Optional) Replica resize type. Must set when change `replication_info`. Valid values: `add`, `delete`.
* `reservation_length` - (Optional) Prepaid billing reservation length, only useful when `payment_timing` is `Prepaid`. Valid values: `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `12`, `24`, `36`
* `reservation_time_unit` - (Optional) Prepaid billing reservation time unit, only useful when `payment_timing` is `Prepaid`. Only support `month` now.
* `shard_num` - (Optional) The number of instance shard. Defaults to `1`. To learn about supported shard number, see documentation on [Supported Node Types](https://cloud.baidu.com/doc/SCS/s/1jwvxtsh0#%E5%AE%9E%E4%BE%8B%E8%A7%84%E6%A0%BC)
* `store_type` - (Optional) Store type of the instance. Valid values: `0`(high performance memory), `1`(ssd local disk), `3`(capacity storage, only for PegaDB).
* `subnets` - (Optional) Subnets of the instance.
* `tags` - (Optional) Tags, support setting when creating instance, do not support modify
* `vpc_id` - (Optional) ID of the specific VPC
* `security_groups` - (Optional) Security group ids of the scs.


The `billing` object supports the following:

* `payment_timing` - (Required) **Deprecated**. Payment timing of billing, which can be Prepaid or Postpaid. The default is Postpaid.
* `reservation` - (Optional) **Deprecated**. Reservation of the Scs.

The `reservation` object supports the following:

* `reservation_length` - (Required) **Deprecated**. The reservation length that you will pay for your resource. It is valid when payment_timing is Prepaid. Valid values: [1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36].
* `reservation_time_unit` - (Required) **Deprecated**. The reservation time unit that you will pay for your resource. It is valid when payment_timing is Prepaid. The value can only be month currently, which is also the default value.

The `replication_info` object supports the following:

* `availability_zone` - (Required) Availability zone of the replica. e.g. `cn-bj-a`.
* `is_master` - (Required) Whether the replica is master node. Valid values: `1`(master node), `0`(slave node).
* `subnet_id` - (Required) Subnet id of the replica.

The `subnets` object supports the following:

* `subnet_id` - (Optional) ID of the subnet.
* `zone_name` - (Optional) Zone name of the subnet. e.g. `cn-bj-a`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `auto_renew_time_length` - The time length of automatic renewal. It is valid when payment_timing is Prepaid, and the value should be 1-9 when the auto_renew_time_unit is month and 1-3 when the auto_renew_time_unit is year.
* `auto_renew_time_unit` - Time unit of automatic renewal, the value can be month or year. The default value is empty, indicating no automatic renewal. It is valid only when the payment_timing is Prepaid.
* `auto_renew` - Whether to automatically renew.
* `capacity` - Memory capacity(GB) of the instance.
* `create_time` - Create time of the instance.
* `domain` - Domain of the instance.
* `expire_time` - Expire time of the instance.
* `instance_id` - ID of the instance.
* `instance_status` - Status of the instance.
* `used_capacity` - The amount of memory(GB) used by the instance.
* `v_net_ip` - The internal ip used to access a instance.
* `zone_names` - Zone name list


## Import

SCS can be imported, e.g.

```hcl
$ terraform import baiducloud_scs.default id
```

