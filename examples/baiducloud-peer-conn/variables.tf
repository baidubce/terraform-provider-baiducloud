variable "local_vpc_name" {
  default = "local-vpc"
}

variable "local_vpc_cidr" {
  default = "172.16.0.0/16"
}

variable "peer_vpc_name" {
  default = "peer-vpc"
}

variable "peer_vpc_cidr" {
  default = "172.17.0.0/16"
}

variable "peer_subnet_name" {
  default = "terraform-subnet"
}

variable "peer_subnet_cidr" {
  default = "172.17.1.0/24"
}