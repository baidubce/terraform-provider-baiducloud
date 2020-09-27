---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_ccev2_cluster"
sidebar_current: "docs-baiducloud-resource-ccev2_cluster"
description: |-
  Use this resource to create a CCEv2 cluster.
---

# baiducloud_ccev2_cluster

Use this resource to create a CCEv2 cluster.

## Example Usage

```hcl
resource "baiducloud_ccev2_cluster" "default_managed" {
  cluster_spec  {
    cluster_name = var.cluster_name
    cluster_type = "normal"
    k8s_version = "1.16.8"
    runtime_type = "docker"
    vpc_id = baiducloud_vpc.default.id
    plugins = ["core-dns", "kube-proxy"]
    master_config {
      master_type = "managed"
      cluster_ha = 2
      exposed_public = false
      cluster_blb_vpc_subnet_id = baiducloud_subnet.defaultA.id
      managed_cluster_master_option {
        master_vpc_subnet_zone = "zoneA"
      }
    }
    container_network_config  {
      mode = "kubenet"
      lb_service_vpc_subnet_id = baiducloud_subnet.defaultA.id
      node_port_range_min = 30000
      node_port_range_max = 32767
      max_pods_per_node = 64
      cluster_pod_cidr = var.cluster_pod_cidr
      cluster_ip_service_cidr = var.cluster_ip_service_cidr
      ip_version = "ipv4"
      kube_proxy_mode = "iptables"
    }
    cluster_delete_option {
      delete_resource = true
      delete_cds_snapshot = true
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_spec` - (Required, ForceNew) Specification of the cluster
* `master_specs` - (Optional, ForceNew) Specification of master nodes cluster

The `cluster_spec` object supports the following:

* `cluster_delete_option` - (Optional) Cluster Delete Option
* `cluster_name` - (Optional) Cluster Name
* `cluster_type` - (Optional) Cluster Type. Available Value: [normal].
* `container_network_config` - (Optional) Container Network Config
* `description` - (Optional) Cluster Description
* `k8s_custom_config` - (Optional) Cluster k8s custom config
* `k8s_version` - (Optional) Kubernetes Version. Available Value: [1.13.10, 1.16.8].
* `master_config` - (Optional) Cluster Master Config
* `plugins` - (Optional) Plugin List
* `runtime_type` - (Optional) Container Runtime Type. Available Value: [docker].
* `runtime_version` - (Optional) Container Runtime Version
* `vpc_cidr_ipv6` - (Optional) VPC CIDR IPv6
* `vpc_cidr` - (Optional) VPC CIDR
* `vpc_id` - (Optional) VPC ID

The `cluster_delete_option` object supports the following:

* `delete_cds_snapshot` - (Optional) Whether to delete CDS snapshot
* `delete_resource` - (Optional) Whether to delete resources

The `container_network_config` object supports the following:

* `cluster_ip_service_cidr_ipv6` - (Optional) Cluster Service ClusterIP CIDR IPv6
* `cluster_ip_service_cidr` - (Optional) Cluster Service ClusterIP CIDR 
* `cluster_pod_cidr_ipv6` - (Optional) Cluster Pod IP CIDR IPv6
* `cluster_pod_cidr` - (Optional) Cluster Pod IP CIDR
* `eni_security_group_id` - (Optional) ENI Security Group ID
* `eni_vpc_subnet_ids` - (Optional) ENI VPC Subnet ID
* `ip_version` - (Optional) IP Version. Available Value: [ipv4, ipv6, dualStack].
* `kube_proxy_mode` - (Optional) KubeProxy Mode. Available Value: [iptables, ipvs].
* `lb_service_vpc_subnet_id` - (Optional) LB Service VPC Sunnet ID
* `max_pods_per_node` - (Optional) Max pod number in one node 
* `mode` - (Optional) Network Mode. Available Value: [kubenet, vpc-cni, vpc-route-veth, vpc-route-ipvlan, vpc-route-auto-detect, vpc-secondary-ip-veth, vpc-secondary-ip-ipvlan, vpc-secondary-ip-auto-detect].
* `node_port_range_max` - (Optional) Node Port Service Port Range Max
* `node_port_range_min` - (Optional) Node Port Service Port Range Min

