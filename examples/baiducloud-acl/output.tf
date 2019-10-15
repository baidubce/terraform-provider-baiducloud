output "acl_rule_id" {
  value = "${baiducloud_acl.default.id}"
}

output "acls" {
  value = "${data.baiducloud_acls.default.acls}"
}