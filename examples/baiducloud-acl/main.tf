provider "baiducloud" {
  # option config, you can use assume role as the operation account
  #assume_role {
  #  account_id = "your-account-id"
  #  role_name = "your-role-name"
  #}
}


data "baiducloud_zones" "default" {}

resource "baiducloud_vpc" "default" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "baiducloud_subnet" "default" {
  name      = var.subnet_name
  zone_name = data.baiducloud_zones.default.zones.0.zone_name
  cidr      = var.subnet_cidr
  vpc_id    = baiducloud_vpc.default.id
}

resource "baiducloud_acl" "default" {
  subnet_id = baiducloud_subnet.default.id

  # support all/tcp/udp/icmp
  protocol               = "tcp"
  source_ip_address      = var.source_address
  destination_ip_address = var.destination_address
  source_port            = var.source_port
  destination_port       = var.destination_port
  position               = 20

  # support ingress/egress
  direction = "ingress"

  # support allow/deny
  action = "allow"
}

data "baiducloud_acls" "default" {
  subnet_id = baiducloud_subnet.default.id
  acl_id    = baiducloud_acl.default.id
}
