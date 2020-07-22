provider "baiducloud" {
  # option config, you can use assume role as the operation account
  #assume_role {
  #  account_id = "your-account-id"
  #  role_name = "your-role-name"
  #}
}

resource "baiducloud_rds_instance" "default" {
  billing = {
    payment_timing = var.payment_timing
  }
  engine_version            = "5.6"
  engine                    = "MySQL"
  cpu_count                 = 1
  memory_capacity           = 1
  volume_capacity           = 5
}

resource "baiducloud_rds_readonly_instance" "default" {
  billing = {
    payment_timing = var.payment_timing
  }
  source_instance_id        = baiducloud_rds_instance.default.instance_id
  cpu_count                 = 1
  memory_capacity           = 1
  volume_capacity           = 5
}

resource "baiducloud_rds_account" "default" {
  instance_id       = baiducloud_rds_instance.default.instance_id
  account_name      = "mysqlaccount"
  password          = "password12"
  desc              = "test"
}

data "baiducloud_rdss" "default" {
  filter {
    name            = "memory_capacity"
    values          = [baiducloud_rds_instance.default.memory_capacity]
  }
}