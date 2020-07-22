provider "baiducloud" {
  # option config, you can use assume role as the operation account
  #assume_role {
  #  account_id = "your-account-id"
  #  role_name = "your-role-name"
  #}
}

data "baiducloud_scs_specs" "default" {
  cluster_type = "master_slave"
  node_capacity = 1
}

data "baiducloud_zones" "default" {}

resource "baiducloud_vpc" "default" {
  name = "terraform-vpc"
  cidr = "192.168.0.0/16"
}

resource "baiducloud_subnet" "default" {
  name = "terraform-subnet"
  zone_name = data.baiducloud_zones.default.zones.0.zone_name
  cidr = "192.168.1.0/24"
  vpc_id = baiducloud_vpc.default.id
}

resource "baiducloud_scs" "default" {
  instance_name = "${var.redis_name}-${format(var.instance_format, var.number)}"
  billing = {
    payment_timing = var.payment_timing
  }
  vpc_id = baiducloud_vpc.default.id
  subnets {
    subnet_id = baiducloud_subnet.default.id
    zone_name = baiducloud_subnet.default.zone_name
  }
  port = 6379
  engine_version = "3.2"
  node_type = data.baiducloud_scs_specs.default.specs.0.node_type
  cluster_type = "master_slave"
  replication_num = 1
}

data "baiducloud_scss" "default" {
  name_regex = "terraform"
  filter {
    name = "cluster_type"
    values = ["master_slave"]
  }
}