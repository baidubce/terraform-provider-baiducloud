---
layout: "baiducloud"
subcategory: "Virtual private Cloud (VPC)"
page_title: "BaiduCloud: baiducloud_subnets"
sidebar_current: "docs-baiducloud-datasource-subnets"
description: |-
  Use this data source to query subnet list.
---

# baiducloud_subnets

Use this data source to query subnet list.

## Example Usage

```hcl
data "baiducloud_subnets" "default" {
 vpc_id = "vpc-y4p102r3mz6m"
}

output "subnets" {
 value = "${data.baiducloud_subnets.default.subnets}"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `output_file` - (Optional, ForceNew) Output file for saving result.
* `subnet_id` - (Optional, ForceNew) ID of the subnet.
* `subnet_type` - (Optional, ForceNew) Specify the subnet type for subnets.
* `vpc_id` - (Optional, ForceNew) VPC ID for subnets to retrieve.
* `zone_name` - (Optional, ForceNew) Specify the zone name for subnets.

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `subnets` - Result of the subnets.
  * `available_ip` - Available IP address of the subnet.
  * `cidr` - CIDR block of the subnet.
  * `description` - Description of the subnet.
  * `name` - Name of the subnet.
  * `subnet_id` - ID of the subnet.
  * `subnet_type` - Type of the subnet.
  * `tags` - Tags
  * `vpc_id` - VPC ID of the subnet.
  * `zone_name` - Zone name of the subnet.


