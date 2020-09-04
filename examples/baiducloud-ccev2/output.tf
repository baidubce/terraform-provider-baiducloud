
# Print detail of master nodes of the cluster
output "ccev2-managed-cluster-masters" {
  value = baiducloud_ccev2_cluster.default_custom.masters
}

# Print detail of worker(slave) nodes of the cluster
output "ccev2-managed-cluster-nodes" {
  value = baiducloud_ccev2_cluster.default_custom.nodes
}

# Print detail of worker(slave) nodes of a instance group
output "ccev2-instance-group-1-nodes" {
  value = baiducloud_ccev2_instance_group.ccev2_instance_group_1.nodes
}

# Print ClusterNodesDataSource
output "data-cluster-nodes" {
  value = data.baiducloud_ccev2_cluster_instances.default
}

# Print InstanceGroupDataSource
output "data-instance-group-nodes" {
  value = data.baiducloud_ccev2_instance_group_instances.default
}

# Tips: cluster.nodes = {instance_group_1.nodes, instance_group_2.nodes, instance_group_3.nodes}