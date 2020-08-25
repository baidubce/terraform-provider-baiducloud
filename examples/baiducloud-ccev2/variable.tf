variable "cluster_name" {
  default = "ccev2_cluster_1"
}

//variable "vpc_id" {
//  default = "vpc-crgbjnaehhhk"
//}
//
//variable "security_group_id" {
//  default = "g-xh04bcdkq5n6"
//}
//
//variable "vpc_subnet_id" {
//  default = "sbn-xdbj15z8v7au"
//}


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
  default = 2
}

variable "instance_group_replica_2" {
  default = 2
}