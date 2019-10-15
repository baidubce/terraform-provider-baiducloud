---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_certs"
sidebar_current: "docs-baiducloud-datasource-certs"
description: |-
  Use this data source to query CERT list.
---

# baiducloud_certs

Use this data source to query CERT list.

## Example Usage

```hcl
data "baiducloud_certs" "default" {
  name = "testCert"
}

output "certs" {
 value = "${data.baiducloud_certs.default.certs}"
}
```

## Argument Reference

The following arguments are supported:

* `cert_name` - (Optional, ForceNew) Name of the Cert to be queried
* `output_file` - (Optional, ForceNew) Certs search result output file

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `certs` - A list of Cert
  * `cert_common_name` - Cert's common name
  * `cert_create_time` - Cert's create time
  * `cert_id` - Cert's ID
  * `cert_name` - Cert's name
  * `cert_start_time` - Cert's start time
  * `cert_stop_time` - Cert's stop time
  * `cert_type` - Cert's type
  * `cert_update_time` - Cert's update time


