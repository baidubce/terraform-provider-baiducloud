---
layout: "baiducloud"
subcategory: "Cloud Function Compute (CFC)"
page_title: "BaiduCloud: baiducloud_cfc_function"
sidebar_current: "docs-baiducloud-resource-cfc_function"
description: |-
  Provide a resource to create an CFC Function.
---

# baiducloud_cfc_function

Provide a resource to create an CFC Function.

## Example Usage

```hcl
resource "baiducloud_cfc_function" "default" {
  function_name  = "terraform-cfc"
  description    = "terraform create"
  handler        = "index.handler"
  memory_size    = 256
  runtime        = "nodejs12"
  time_out       = 20
  code_file_name = "../../baiducloud/testFiles/cfcTestCode.zip"
}
```

## Argument Reference

The following arguments are supported:

* `function_name` - (Required, ForceNew) CFC function name, length must be between 1 and 64 bytes
* `handler` - (Required) CFC Function execution handler
* `runtime` - (Required) CFC Function runtime
* `time_out` - (Required) Function time out, support [1, 300]s
* `code_bos_bucket` - (Optional) CFC Function Code storage bos bucket name
* `code_bos_object` - (Optional) CFC Function Code storage bos object key
* `code_file_dir` - (Optional) CFC Function Code local file dir
* `code_file_name` - (Optional) CFC Function Code local zip file name
* `description` - (Optional) Function description
* `environment` - (Optional) CFC Function environment variables
* `log_bos_dir` - (Optional) Log save dir if log type is bos
* `log_type` - (Optional) Log save type, support bos/none
* `memory_size` - (Optional) CFC Function memory size, should be an integer multiple of 128
* `reserved_concurrent_executions` - (Optional) Function reserved concurrent executions, support [0-90]

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `code_sha256` - Function code sha256
* `code_size` - Function code size
* `code_storage` - CFC Code storage information
* `commit_id` - Function commit id
* `function_arn` - The same as function brn
* `function_brn` - Function brn
* `last_modified` - The same as update_time
* `region` - Function region
* `role` - Function exec role
* `source_tag` - CFC Function source tag
* `uid` - Function user uid
* `update_time` - Last update time
* `version` - Function version, should only be $LATEST


## Import

CFC can be imported, e.g.

```hcl
$ terraform import baiducloud_cfc_function.default functionName
```

