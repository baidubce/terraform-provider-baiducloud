output "cfc_functions" {
  value = "${baiducloud_cfc_function.default}"
}

output "cfc_versions" {
  value = "${baiducloud_cfc_version.default}"
}

output "cfc_alias" {
  value = "${baiducloud_cfc_alias.default}"
}

output "http_trigger" {
  value = "${baiducloud_cfc_trigger.http-trigger}"
}

output "bos_trigger" {
  value = "${baiducloud_cfc_trigger.bos-trigger}"
}

output "crontab_trigger" {
  value = "${baiducloud_cfc_trigger.crontab-trigger}"
}

output "dueros_trigger" {
  value = "${baiducloud_cfc_trigger.dueros-trigger}"
}

output "duedge_trigger" {
  value = "${baiducloud_cfc_trigger.duedge-trigger}"
}

output "cdn_trigger" {
  value = "${baiducloud_cfc_trigger.cdn-trigger}"
}

output "data_version" {
  value = "${data.baiducloud_cfc_function.default}"
}