output "cce-managed" {
  value = baiducloud_cce_cluster.default_managed
}

#output "cce-independent" {
#  value = baiducloud_cce_cluster.default_independent
#}

output "cce-instance-list" {
  value = data.baiducloud_cce_cluster_nodes.default
}