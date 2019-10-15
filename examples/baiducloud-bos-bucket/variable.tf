variable "my_bucket" {
  default = "test-terraform02"
}

variable "bucket" {
  default = "test-terraform"
}

variable "acl" {
  default = "public-read-write"
}

variable "rc_id" {
  default = "test-rc"
}

variable "logging_prefix" {
  default = "logs/"
}

variable "lr_id" {
  default = "test-id"
}

variable "lr_id01" {
  default = "test-id01"
}

variable "lr_id02" {
  default = "test-id02"
}

variable "date_greater_than_date" {
  default = "2019-09-07T00:00:00Z"
}

variable "date_greater_than_days" {
  default = "$(lastModified)+P7D"
}

variable "action" {
  default = "DeleteObject"
}

variable "storage_class" {
  default = "COLD"
}

variable "server_side_encryption_rule" {
  default = "AES256"
}

variable "index" {
   default = "index.html"
}

variable "err" {
  default = "err.html"
}

variable "resource" {
  default = "test-terraform/*"
}

variable "allowed_origins" {
  default = "https://www.baidu.com"
}

variable "allowed_methods" {
  default = "GET"
}
