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
    k8s_version = "1.16.8"
    runtime_type = "docker"
    vpc_id = var.vpc_id

    master_config {
      master_type = "managed"
      cluster_ha = 1
      exposed_public = false
      cluster_blb_vpc_subnet_id = var.vpc_subnet_id
      managed_cluster_master_option {
        master_vpc_subnet_zone = "zoneA"
      }
    }
    container_network_config  {
      mode = "kubenet"
      lb_service_vpc_subnet_id = var.vpc_subnet_id
      cluster_pod_cidr = var.cluster_pod_cidr
      cluster_ip_service_cidr = var.cluster_ip_service_cidr
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

The `cluster_spec` object supports the following:

* `cluster_delete_option` - (Optional) Cluster Delete Option
* `cluster_name` - (Optional) Cluster Name
* `cluster_type` - (Optional) Cluster Type
* `container_network_config` - (Optional) Container Network Config
* `description` - (Optional) Cluster Description
* `k8s_custom_config` - (Optional) Cluster k8s custom config
* `k8s_version` - (Optional) Kubernetes Version
* `master_config` - (Optional) Cluster Master Config
* `plugins` - (Optional) Plugin List
* `runtime_type` - (Optional) Container Runtime Type
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
* `ip_version` - (Optional) IP Version
* `kube_proxy_mode` - (Optional) KubeProxy Mode
* `lb_service_vpc_subnet_id` - (Optional) LB Service VPC Sunnet ID
* `max_pods_per_node` - (Optional) Max pod number in one node 
* `mode` - (Optional) Network MOde
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
* `cluster_ha` - (Optional) Number of master nodes
* `exposed_public` - (Optional) Whether exposed to public network
* `managed_cluster_master_option` - (Optional) Managed cluster master option
* `master_type` - (Optional) Master Type

The `managed_cluster_master_option` object supports the following:

* `master_vpc_subnet_zone` - (Optional) Master VPC Sunbet Zone

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cluster_status` - Statue of the cluster
* `created_at` - Create time of the cluster
* `masters` - Master machines of the cluster
  * `created_at` - Instance create time
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


