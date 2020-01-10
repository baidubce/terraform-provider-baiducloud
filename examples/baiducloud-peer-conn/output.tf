output "peer_conns" {
  value = data.baiducloud_peer_conns.default.peer_conns
}

output "peer_conn_acceptor" {
  value = baiducloud_peer_conn_acceptor.default
}