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

resource "baiducloud_eip" "default" {
  name              = var.eip_name
  bandwidth_in_mbps = var.eip_bandwidth

  # support Prepaid/Postpaid
  payment_timing = "Postpaid"

  # support Bytraffic/ByBandwith
  billing_method = "ByTraffic"
}

resource "baiducloud_instance" "my-server" {
  image_id              = data.baiducloud_images.default.images.0.id
  name                  = var.instance_name
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  billing = {
    payment_timing = var.payment_timing
  }

  related_release_flag     = true
  delete_cds_snapshot_flag = true

  tags = {
    "testKey" = "testValue"
  }
}

resource "baiducloud_eip_association" "default" {
  eip = baiducloud_eip.default.id

  # support BCC/BLB/NAT/VPN
  instance_type = "BCC"
  instance_id   = baiducloud_instance.my-server.id
}
