provider "baiducloud" {}

data "baiducloud_zones" "default" {}

resource "baiducloud_vpc" "default" {
  name        = "${var.vpc_name}"
  description = "terraform create"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = "${var.subnet_name}"
  zone_name   = "${data.baiducloud_zones.default.zones.0.zone_name}"
  cidr        = "192.168.0.0/24"
  vpc_id      = "${baiducloud_vpc.default.id}"
  description = "terraform create"
}

resource "baiducloud_appblb" "default" {
  name        = "${var.appblb_name}"
  description = "terraform create"
  vpc_id      = "${baiducloud_vpc.default.id}"
  subnet_id   = "${baiducloud_subnet.default.id}"

  tags {
    tag_key   = "testKey"
    tag_value = "testValue"
  }
}

resource "baiducloud_appblb_server_group" "default" {
  name        = "${var.servergroup_name}"
  description = "terraform create"
  blb_id      = "${baiducloud_appblb.default.id}"

  port_list {
    port = 66
    type = "TCP"
  }
}

data "baiducloud_appblbs" "default" {}

data "baiducloud_appblb_server_groups" "default" {
  blb_id = "${data.baiducloud_appblbs.default.appblbs.0.blb_id}"
}