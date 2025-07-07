---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_bls_log_store"
subcategory: "Baidu Log Service (BLS)"
sidebar_current: "docs-baiducloud-resource-bls_log_store"
description: |-
  Provide a resource to create an BLS LogStore.
---

# baiducloud_bls_log_store

Provide a resource to create an BLS LogStore.

## Example Usage

```hcl
resource "baiducloud_bls_log_store" "default" {
  log_store_name   = "MyTest"
  retention        = 10

}
```

## Argument Reference

The following arguments are supported:

* `log_store_name` - (Required, ForceNew) name of log store
* `retention` - (Required) retention days of log store
* `creation_date_time` - (Computed) log store create date time
* `last_modified_time` - (Computed) log store last modified time


## Import

BLS LogStore can be imported, e.g.

```hcl
$ terraform import baiducloud_bls_log_store.default id
```

