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

variable "source_address" {
  default = "192.168.0.0/24"
}

variable "destination_address" {
  default = "192.168.1.0/24"
}

variable "name" {
  default = "terraform-BCC"
}

variable "payment_timing" {
  default = "Postpaid"
}