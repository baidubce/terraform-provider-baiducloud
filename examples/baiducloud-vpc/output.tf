output "vpc_id" {
  value = "${baiducloud_vpc.default.id}"
}

output "vpcs" {
  value = "${data.baiducloud_vpcs.default.vpcs}"
}