output "appblb-listener" {
  value = "${baiducloud_appblb_listener.my-appblb-listener}"
}

output "appblb-listeners" {
  value = "${data.baiducloud_appblb_listeners.default.listeners}"
}

output "certs" {
  value = "${data.baiducloud_certs.default}"
}