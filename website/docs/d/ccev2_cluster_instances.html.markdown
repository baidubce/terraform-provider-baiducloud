---
layout: "baiducloud"
subcategory: "Cloud Container Engine v2 (CCEv2)"
page_title: "BaiduCloud: baiducloud_ccev2_cluster_instances"
sidebar_current: "docs-baiducloud-datasource-ccev2_cluster_instances"
description: |-
  Use this data source to list instances of a cluster.
---

# baiducloud_ccev2_cluster_instances

Use this data source to list instances of a cluster.

## Example Usage

```hcl
data "baiducloud_ccev2_cluster_instances" "default" {
  cluster_id = baiducloud_ccev2_cluster.default_custom.id
  keyword_type = "instanceName"
  keyword = ""
  order_by = "instanceName"
  order = "ASC"
  page_no = 0
  page_size = 0
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) CCEv2 Cluster ID
* `keyword_type` - (Optional, ForceNew) Keyword type. Available Value: [instanceName, instanceID].
* `keyword` - (Optional, ForceNew) The search keyword
* `order_by` - (Optional, ForceNew) The field that used to order the list. Available Value: [instanceName, instanceID, createdAt].
* `order` - (Optional, ForceNew) Ascendant or descendant order. Available Value: [ASC, DESC].
* `page_no` - (Optional, ForceNew) Page number of query result
* `page_size` - (Optional, ForceNew) The size of every page

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `master_list` - The search result
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
    * `ehc_cluster_id` - EHC Cluster ID for instances
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
      * `purchase_time_unit` - Time unit for purchase
      * `purchase_time` - Time of purchase
    * `instance_resource` - Instance Resource Config
      * `cds_list` - CDS List
        * `cds_size` - CDS Size
        * `path` - CDS path
        * `snapshot_id` - Snap shot ID
        * `storage_type` - Storage Type. Available Value: [std1, hp1, cloud_hp1, local, sata, ssd, hdd].
      * `cpu` - CPU cores
      * `ephemeral_disk_list` - Ephemeral Disk List for instances
        * `disk_path` - Custom disk mount path for local disks
      * `gpu_count` - GPU Number
      * `gpu_type` - GPU Type. Available Value: [V100-32, V100-16, P40, P4, K40, DLCard].
      * `local_disk_size` - Local disk size
      * `machine_spec` - Machine specification for instances, e.g., 'llama_7B_train/10k'
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
* `nodes_list` - The search result
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
    * `ehc_cluster_id` - EHC Cluster ID for instances
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
      * `purchase_time_unit` - Time unit for purchase
      * `purchase_time` - Time of purchase
    * `instance_resource` - Instance Resource Config
      * `cds_list` - CDS List
        * `cds_size` - CDS Size
        * `path` - CDS path
        * `snapshot_id` - Snap shot ID
        * `storage_type` - Storage Type. Available Value: [std1, hp1, cloud_hp1, local, sata, ssd, hdd].
      * `cpu` - CPU cores
      * `ephemeral_disk_list` - Ephemeral Disk List for instances
        * `disk_path` - Custom disk mount path for local disks
      * `gpu_count` - GPU Number
      * `gpu_type` - GPU Type. Available Value: [V100-32, V100-16, P40, P4, K40, DLCard].
      * `local_disk_size` - Local disk size
      * `machine_spec` - Machine specification for instances, e.g., 'llama_7B_train/10k'
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
* `total_count` - The total count of the result


