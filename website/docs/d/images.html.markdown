---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_images"
sidebar_current: "docs-baiducloud-datasource-images"
description: |-
  Use this data source to query image list.
---

# baiducloud_images

Use this data source to query image list.

## Example Usage

```hcl
data "baiducloud_images" "default" {}

output "images" {
  value = "${data.baiducloud_images.default.images}"
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
  * `create_time` - Image create time
  * `description` - Image description
  * `id` - Image id
  * `name` - Image name
  * `os_arch` - Image os arch
  * `os_build` - Image os build
  * `os_name` - Image os name
  * `os_type` - Image os type
  * `os_version` - Image os version
  * `special_version` - Image special version
  * `status` - Image status
  * `type` - Image type


