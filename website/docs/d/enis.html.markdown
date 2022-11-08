---
layout: "baiducloud"
subcategory: "Elastic Network Interface (ENI)"
page_title: "BaiduCloud: baiducloud_enis"
sidebar_current: "docs-baiducloud-datasource-enis"
description: |-
  Use this data source to query ENI list.
---

# baiducloud_enis

Use this data source to query ENI list.

## Example Usage

```hcl
data "baiducloud_enis" "default" {
  vpc_id      = "vpc-xxxxxx"
}

output "enis" {
 value = "${data.baiducloud_enis.default.enis}"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required) Vpc id which ENI belong to
* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `instance_id` - (Optional) Instance id the ENI bind
* `name` - (Optional) Name of ENI
* `output_file` - (Optional, ForceNew) ENI list result output file
* `private_ip_address` - (Optional) Eni private IP address

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `enis` - ENI list
  * `created_time` - ENI create time
  * `description` - Description of ENI
  * `eni_id` - ENI ID
  * `enterprise_security_group_ids` - ENI enterprise security group IDs
  * `instance_id` - Instance id which ENI bind
  * `mac_address` - ENI Mac Address
  * `name` - Name of ENI
  * `private_ip_set` - ENI private ip set
    * `primary` - True or false, true mean it is primary IP, it's private IP address can not modify, only one primary IP in a ENI
    * `private_ip_address` - Private IP address
    * `public_ip_address` - Public IP address
  * `security_group_ids` - ENI security group IDs
  * `status` - Status of ENI
  * `subnet_id` - Subnet ID which ENI belong to
  * `zone_name` - ENI Availability Zone Name


