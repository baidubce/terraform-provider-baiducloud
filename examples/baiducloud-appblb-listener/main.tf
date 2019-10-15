provider "baiducloud" {}

data "baiducloud_zones" "default" {}

data "baiducloud_certs" "default" {}

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

resource "baiducloud_appblb_listener" "default" {
  blb_id               = "${baiducloud_appblb.default.id}"
  listener_port        = 130
  protocol             = "HTTPS"
  scheduler            = "LeastConnection"
  keep_session         = true
  cert_ids             = ["${data.baiducloud_certs.default.certs.0.cert_id}"]
  encryption_protocols = ["sslv3", "tlsv10", "tlsv11"]
  encryption_type      = "userDefind"
}


data "baiducloud_appblb_listeners" "default" {
  blb_id = "${baiducloud_appblb.default.id}"
}
