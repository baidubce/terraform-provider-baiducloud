---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_cfc_version"
sidebar_current: "docs-baiducloud-resource-cfc_version"
description: |-
  Provide a resource to publish a CFC Function Version.
---

# baiducloud_cfc_version

Provide a resource to publish a CFC Function Version.

## Example Usage

```hcl
resource "baiducloud_cfc_version" "default" {
  function_name       = "terraform-cfc"
  version_description = "terraformVersion"
}
```

```

## Argument Reference

The following arguments are supported:

* `function_name` - (Required, ForceNew) CFC function name, length must be between 1 and 64 bytes
* `code_sha256` - (Optional) Function code sha256
* `log_bos_dir` - (Optional) Log save dir if log type is bos
* `log_type` - (Optional) Log save type, support bos/none
* `version_description` - (Optional, ForceNew) Function version description

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `code_size` - Function code size
* `commit_id` - Function commit id
* `description` - Function description
* `environment` - CFC Function environment variables
* `function_arn` - The same as function brn
* `function_brn` - Function brn
* `handler` - CFC Function execution handler
* `last_modified` - The same as update_time
* `memory_size` - CFC Function memory size
* `region` - CFC Function bined VPC Subnet id list
* `role` - Function exec role
* `runtime` - CFC Function runtime
* `time_out` - Function time out
* `uid` - Function user uid
* `update_time` - Last update time
* `version` - Function version


