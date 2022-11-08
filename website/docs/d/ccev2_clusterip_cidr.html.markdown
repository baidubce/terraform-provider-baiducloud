---
layout: "baiducloud"
subcategory: "Cloud Container Engine v2 (CCEv2)"
page_title: "BaiduCloud: baiducloud_ccev2_clusterip_cidr"
sidebar_current: "docs-baiducloud-datasource-ccev2_clusterip_cidr"
description: |-
  Use this data source to recommend ccev2 cluster IP CIDR.
---

# baiducloud_ccev2_clusterip_cidr

Use this data source to recommend ccev2 cluster IP CIDR.

## Example Usage

```hcl
data "baiducloud_ccev2_clusterip_cidr" "default" {
  vpc_cidr = var.vpc_cidr
  container_cidr = var.container_cidr
  cluster_max_service_num = 32
  private_net_cidrs = ["172.16.0.0/12",]
  ip_version = "ipv4"
  output_file = "${path.cwd}/recommendClusterIPCidr.txt"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_max_service_num` - (Optional) Max service number in the cluster
* `container_cidr_ipv6` - (Optional) Container CIDR IPv6
* `container_cidr` - (Optional) Container CIDR
* `ip_version` - (Optional) IP Version
* `output_file` - (Optional) Result output file
* `private_net_cidrs_ipv6` - (Optional) Private Net CIDRs IPv6
* `private_net_cidrs` - (Optional) Private Net CIDRs
* `vpc_cidr_ipv6` - (Optional) VPC CIDR IPv6
* `vpc_cidr` - (Optional) VPC CIDR

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `err_msg` - Error message if an error occurs
* `is_success` - Is the recommendation request success
* `recommended_clusterip_cidrs_ipv6` - Recommend Cluster IP CIDR List IPv6
* `recommended_clusterip_cidrs` - Recommend Cluster IP CIDR List
* `request_id` - Request ID


