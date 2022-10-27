---
layout: "baiducloud"
subcategory: "CCEv2"
page_title: "BaiduCloud: baiducloud_ccev2_container_cidr"
sidebar_current: "docs-baiducloud-datasource-ccev2_container_cidr"
description: |-
  Use this data source to recommend ccev2 container CIDR.
---

# baiducloud_ccev2_container_cidr

Use this data source to recommend ccev2 container CIDR.

## Example Usage

```hcl
data "baiducloud_ccev2_container_cidr" "default" {
  vpc_id = var.vpc_id
  vpc_cidr = var.vpc_cidr
  cluster_max_node_num = 16
  max_pods_per_node = 32
  private_net_cidrs = ["172.16.0.0/12",]
  k8s_version = "1.16.8"
  ip_version = "ipv4"
  output_file = "${path.cwd}/recommendContainerCidr.txt"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_max_node_num` - (Optional) Max node number in a cluster
* `ip_version` - (Optional) IP version
* `k8s_version` - (Optional) K8s Version
* `max_pods_per_node` - (Optional) Max pod number in a node
* `output_file` - (Optional) Eips search result output file
* `private_net_cidrs_ipv6` - (Optional) Private Net CIDR List IPv6
* `private_net_cidrs` - (Optional) Private Net CIDR List
* `vpc_cidr_ipv6` - (Optional) VPC CIDR IPv6
* `vpc_cidr` - (Optional) VPC CIDR
* `vpc_id` - (Optional) VPC ID

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `err_msg` - Error message if an error occures
* `is_success` - Whether the recommendation success
* `recommended_container_cidrs_ipv6` - Recomment Container CIDRs IPv6
* `recommended_container_cidrs` - Recomment Container CIDR
* `request_id` - Request ID


