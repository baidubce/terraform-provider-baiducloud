provider "baiducloud" {
  # option config, you can use assume role as the operation account
  #assume_role {
  #  account_id = "your-account-id"
  #  role_name = "your-role-name"
  #}
}

resource "baiducloud_vpc" "vpc" {
  name = "terraform_vpc"
  cidr = "172.16.0.0/16"
}
resource "baiducloud_subnet" "subnet" {
  name        = "terraform_subnet-c"
  zone_name   = "cn-bj-c"
  cidr        = "172.16.128.0/24"
  vpc_id      = baiducloud_vpc.vpc.id
  description = "terraform test subnet"
}
resource "baiducloud_subnet" "subnet2" {
  name        = "terraform_subnet-d"
  zone_name   = "cn-bj-d"
  cidr        = "172.16.64.0/24"
  vpc_id      = baiducloud_vpc.vpc.id
  description = "terraform test subnet"
}

resource "baiducloud_cfs" "default" {
  name = "terraform_test"
  zone = "zoneD"
}

resource "baiducloud_cfs_mount_target" "default" {
  fs_id = baiducloud_cfs.default.id
  subnet_id = baiducloud_subnet.subnet.id
  vpc_id = baiducloud_vpc.vpc.id
}
resource "baiducloud_cfs_mount_target" "default2" {
  fs_id = baiducloud_cfs.default.id
  subnet_id = baiducloud_subnet.subnet2.id
  vpc_id = baiducloud_vpc.vpc.id
}

data "baiducloud_cfs_mount_targets" "default" {
  fs_id = baiducloud_cfs.default.id
  filter{
    name = "mount_id"
    values = [baiducloud_cfs_mount_target.default.id]
  }
}

data "baiducloud_cfss" "default" {
  filter{
    name = "fs_id"
    values = [baiducloud_cfs.default.id]
  }
}