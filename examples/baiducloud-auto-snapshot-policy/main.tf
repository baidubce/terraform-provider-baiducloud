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
  name                  = var.name
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
  name            = var.name
  description     = ""
  disk_size_in_gb = 5
  payment_timing  = "Postpaid"
  count           = var.number
}

resource "baiducloud_cds_attachment" "default" {
  count       = var.number
  cds_id      = baiducloud_cds.default.*.id[count.index]
  instance_id = baiducloud_instance.default.id
}

resource "baiducloud_auto_snapshot_policy" "my-asp" {
  name            = var.name
  time_points     = [0, 22]
  repeat_weekdays = [0, 3]
  retention_days  = -1

  # cds must be in-use
  depends_on = [baiducloud_cds_attachment.default]
  volume_ids = baiducloud_cds.default.*.id
}

data "baiducloud_auto_snapshot_policies" "default" {
  asp_name    = baiducloud_auto_snapshot_policy.my-asp.name
  volume_name = baiducloud_cds.default.0.name
}
