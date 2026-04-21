resource "baiducloud_hpas_reserved_instance" "example" {
  name                  = "example-reserved"
  zone_name             = "cn-bj-a"
  app_type              = "llama2_7B_train"
  app_performance_level = "10k"
  payment_timing        = "NoPrepay"
}
