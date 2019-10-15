output "snapshot" {
  value = "${baiducloud_snapshot.my-snapshot}"
}

output "snapshots" {
  value = "${data.baiducloud_snapshots.default.snapshots}"
}