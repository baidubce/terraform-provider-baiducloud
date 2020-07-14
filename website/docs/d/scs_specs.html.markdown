---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_scs_specs"
sidebar_current: "docs-baiducloud-datasource-scs_specs"
description: |-
  Use this data source to query scs specs list.
---

# baiducloud_scs_specs

Use this data source to query scs specs list.

## Example Usage

```hcl
data "data.baiducloud_scs_specs" "default" {}

output "specs" {
  value = "${data.baiducloud_scs_specs.default.specs}"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_type` - (Required, ForceNew) Type of the instance,  Available values are cluster, master_slave.
* `node_capacity` - (Required, ForceNew) Memory capacity(GB) of the instance node.
* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `output_file` - (Optional, ForceNew) Output file for saving result.

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `specs` - Useful spec list, when create a scs instance, suggest use node_type/cpu_num/instance_flavor/allowed_nodeNum_list as scs instance parameters
  * `node_capacity` - Memory capacity(GB) of the instance node.
  * `node_type` - Useful node type


