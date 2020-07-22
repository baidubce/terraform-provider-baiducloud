provider "baiducloud" {
  # option config, you can use assume role as the operation account
  #assume_role {
  #  account_id = "your-account-id"
  #  role_name = "your-role-name"
  #}
}

data "baiducloud_specs" "default" {
  # for more detailed conf, please refer to https://cloud.baidu.com/doc/BCC/s/6jwvyo0q2#%E5%8C%BA%E5%9F%9F%E6%9C%BA%E5%9E%8B%E4%BB%A5%E5%8F%8A%E5%8F%AF%E9%80%89%E9%85%8D%E7%BD%AE

  # support General/memory/cpu
  #instance_type     = "General"
  #name_regex        = "bcc.g1.tiny"
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

resource "baiducloud_security_group_rule" "default2" {
  security_group_id = baiducloud_security_group.default.id
  remark            = "remark"
  protocol          = "tcp"
  port_range        = "22"
  direction         = "ingress"
}

resource "baiducloud_eip" "default" {
  count             = var.number
  name              = var.eip_name
  bandwidth_in_mbps = var.eip_bandwidth
  payment_timing    = var.payment_timing
  billing_method    = "ByTraffic"
}

resource "baiducloud_cds" "default" {
  name            = var.cds_name
  disk_size_in_gb = 5
  payment_timing  = "Postpaid"
  storage_type    = "hdd"
  zone_name       = data.baiducloud_zones.default.zones.0.zone_name

  depends_on = [baiducloud_instance.my-server]
  count      = var.number
}

resource "baiducloud_instance" "my-server" {
  count                 = var.number
  image_id              = data.baiducloud_images.default.images.0.id
  name                  = "${var.instance_short_name}-${var.instance_role}-${format(var.instance_format, count.index + 1)}"
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  billing = {
    payment_timing = var.payment_timing
  }

  subnet_id       = baiducloud_subnet.default.id
  security_groups = [baiducloud_security_group.default.id]

  related_release_flag     = true
  delete_cds_snapshot_flag = true

  // The action is optional, which can be start or stop, default is start.
  action = "start"

  // option parameter, please set your keypair id
  #keypair_id = "k-xxxxxx"

  cds_disks {
    cds_size_in_gb = 50
    storage_type   = "cloud_hp1"
  }

  cds_disks {
    cds_size_in_gb = 60
    storage_type   = "hp1"
  }

  tags = {
    "testKey"  = "testValue"
    "testKey2" = "testValue2"
  }
}

resource "baiducloud_eip_association" "default" {
  count         = var.number
  eip           = baiducloud_eip.default.*.id[count.index]
  instance_type = "BCC"
  instance_id   = baiducloud_instance.my-server.*.id[count.index]
}

resource "baiducloud_cds_attachment" "default" {
  count       = var.number
  cds_id      = baiducloud_cds.default.*.id[count.index]
  instance_id = baiducloud_instance.my-server.*.id[count.index]
}
