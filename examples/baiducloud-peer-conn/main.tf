provider "baiducloud" {}

data "baiducloud_zones" "default" {}

resource "baiducloud_vpc" "local-vpc" {
  name = "${var.local_vpc_name}"
  cidr = "${var.local_vpc_cidr}"
}

resource "baiducloud_subnet" "default" {
  name = "${var.peer_subnet_name}"
  zone_name = "${data.baiducloud_zones.default.zones.1.zone_name}"
  cidr = "${var.peer_subnet_cidr}"
  vpc_id = "${baiducloud_vpc.peer-vpc.id}"
}

resource "baiducloud_vpc" "peer-vpc" {
  name = "${var.peer_vpc_name}"
  cidr = "${var.peer_vpc_cidr}"
}

resource "baiducloud_peer_conn" "default" {
  bandwidth_in_mbps = 20
  local_vpc_id = "${baiducloud_vpc.local-vpc.id}"
  peer_vpc_id = "${baiducloud_vpc.peer-vpc.id}"
  peer_region = split("-", "${baiducloud_subnet.default.zone_name}").1
  billing = {
    payment_timing = "Postpaid"
  }
}

data "baiducloud_peer_conns" "default" {
  peer_conn_id = "${baiducloud_peer_conn.default.id}"
}
