---
layout: "baiducloud"
subcategory: "Baidu Object Storage (BOS)"
page_title: "BaiduCloud: baiducloud_bos_bucket"
sidebar_current: "docs-baiducloud-resource-bos_bucket"
description: |-
  Provide a resource to create a BOS Bucket.
---

# baiducloud_bos_bucket

Provide a resource to create a BOS Bucket.

## Example Usage

```hcl
Versioning bucket
resource "baiducloud_bos_bucket" "default" {
  bucket = "${var.bucket}"
  versioning_status = "enabled"
}
```

```hcl
Private bucket
resource "baiducloud_bos_bucket" "default" {
  bucket = "${var.bucket}"
  acl = "private"
}
```

Using replication configuration
```hcl
resource "baiducloud_bos_bucket" "default" {
  bucket = "${var.bucket}"

  replication_configuration {
    id = "test-rc"
    status = "enabled"
    resource = ["test-terraform/abc*"]
    destination {
      bucket = "test-terraform"
    }
    replicate_deletes = "disabled"
  }
}
```

Using logging
```hcl
resource "baiducloud_bos_bucket" "default" {
  bucket = "${var.bucket}"

  logging {
    target_bucket = "test-terraform"
    target_prefix = "logs/"
  }
}
```

Using lifecycle rule
```hcl
resource "baiducloud_bos_bucket" "default" {
  bucket = "${var.bucket}"

  lifecycle_rule {
	id = "test-id01"
	status =  "enabled"
	resource = ["test-terraform/abc*"]
	condition {
	  time {
	   date_greater_than = "2019-09-07T00:00:00Z"
	  }
	}
	action {
	  name = "DeleteObject"
	}
  }
  lifecycle_rule {
	id = "test-id02"
	status =  "enabled"
	resource = ["test-terraform/def*"]
	condition {
	  time {
	   date_greater_than = "$(lastModified)+P7D"
	  }
	}
	action {
	  name = "DeleteObject"
	}
  }
}
```

Using website
```hcl
resource "baiducloud_bos_bucket" "default" {
  bucket = "${var.bucket}"

  website{
    index_document = "index.html"
    error_document = "err.html"
  }
}
```

Using cors rule
```hcl
resource "baiducloud_bos_bucket" "default" {
  bucket = "${var.bucket}"

  cors_rule {
    allowed_origins = ["https://www.baidu.com"]
    allowed_methods = ["GET"]
    max_age_seconds = 1800
  }
}
```

