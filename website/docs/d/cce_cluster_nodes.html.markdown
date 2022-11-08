---
layout: "baiducloud"
subcategory: "Cloud Container Engine (CCE)"
page_title: "BaiduCloud: baiducloud_cce_cluster_nodes"
sidebar_current: "docs-baiducloud-datasource-cce_cluster_nodes"
description: |-
  Use this data source to get cce cluster nodes.
---

# baiducloud_cce_cluster_nodes

Use this data source to get cce cluster nodes.

## Example Usage

```hcl
data "baiducloud_cce_cluster_nodes" "default" {
   cluster_uuid	 = "c-NqYwWEhu"
}

output "nodes" {
 value = "${data.baiducloud_cce_cluster_nodes.default.nodes}"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_uuid` - (Required, ForceNew) UUID of the cce cluster.
* `available_zone` - (Optional, ForceNew) Available zone of the cluster node.
* `instance_id` - (Optional, ForceNew) ID of the search instance.
* `instance_name_regex` - (Optional, ForceNew) Regex pattern of the search spec name.
* `instance_type` - (Optional, ForceNew) Type of the search instance.
* `subnet_id` - (Optional, ForceNew) ID of the subnet.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `nodes` - Result of the cluster nodes list.
  * `available_zone` - CDS disk size, should in [1, 32765], when snapshot_id not set, this parameter is required.
  * `blb` - BLB address of the node.
  * `cpu` - Number of cpu cores.
  * `create_time` - Create time of the instance.
  * `delete_time` - Delete time of the instance.
  * `disk_size` - Local disk size of the node.
  * `eip_bandwidth` - Eip bandwidth(Mbps) of the instance.
  * `eip` - Eip of the instance.
  * `expire_time` - Expire time of the instance.
  * `fix_ip` - Fix ip of the node, which is assigned in VPC.
  * `floating_ip` - Floating ip of the node.
  * `instance_id` - ID of the instance.
  * `instance_name` - Name of the instance.
  * `instance_type` - Type of the instance.
  * `instance_uuid` - UUID of the instance.
  * `memory` - Memory capacity(GB) of the instance.
  * `payment_method` - Payment method of the node.
  * `runtime_version` - Version of the instance runtime.
  * `status` - Status of the instance.
  * `subnet_id` - Subnet id of the instance.
  * `subnet_type` - Subnet type of the instance.
  * `sys_disk` - System disk size of the node.
  * `vpc_cidr` - VPC cidr of the instance.
  * `vpc_id` - VPC id of the instance.


