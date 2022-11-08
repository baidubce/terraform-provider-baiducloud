---
layout: "baiducloud"
subcategory: "Simple Message Service (SMS)"
page_title: "BaiduCloud: baiducloud_sms_template"
sidebar_current: "docs-baiducloud-resource-sms_template"
description: |-
  Provide a resource to create an SMS template.
---

# baiducloud_sms_template

Provide a resource to create an SMS template.

## Example Usage

```hcl
resource "baiducloud_sms_template" "default" {
  name	         = "My test template"
  content        = "Test content"
  sms_type       = "CommonNotice"
  country_type   = "GLOBAL"
  description    = "this is a test sms template"

}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, ForceNew) Template content of sms
* `country_type` - (Required, ForceNew) Template type of country
* `description` - (Required, ForceNew) LoadBalance instance's status
* `name` - (Required, ForceNew) Template name of sms
* `sms_type` - (Required, ForceNew) Type of sms
* `template_id` - (Optional) base64 of Template file

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `review` - commit review
* `status` - status
* `user_id` - User id


## Import

SMS template can be imported, e.g.

```hcl
$ terraform import baiducloud_sms_template.default id
```

