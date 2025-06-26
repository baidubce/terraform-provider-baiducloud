resource "baiducloud_eipgroup_detachment" "example" {

  eip_group_id = "eg-example"
  move_out_eips {
    eip               = "100.88.2.121"
    bandwidth_in_mbps = 10
    payment_timing    = "Postpaid"
    billing_method    = "ByTraffic"
  }
  move_out_eips {
    eip = "240c:4082:ffff:ff01:0:4:0:159"
  }

}
