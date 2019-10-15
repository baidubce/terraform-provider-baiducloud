output "nat_gateway_id" {
  value = "${baiducloud_nat_gateway.default.id}"
}

output "nat_gateways" {
  value = "${data.baiducloud_nat_gateways.default.nat_gateways}"
}