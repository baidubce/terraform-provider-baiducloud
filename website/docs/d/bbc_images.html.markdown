---
layout: "baiducloud"
subcategory: "Baidu Baremetal Compute (BBC)"
page_title: "BaiduCloud: baiducloud_bbc_images"
sidebar_current: "docs-baiducloud-datasource-bbc_images"
description: |-
  Use this data source to get bbc images list.
---

# baiducloud_bbc_images

Use this data source to get bbc images list.

## Example Usage

```hcl
data "baiducloud_bbc_images" "bbc_images" {
  image_type = "BbcSystem"
  os_name    = "CentOS"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `image_type` - (Optional, ForceNew) Image type of the images to be queried, support ALL/System/Custom/Integration/Sharing/GpuBccSystem/GpuBccCustom/FpgaBccSystem/FpgaBccCustom
* `name_regex` - (Optional, ForceNew) Regex pattern of the search image name
* `os_name` - (Optional, ForceNew) Search image OS Name
* `output_file` - (Optional, ForceNew) Images search result output file

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `images` - Image list
  * `create_time` - The creation time of the image, in a date format that conforms to the BCE specification
  * `desc` - Image description
  * `id` - Image id.Computed after apply
  * `os_arch` - Image os arch.
  * `os_build` - Image os build.
  * `os_name` - Image os name.
  * `os_type` - Image os type.CentOS, Windows, etc.
  * `os_version` - Image os version.
  * `status` - Image status.Creating, CreatedFailed, Available, NotAvailable, Error
  * `type` - Image type.System, Custom, Integration


