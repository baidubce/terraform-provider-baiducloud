---
layout: "baiducloud"
subcategory: "Baidu Cloud Compute (BCC)"
page_title: "BaiduCloud: baiducloud_deploysets"
sidebar_current: "docs-baiducloud-datasource-deploysets"
description: |-
  Use this data source to query deploy set list.
---

# baiducloud_deploysets

Use this data source to query deploy set list.

## Example Usage

```hcl
data "baiducloud_deploysets" "default" {}

output "deploysets" {
 value = "${data.baiducloud_deploysets.default.deploysets}"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `output_file` - (Optional, ForceNew) deployset search result output file

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `deploy_sets` - Image list
  * `az_intstance_statis_list` - Availability Zone Instance Statistics List.
    * `bbc_instance_cnt` - Count of BBC instance which is in the deployset.
    * `bcc_instance_cnt` - Count of BCC instance which is in the deployset.
    * `instance_count` - Count of instance which is in the deployset.
    * `instance_total` - Total of instance which is in the deployset.
    * `zone_name` - Zone name of deployset.
  * `concurrency` - concurrency of deployset.
  * `deployset_id` - Id of deployset.
  * `desc` - Description of the deployset.
  * `name` - Name of the deployset.
  * `strategy` - Strategy of deployset.Available values are HOST_HA, RACK_HA and TOR_HA


