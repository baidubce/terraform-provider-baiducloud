output "security_group" {
  value = "${baiducloud_security_group.my-sg}"
}

output "security_group_rule" {
  value = "${baiducloud_security_group_rule.default}"
}