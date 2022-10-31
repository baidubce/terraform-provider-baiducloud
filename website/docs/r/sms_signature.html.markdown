---
layout: "baiducloud"
subcategory: "SMS"
page_title: "BaiduCloud: baiducloud_sms_signature"
sidebar_current: "docs-baiducloud-resource-sms_signature"
description: |-
  Provide a resource to create an SMS signature.
---

# baiducloud_sms_signature

Provide a resource to create an SMS signature.

## Example Usage

```hcl
resource "baiducloud_sms_signature" "default" {
  content        = "baidu"
  description    = "this is a test sms signature"
  content_type   = "Enterprise"
  country_type   = "DOMESTIC"

}
```

## Argument Reference

The following arguments are supported:

* `content_type` - (Required, ForceNew) type of content
* `content` - (Required, ForceNew) signature content of sms
* `country_type` - (Optional, ForceNew) signature type of country
* `description` - (Optional) LoadBalance instance's status
* `signature_file_base64` - (Optional) base64 of signature file
* `signature_file_format` - (Optional) Format of signature file

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `review` - commit review
* `status` - status
* `user_id` - User id


## Import

SMS signature can be imported, e.g.

```hcl
$ terraform import baiducloud_sms_signature.default id
```

