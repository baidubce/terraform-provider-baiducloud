provider "baiducloud" {
  # option config, you can use assume role as the operation account
  #assume_role {
  #  account_id = "your-account-id"
  #  role_name = "your-role-name"
  #}
}

resource "baiducloud_vpc" "default" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

data "baiducloud_zones" "default" {}

resource "baiducloud_subnet" "default" {
  name      = var.subnet_name
  zone_name = data.baiducloud_zones.default.zones.0.zone_name
  cidr      = var.subnet_cidr
  vpc_id    = baiducloud_vpc.default.id
}

resource "baiducloud_eip" "default" {
  name              = "terraform-eip"
  bandwidth_in_mbps = 10
  payment_timing    = "Postpaid"
  billing_method    = "ByTraffic"
}

resource "baiducloud_nat_gateway" "default" {
  name   = var.nat_name
  vpc_id = baiducloud_vpc.default.id
  spec   = var.spec
  billing = {
    payment_timing = "Postpaid"
  }
  depends_on = [baiducloud_subnet.default]
}

resource "baiducloud_eip_association" "default" {
  eip           = baiducloud_eip.default.id
  instance_type = "NAT"
  instance_id   = baiducloud_nat_gateway.default.id
}

data "baiducloud_nat_gateways" "default" {
  nat_id = baiducloud_nat_gateway.default.id
}
