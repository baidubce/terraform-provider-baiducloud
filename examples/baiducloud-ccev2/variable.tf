variable "cluster_name" {
  default = "ccev2_cluster_1"
}

variable "vpc_cidr" {
  default = "192.168.0.0/16"
}

variable "container_cidr" {
  default = "172.28.0.0/16"
}

variable "cluster_pod_cidr" {
  default = "172.28.0.0/16"
}

variable "cluster_ip_service_cidr" {
  default = "172.31.0.0/16"
}

variable "instance_group_replica_1" {
  default = 1
}

variable "instance_group_replica_2" {
  default = 0
}

variable "instance_group_replica_3" {
  default = 0
}