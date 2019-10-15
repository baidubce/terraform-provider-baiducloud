---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_bos_bucket_object"
sidebar_current: "docs-baiducloud-resource-bos_bucket_object"
description: |-
  Provide a resource to create a BOS bucket object.
---

# baiducloud_bos_bucket_object

Provide a resource to create a BOS bucket object.

## Example Usage

```hcl
resource "baiducloud_bos_bucket_object" "default" {
  bucket = "my-bucket"
  key = "test-key"
  source = "/tmp/test-file"
  acl = "public-read"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, ForceNew) Name of the bucket to put the file in.
* `key` - (Required, ForceNew) Name of the object once it is in the bucket.
* `acl` - (Optional) Canned ACL of the object, which can be private or public-read. If the value is not set, the object permission will be empty by default, and then the bucket permission as default.
* `cache_control` - (Optional) The caching behavior along the request/reply chain. Valid values are private, no-cache, max-age and must-revalidate. If not set, the value is empty.
* `content_crc32` - (Optional) Crc(cyclic redundancy check code) value of the object.
* `content_disposition` - (Optional) Specifies presentational information for the object, which can be inline or attachment. If not set, the value is empty.
* `content_length` - (Optional) Length of the content to be uploaded.
* `content_md5` - (Optional) MD5 digest of the HTTP request content defined in RFC2616 can be carried by the field to verify whether the file saved on the BOS side is consistent with the file expected by the user.
* `content_sha256` - (Optional) Sha256 value of the object, which is used to verify whether the file saved on the BOS side is consistent with the file expected by the user, the sha256 has higher verification accuracy, and the sha256 value of the transmitted data must match this, otherwise the object uploaded fails.
* `content_type` - (Optional) Type to describe the format of the object data.
* `content` - (Optional, ForceNew) The literal string value that will be uploaded as the object content.
* `expires` - (Optional) The expire date is used to set the cache expiration time when downloading object. If it is not set, the BOS will set the cache expiration time to three days by default.
* `source` - (Optional, ForceNew) The file path that will be read and uploaded as raw bytes for the object content.
* `storage_class` - (Optional) Storage class of the object, which can be COLD, STANDARD_IA, STANDARD or ARCHIVE. Default to STANDARD.
* `user_meta` - (Optional) The mapping of key/values to to provision metadata, which will be automatically prefixed by x-bce-meta-.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `content_encoding` - Encoding of the object.
* `etag` - Etag generated of the object.
* `last_modified` - Last modified date of the object.


