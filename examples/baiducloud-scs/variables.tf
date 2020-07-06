variable "number" {
  default = 16
}

variable "redis_name" {
  default = "terraform-redis"
}

variable "instance_format" {
  default = "%04d"
}

variable "payment_timing" {
  default = "Postpaid"
}

variable "Prepaid" {
  default = "Prepaid"
}