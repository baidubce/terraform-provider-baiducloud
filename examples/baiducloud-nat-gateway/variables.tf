variable "nat_name" {
  default = "terraform-nat-gateway"
}

variable "spec" {
  default = "medium"
}

variable "vpc_name" {
  default = "my-vpc"
}

variable "vpc_cidr" {
  default = "192.168.0.0/16"
}

variable "subnet_name" {
  default = "terraform-subnet"
}

variable "subnet_cidr" {
  default = "192.168.1.0/24"
}