output "cds" {
  value = "${baiducloud_cds.my-cds}"
}

output "cdss" {
  value = "${data.baiducloud_cdss.default.cdss}"
}