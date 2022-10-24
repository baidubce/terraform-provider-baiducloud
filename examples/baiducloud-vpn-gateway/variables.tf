variable "vpc_name" {
  default = "terraform-vpc"
}
variable "vpn_name" {
  default = "terraform-vpn"
}

variable "description" {
  default = "this is created by terraform"
}

variable "cidr" {
  default = "172.16.32.0/24"
}