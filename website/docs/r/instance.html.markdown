---
layout: "baiducloud"
subcategory: "Baidu Cloud Compute (BCC)"
page_title: "BaiduCloud: baiducloud_instance"
sidebar_current: "docs-baiducloud-resource-instance"
description: |-
  Use this resource to create bcc instance.
---

# baiducloud_instance

Use this resource to create a BCC instance.

~> **NOTE:** The terminate operation of bcc does NOT take effect immediately，maybe takes for several minites.

## Example Usage

### Create instance
```hcl
resource "baiducloud_instance" "my-server" {
  image_id = "m-A4jJpFzi"
  name = "my-instance"
  availability_zone = "cn-bj-a"
  cpu_count = "2"
  memory_capacity_in_gb = "8"
  payment_timing = "Postpaid"
#  if create prepaid instance, please refer below
#  payment_timing = "Prepaid"
#  reservation = {
#    reservation_length = 1
#    reservation_time_unit =  "Month"
#  }
}
```

### Create instance by spec
> Use parameter *instance_spec* to replace instance_type, cpu_count, memory_capacity_in_gb, gpu_card, fpga_card, card_count, ephemeral_disks parameters，if the parameter *instance_spec* is specified, the above parameters are invalid
```hcl
resource "baiducloud_instance" "my-server" {
  image_id = "m-pUgPC9sJ"
  name = "my-instance"
  availability_zone = "cn-bj-d"
  instance_spec = "bcc.gr1.c1m4"
  payment_timing = "Postpaid"
#  if create prepaid instance, please refer below
#  payment_timing = "Prepaid"
#  reservation = {
#    reservation_length = 1
#    reservation_time_unit =  "Month"
#  }
}
```

## Argument Reference

The following arguments are supported:

