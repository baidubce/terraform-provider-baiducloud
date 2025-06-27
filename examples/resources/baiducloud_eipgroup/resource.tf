resource "baiducloud_eipgroup" "example" {

  name               = "example-group"
  route_type         = "BGP"
  eip_count          = 2
  eipv6_count        = 2
  bandwidth_in_mbps  = 20
  payment_timing     = "Prepaid"
  billing_method     = "ByBandwidth"
  reservation_length = 1
  tags = {
    key1 = "value1"
    key2 = "value2"
  }

}