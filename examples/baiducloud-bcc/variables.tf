variable "number" {
  default = 1
}

variable "eip_name" {
  default = "terraform-eip"
}

variable "eip_bandwidth" {
  default = 2
}

variable "cds_name" {
  default = "terraform-cds"
}

variable "vpc_name" {
  default = "terraform-vpc"
}

variable "subnet_name" {
  default = "terraform-subnet"
}

variable "security_group_name" {
  default = "terraform-sg"
}

variable "instance_name" {
  default = "terraform-bcc"
}

variable "payment_timing" {
  default = "Postpaid"
}