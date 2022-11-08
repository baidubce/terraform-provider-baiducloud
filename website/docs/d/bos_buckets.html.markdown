---
layout: "baiducloud"
subcategory: "Baidu Object Storage (BOS)"
page_title: "BaiduCloud: baiducloud_bos_buckets"
sidebar_current: "docs-baiducloud-datasource-bos_buckets"
description: |-
  Use this data source to query BOS bucket list.
---

# baiducloud_bos_buckets

Use this data source to query BOS bucket list.

## Example Usage

```hcl
data "baiducloud_bos_buckets" "default" {}

output "buckets" {
 value = "${data.baiducloud_bos_buckets.default.buckets}"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Optional) Name of the bucket to retrieve.
* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `output_file` - (Optional, ForceNew) Output file for saving result.

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `buckets` - List of buckets.
  * `acl` - Acl of the bucket.
  * `bucket` - Name of the bucket.
  * `copyright_protection` - Configuration of the copyright protection.
    * `resource` - Resources to be protected for copyright.
  * `cors_rule` - Configuration of the Cross-Origin Resource Sharing.
    * `allowed_expose_headers` - Indicate which expose headers are allowed.
    * `allowed_headers` - Indicate which headers are allowed.
    * `allowed_methods` - Indicate which methods are allowed.
    * `allowed_origins` - Indicate which origins are allowed.
    * `max_age_seconds` - Indicate time in seconds that browser can cache the response for a preflight request.
  * `creation_date` - Creation date of the bucket.
  * `lifecycle_rule` - Configuration of object lifecycle management.
    * `action` - Action of the lifecycle rule.
      * `name` - Action name of the lifecycle rule.
      * `storage_class` - Storage class of the action.
    * `condition` - Condition of the lifecycle rule.
      * `time` - Condition time, implemented by the date_greater_than.
        * `date_greater_than` - Support absolute time date and relative time days.
    * `id` - ID of the lifecycle rule.
    * `resource` - Resource of the lifecycle rule.
    * `status` - Status of the lifecycle rule.
  * `location` - Bucket location of the bucket.
  * `logging` - Logging of the bucket.
    * `target_bucket` - Target bucket name of the logging.
    * `target_prefix` - Target prefix of the logging.
  * `owner_id` - Owner id of the bucket.
  * `owner_name` - Owner name of the bucket.
  * `replication_configuration` - Replication configuration of the bucket.
    * `destination` - Destination of the replication configuration.
      * `bucket` - Destination bucket name of the replication configuration.
      * `storage_class` - Destination storage class of the replication configuration.
    * `id` - ID of the replication configuration.
    * `replicate_deletes` - Whether to enable the delete synchronization.
    * `replicate_history` - Configuration of the replicate history.
      * `bucket` - Destination bucket name of the replication configuration.
      * `storage_class` - Destination storage class of the replication configuration.
    * `resource` - Resource of the replication configuration.
    * `status` - Status of the replication configuration.
  * `server_side_encryption_rule` - Encryption of the bucket.
  * `storage_class` - Storage class of the bucket.
  * `website` - Website of the BOS bucket.
    * `error_document` - An absolute path to the document to return in case of a 404 error.
    * `index_document` - Baiducloud BOS returns this index document when requests are made to the root domain or any of the subfolders.


