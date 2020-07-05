---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_cce_cluster"
sidebar_current: "docs-baiducloud-resource-cce_cluster"
description: |-
  Use this resource to get information about a CCE Cluster.
---

# baiducloud_cce_cluster

Use this resource to get information about a CCE Cluster.

~> **NOTE:** The terminate operation of cce does NOT take effect immediately，maybe takes for several minites.

## Example Usage

```hcl
resource "baiducloud_cce_cluster" "my-cluster" {
  cluster_name        = "test-cce-cluster"
  main_available_zone = "zoneA"
  container_net       = "172.16.0.0/16"
  deploy_mode		  = "BCC"
  master_config {
    instance_type     = "10"
    cpu               = 4
    memory            = 8
    image_type        = "common"
    logical_zone      = "zoneA"
    subnet_uuid       = baiducloud_subnet.defaultA.id
    security_group_id = baiducloud_security_group.defualt.id
    product_type      = "postpay"
    image_id          = data.baiducloud_images.default.images.0.id
  }
  worker_config {
    count = {
      "zoneA" : 2
    }
    instance_type = "10"
    cpu           = 1
    memory        = 2
    subnet_uuid = {
      "zoneA" : baiducloud_subnet.defaultA.id
      "zoneB" : baiducloud_subnet.defaultB.id
    }
    security_group_id = baiducloud_security_group.defualt.id
    product_type      = "postpay"
    image_id          = data.baiducloud_images.default.images.0.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_name` - (Required, ForceNew) Name of the Cluster. Support for uppercase and lowercase letters, numbers, Chinese and special characters, such as "-","_","/",".", the value must start with a letter, length 1-65.
* `container_net` - (Required, ForceNew) Container network type of the cce cluster.
* `worker_config` - (Required) Worker node config of the cce cluster.
* `advanced_options` - (Optional, ForceNew) Advanced options configuration of the cce cluster.
* `comment` - (Optional, ForceNew) Comment information of the cce cluster.
* `delete_eip_cds` - (Optional) Whether to delete the eip and cds, default to true.
* `delete_snapshots` - (Optional) Whether to delete the snapshots, default to true.
* `deploy_mode` - (Optional, ForceNew) Deployment mode of the cce cluster, which can only be BCC.
* `main_available_zone` - (Optional, ForceNew) Main available zone of the cce cluster, support zoneA, zoneB, etc.
* `master_config` - (Optional, ForceNew) Master config of the cce cluster.
* `version` - (Optional, ForceNew) Kubernetes version of the cce cluster.

The `advanced_options` object supports the following:

* `cni_mode` - (Optional, ForceNew) Mode of the container network interface, which can only be cni or kubenet.
* `cni_type` - (Optional, ForceNew) Type of the container network interface, which can be VPC_ROUTE_AUTODETECT, VPC_SECONDARY_IP_VETH.
* `dns_mode` - (Optional, ForceNew) Mode of the dns, which can be coreDNS or kubeDNS.
* `kube_proxy_mode` - (Optional, ForceNew) Mode of kube-proxy, which can only be iptables or ipvs.
* `max_pod_num` - (Optional, ForceNew) Maximum number of pods in a node.

The `master_config` object supports the following:

* `cpu` - (Required, ForceNew) Number of cpu cores.
* `image_id` - (Required, ForceNew) Image id of the master node.
* `image_type` - (Required, ForceNew) Image type of the master node.
* `instance_type` - (Required, ForceNew) Instance type of the master node.
* `logical_zone` - (Required, ForceNew) Logical zone of the master node.
* `memory` - (Required, ForceNew) Memory capacity(GB) of the master node.
* `security_group_id` - (Required, ForceNew) ID of the security group.
* `subnet_uuid` - (Required, ForceNew) Subnet uuid of the master node.
* `admin_pass` - (Optional, ForceNew) Password of the worker node.
* `auto_renew_time_unit` - (Optional, ForceNew) Time unit of automatic renewal, the default value is month, It is valid only when the product_type is prepay and auto_renew is true.
* `auto_renew_time` - (Optional, ForceNew) The time length of automatic renewal. It is valid only when the product_type is prepay and auto_renew is true.
* `auto_renew` - (Optional, ForceNew) Whether the master is automatically renewed.
* `gpu_card` - (Optional, ForceNew) Gpu card of the master node.
* `gpu_count` - (Optional, ForceNew) Count of gpu card.
* `product_type` - (Optional, ForceNew) Product type of the master node, which can be postpay or prepay.
* `purchase_length` - (Optional, ForceNew) Purchase duration of the master node.
* `root_disk_size_in_gb` - (Optional, ForceNew) System disk size(GB) of the master node.
* `root_disk_storage_type` - (Optional, ForceNew) System disk storage type of the master node.

