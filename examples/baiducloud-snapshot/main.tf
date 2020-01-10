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


resource "baiducloud_instance" "default" {
  name                  = var.instance-name
  image_id              = data.baiducloud_images.default.images.0.id
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  billing = {
    payment_timing = "Postpaid"
  }
}

resource "baiducloud_cds" "default" {
  depends_on      = [baiducloud_instance.default]
  name            = var.cds-name
  description     = var.description
  disk_size_in_gb = 5
  payment_timing  = "Postpaid"
}

resource "baiducloud_snapshot" "my-snapshot" {
  name        = var.sp-name
  description = var.description
  volume_id   = baiducloud_cds.default.id
}

data "baiducloud_snapshots" "default" {
  volume_id = baiducloud_snapshot.my-snapshot.id
}