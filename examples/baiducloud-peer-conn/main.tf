provider "baiducloud" {
  alias      = "local"
  access_key = "xxx"
  secret_key = "xxx"
  region     = "su"
}

provider "baiducloud" {
  alias      = "peer"
  access_key = "xxx"
  secret_key = "xxx"
  region     = "su"
}

resource "baiducloud_vpc" "local-vpc" {
  provider = baiducloud.local
  name     = var.local_vpc_name
  cidr     = var.local_vpc_cidr
}

resource "baiducloud_vpc" "peer-vpc" {
  provider = baiducloud.peer
  name     = var.peer_vpc_name
  cidr     = var.peer_vpc_cidr
}

resource "baiducloud_peer_conn" "default" {
  provider          = baiducloud.local
  bandwidth_in_mbps = 20
  local_vpc_id      = baiducloud_vpc.local-vpc.id
  peer_vpc_id       = baiducloud_vpc.peer-vpc.id
  peer_region       = var.region
  peer_if_name      = "peer-interface"
  description       = var.description
  local_if_name     = "local-interface-update"
  dns_sync          = false
  billing = {
    payment_timing = "Postpaid"
  }
}

resource "baiducloud_peer_conn_acceptor" "default" {
  provider     = baiducloud.peer
  peer_conn_id = baiducloud_peer_conn.default.id
  auto_accept  = true
  dns_sync     = false
}

data "baiducloud_peer_conns" "default" {
  provider     = baiducloud.local
  peer_conn_id = baiducloud_peer_conn.default.id
}
