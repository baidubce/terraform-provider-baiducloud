---
layout: "baiducloud"
subcategory: "BCC"
page_title: "BaiduCloud: baiducloud_specs"
sidebar_current: "docs-baiducloud-datasource-specs"
description: |-
  Use this data source to query spec list.
---

# baiducloud_specs

Use this data source to query spec list.

~> **NOTE:** Since v1.16.2, the update of this datasource is not compatible with the old version, please read the following documents carefully, if your provider version >= v1.16.2, the datasource configuration needs to be updated accordingly in your .tf files
## Example Usage

```hcl
data "baiducloud_bcc_specs" "default" {
  zone_name = "cn-bj-d"
  output_file = "specs.json"

  filter {
    name = "cpu_count"
    values = ["^([1])$"]
  }
}

output "spec" {
  value = "${data.baiducloud_specs.default.specs}"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `output_file` - (Optional, ForceNew) BCC Flavor search result output file
* `zone_name` - (Optional, ForceNew) Zone name

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `specs` - Specs list
  * `cpu_count` - CPU count
  * `cpu_ghz` - CPU frequency
  * `cpu_model` - CPU model name
  * `ephemeral_disk_count` - Count of ephemeral disk
  * `ephemeral_disk_in_gb` - Ephemeral disk size in gb
  * `ephemeral_disk_type` - Type of ephemeral disk
  * `fpga_card_count` - Count of FPGA card
  * `fpga_card_type` - Type of FPGA card
  * `gpu_card_count` - Count of gpu card
  * `gpu_card_type` - Type of gpu card
  * `group_id` - Group id
  * `memory_capacity_in_gb` - Memory capacity in GB
  * `network_bandwidth` - Network bandwidth
  * `network_package` - Network package
  * `product_type` - Product type
  * `spec_id` - Spec id
  * `spec` - Spec name
  * `zone_name` - Zone name