The `worker_config` object supports the following:

* `count` - (Required) Count of the worker node.
* `cpu` - (Required, ForceNew) Number of cpu cores for the worker node.
* `image_id` - (Required, ForceNew) Image id of the worker node.
* `image_type` - (Required, ForceNew) Image type of the worker node, which can be common, custom, gpuBccImage, gpuBccCustom, sharing.
* `instance_type` - (Required, ForceNew) Instance type of the worker node.
* `memory` - (Required, ForceNew) Memory size of the worker node.
* `security_group_id` - (Required, ForceNew) ID of the security group.
* `subnet_uuid` - (Required) Subnet uuid of the worker node.
* `admin_pass` - (Optional, ForceNew) Password of the worker node.
* `auto_renew_time_unit` - (Optional, ForceNew) Time unit of automatic renewal, the default value is month, It is valid only when the product_type is prepay and auto_renew is true.
* `auto_renew_time` - (Optional, ForceNew) The time length of automatic renewal. It is valid only when the product_type is prepay and auto_renew is true.
* `auto_renew` - (Optional, ForceNew) Whether the worker is automatically renewed.
* `cds_disks` - (Optional, ForceNew) CDS disks of the worker node.
* `eip` - (Optional, ForceNew) Eip of the worker node.
* `gpu_card` - (Optional, ForceNew) Gpu card of the worker node.
* `gpu_count` - (Optional, ForceNew) Gpu count of the worker node.
* `product_type` - (Optional, ForceNew) Product type of the worker node, which can be postpay or prepay.
* `purchase_length` - (Optional, ForceNew) Purchase duration of the worker node.
* `root_disk_size_in_gb` - (Optional, ForceNew) System disk size(GB) of the worker node.
* `root_disk_storage_type` - (Optional, ForceNew) System disk storage type of the worker node.

The `cds_disks` object supports the following:

* `disk_size_in_gb` - (Required, ForceNew) Volume of disk in GB. Default is 0.
* `volume_type` - (Required, ForceNew) Types of disk，available values: CLOUD_PREMIUM and CLOUD_SSD.
* `snapshot_id` - (Optional, ForceNew) Data disk snapshot ID.

The `eip` object supports the following:

* `bandwidth_in_mbps` - (Required, ForceNew) Eip bandwidth(Mbps) of the worker node.
* `sub_product_type` - (Required, ForceNew) Eip product type of the worker node, which can be bandwidth or netraffic.
* `eip_name` - (Optional, ForceNew) Eip name of the worker node.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cluster_node_num` - Number of nodes in the cluster.
* `cluster_uuid` - UUID of cce cluster.
* `create_start_time` - Create time of the cce cluster.
* `delete_time` - Delete time of the cce cluster.
* `has_prepay` - Whether to include prepaid nodes.
* `instance_mode` - Instance mode of the cce cluster.
* `master_vm_count` - Number of virtual machines in the master node of the cce cluster.
* `master_zone_subnet_map` - Availability zone of master node.
* `region` - Region of the cce cluster.
* `status` - Status of the cce cluster.
* `vpc_cidr` - VPC cidr of the cce cluster.
* `vpc_id` - VPC id of the cce cluster.
* `vpc_uuid` - VPC uuid of the cce cluster.
* `worker_instances_list` - List of the worker instances.
  * `available_zone` - Available zone of the instance.
  * `eip` - Eip of the instance.
  * `instance_id` - ID of the instance.
  * `status` - Status of the instance.
* `zone_subnet_map` - Subnet of the zone.


