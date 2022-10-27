---
layout: "baiducloud"
subcategory: "APPBLB"
page_title: "BaiduCloud: baiducloud_appblbs"
sidebar_current: "docs-baiducloud-datasource-appblbs"
description: |-
  Use this data source to query APPBLB list.
---

# baiducloud_appblbs

Use this data source to query APPBLB list.

## Example Usage

```hcl
data "baiducloud_appblbs" "default" {
 name = "myLoadBalance"
}

output "blbs" {
 value = "${data.baiducloud_appblbs.default.blbs}"
}
```

## Argument Reference

The following arguments are supported:

* `address` - (Optional) Address ip of the LoadBalance instance to be queried
* `bcc_id` - (Optional) ID of the BCC instance bound to the LoadBalance
* `blb_id` - (Optional) ID of the LoadBalance instance to be queried
* `exactly_match` - (Optional) Whether the query condition is an exact match or not, default false
* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `name` - (Optional) Name of the LoadBalance instance to be queried
* `output_file` - (Optional, ForceNew) Query result output file path

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `appblbs` - A list of Application LoadBalance Instance
  * `address` - LoadBalance instance's service IP, instance can be accessed through this IP
  * `blb_id` - LoadBalance instance's ID
  * `cidr` - Cidr of the network where the LoadBalance instance reside
  * `create_time` - LoadBalance instance's create time
  * `description` - LoadBalance instance's description
  * `listener` - List of listeners mounted under the instance
    * `port` - Listening port
    * `type` - Listening protocol type
  * `name` - LoadBalance instance's name
  * `public_ip` - LoadBalance instance's public ip
  * `release_time` - LoadBalance instance's auto release time
  * `status` - LoadBalance instance's status
  * `subnet_cidr` - Cidr of the subnet which the LoadBalance instance belongs
  * `subnet_id` - The subnet ID to which the LoadBalance instance belongs
  * `subnet_name` - The subnet name to which the LoadBalance instance belongs
  * `tags` - Tags
  * `vpc_id` - The VPC short ID to which the LoadBalance instance belongs
  * `vpc_name` - The VPC name to which the LoadBalance instance belongs


