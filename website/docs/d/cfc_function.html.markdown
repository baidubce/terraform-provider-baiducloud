---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_cfc_function"
sidebar_current: "docs-baiducloud-datasource-cfc_function"
description: |-
  Use this data source to get a function.
---

# baiducloud_cfc_function

Use this data source to get a function.

## Example Usage

```hcl
data "baiducloud_cfc_function" "default" {
   function_name = "terraform-create"
}

output "function" {
 value = "${data.baiducloud_cfc_function.default}"
}
```

## Argument Reference

The following arguments are supported:

* `function_name` - (Required, ForceNew) CFC function name, length must be between 1 and 64 bytes
* `qualifier` - (Optional) Function search qualifier

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `code_id` - CFC Function code id
* `code_sha256` - Function code sha256
* `code_size` - Function code size
* `code_storage` - CFC Code storage information
* `code_zip_file` - CFC Function Code base64-encoded data
* `commit_id` - Function commit id
* `description` - Function description
* `environment` - CFC Function environment variables
* `function_arn` - The same as function brn
* `function_brn` - Function brn
* `handler` - CFC Function execution handler
* `last_modified` - The same as update_time
* `log_bos_dir` - Log save dir if log type is bos
* `log_type` - Log save type, support bos/none
* `memory_size` - CFC Function memory size, should be an integer multiple of 128
* `region` - Function region
* `reserved_concurrent_executions` - Function reserved concurrent executions, support [0-90]
* `role` - Function exec role
* `runtime` - CFC Function runtime
* `source_tag` - CFC Function source tag
* `time_out` - Function time out, support [1, 300]s
* `uid` - Function user uid
* `update_time` - Last update time
* `version` - Function version, should only be $LATEST
* `vpc_config` - Function VPC Config
  * `security_group_ids` - CFC Function binded Security group list
  * `subnet_ids` - CFC Function bined VPC Subnet id list


