output "instance_id" {
  value = "${baiducloud_instance.my-server.id}"
}

output "instacnes" {
  value = "${data.baiducloud_instances.default.instances}"
}

output "eip_association" {
  value = "${baiducloud_eip_association.default}"
}