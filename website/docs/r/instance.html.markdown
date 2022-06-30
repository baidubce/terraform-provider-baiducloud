---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_instance"
sidebar_current: "docs-baiducloud-resource-instance"
description: |-
  Use this resource to get information about a BCC instance.
---

# baiducloud_instance

Use this resource to get information about a BCC instance.

~> **NOTE:** The terminate operation of bcc does NOT take effect immediately，maybe takes for several minites.

~> **NOTE:** It is recommended to set the maximum parallelism number to 18, otherwise it may cause errors ("There are too many connections").

## Example Usage

```hcl
resource "baiducloud_instance" "my-server" {
  image_id = "m-A4jJpFzi"
  name = "my-instance"
  availability_zone = "cn-bj-a"
  cpu_count = "2"
  memory_capacity_in_gb = "8"
  billing = {
    payment_timing = "Postpaid"
  }
}
```

## Argument Reference

The following arguments are supported:

* `billing` - (Required) Billing information of the instance.
* `cpu_count` - (Required) Number of CPU cores to be created for the instance.
* `image_id` - (Required) ID of the image to be used for the instance.
* `memory_capacity_in_gb` - (Required) Memory capacity(GB) of the instance to be created.
* `action` - (Optional) Start or stop the instance, which can only be start or stop, default start.
* `admin_pass` - (Optional) Password of the instance to be started. This value should be 8-16 characters, and English, numbers and symbols must exist at the same time. The symbols is limited to "!@#$%^*()".
* `auto_renew_time_length` - (Optional, ForceNew) The time length of automatic renewal. It is valid when payment_timing is Prepaid, and the value should be 1-9 when the auto_renew_time_unit is month and 1-3 when the auto_renew_time_unit is year. Default to 1.
* `auto_renew_time_unit` - (Optional, ForceNew) Time unit of automatic renewal, the value can be month or year. The default value is empty, indicating no automatic renewal. It is valid only when the payment_timing is Prepaid.
* `availability_zone` - (Optional, ForceNew) Availability zone to start the instance in.
* `card_count` - (Optional, ForceNew) Count of the GPU cards or FPGA cards to be carried for the instance to be created, it is valid only when the gpu_card or fpga_card field is not empty.
* `cds_auto_renew` - (Optional, ForceNew) Whether the cds is automatically renewed. It is valid when payment_timing is Prepaid. Default to false.
* `cds_disks` - (Optional) CDS disks of the instance.
* `dedicate_host_id` - (Optional, ForceNew) The ID of dedicated host.
* `delete_cds_snapshot_flag` - (Optional, ForceNew) Whether to release the cds disk snapshots, default to false. It is effective only when the related_release_flag is true.
* `description` - (Optional) Description of the instance.
* `ephemeral_disks` - (Optional) Ephemeral disks of the instance.
* `fpga_card` - (Optional, ForceNew) FPGA card of the instance.
* `gpu_card` - (Optional, ForceNew) GPU card of the instance.
* `instance_type` - (Optional, ForceNew) Type of the instance to start. Available values are N1, N2, N3, N4, N5, C1, C2, S1, G1, F1. Default to N3.
* `keypair_id` - (Optional, ForceNew) Key pair id of the instance.
* `name` - (Optional) Name of the instance. Support for uppercase and lowercase letters, numbers, Chinese and special characters, such as "-","_","/",".", the value must start with a letter, length 1-65.
* `related_release_flag` - (Optional, ForceNew) Whether to release the eip and data disks mounted by the current instance. Can only be released uniformly or not. Default to false.
* `relation_tag` - (Optional, ForceNew) The new instance associated with existing Tags or not, default false. The Tags should already exit if set true
* `root_disk_size_in_gb` - (Optional, ForceNew) System disk size(GB) of the instance to be created. The value range is [40,500]GB, Default to 40GB, and more than 40GB is charged according to the cloud disk price. Note that the specified system disk size needs to meet the minimum disk space limit of the mirror used.
* `root_disk_storage_type` - (Optional, ForceNew) System disk storage type of the instance. Available values are std1, hp1, cloud_hp1, local, sata, ssd. Default to cloud_hp1.
* `security_groups` - (Optional) Security groups of the instance.
* `subnet_id` - (Optional) The subnet ID of VPC. The default subnet will be used when it is empty. The instance will restart after changing the subnet.
* `tags` - (Optional, ForceNew) Tags, do not support modify
* `user_data` - (Optional) User Data

The `billing` object supports the following:

* `payment_timing` - (Required) Payment timing of billing, which can be Prepaid or Postpaid. The default is Postpaid.
* `reservation` - (Optional) Reservation of the instance.

The `reservation` object supports the following:

* `reservation_length` - (Required) The reservation length that you will pay for your resource. It is valid when payment_timing is Prepaid. Valid values: [1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36].
* `reservation_time_unit` - (Required) The reservation time unit that you will pay for your resource. It is valid when payment_timing is Prepaid. The value can only be month currently, which is also the default value.

The `cds_disks` object supports the following:

* `cds_size_in_gb` - (Optional, ForceNew) The size(GB) of CDS.
* `snapshot_id` - (Optional, ForceNew) Snapshot ID of CDS.
* `storage_type` - (Optional, ForceNew) Storage type of the CDS.

The `ephemeral_disks` object supports the following:

* `size_in_gb` - (Optional, ForceNew) Size(GB) of the ephemeral disk.
* `storage_type` - (Optional, ForceNew) Storage type of the ephemeral disk. Available values are std1, hp1, cloud_hp1, local, sata, ssd. Default to cloud_hp1.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `auto_renew` - Whether to automatically renew.
* `create_time` - Create time of the instance.
* `expire_time` - Expire time of the instance.
* `internal_ip` - Internal IP assigned to the instance.
* `keypair_name` - Key pair name of the instance.
* `network_capacity_in_mbps` - Public network bandwidth(Mbps) of the instance.
* `placement_policy` - The placement policy of the instance, which can be default or dedicatedHost.
* `public_ip` - Public IP
* `status` - Status of the instance.
* `vpc_id` - VPC ID of the instance.


## Import

BCC instance can be imported, e.g.

```hcl
$ terraform import baiducloud_instance.my-server id
```

