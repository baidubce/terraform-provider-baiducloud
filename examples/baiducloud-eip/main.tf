provider "baiducloud" {}

resource "baiducloud_eip" "my-eip" {
  name              = var.name
  bandwidth_in_mbps = 100
  payment_timing    = "Postpaid"

  # support ByTraffic/ByBandwidth
  billing_method = "ByTraffic"

  tags = {
    "testKey" = "testValue"
  }
}

data "baiducloud_eips" "default" {
  eip = baiducloud_eip.my-eip.id
}