The `eni_vpc_subnet_ids` object supports the following:

* `zone_and_id` - (Optional) Available Zone and ENI ID

The `k8s_custom_config` object supports the following:

* `admission_plugins` - (Optional) custom Admission Plugins
* `etcd_data_path` - (Optional) etcd data directory
* `kube_api_burst` - (Optional) custom Kube API Burst
* `kube_api_qps` - (Optional) custom Kube API QPS
* `master_feature_gates` - (Optional) custom master Feature Gates
* `node_feature_gates` - (Optional) custom node Feature Gates
* `pause_image` - (Optional) custom PauseImage
* `scheduler_predicated` - (Optional) custom Scheduler Predicates
* `scheduler_priorities` - (Optional) custom SchedulerPriorities

The `master_config` object supports the following:

* `cluster_blb_vpc_subnet_id` - (Optional) Cluster BLB VPC Subnet ID
* `cluster_ha` - (Optional) Number of master nodes. Available Value: [1, 3, 5, 2(for serverless)].
* `exposed_public` - (Optional) Whether exposed to public network
* `managed_cluster_master_option` - (Optional) Managed cluster master option
* `master_type` - (Optional) Master Type. Available Value: [managed, custom, serverless].

The `managed_cluster_master_option` object supports the following:

* `master_vpc_subnet_zone` - (Optional) Master VPC Subnet Zone. Available Value: [zoneA, zoneB, zoneC, zoneD, zoneE, zoneF].

The `master_specs` object supports the following:

* `count` - (Required) Count of this type master
* `master_spec` - (Required) Count of this type master

The `master_spec` object supports the following:

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
* `image_id` - (Optional) Image ID
* `instance_charging_type` - (Optional) Instance charging type. Available Value: [Prepaid, Postpaid, bidding].
* `instance_group_id` - (Optional) Instance Group ID of this Instance
* `instance_group_name` - (Optional) Name of Instance Group
* `instance_name` - (Optional) Instance Name
* `instance_os` - (Optional) OS Config of the instance
* `instance_precharging_option` - (Optional) Instance Pre-charging Option
* `instance_resource` - (Optional) Instance Resource Config
* `instance_taints` - (Optional) Taint List
* `instance_type` - (Optional) Instance Type Available Value: [N1, N2, N3, N4, N5, C1, C2, S1, G1, F1].
* `labels` - (Optional) Labels List
* `machine_type` - (Optional) Machine Type. Available Value: [BCC, BBC, Metal].
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
* `vpc_id` - (Optional) VPC ID
* `vpc_subnet_cidr_ipv6` - (Optional) VPC Sunbet CIDR IPv6
* `vpc_subnet_cidr` - (Optional) VPC Subnet CIDR
* `vpc_subnet_id` - (Optional) VPC Subnet ID
* `vpc_subnet_type` - (Optional) VPC Subnet type. Available Value: [BCC, BCC_NAT, BBC].

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cluster_status` - Statue of the cluster
  * `cluster_blb` - Cluster BLB
  * `cluster_phase` - Cluster Phase
  * `node_num` - Cluster Node Number
* `created_at` - Create time of the cluster
* `masters` - Master machines of the cluster
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
    * `instance_type` - Instance Type Available Value: [N1, N2, N3, N4, N5, C1, C2, S1, G1, F1].
    * `labels` - Labels List
    * `machine_type` - Machine Type. Available Value: [BCC, BBC, Metal].
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
* `nodes` - Slave machines of the cluster
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
    * `instance_type` - Instance Type Available Value: [N1, N2, N3, N4, N5, C1, C2, S1, G1, F1].
    * `labels` - Labels List
    * `machine_type` - Machine Type. Available Value: [BCC, BBC, Metal].
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
* `updated_at` - Update time of the cluster


