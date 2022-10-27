---
layout: "baiducloud"
subcategory: "BBC"
page_title: "BaiduCloud: baiducloud_bbc_flavors"
sidebar_current: "docs-baiducloud-datasource-bbc_flavors"
description: |-
  Use this data source to query BBC flavors list.
---

# baiducloud_bbc_flavors

Use this data source to query BBC flavors list.

## Example Usage

```hcl
data "baiducloud_bbc_flavors" "bbc_flavors" {

}

output "flavors" {
 value = "${data.baiducloud_bbc_flavors.bbc_flavors.flavors}"
}

```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `output_file` - (Optional, ForceNew) Flavor search result output file

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `flavors` - flavor list
  * `cpu_count` - cpu count
  * `cpu_type` - cpu type
  * `disk` - Disk information, including SSD and SATA disks
  * `flavor_id` - flavor id
  * `memory_capacity_in_gb` - Memory capacity in GB
  * `network_card` - Network device information
  * `others` - Additional information included in the flavor


