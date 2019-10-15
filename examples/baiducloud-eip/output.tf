output "eip_id" {
  value = "${baiducloud_eip.my-eip}"
}

output "eips" {
  value = "${data.baiducloud_eips.default}"
}
