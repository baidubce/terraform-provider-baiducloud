---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_scss"
subcategory: "SCS"
sidebar_current: "docs-baiducloud-datasource-scss"
description: |-
  Use this data source to query SCS list.
---

# baiducloud_scss

Use this data source to query SCS list.

## Example Usage

```hcl
data "baiducloud_scss" "default" {}

output "scss" {
 value = "${data.baiducloud_scss.default.scss}"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `name_regex` - (Optional, ForceNew) Regex pattern of the search name of scs instance
* `output_file` - (Optional, ForceNew) Output file of the instances search result

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `scss` - The result of the instances list.
  * `auto_renew` - Whether to automatically renew.
  * `capacity` - Memory capacity(GB) of the instance.
  * `cluster_type` - Type of the instance,  Available values are cluster, master_slave.
  * `create_time` - Create time of the instance.
  * `domain` - Domain of the instance.
  * `engine_version` - Engine version of the instance. Available values are 3.2, 4.0.
  * `engine` - Engine of the instance. Available values are redis, memcache.
  * `expire_time` - Expire time of the instance.
  * `instance_id` - ID of the instance.
  * `instance_name` - Name of the instance. Support for uppercase and lowercase letters, numbers, Chinese and special characters, such as "-","_","/",".", the value must start with a letter, length 1-65.
  * `instance_status` - Status of the instance.
  * `node_type` - Type of the instance. Available values are cache.n1.micro, cache.n1.small, cache.n1.medium...cache.n1hs3.4xlarge.
  * `payment_timing` - SCS payment timing
  * `port` - The port used to access a instance.
  * `proxy_num` - The number of instance proxy.
  * `replication_num` - The number of instance copies.
  * `shard_num` - The number of instance shard. IF cluster_type is cluster, support 2/4/6/8/12/16/24/32/48/64/96/128, if cluster_type is master_slave, support 1.
  * `tags` - Tags
  * `used_capacity` - Memory capacity(GB) of the instance to be used.
  * `v_net_ip` - ID of the specific vnet.
  * `zone_names` - Zone name list


