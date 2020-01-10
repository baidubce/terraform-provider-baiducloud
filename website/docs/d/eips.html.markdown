---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_eips"
sidebar_current: "docs-baiducloud-datasource-eips"
description: |-
  Use this data source to query EIP list.
---

# baiducloud_eips

Use this data source to query EIP list.

## Example Usage

```hcl
data "baiducloud_eips" "default" {}

output "eips" {
 value = "${data.baiducloud_eips.default.eips}"
}
```

## Argument Reference

The following arguments are supported:

* `eip` - (Optional, ForceNew) Eip address
* `instance_id` - (Optional, ForceNew) Eip bind instance id
* `instance_type` - (Optional, ForceNew) Eip bind instance type
* `output_file` - (Optional, ForceNew) Eips search result output file
* `status` - (Optional, ForceNew) Eip status

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `eips` - Eip list
  * `bandwidth_in_mbps` - Eip bandwidth(Mbps)
  * `billing_method` - Eip billing method
  * `create_time` - Eip create time
  * `eip_instance_type` - Eip instance type
  * `eip` - Eip address
  * `expire_time` - Eip expire time
  * `name` - Eip name
  * `payment_timing` - Eip payment timing
  * `share_group_id` - Eip share group id
  * `status` - Eip status
  * `tags` - Tags


