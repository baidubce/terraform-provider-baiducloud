provider "baiducloud" {
  # option config, you can use assume role as the operation account
  #assume_role {
  #  account_id = "your-account-id"
  #  role_name = "your-role-name"
  #}
}

data "baiducloud_bbc_images" "bbc_images" {
  image_type = "BbcSystem"
  os_name    = "CentOS"
}
data "baiducloud_security_groups" "sg" {
  filter {
    name   = "name"
    values = ["terraform-test"]
  }
}
data "baiducloud_subnets" "subnets" {
  filter {
    name   = "zone_name"
    values = ["cn-bj-d"]
  }
  filter {
    name   = "name"
    values = ["系统预定义子网D"]
  }
}
data "baiducloud_bbc_flavors" "bbc_flavors" {
  filter {
    name   = "flavor_id"
    values = ["BBC-I4-01S"]
  }
}
resource "baiducloud_bbc_instance" "bbc_instance2" {
  action               = "start"
  payment_timing       = "Postpaid"
#  If you want to create a prepaid BBC, use following properties
#  payment_timing = "Prepaid"
#  reservation             = {
#    reservation_length    = 1
#    reservation_time_unit = "Month"
#  }
  flavor_id            = "${data.baiducloud_bbc_flavors.bbc_flavors.flavors.0.flavor_id}"
  image_id             = "${data.baiducloud_bbc_images.bbc_images.images.0.id}"
  name                 = "terraform_test1"
  purchase_count       = 1
  raid                 = "Raid5"
  zone_name            = "cn-bj-d"
  root_disk_size_in_gb = 40
  security_groups      = [
    "${data.baiducloud_security_groups.sg.security_groups.0.id}",
    "${data.baiducloud_security_groups.sg.security_groups.1.id}",
  ]
  tags = {
    "testKey" = "terraform_test"
  }
  description = "terraform test"
  admin_pass  = "terra12345"
}