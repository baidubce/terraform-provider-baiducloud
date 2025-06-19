data "baiducloud_hpas_instances" "example" {

  name        = "instance"
  hpas_status = "Active"
  app_type    = "llama2_7B_train"
  
}