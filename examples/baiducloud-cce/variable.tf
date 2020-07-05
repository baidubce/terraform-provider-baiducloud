variable "cce-name" {
  default = "terraform-cce"
}

variable "vpc-name" {
  default = "terraform-vpc"
}

variable "description" {
  default = "terraform create"
}

variable "subnet_name_a" {
  default = "terraform-subnet-a"
}

variable "subnet_name_b" {
  default = "terraform-subnet-b"
}

variable "subnet_cidr_a" {
  default = "192.168.1.0/24"
}

variable "subnet_cidr_b" {
  default = "192.168.2.0/24"
}

variable "scurity-group-name" {
  default = "terraform-sg"
}