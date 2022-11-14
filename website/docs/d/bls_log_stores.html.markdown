---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_bls_log_stores"
sidebar_current: "docs-baiducloud-datasource-bls_log_stores"
description: |-
  Use this data source to query bls log stores .
---

# baiducloud_bls_log_stores

Use this data source to query bls log stores .

## Example Usage

```hcl
data "baiducloud_bls_log_stores" "default" {

}

output "log_stores" {
 	value = "${data.baiducloud_bls_log_stores.default.log_stores}"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `name_pattern` - (Optional, ForceNew) Log store namePattern
* `order_by` - (Optional, ForceNew) order field
* `order` - (Optional, ForceNew) search order
* `output_file` - (Optional, ForceNew) log stores search result output file
* `page_no` - (Optional, ForceNew) number of page 
* `page_size` - (Optional, ForceNew) size of page
* `total_count` - (Optional) Total number of items

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `log_stores` - log store list
  * `log_store_name` - name of log store
  * `retention` - retention days of log store
  * `creation_date_time` - log store create date time
  * `last_modified_time` - log store last modified time


