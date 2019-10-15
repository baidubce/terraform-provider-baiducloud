output "bucket_id" {
  value = "${baiducloud_bos_bucket.default.id}"
}

output "buckets" {
  value = "${data.baiducloud_bos_buckets.default.buckets}"
}