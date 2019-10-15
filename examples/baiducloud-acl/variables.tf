variable "source_address" {
  default = "192.168.0.0/24"
}

variable "destination_address" {
  default = "192.168.1.0/24"
}

variable "source_port" {
  default = "8888"
}

variable "destination_port" {
  default = "9999"
}

variable "subnet_name" {
  default = "test-subnet"
}

variable "subnet_cidr" {
  default = "192.168.1.0/24"
}

variable "vpc_name" {
  default = "test-vpc"
}

variable "vpc_cidr" {
  default = "192.168.0.0/16"
}