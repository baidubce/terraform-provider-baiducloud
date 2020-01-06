---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_cfc_trigger"
sidebar_current: "docs-baiducloud-resource-cfc_trigger"
description: |-
  Provide a resource to create a CFC Function Trigger.
---

# baiducloud_cfc_trigger

Provide a resource to create a CFC Function Trigger.

## Example Usage

```hcl
resource "baiducloud_cfc_trigger" "http-trigger" {
  source_type   = "cfc-http-trigger/v1/CFCAPI"
  target        = "function_brn"
  resource_path = "/aaabbs"
  method        = ["GET","PUT"]
  auth_type     = "iam"
}
```

```

## Argument Reference

The following arguments are supported:

* `source_type` - (Required, ForceNew) CFC Funtion Trigger source type, support bos/http/crontab/dueros/duedge/cdn
* `target` - (Required, ForceNew) CFC Function Trigger target, it should be function brn
* `auth_type` - (Optional) CFC Function Trigger auth type if source_type is http, support anonymous or iam
* `bos_event_type` - (Optional) CFC Function Trigger bos event type
* `bucket` - (Optional, ForceNew) CFC Function Trigger source bucket if source_type is bos
* `cdn_event_type` - (Optional) CFC Function Trigger cdn event type
* `domain` - (Optional) CFC Function Trigger domain if source_type is cdn
* `enabled` - (Optional) CFC Function Trigger enabled if source_type is crontab
* `input` - (Optional) CFC Function Trigger input if source_type is crontab
* `method` - (Optional) CFC Function Trigger method if source_type is http
* `name` - (Optional) CFC Function Trigger name if source_type is bos or crontab
* `remark` - (Optional) CFC Function Trigger remark if source_type is cdn
* `resource_path` - (Optional) CFC Function Trigger resource path if source_type is http
* `resource` - (Optional) CFC Function Trigger resource if source_type is bos
* `schedule_expression` - (Optional) CFC Function Trigger schedule expression if source_type is crontab
* `status` - (Optional) CFC Funtion Trigger status if source_type is bos or cdn

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `relation_id` - CFC Function Trigger relation id
* `sid` - CFC Funtion Trigger sid


