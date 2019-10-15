---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_cert"
sidebar_current: "docs-baiducloud-resource-cert"
description: |-
  Provide a resource to Upload a cert.
---

# baiducloud_cert

Provide a resource to Upload a cert.

## Example Usage

```hcl
resource "baiducloud_cert" "cert" {
  cert_name         = "testCert"
  cert_server_data  = ""
  cert_private_data = ""
}
```

## Argument Reference

The following arguments are supported:

* `cert_name` - (Required) Cert Name
* `cert_private_data` - (Required) Cert private key data, base64 encode
* `cert_server_data` - (Required) Server Cert data, base64 encode
* `cert_link_data` - (Optional) Cert lint data, base64 encode
* `cert_type` - (Optional) Cert type

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cert_common_name` - Cert common name
* `cert_create_time` - Cert create time
* `cert_start_time` - Cert start time
* `cert_stop_time` - Cert stop time
* `cert_update_time` - Cert update time


