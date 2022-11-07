---
layout: "baiducloud"
subcategory: "SMS"
page_title: "BaiduCloud: baiducloud_sms_template"
sidebar_current: "docs-baiducloud-datasource-sms_template"
description: |-
  Use this data source to query sms template .
---

# baiducloud_sms_template

Use this data source to query sms template .

## Example Usage

```hcl
data "baiducloud_sms_template" "default" {
	template_id = "xxxxxx"
}

output "template_info" {
 	value = "${data.baiducloud_sms_template.default.template_info}"
}
```

## Argument Reference

The following arguments are supported:

* `template_id` - (Required, ForceNew) template id
* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `output_file` - (Optional, ForceNew) Query result output file path

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `template_info` - template content of sms
  * `name` - name of template
  * `content` - content of template
  * `sms_type` - type of sms
  * `country_type` - Template country type
  * `description` - description of template
  * `template_id` - Template id
  * `review` - commit review
  * `status` - status
  * `user_id` - User id


