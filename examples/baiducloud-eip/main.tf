provider "baiducloud" {}

resource "baiducloud_eip" "my-eip" {
  name              = "${var.name}"
  bandwidth_in_mbps = 100
  payment_timing    = "Postpaid"
  billing_method    = "ByTraffic"
}

data "baiducloud_eips" "default" {}
