---
layout: "baiducloud"
subcategory: "Cloud Container Engine v2 (CCEv2)"
page_title: "BaiduCloud: baiducloud_ccev2_instance_group"
sidebar_current: "docs-baiducloud-resource-ccev2_instance_group"
description: |-
  Use this resource to create a CCEv2 InstanceGroup.
---

# baiducloud_ccev2_instance_group

Use this resource to create a CCEv2 InstanceGroup.

~> **NOTE:** The create/update/delete operation of ccev2 does NOT take effect immediately，maybe takes for several minutes.

## Example Usage

```hcl
resource "baiducloud_ccev2_instance_group" "ccev2_instance_group_1" {
  spec {
    cluster_id = baiducloud_ccev2_cluster.default_custom.id
    replicas = var.instance_group_replica_1
    instance_group_name = "ig_1"
    instance_template {
      cce_instance_id = ""
      instance_name = "tf_ins_ig_1"
      cluster_role = "node"
      existed = false
      machine_type = "BCC"
      instance_type = "N3"
      vpc_config {
        vpc_id = baiducloud_vpc.default.id
        vpc_subnet_id = baiducloud_subnet.defaultA.id
        security_group_id = baiducloud_security_group.default.id
        available_zone = "zoneA"
      }
      deploy_custom_config {
        pre_user_script  = "ls"
        post_user_script = "date"
      }
      instance_resource {
        cpu = 4
        mem = 8
        root_disk_size = 40
        local_disk_size = 0
      }
      image_id = data.baiducloud_images.default.images.0.id
      instance_os {
        image_type = "System"
      }
      need_eip = false
      admin_password = "test123!YT"
      ssh_key_id = ""
      instance_charging_type = "Postpaid"
      runtime_type = "docker"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `spec` - (Required) Instance Group Spec

The `spec` object supports the following:

* `cluster_id` - (Required, ForceNew) Cluster ID of Instance Group
* `instance_group_name` - (Required, ForceNew) Name of Instance Group
* `instance_template` - (Required, ForceNew) Instance Spec of Instances in this Instance Group
* `replicas` - (Required) Number of instances in this Instance Group

The `instance_template` object supports the following:

* `admin_password` - (Optional) Admin Password
* `bbc_option` - (Optional) BBC Option
* `cce_instance_id` - (Optional) Instance ID
* `cluster_id` - (Optional) Cluster ID of this Instance
* `cluster_role` - (Optional) Cluster Role of Instance, Master or Nodes. Available Value: [master, node].
* `delete_option` - (Optional) Delete Option
* `deploy_custom_config` - (Optional) Deploy Custom Option
* `eip_option` - (Optional) EIP Option
* `existed_option` - (Optional) Existed Instance Option
* `existed` - (Optional) Is the instance existed
* `hpas_option` - (Optional) HPAS Option
* `image_id` - (Optional) Image ID
* `instance_charging_type` - (Optional) Instance charging type. Available Value: [Prepaid, Postpaid, bidding].
* `instance_group_id` - (Optional) Instance Group ID of this Instance
* `instance_group_name` - (Optional) Name of Instance Group
* `instance_name` - (Optional) Instance Name
* `instance_os` - (Optional) OS Config of the instance
* `instance_precharging_option` - (Optional) Instance Pre-charging Option
* `instance_resource` - (Optional) Instance Resource Config
* `instance_taints` - (Optional) Taint List
* `instance_type` - (Optional) Instance Type. Available Values: [N1, N2, N3, N4, N5, C1, C2, S1, G1, F1, HPAS].
* `labels` - (Optional) Labels List
* `machine_type` - (Optional) Machine Type. Available Values: [BCC, BBC, EBC, HPAS].
* `master_type` - (Optional) Master Type. Available Value: [managed, custom, serverless].
* `need_eip` - (Optional) Whether the instance need a EIP
* `runtime_type` - (Optional) Container Runtime Type. Available Value: [docker].
* `runtime_version` - (Optional) Container Runtime Version
* `ssh_key_id` - (Optional) SSH Key ID
* `tag_list` - (Optional) Tag List
* `vpc_config` - (Optional) VPC Config
* `cce_instance_priority` - Priority of this instance.

The `bbc_option` object supports the following:

* `raid_id` - (Optional) Disk Raid ID
* `reserve_data` - (Optional) Whether reserve data
* `sys_disk_size` - (Optional) System Disk Size

The `delete_option` object supports the following:

* `delete_cds_snapshot` - (Optional) Whether delete CDS snapshot
* `delete_resource` - (Optional) Whether delete resources
* `move_out` - (Optional) Whether move out the instance

The `deploy_custom_config` object supports the following:

* `docker_config` - (Optional) Docker Config Info
* `enable_cordon` - (Optional) Whether enable cordon
* `enable_resource_reserved` - (Optional) Whether to Enable Resource Quota
* `kube_reserved` - (Optional) Resource Quota
* `kubelet_root_dir` - (Optional) kubelet Data Directory
* `post_user_script` - (Optional) Script after deployment, base64 encoded
* `pre_user_script` - (Optional) Script before deployment, base64 encoded

The `docker_config` object supports the following:

* `bip` - (Optional) docker0 Network Bridge Network Segment
* `docker_data_root` - (Optional) Customized Docker Data Directory
* `docker_log_max_file` - (Optional) docker Log Max File
* `docker_log_max_size` - (Optional) docker Log Max Size
* `insecure_registries` - (Optional) Customized InsecureRegistries
* `registry_mirrors` - (Optional) Customized RegistryMirrors

The `eip_option` object supports the following:

* `eip_bandwidth` - (Optional) EIP Bandwidth
* `eip_charging_type` - (Optional) EIP Charging Type. Available Value: [ByTraffic, ByBandwidth].
* `eip_name` - (Optional) EIP Name

The `existed_option` object supports the following:

* `existed_instance_id` - (Optional) Existed Instance ID
* `rebuild` - (Optional) Whether re-install OS

The `hpas_option` object supports the following:

* `app_performance_level` - (Required) Performance level of the application. e.g., `10k`.
* `app_type` - (Required) Application type of the HPAS instance. e.g., `llama2_7B_train`.

The `instance_os` object supports the following:

* `image_name` - (Optional) Image Name
* `image_type` - (Optional) Image type. Available Value: [Integration, System, All, Custom, Sharing, GpuBccSystem, GpuBccCustom, BbcSystem, BbcCustom].
* `os_arch` - (Optional) OS arch
* `os_build` - (Optional) OS Build Time
* `os_name` - (Optional) OS name. Available Value: [CentOS, Ubuntu, Windows Server, Debian, opensuse].
* `os_type` - (Optional) OS type. Available Value: [linux, windows].
* `os_version` - (Optional) OS version

The `instance_precharging_option` object supports the following:

* `auto_renew_time_unit` - (Optional) Time unit for auto renew
* `auto_renew_time` - (Optional) Number of time unit for auto renew
* `auto_renew` - (Optional) Is Auto Renew
* `purchase_time` - (Optional) Time of purchase

The `instance_resource` object supports the following:

* `cds_list` - (Optional) CDS List
* `cpu` - (Optional) CPU cores
* `gpu_count` - (Optional) GPU Number
* `gpu_type` - (Optional) GPU Type. Available Value: [V100-32, V100-16, P40, P4, K40, DLCard].
* `local_disk_size` - (Optional) Local disk size
* `mem` - (Optional) memory GB
* `node_cpu_quota` - (Optional) Node cpu quota
* `node_mem_quota` - (Optional) Node memory quota
* `root_disk_size` - (Optional) Root disk size
* `root_disk_type` - (Optional) Root disk type. Available Value: [std1, hp1, cloud_hp1, local, sata, ssd, hdd].

The `cds_list` object supports the following:

* `cds_size` - (Optional) CDS Size
* `path` - (Optional) CDS path
* `snapshot_id` - (Optional) Snap shot ID
* `storage_type` - (Optional) Storage Type. Available Value: [std1, hp1, cloud_hp1, local, sata, ssd, hdd].

The `instance_taints` object supports the following:

* `effect` - (Optional) Taint Effect. Available Value: [NoSchedule, PreferNoSchedule, NoExecute].
* `key` - (Optional) Taint Key
* `time_added` - (Optional) Taint Added Time. Format RFC3339
* `value` - (Optional) Taint Value

The `tag_list` object supports the following:

* `tag_key` - (Optional) Tag Key
* `tag_value` - (Optional) Tag Value

The `vpc_config` object supports the following:

* `available_zone` - (Optional) Available Zone. Available Value: [zoneA, zoneB, zoneC, zoneD, zoneE, zoneF].
* `security_group_id` - (Optional) Security Group ID
* `security_group_type` - (Optional) Security Group type. Available Values: [normal, enterprise]. Default: `normal`
* `vpc_id` - (Optional) VPC ID
* `vpc_subnet_cidr_ipv6` - (Optional) VPC Sunbet CIDR IPv6
* `vpc_subnet_cidr` - (Optional) VPC Subnet CIDR
* `vpc_subnet_id` - (Optional) VPC Subnet ID
* `vpc_subnet_type` - (Optional) VPC Subnet type. Available Value: [BCC, BCC_NAT, BBC].

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `nodes` - All detail info of nodes in this instance group
  * `created_at` - Instance create time
  * `instance_spec` - Instance specification
    * `admin_password` - Admin Password
    * `bbc_option` - BBC Option
      * `raid_id` - Disk Raid ID
      * `reserve_data` - Whether reserve data
      * `sys_disk_size` - System Disk Size
    * `cce_instance_id` - Instance ID
    * `cce_instance_priority` - Priority of this instance.
    * `cluster_id` - Cluster ID of this Instance
    * `cluster_role` - Cluster Role of Instance, Master or Nodes. Available Value: [master, node].
    * `delete_option` - Delete Option
      * `delete_cds_snapshot` - Whether delete CDS snapshot
      * `delete_resource` - Whether delete resources
      * `move_out` - Whether move out the instance
    * `deploy_custom_config` - Deploy Custom Option
      * `docker_config` - Docker Config Info
        * `bip` - docker0 Network Bridge Network Segment
        * `docker_data_root` - Customized Docker Data Directory
        * `docker_log_max_file` - docker Log Max File
        * `docker_log_max_size` - docker Log Max Size
        * `insecure_registries` - Customized InsecureRegistries
        * `registry_mirrors` - Customized RegistryMirrors
      * `enable_cordon` - Whether enable cordon
      * `enable_resource_reserved` - Whether to Enable Resource Quota
      * `kube_reserved` - Resource Quota
      * `kubelet_root_dir` - kubelet Data Directory
      * `post_user_script` - Script after deployment, base64 encoded
      * `pre_user_script` - Script before deployment, base64 encoded
    * `eip_option` - EIP Option
      * `eip_bandwidth` - EIP Bandwidth
      * `eip_charging_type` - EIP Charging Type. Available Value: [ByTraffic, ByBandwidth].
      * `eip_name` - EIP Name
    * `existed_option` - Existed Instance Option
      * `existed_instance_id` - Existed Instance ID
      * `rebuild` - Whether re-install OS
    * `existed` - Is the instance existed
    * `hpas_option` - HPAS Option
    * `image_id` - Image ID
    * `instance_charging_type` - Instance charging type. Available Value: [Prepaid, Postpaid, bidding].
    * `instance_group_id` - Instance Group ID of this Instance
    * `instance_group_name` - Name of Instance Group
    * `instance_name` - Instance Name
    * `instance_os` - OS Config of the instance
      * `image_name` - Image Name
      * `image_type` - Image type. Available Value: [Integration, System, All, Custom, Sharing, GpuBccSystem, GpuBccCustom, BbcSystem, BbcCustom].
      * `os_arch` - OS arch
      * `os_build` - OS Build Time
      * `os_name` - OS name. Available Value: [CentOS, Ubuntu, Windows Server, Debian, opensuse].
      * `os_type` - OS type. Available Value: [linux, windows].
      * `os_version` - OS version
    * `instance_precharging_option` - Instance Pre-charging Option
      * `auto_renew_time_unit` - Time unit for auto renew
      * `auto_renew_time` - Number of time unit for auto renew
      * `auto_renew` - Is Auto Renew
      * `purchase_time` - Time of purchase
    * `instance_resource` - Instance Resource Config
      * `cds_list` - CDS List
        * `cds_size` - CDS Size
        * `path` - CDS path
        * `snapshot_id` - Snap shot ID
        * `storage_type` - Storage Type. Available Value: [std1, hp1, cloud_hp1, local, sata, ssd, hdd].
      * `cpu` - CPU cores
      * `gpu_count` - GPU Number
      * `gpu_type` - GPU Type. Available Value: [V100-32, V100-16, P40, P4, K40, DLCard].
      * `local_disk_size` - Local disk size
      * `mem` - memory GB
      * `node_cpu_quota` - Node cpu quota
      * `node_mem_quota` - Node memory quota
      * `root_disk_size` - Root disk size
      * `root_disk_type` - Root disk type. Available Value: [std1, hp1, cloud_hp1, local, sata, ssd, hdd].
    * `instance_taints` - Taint List
      * `effect` - Taint Effect. Available Value: [NoSchedule, PreferNoSchedule, NoExecute].
      * `key` - Taint Key
      * `time_added` - Taint Added Time. Format RFC3339
      * `value` - Taint Value
    * `instance_type` - Instance Type. Available Values: [N1, N2, N3, N4, N5, C1, C2, S1, G1, F1, HPAS].
    * `labels` - Labels List
    * `machine_type` - Machine Type. Available Values: [BCC, BBC, EBC, HPAS].
    * `master_type` - Master Type. Available Value: [managed, custom, serverless].
    * `need_eip` - Whether the instance need a EIP
    * `runtime_type` - Container Runtime Type. Available Value: [docker].
    * `runtime_version` - Container Runtime Version
    * `ssh_key_id` - SSH Key ID
    * `tag_list` - Tag List
      * `tag_key` - Tag Key
      * `tag_value` - Tag Value
    * `vpc_config` - VPC Config
      * `available_zone` - Available Zone. Available Value: [zoneA, zoneB, zoneC, zoneD, zoneE, zoneF].
      * `security_group_id` - Security Group ID
      * `vpc_id` - VPC ID
      * `vpc_subnet_cidr_ipv6` - VPC Sunbet CIDR IPv6
      * `vpc_subnet_cidr` - VPC Subnet CIDR
      * `vpc_subnet_id` - VPC Subnet ID
      * `vpc_subnet_type` - VPC Subnet type. Available Value: [BCC, BCC_NAT, BBC].
  * `instance_status` - Instance status
    * `instance_phase` - Instance Phase
    * `machine_status` - Machine status
    * `machine` - Machine info
      * `eip` - EIP
      * `instance_id` - Instance ID
      * `mount_list` - Mount List of Machine
      * `order_id` - Order ID
      * `vpc_ip_ipv6` - VPC IPv6
      * `vpc_ip` - VPC IP
  * `updated_at` - Instance update time
* `status` - Instance Group Status
  * `ready_replicas` - Number of instances in RUNNING


