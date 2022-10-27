---
layout: "baiducloud"
subcategory: "BOS"
page_title: "BaiduCloud: baiducloud_bos_bucket_objects"
sidebar_current: "docs-baiducloud-datasource-bos_bucket_objects"
description: |-
  Use this data source to query BOS bucket object list.
---

# baiducloud_bos_bucket_objects

Use this data source to query BOS bucket object list.

## Example Usage

```hcl
data "baiducloud_bos_bucket_objects" "default" {
  bucket = "my-bucket"
}

output "objects" {
  value = "${data.baiducloud_bos_bucket_objects.default.objects}"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) Bucket name of the objects to retrieve.
* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `output_file` - (Optional, ForceNew) Output file for saving result.
* `prefix` - (Optional) Prefix of the objects.

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `objects` - List of the objects.
  * `acl` - Acl of the object.
  * `bucket` - Bucket of the object.
  * `cache_control` - Caching behavior of the object.
  * `content_crc32` - Crc(cyclic redundancy check code) value of the object.
  * `content_disposition` - Content disposition of the object.
  * `content_encoding` - Encoding of the object.
  * `content_length` - Content length of the object.
  * `content_md5` - MD5 value of the object content defined in RFC2616.
  * `content_sha256` - Sha256 value of the object.
  * `content_type` - Content type of the object data.
  * `etag` - Etag of the object.
  * `expires` - Expire date of the object.
  * `key` - Key of the object.
  * `last_modified` - Last modifyed time of the object.
  * `owner_id` - Owner id of the object.
  * `owner_name` - Owner name of the object.
  * `size` - Size of the object.
  * `storage_class` - Storage class of the object.
  * `user_meta` - Metadata of the object.


