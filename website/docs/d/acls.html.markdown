---
layout: "baiducloud"
subcategory: "Virtual private Cloud (VPC)"
page_title: "BaiduCloud: baiducloud_acls"
sidebar_current: "docs-baiducloud-datasource-acls"
description: |-
  Use this data source to query ACL list.
---

# baiducloud_acls

Use this data source to query ACL list.

## Example Usage

```hcl
data "baiducloud_acls" "default" {
 vpc_id = "vpc-y4p102r3mz6m"
}

output "acls" {
 value = "${data.baiducloud_acls.default.acls}"
}
```

## Argument Reference

The following arguments are supported:

* `acl_id` - (Optional) ID of the ACL to retrieve.
* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `output_file` - (Optional, ForceNew) Output file for saving result.
* `subnet_id` - (Optional) Subnet ID of the ACLs to retrieve.
* `vpc_id` - (Optional) VPC ID of the ACLs to retrieve.

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `acls` - List of the ACLs.
  * `acl_id` - ID of the ACL.
  * `action` - Action of the ACL.
  * `description` - Description of the ACL.
  * `destination_ip_address` - Destination IP address of the ACL.
  * `destination_port` - Destination port of the ACL.
  * `direction` - Direction of the ACL.
  * `position` - Position of the ACL.
  * `protocol` - Protocol of the ACL.
  * `source_ip_address` - Source IP address of the ACL.
  * `source_port` - Source port of the ACL.
  * `subnet_id` - Subnet ID of the ACL.


