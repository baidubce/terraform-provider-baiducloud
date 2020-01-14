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

variable "instance_role" {
  default = "terraform"
}

variable "instance_short_name" {
  default = "short"
}

variable "instance_format" {
  default = "%02d"
}

variable "payment_timing" {
  default = "Postpaid"
}