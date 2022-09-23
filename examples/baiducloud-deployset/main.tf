provider "baiducloud" {
  # option config, you can use assume role as the operation account
  #assume_role {
  #  account_id = "your-account-id"
  #  role_name = "your-role-name"
  #}
}
resource "baiducloud_deployset" "default" {
  name     = var.deployset-name
  desc     = "test desc"
  strategy = "HOST_HA"
}
data "baiducloud_images" "default" {
  image_type = "System"
  name_regex = "8.4 aarch"
  os_name    = "CentOS"
}
resource "baiducloud_instance" "default" {
  billing = {
    payment_timing = "Postpaid"
  }
  instance_spec = "bcc.gr1.c1m4"
  image_id      = data.baiducloud_images.default.images.0.id
  tags          = {
    "use"  = "terraform-bcc"
  }
  availability_zone = "cn-bj-d"
  deploy_set_ids = [
    baiducloud_deployset.default.id
  ]
}
data "baiducloud_deploysets" "default" {

}