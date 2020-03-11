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

  filter{
    name = "os_name"
    values = ["CentOS"]
  }
}

resource "baiducloud_cds" "default" {
  name            = var.cds_name
  disk_size_in_gb = 5
  payment_timing  = "Postpaid"

  # support hp1/std1/cloud_hp1/hdd
  storage_type = "hdd"
  zone_name    = data.baiducloud_zones.default.zones.0.zone_name

  depends_on = [baiducloud_instance.my-server]
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

  cds_disks {
    cds_size_in_gb = 50
    storage_type   = "cloud_hp1"
  }

  cds_disks {
    cds_size_in_gb = 60
    storage_type   = "hp1"
  }

  tags = {
    "testKey" = "testValue"
  }
}

resource "baiducloud_cds_attachment" "default" {
  cds_id      = baiducloud_cds.default.id
  instance_id = baiducloud_instance.my-server.id
}
