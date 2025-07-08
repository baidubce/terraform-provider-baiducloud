resource "baiducloud_eip_ddos_protection" "example" {

  ip             = "1.2.3.4"
  threshold_type = "manual"
  ip_clean_mbps  = 120
  ip_clean_pps   = 100000

}