provider "baiducloud" {}

data "baiducloud_specs" "default" {}

data "baiducloud_images" "default" {}

data "baiducloud_zones" "default" {}

data "baiducloud_security_groups" "default" {
  vpc_id = "${baiducloud_vpc.default.id}"
}

resource "baiducloud_vpc" "default" {
  name = "${var.vpc_name}"
  cidr = "${var.vpc_cidr}"
}

resource "baiducloud_subnet" "default" {
  name = "${var.subnet_name}"
  zone_name = "${data.baiducloud_zones.default.zones.1.zone_name}"
  cidr = "${var.subnet_cidr}"
  description = "subnet created by terraform"
  vpc_id = "${baiducloud_vpc.default.id}"
}

resource "baiducloud_instance" "default" {
  image_id = "${data.baiducloud_images.default.images.0.id}"
  name = "${var.name}"
  cpu_count = "${data.baiducloud_specs.default.specs.0.cpu_count}"
  memory_capacity_in_gb = "${data.baiducloud_specs.default.specs.0.memory_size_in_gb}"
  availability_zone = "${data.baiducloud_zones.default.zones.1.zone_name}"
  subnet_id = "${baiducloud_subnet.default.id}"
  security_groups = ["${data.baiducloud_security_groups.default.security_groups.0.id}"]
  billing = {
    payment_timing = "${var.payment_timing}"
  }
}

resource "baiducloud_route_rule" "default" {
  route_table_id = "${baiducloud_vpc.default.route_table_id}"
  source_address = "${var.source_address}"
  destination_address = "${var.destination_address}"
  next_hop_type = "custom"
  next_hop_id = "${baiducloud_instance.default.id}"
  description = "route rule created by terraform"
}

data "baiducloud_route_rules" "default" {
  vpc_id = "${baiducloud_vpc.default.id}"
  route_rule_id = "${baiducloud_route_rule.default.id}"
}
