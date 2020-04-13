provider "baiducloud" {}

data "baiducloud_zones" "default" {}

resource "baiducloud_vpc" "default" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "baiducloud_subnet" "default" {
  name        = var.subnet_name
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = var.subnet_cidr
  vpc_id      = baiducloud_vpc.default.id
  description = var.description

  tags = {
    "testKey"  = "testValue"
    "testKey2" = "testValue2"
  }
}

data "baiducloud_subnets" "default" {
  #subnet_id = baiducloud_subnet.default.id
  vpc_id = baiducloud_vpc.default.id

  filter {
    name = "cidr"
    values = ["192.168.1.0/26", ".*/24"]
  }
}