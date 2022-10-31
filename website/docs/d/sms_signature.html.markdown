---
layout: "baiducloud"
subcategory: "SMS"
page_title: "BaiduCloud: baiducloud_sms_signature"
sidebar_current: "docs-baiducloud-datasource-sms_signature"
description: |-
  Use this data source to query sms signature .
---

# baiducloud_sms_signature

Use this data source to query sms signature .

## Example Usage

```hcl
data "baiducloud_sms_signature" "default" {
	signature_id = "xxxxxx"
}

output "signature_info" {
 	value = "${data.baiducloud_sms_signature.default.signature_info}"
}
```

## Argument Reference

The following arguments are supported:

* `signature_id` - (Required, ForceNew) signature id
* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `output_file` - (Optional, ForceNew) Query result output file path

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `signature_info` - signature content of sms
  * `content_type` - type of content
  * `content` - signature content of sms
  * `country_type` - signature type of country
  * `review` - commit review
  * `signature_id` - signature id
  * `status` - status
  * `user_id` - User id


