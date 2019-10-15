provider "baiducloud" {}

resource "baiducloud_security_group" "my-sg" {
  name        = "${var.name}"
  description = "${var.description}"
}

resource "baiducloud_security_group_rule" "default" {
  security_group_id = "${baiducloud_security_group.my-sg.id}"
  remark            = "remark"
  protocol          = "udp"
  port_range        = "1-65523"
  direction         = "ingress"
}