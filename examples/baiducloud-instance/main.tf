provider "baiducloud" {}

data "baiducloud_specs" "default" {}

data "baiducloud_zones" "default" {}

data "baiducloud_images" "default" {}

resource "baiducloud_eip" "default" {
  name              = "${var.eip_name}"
  bandwidth_in_mbps = var.eip_bandwidth
  payment_timing    = "Postpaid"
  billing_method    = "ByTraffic"
}

resource "baiducloud_instance" "my-server" {
  image_id              = "${data.baiducloud_images.default.images.0.id}"
  name                  = "${var.instance_name}"
  availability_zone     = "${data.baiducloud_zones.default.zones.1.zone_name}"
  cpu_count             = "${data.baiducloud_specs.default.specs.0.cpu_count}"
  memory_capacity_in_gb = "${data.baiducloud_specs.default.specs.0.memory_size_in_gb}"
  billing = {
    payment_timing = "${var.payment_timing}"
  }

  related_release_flag = true
  delete_cds_snapshot_flag = true

  cds_disks {
    cds_size_in_gb       = 50
    storage_type         = "cloud_hp1"
  }
}

resource "baiducloud_eip_association" "default" {
  eip           = "${baiducloud_eip.default.id}"
  instance_type = "BCC"
  instance_id   = "${baiducloud_instance.my-server.id}"
}

data "baiducloud_instances" "default" {}
