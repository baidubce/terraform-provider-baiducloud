---
layout: "baiducloud"
subcategory: "Baidu Cloud Compute (BCC)"
page_title: "BaiduCloud: baiducloud_cds_attachment"
sidebar_current: "docs-baiducloud-resource-cds_attachment"
description: |-
  Provide a resource to create a CDS attachment, can attach a CDS volume with instance.
---

# baiducloud_cds_attachment

Provide a resource to create a CDS attachment, can attach a CDS volume with instance.

## Example Usage

```hcl
resource "baiducloud_cds_attachment" "default" {
  cds_id      = "v-FJjJeTiG"
  instance_id = "i-tgZhS50C"
}
```

## Argument Reference

The following arguments are supported:

* `cds_id` - (Required, ForceNew) CDS volume ID
* `instance_id` - (Required, ForceNew) The ID of Instance which will attach CDS volume

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `attachment_device` - CDS mount device path
* `attachment_serial` - CDS serial


## Import

CDS attachment can be imported, e.g.

```hcl
$ terraform import baiducloud_cds_attachment.default id
```

