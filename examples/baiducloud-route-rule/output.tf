output "route_rule_id" {
  value = "${baiducloud_route_rule.default.id}"
}

output "route_rules" {
  value = "${data.baiducloud_route_rules.default.route_rules}"
}