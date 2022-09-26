provider "baiducloud" {
  # option config, you can use assume role as the operation account
  #assume_role {
  #  account_id = "your-account-id"
  #  role_name = "your-role-name"
  #}
}


data "baiducloud_images" "default" {
  image_type = "System"
  name_regex = "8.4 aarch"
  os_name    = "CentOS"
}

resource "baiducloud_instance" "default1" {
  billing = {
    payment_timing = "Postpaid"
  }
  instance_spec = "bcc.gr1.c1m4"
  image_id      = data.baiducloud_images.default.images.0.id
  tags          = {
    "use"  = "xx-bcc"
  }
  availability_zone = "cn-bj-d"
}

resource "baiducloud_blb" "default2" {
  name        = "var.blb_name"
  description = "created by terraform"
  vpc_id      = "${baiducloud_instance.default1.vpc_id}"
  subnet_id   = "${baiducloud_instance.default1.subnet_id}"
}

resource "baiducloud_blb_backend_server" "default" {
  blb_id       = "${baiducloud_blb.default2.id}"
  backend_server_list {
    instance_id = "${baiducloud_instance.default1.id}"
    weight      = 40
  }

}
data "baiducloud_blb_backend_servers" "default" {
    depends_on  = [baiducloud_blb_backend_server.default]
    blb_id ="${baiducloud_blb.default2.id}"
}