* `image_id` - (Required) ID of the image to be used for the instance.
* `action` - (Optional) Start or stop the instance, which can only be start or stop, default start.
* `admin_pass` - (Optional, Sensitive) Password of the instance to be started. This value should be 8-16 characters, and English, numbers and symbols must exist at the same time. The symbols is limited to "!@#$%^*()".
* `auto_renew_time_length` - (Optional) The time length of automatic renewal. Effective only when `payment_timing` is `Prepaid`. Valid values are `1–9` when `auto_renew_time_unit` is `month` and `1–3` when it is `year`. Defaults to `1`. Due to API limitations, modifying this parameter after the auto-renewal rule is created will first delete the existing rule and then recreate it.
* `auto_renew_time_unit` - (Optional) Time unit of automatic renewal, the value can be month or year. The default value is empty, indicating no automatic renewal. It is valid only when the payment_timing is Prepaid.
* `availability_zone` - (Optional, ForceNew) Availability zone to start the instance in.
* `card_count` - (Optional) Count of the GPU cards or FPGA cards to be carried for the instance to be created, it is valid only when the gpu_card or fpga_card field is not empty.
* `cds_auto_renew` - (Optional) **This parameter is deprecated as CDS auto-renewal now aligns with the BCC instance.** Whether the cds is automatically renewed. It is valid when payment_timing is Prepaid. Default to false.
* `cds_disks` - (Optional) CDS disks of the instance.
* `cpu_count` - (Optional) Number of CPU cores to be created for the instance.
* `dedicate_host_id` - (Optional, ForceNew) The ID of dedicated host.
* `delete_cds_snapshot_flag` - (Optional, ForceNew) Whether to release the cds disk snapshots, default to false. It is effective only when the related_release_flag is true.
* `deploy_set_ids` - (Optional) Deploy set ids the instance belong to
* `description` - (Optional) Description of the instance.
* `enterprise_security_groups` - (Optional) Enterprise security group ids of the instance.
* `ephemeral_disks` - (Optional) Ephemeral disks of the instance.
* `fpga_card` - (Optional, ForceNew) FPGA card of the instance.
* `gpu_card` - (Optional, ForceNew) GPU card of the instance.
* `hostname` - (Optional) Hostname of the instance.
* `instance_spec` - (Optional) spec
* `instance_type` - (Optional, ForceNew) Type of the instance to start. Available values are N1, N2, N3, N4, N5, C1, C2, S1, G1, F1. Default to N3.
* `is_open_hostname_domain` - (Optional) Whether to automatically generate hostname domain.
* `is_open_ipv6` - (Optional) Whether to enable IPv6 for the instance to be created. It can be enabled only when both the image and the subnet support IPv6. True means enabled, false means disabled, undefined means automatically adapting to the IPv6 support of the image and subnet.
* `keypair_id` - (Optional, ForceNew) Key pair id of the instance.
* `memory_capacity_in_gb` - (Optional) Memory capacity(GB) of the instance to be created.
* `name` - (Optional) Name of the instance. Support for uppercase and lowercase letters, numbers, Chinese and special characters, such as "-","_","/",".", the value must start with a letter, length 1-65.
* `payment_timing` - (Optional) Payment timing of billing, which can be Prepaid or Postpaid. The default is Postpaid. When switching to Prepaid, reservation length must be set. Switching to Postpaid takes effect immediately.
* `related_release_flag` - (Optional, ForceNew) Whether to release the eip and data disks mounted by the current instance. Can only be released uniformly or not. Default to false.
* `relation_tag` - (Optional, ForceNew) The new instance associated with existing Tags or not, default false. The Tags should already exit if set true
* `reservation` - (Optional) Reservation of the instance.
* `resource_group_id` - (Optional) Resource group Id of the instance.
* `root_disk_size_in_gb` - (Optional, ForceNew) System disk size(GB) of the instance to be created. The value range is [40,2048]GB, Default to 40GB, and more than 40GB is charged according to the cloud disk price. Note that the specified system disk size needs to meet the minimum disk space limit of the mirror used.
* `root_disk_storage_type` - (Optional, ForceNew) System disk storage type of the instance. Available values are enhanced_ssd_pl1, enhanced_ssd_pl2, cloud_hp1, premium_ssd, hp1, ssd, sata, hdd, local, sata, local-ssd, local-hdd, local-nvme. Default to cloud_hp1.
* `security_groups` - (Optional) Security group ids of the instance.
* `stop_with_no_charge` - (Optional) Whether to enable stopping charging after shutdown for postpaid instance without local disks. Defaults to false.
* `subnet_id` - (Optional) The subnet ID of VPC. The default subnet will be used when it is empty. The instance will restart after changing the subnet.
* `sync_eip_auto_renew_rule` - (Optional) Whether to synchronize the EIP's auto-renewal rule with that of the associated BCC instance. This setting applies during both the creation and deletion of the BCC's auto-renewal rule. Modifying this parameter alone does not trigger any change to the EIP's auto-renewal rule. Effective only when `payment_timing` is `Prepaid`. Defaults to `true`.
* `tags` - (Optional, ForceNew) Tags, do not support modify
* `user_data` - (Optional) User Data

The `cds_disks` object supports the following:

* `cds_size_in_gb` - (Optional, ForceNew) The size(GB) of CDS.
* `snapshot_id` - (Optional, ForceNew) Snapshot ID of CDS.
* `storage_type` - (Optional, ForceNew) Storage type of the CDS.

The `ephemeral_disks` object supports the following:

* `size_in_gb` - (Optional, ForceNew) Size(GB) of the ephemeral disk.
* `storage_type` - (Optional, ForceNew) Storage type of the ephemeral disk. Available values are std1, hp1, cloud_hp1, local, sata, ssd. Default to cloud_hp1.

The `reservation` object supports the following:

* `reservation_length` - (Required) The reservation length that you will pay for your resource. It is valid when payment_timing is Prepaid. Valid values: [1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36].
* `reservation_time_unit` - (Required) The reservation time unit that you will pay for your resource. It is valid when payment_timing is Prepaid. The value can only be month currently, which is also the default value.

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

