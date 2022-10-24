provider "baiducloud" {
  # option config, you can use assume role as the operation account
  #assume_role {
  #  account_id = "your-account-id"
  #  role_name = "your-role-name"
  #}
}

data "baiducloud_vpn_gateways" "default" {
  vpc_id = baiducloud_vpc.vpc.id
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
resource "baiducloud_eip" "eip" {
  bandwidth_in_mbps = 1
  billing_method    = "ByBandwidth"
  payment_timing    = "Postpaid"
}

resource "baiducloud_vpn_gateway" "vpn_gateway" {
  depends_on     = [baiducloud_eip.eip]
  vpn_name       = "test_vpn_gateway"
  vpc_id         = baiducloud_vpc.vpc.id
  description    = "test desc"
  payment_timing = "Postpaid"
  eip            = baiducloud_eip.eip.eip
}

resource "baiducloud_vpn_conn" "vpn_conn" {
  vpn_id        = baiducloud_vpn_gateway.vpn_gateway.id
  secret_key    = "ddd22@www"
  local_subnets = [
    baiducloud_subnet.subnet.cidr
  ]
  remote_ip      = "123.11.11.11"
  remote_subnets = [
    "10.24.0.0/24"
  ]
  description   = "111"
  vpn_conn_name = "vpnconn1"
  ike_config {
    ike_version   = "v1"
    ike_mode      = "main"
    ike_enc_alg   = "aes"
    ike_auth_alg  = "sha1"
    ike_pfs       = "group2"
    ike_life_time = 28800
  }
  ipsec_config {
    ipsec_enc_alg   = "aes"
    ipsec_auth_alg  = "sha1"
    ipsec_pfs       = "group2"
    ipsec_life_time = 28800
  }
}

resource "baiducloud_route_rule" "route_rule" {
  route_table_id      = baiducloud_vpc.vpc.route_table_id
  source_address      = baiducloud_subnet.subnet.cidr
  destination_address = "10.24.0.0/24"
  next_hop_id         = baiducloud_vpn_gateway.vpn_gateway.id
  next_hop_type       = "vpn"
  description         = "created by terraform"
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
resource "baiducloud_security_group_rule" "sgr1_out" {
  security_group_id = baiducloud_security_group.sg.id
  remark            = "remark"
  protocol          = "icmp"
  port_range        = ""
  direction         = "egress"
}
