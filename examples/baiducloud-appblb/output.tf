output "appblb" {
  value = "${data.baiducloud_appblbs.default.blbs}"
}

output "appblb-id" {
  value = "${baiducloud_appblb.default}"
}