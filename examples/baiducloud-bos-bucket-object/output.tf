output "object_key" {
  value = "${baiducloud_bos_bucket_object.default.key}"
}

output "objects" {
  value = "${data.baiducloud_bos_bucket_objects.default.objects}"
}