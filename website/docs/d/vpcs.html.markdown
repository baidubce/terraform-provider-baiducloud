---
layout: "baiducloud"
subcategory: "VPC"
page_title: "BaiduCloud: baiducloud_vpcs"
sidebar_current: "docs-baiducloud-datasource-vpcs"
description: |-
  Use this data source to query vpc list.
---

# baiducloud_vpcs

Use this data source to query vpc list.

## Example Usage

```hcl
data "baiducloud_vpcs" "default" {
    name="test-vpc"
}

output "cidr" {
  value = "${data.baiducloud_vpcs.default.vpcs.0.cidr}"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `name` - (Optional, ForceNew) Name of the specific VPC to retrieve.
* `output_file` - (Optional, ForceNew) Output file for saving result.
* `vpc_id` - (Optional, ForceNew) ID of the specific VPC to retrieve.

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `vpcs` - Result of VPCs.
  * `cidr` - CIDR block of the VPC.
  * `description` - Description of the VPC.
  * `is_default` - Specify if it is the default VPC.
  * `name` - Name of the VPC.
  * `route_table_id` - Route table ID of the VPC.
  * `secondary_cidrs` - The secondary cidr list of the VPC. They will not be repeated.
  * `tags` - Tags
  * `vpc_id` - ID of the VPC.


