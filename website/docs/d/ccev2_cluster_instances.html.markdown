---
layout: "baiducloud"
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


