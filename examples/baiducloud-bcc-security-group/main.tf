provider "baiducloud" {}

data "baiducloud_specs" "default" {
  #name_regex        = "bcc.g1.tiny"
  #instance_type     = "General"
  cpu_count         = 1
  memory_size_in_gb = 4
}

data "baiducloud_zones" "default" {
  name_regex = ".*a$"
}

data "baiducloud_images" "default" {
  image_type = "System"
  name_regex = "7.5.*"
  os_name    = "CentOS"
}

resource "baiducloud_vpc" "default" {
  name = var.vpc_name
  cidr = "192.168.0.0/16"
}

resource "baiducloud_subnet" "default" {
  name      = var.subnet_name
  zone_name = data.baiducloud_zones.default.zones.0.zone_name
  cidr      = "192.168.1.0/24"
  vpc_id    = baiducloud_vpc.default.id
}

resource "baiducloud_security_group" "default" {
  name        = var.security_group_name
  description = "security group created by terraform"
  vpc_id      = baiducloud_vpc.default.id
}

resource "baiducloud_security_group_rule" "default" {
  security_group_id = baiducloud_security_group.default.id
  remark            = "remark"
  protocol          = "udp"
  port_range        = "1-65523"
  direction         = "ingress"
}

resource "baiducloud_security_group_rule" "default" {
  security_group_id = baiducloud_security_group.default.id
  remark            = "remark"
  protocol          = "tcp"
  port_range        = "22"
  direction         = "ingress"
}

resource "baiducloud_instance" "default" {
  image_id              = data.baiducloud_images.default.images.0.id
  name                  = var.instance_name
  description           = var.instance_description
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  billing = {
    payment_timing = "Postpaid"
  }

  subnet_id       = baiducloud_subnet.default.id
  security_groups = [baiducloud_security_group.default.id]

  related_release_flag     = true
  delete_cds_snapshot_flag = true

  cds_disks {
    cds_size_in_gb = 50
    storage_type   = "cloud_hp1"
  }

  tags = {
    "testKey" = "testValue"
  }
}