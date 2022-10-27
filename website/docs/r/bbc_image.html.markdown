---
layout: "baiducloud"
subcategory: "BBC"
page_title: "BaiduCloud: baiducloud_bbc_image"
sidebar_current: "docs-baiducloud-resource-bbc_image"
description: |-
  Use this resource to create BBC custom image.
---

# baiducloud_bbc_image

Use this resource to create BBC custom image.

## Example Usage

```hcl
resource "baiducloud_bbc_image" "test-image" {
  image_name = "terrform-bbc-image-test"
  instance_id = "i-qwIq4vKi"
}
```

## Argument Reference

The following arguments are supported:

* `image_name` - (Required, ForceNew) Image name.
* `instance_id` - (Required, ForceNew) The id of the image source instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The creation time of the image, in a date format that conforms to the BCE specification
* `desc` - Image description
* `image_id` - Image id.Computed after apply
* `os_arch` - Image os arch.
* `os_build` - Image os build.
* `os_name` - Image os name.
* `os_type` - Image os type.CentOS, Windows, etc.
* `os_version` - Image os version.
* `status` - Image status.Creating, CreatedFailed, Available, NotAvailable, Error
* `type` - Image type.System, Custom, Integration