Using copyright protection
```hcl
resource "baiducloud_bos_bucket" "default" {
  bucket = "${var.bucket}"

  copyright_protection {
    resource = ["test-terraform/abc*"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, ForceNew) Name of the bucket.
* `acl` - (Optional) Canned ACL to apply, available values are private, public-read and public-read-write. Default to private.
* `copyright_protection` - (Optional) Configuration of the copyright protection.
* `cors_rule` - (Optional) Configuration of the Cross-Origin Resource Sharing. Up to 100 rules are allowed per bucket, if there are multiple configurations, the execution order is from top to bottom.
* `enable_multi_az` - (Optional, ForceNew) Whether to enable multi-az replication for the bucket. Default to false.
* `force_destroy` - (Optional) Whether to force delete the bucket and related objects when the bucket is not empty. Default to false.
* `lifecycle_rule` - (Optional) Configuration of object lifecycle management.
* `logging` - (Optional) Settings of the bucket logging.
* `replication_configuration` - (Optional) Replication configuration of the BOS bucket.
* `resource_group` - (Optional, ForceNew) resource group of bucket.
* `server_side_encryption_rule` - (Optional) Encryption rule for the server side, which can only be AES256 currently.
* `storage_class` - (Optional) Storage class of the BOS bucket, available values are STANDARD, STANDARD_IA, MAZ_STANDARD, MAZ_STANDARD_IA, COLD or ARCHIVE.
* `tags` - (Optional, ForceNew) Tags, do not support modify
* `versioning_status` - (Optional) Versioning status of the BOS bucket.
* `website` - (Optional) Website of the BOS bucket.

The `copyright_protection` object supports the following:

* `resource` - (Required) The resources to be protected for copyright.

The `cors_rule` object supports the following:

* `allowed_methods` - (Required) Specifies which methods are allowed. Can be GET,PUT,DELETE,POST or HEAD.
* `allowed_origins` - (Required) Specifies which origins are allowed, containing up to one * wildcard.
* `allowed_expose_headers` - (Optional) Specifies which expose headers are allowed.
* `allowed_headers` - (Optional) Specifies which headers are allowed.
* `max_age_seconds` - (Optional) Specifies time in seconds that browser can cache the response for a preflight request.

The `lifecycle_rule` object supports the following:

* `action` - (Required) Action of the lifecycle rule.
* `condition` - (Required) Condition of the lifecycle rule, only the time form is supported currently.
* `resource` - (Required) Resource of the lifecycle rule. For example, samplebucket/prefix/* will be valid for the object prefixed with prefix/ in samplebucket; samplebucket/* will be valid for all objects in samplebucket.
* `status` - (Required) Status of the lifecycle rule, which can be enabled, disabled. The rule cannot take effect when the status is disabled.
* `id` - (Optional) ID of the lifecycle rule. The id must be unique and cannot be repeated in the same bucket. The system will automatically generate one when it is not specified.

The `action` object supports the following:

* `name` - (Required) Action name of the lifecycle rule, which can be Transition, DeleteObject and AbortMultipartUpload.
* `storage_class` - (Optional) When the action is Transition, it can be set to STANDARD_IA or COLD or ARCHIVE, indicating that it is changed from the original storage type to low frequency storage or cold storage or archive storage.

The `condition` object supports the following:

* `time` - (Required) The condition time, implemented by the date_greater_than.

The `time` object supports the following:

* `date_greater_than` - (Required) Support absolute time date and relative time days. The absolute time date format is yyyy-mm-ddThh:mm:ssZ,eg. 2019-09-07T00:00:00Z. The absolute time is UTC time, which must end at 00:00:00(UTC 0 point); the description of relative time days follows ISO8601, and the minimum time granularity supported is days, such as: $(lastModified)+P7D indicates the time of object 7 days after last-modified.

The `logging` object supports the following:

* `target_bucket` - (Required) Target bucket name that will receive the log data.
* `target_prefix` - (Optional) Target prefix for the log data.

The `replication_configuration` object supports the following:

* `destination` - (Required) Destination of the replication configuration.
* `id` - (Required) ID of the replication configuration.
* `replicate_deletes` - (Required) Whether to enable the delete synchronization, which can be enabled, disabled.
* `resource` - (Required) Resource of the replication configuration. The configuration format of the resource is {$bucket_name/<effective object prefix>}, which must start with "$bucket_name"+"/"
* `status` - (Required) Status of the replication configuration. Valid values are enabled and disabled.
* `replicate_history` - (Optional) Configuration of the replicate history. The bucket name in replicate history needs to be the same as the bucket name in the destination above. After the history file is copied, all the objects of the inventory are copied to the destination bucket synchronously. The history file copy range is not referenced to the resource.

The `destination` object supports the following:

* `bucket` - (Required) Destination bucket name of the replication configuration.
* `storage_class` - (Optional) Destination storage class of the replication configuration, the parameter does not need to be configured if it is consistent with the storage class of the source bucket, if you need to specify the storage class separately, it can be COLD, STANDARD, STANDARD_IA, MAZ_STANDARD, MAZ_STANDARD_IA.

The `replicate_history` object supports the following:

* `bucket` - (Required) Destination bucket name of the replication configuration.
* `storage_class` - (Optional) Destination storage class of the replication configuration, the parameter does not need to be configured if it is consistent with the storage class of the source bucket, if you need to specify the storage class separately, it can be COLD, STANDARD, STANDARD_IA, MAZ_STANDARD, MAZ_STANDARD_IA.

The `website` object supports the following:

* `error_document` - (Optional) An absolute path to the document to return in case of a 404 error.
* `index_document` - (Optional) Baiducloud BOS returns this index document when requests are made to the root domain or any of the subfolders.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `creation_date` - Creation date of the BOS bucket.
* `location` - Location of the BOS bucket.
* `owner_id` - Owner ID of the BOS bucket.
* `owner_name` - Owner name of the BOS bucket.


## Import

BOS bucket can be imported, e.g.

```hcl
$ terraform import baiducloud_bos_bucket.default bucket_id
```

