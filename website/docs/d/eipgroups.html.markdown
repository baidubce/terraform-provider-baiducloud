---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_eipgroups"
subcategory: "Elastic IP (EIP)"
sidebar_current: "docs-baiducloud-datasource-eipgroups"
description: |-
  Use this data source to query EIP group list.
---

# baiducloud_eipgroups

Use this data source to query EIP group list.

## Example Usage

```hcl
data "baiducloud_eipgroups" "default" {}

output "eip_groups" {
 value = "${data.baiducloud_eipgroups.default.eip_groups}"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `group_id` - (Optional, ForceNew) Id of Eip group
* `name` - (Optional, ForceNew) name of Eip group
* `output_file` - (Optional, ForceNew) Eipgroups search result output file
* `status` - (Optional, ForceNew) Eip group status

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `eip_groups` - Eip group list
  * `band_width_in_mbps` - Eip group band width in mbps
  * `billing_method` - Eip group billing method
  * `bw_bandwidth_in_mbps` - Eip group bw bandwidth in mbps
  * `bw_short_id` - Eip group bw short id
  * `create_time` - Eip group create time
  * `default_domestic_bandwidth` - Eip group default domestic bandwidth
  * `domestic_bw_bandwidth_in_mbps` - Eip group domestic bw bandwidth in mbps
  * `domestic_bw_short_id` - Eip group domestic bw short id
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
  * `expire_time` - Eip group expire time
  * `id` - Eip group id
  * `name` - Eip group name
  * `payment_timing` - Eip group payment timing
  * `region` - Eip group region
  * `route_type` - Eip group route type
  * `status` - Eip group status
  * `tags` - Tags


