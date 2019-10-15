output "peer_conn_id" {
  value = "${baiducloud_peer_conn.default.id}"
}

output "peer_conns" {
  value = "${data.baiducloud_peer_conns.default.peer_conns}"
}