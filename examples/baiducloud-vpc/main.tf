provider "baiducloud" {
  # option config, you can use assume role as the operation account
  #assume_role {
  #  account_id = "your-account-id"
  #  role_name = "your-role-name"
  #}
}

resource "baiducloud_vpc" "default" {
  name        = var.name
  description = var.description
  cidr        = var.cidr

  tags = {
    "testKey"  = "testValue"
    "testKey2" = "testValue2"
  }
}

data "baiducloud_vpcs" "default" {
  vpc_id = baiducloud_vpc.default.id
}
