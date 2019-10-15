---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_instances"
sidebar_current: "docs-baiducloud-datasource-instances"
description: |-
  Use this data source to query BCC Instance list.
---

# baiducloud_instances

Use this data source to query BCC Instance list.

## Example Usage

```hcl
data "baiducloud_instances" "default" {}

output "instances" {
 value = "${data.baiducloud_instances.default.instances}"
}
```

## Argument Reference

The following arguments are supported:

* `dedicated_host_id` - (Optional) Dedicated host id of the instance to retrieve.
* `internal_ip` - (Optional) Internal ip address of the instance to retrieve.
* `output_file` - (Optional, ForceNew) Output file of the instances search result
* `zone_name` - (Optional) Name of the available zone to which the instance belongs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instances` - The result of the instances list.
  * `auto_renew` - Whether to automatically renew.
  * `card_count` - The card count of the instance.
  * `cpu_count` - The cpu count of the instance.
  * `create_time` - The creation time of the instance.
  * `dedicated_host_id` - The dedicated host id of the instance.
  * `description` - The description of the instance.
  * `ephemeral_disks` - The ephemeral disks of the instance.
    * `size_in_gb` - The size(GB) of the ephemeral disk.
    * `storage_type` - The storage type of the ephemeral disk.
  * `expire_time` - The expire time of the instance.
  * `fpga_card` - The fgpa card of the instance.
  * `gpu_card` - The gpu card of the instance.
  * `image_id` - The image id of the instance.
  * `instance_id` - The ID of the instance.
  * `instance_type` - The type of the instance.
  * `internal_ip` - The internal ip of the instance.
  * `keypair_id` - The key pair id of the instance.
  * `keypair_name` - The key pair name of the instance.
  * `memory_capacity_in_gb` - The memory capacity in GB of the instance.
  * `name` - The name of the instance.
  * `network_capacity_in_mbps` - The network capacity in Mbps of the instance.
  * `payment_timing` - The payment timing of the instance.
  * `placement_policy` - The placement policy of the instance.
  * `public_ip` - The public ip of the instance.
  * `root_disk_size_in_gb` - The system disk size in GB of the instance.
  * `root_disk_storage_type` - The system disk storage type of the instance.
  * `status` - The status of the instance.
  * `subnet_id` - The subnet ID of the instance.
  * `tags` - Tags
    * `tag_key` - Tag's key
    * `tag_value` - Tag's value
  * `vpc_id` - The VPC ID of the instance.
  * `zone_name` - The zone name of the instance.


