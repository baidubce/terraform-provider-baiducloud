---
layout: "baiducloud"
subcategory: "Cloud Container Engine (CCE)"
page_title: "BaiduCloud: baiducloud_cce_kubeconfig"
sidebar_current: "docs-baiducloud-datasource-cce_kubeconfig"
description: |-
  Use this data source to get cce kubeconfig.
---

# baiducloud_cce_kubeconfig

Use this data source to get cce kubeconfig.

## Example Usage

```hcl
data "baiducloud_cce_kubeconfig" "default" {
	cluster_uuid = "c-NqYwWEhu"
}

output "kubeconfig" {
  value = "${data.baiducloud_cce_kubeconfig.default.data}"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_uuid` - (Required, ForceNew) UUID of the cce cluster.
* `config_type` - (Optional, ForceNew) Config type of the cce cluster.
* `output_file` - (Optional, ForceNew) Output file for saving result.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Data of the cce kubeconfig.


