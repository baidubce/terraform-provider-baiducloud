provider "baiducloud" {}

data "baiducloud_specs" "default" {}

data "baiducloud_zones" "default" {}

data "baiducloud_images" "default" {
  image_type = "System"
}

resource "baiducloud_instance" "default" {
  name                  = "${var.instance_name}"
  image_id              = "${data.baiducloud_images.default.images.0.id}"
  availability_zone     = "${data.baiducloud_zones.default.zones.0.zone_name}"
  cpu_count             = "${data.baiducloud_specs.default.specs.0.cpu_count}"
  memory_capacity_in_gb = "${data.baiducloud_specs.default.specs.0.memory_size_in_gb}"
  billing = {
    payment_timing = "Postpaid"
  }
}

resource "baiducloud_cds" "default" {
  name            = "%s"
  disk_size_in_gb = 5
  payment_timing  = "Postpaid"
}

resource "baiducloud_cds_attachment" "default" {
  cds_id      = "${baiducloud_cds.default.id}"
  instance_id = "${baiducloud_instance.default.id}"
}

data "baiducloud_cdss" "default" {}