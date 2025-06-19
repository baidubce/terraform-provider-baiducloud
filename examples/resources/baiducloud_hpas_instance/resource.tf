resource "baiducloud_hpas_instance" "example" {

  payment_timing        = "Postpaid"
  app_type              = "llama2_7B_train"
  app_performance_level = "10k"
  name                  = "example-instance"
  application_name      = "example-application"
  zone_name             = "cn-bj-a"
  image_id              = "m-example"
  internal_ip           = "192.168.1.100"
  subnet_id             = "sbn-example"
  password              = "1234@password"
  security_group_ids = ["g-example"]
  tags = {
    key1 = "value1"
    key2 = "value2"
  }

}