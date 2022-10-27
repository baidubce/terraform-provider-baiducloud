---
layout: "baiducloud"
subcategory: "LOCALDNS"
page_title: "BaiduCloud: baiducloud_localdns_vpcs"
sidebar_current: "docs-baiducloud-datasource-localdns_vpcs"
description: |-
  Use this data source to query localdns VPCs.
---

# baiducloud_localdns_vpcs

Use this data source to query localdns VPCs.

## Example Usage

```hcl
data "baiducloud_localdns_vpcs" "default" {}

output "vpcs" {
   value = "${data.baiducloud_localdns_vpcs.default.bind_vpcs}"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) zone_id of the DNS privatezone 
* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `output_file` - (Optional, ForceNew) local dns vpc search result output file

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `bind_vpcs` - privatezone bind vpcs
  * `vpc_id` - bind vpc id
  * `vpc_name` - name of vpc
  * `vpc_region` - region of vpc


