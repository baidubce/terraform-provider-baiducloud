provider "baiducloud" {
  # option config, you can use assume role as the operation account
  #assume_role {
  #  account_id = "your-account-id"
  #  role_name = "your-role-name"
  #}
}

data "baiducloud_images" "images" {
  image_type = "System"
  name_regex = "8.4 aarch"
  os_name    = "CentOS"
}
data "baiducloud_enis" "default" {
  vpc_id      = baiducloud_vpc.vpc.id
  output_file = "res.json"
}
resource "baiducloud_vpc" "vpc" {
  name = "terraform_vpc"
  cidr = "172.16.0.0/20"
}
resource "baiducloud_subnet" "subnet" {
  name        = "terraform_subnet"
  zone_name   = "cn-bj-d"
  cidr        = "172.16.0.0/24"
  vpc_id      = baiducloud_vpc.vpc.id
  description = "terraform test subnet"
}
resource "baiducloud_security_group" "sg" {
  name        = "terraform-sg"
  description = "security group created by terraform"
  vpc_id      = baiducloud_vpc.vpc.id
}
resource "baiducloud_security_group_rule" "sgr1_in" {
  security_group_id = baiducloud_security_group.sg.id
  remark            = "remark"
  protocol          = "icmp"
  port_range        = ""
  direction         = "ingress"
  source_ip         = "all"
}
resource "baiducloud_security_group_rule" "sgr2_in" {
  security_group_id = baiducloud_security_group.sg.id
  remark            = "remark"
  protocol          = "tcp"
  port_range        = "22"
  direction         = "ingress"
  source_ip         = "all"
}
resource "baiducloud_security_group_rule" "sgr1_out" {
  security_group_id = baiducloud_security_group.sg.id
  remark            = "remark"
  protocol          = "icmp"
  port_range        = ""
  direction         = "egress"
  dest_ip           = "all"
}
resource "baiducloud_security_group_rule" "sgr2_out" {
  security_group_id = baiducloud_security_group.sg.id
  remark            = "remark"
  protocol          = "all"
  port_range        = ""
  direction         = "egress"
  dest_ip           = "all"
}

resource "baiducloud_instance" "server1" {
  availability_zone = "cn-bj-d"
  instance_spec     = "bcc.gr1.c1m4"
  image_id          = data.baiducloud_images.images.images.0.id
  billing           = {
    payment_timing = "Postpaid"
  }
  admin_pass      = "Eni12345"
  subnet_id       = baiducloud_subnet.subnet.id
  security_groups = [
    baiducloud_security_group.sg.id
  ]
  #  action = "stop"
}
resource "baiducloud_eip" "eip1" {
  bandwidth_in_mbps = 1
  billing_method    = "ByBandwidth"
  payment_timing    = "Postpaid"
}
resource "baiducloud_eip" "eip2" {
  bandwidth_in_mbps = 1
  billing_method    = "ByBandwidth"
  payment_timing    = "Postpaid"
}


resource "baiducloud_eni" "eni" {
  name      = var.name
  subnet_id = baiducloud_subnet.subnet.id

  description        = "terraform test"
  security_group_ids = [
    baiducloud_security_group.sg.id
  ]
  private_ip {
    primary            = true
    private_ip_address = "172.16.0.10"
    public_ip_address  = baiducloud_eip.eip2.eip
  }
  private_ip {
    primary            = false
    private_ip_address = "172.16.0.11"
    public_ip_address  = baiducloud_eip.eip1.eip
  }
  private_ip {
    primary            = false
    private_ip_address = "172.16.0.13"
  }
}
resource "time_sleep" "wait_30_seconds" {
  depends_on      = [baiducloud_instance.server1, baiducloud_eni.eni]
  create_duration = "60s"
}
resource "baiducloud_eni_attachment" "default" {
  depends_on  = [time_sleep.wait_30_seconds]
  eni_id      = baiducloud_eni.eni.id
  instance_id = baiducloud_instance.server1.id
}