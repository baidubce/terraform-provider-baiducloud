output "subnet_id" {
  value = "${baiducloud_subnet.default.id}"
}

output "subnets" {
  value = "${data.baiducloud_subnets.default.subnets}"
}