variable "vpc_name" {
  default = "terraform-vpc"
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

variable "description" {
  default = "this is created by terraform"
}
