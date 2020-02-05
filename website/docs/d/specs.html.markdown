---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_specs"
sidebar_current: "docs-baiducloud-datasource-specs"
description: |-
  Use this data source to query spec list.
---

# baiducloud_specs

Use this data source to query spec list.

## Example Usage

```hcl
data "baiducloud_specs" "default" {}

output "spec" {
  value = "${data.baiducloud_specs.default.specs}"
}
```

## Argument Reference

The following arguments are supported:

* `cpu_count` - (Optional, ForceNew) Useful cpu count of the search spec
* `instance_type` - (Optional, ForceNew) Instance type of the search spec
* `memory_size_in_gb` - (Optional, ForceNew) Useful memory size in GB of the search spec
* `name_regex` - (Optional, ForceNew) Regex pattern of the search spec name
* `output_file` - (Optional, ForceNew) Output file for saving result.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `specs` - Useful spec list, when create a bcc instance, suggest use instance_type/cpu_count/memory_capacity_in_gb as bcc instance parameters
  * `cpu_count` - Useful cpu count
  * `instance_type` - Useful instance type
  * `local_disk_size_in_gb` - Useful local disk size in GB
  * `memory_size_in_gb` - Useful memory size in GB
  * `name` - Spec name


