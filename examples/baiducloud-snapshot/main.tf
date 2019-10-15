provider "baiducloud" {}

resource "baiducloud_snapshot" "my-snapshot" {
  name        = "${var.name}"
  description = "${var.description}"
  volume_id   = ""
}

data "baiducloud_snapshots" "default" {}