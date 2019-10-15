provider "baiducloud" {}

resource "baiducloud_auto_snapshot_policy" "my-asp" {
  name            = "${var.name}"
  time_points     = [0, 22]
  repeat_weekdays = [0, 3]
  retention_days  = -1
  volume_ids      = ["v-Trb3rQXa"]
}

data "baiducloud_auto_snapshot_policies" "default" {}
