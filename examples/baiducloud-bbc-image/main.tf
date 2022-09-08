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
data "baiducloud_bbc_flavors" "bbc_flavors" {
  filter {
    name   = "flavor_id"
    values = ["BBC-I4-01S"]
  }
}
resource "baiducloud_bbc_instance" "my-server" {
  name = "terraform_test"
  flavor_id = "${data.baiducloud_bbc_flavors.bbc_flavors.flavors.0.flavor_id}"
  image_id = "${data.baiducloud_bbc_images.bbc_images.images.0.id}"
  raid = "Raid5"
  root_disk_size_in_gb = 40
  purchase_count = 1
  zone_name = "cn-bj-d"
  security_groups = [
    "${data.baiducloud_security_groups.sg.security_groups.0.id}",
  ]
  billing = {
    payment_timing = "Postpaid"
  }
  tags = {
    "业务"  = "terraform_test"
  }
}
resource "baiducloud_bbc_image" "test-image" {
  image_name = "terraform-bbc-image-test"
  instance_id = baiducloud_bbc_instance.my-server.id
}