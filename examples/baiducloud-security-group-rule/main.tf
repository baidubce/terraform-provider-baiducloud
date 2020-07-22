provider "baiducloud" {
  # option config, you can use assume role as the operation account
  #assume_role {
  #  account_id = "your-account-id"
  #  role_name = "your-role-name"
  #}
}

resource "baiducloud_vpc" "default" {
  name = var.vpc-name
  description = var.description
  cidr = "192.168.0.0/24"
}

resource "baiducloud_security_group" "my-sg" {
  name        = var.name
  description = var.description
  vpc_id      = baiducloud_vpc.default.id
}

resource "baiducloud_security_group_rule" "default" {
  security_group_id = baiducloud_security_group.my-sg.id
  remark            = "remark"
  # support all/tcp/udp/icmp
  protocol          = "udp"
  port_range        = "1-65523"
  # support ingress/egress
  direction         = "ingress"
  # support IPv4/IPv6
  #ether_type        = "IPv4"
}

resource "baiducloud_security_group_rule" "default2" {
  security_group_id = baiducloud_security_group.my-sg.id
  remark            = "remark"
  protocol          = "tcp"
  port_range        = "22"
  direction         = "ingress"
}

data "baiducloud_security_group_rules" "default" {
  security_group_id = baiducloud_security_group_rule.default.security_group_id
  vpc_id            = baiducloud_vpc.default.id
}