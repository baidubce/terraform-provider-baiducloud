provider "baiducloud" {}

data "baiducloud_zones" "default" {}

resource "baiducloud_vpc" "default" {
  name = "${var.vpc_name}"
  cidr = "${var.vpc_cidr}"
}

resource "baiducloud_subnet" "default" {
  name = "${var.subnet_name}"
  zone_name = "${data.baiducloud_zones.default.zones.0.zone_name}"
  cidr = "${var.subnet_cidr}"
  vpc_id = "${baiducloud_vpc.default.id}"
  description = "${var.description}"

  tags {
    tag_key = "tagA"
    tag_value = "tagA"
  }
  tags {
    tag_key = "tagB"
    tag_value = "tagB"
  }
}

data "baiducloud_subnets" "default" {
  subnet_id = "${baiducloud_subnet.default.id}"
}