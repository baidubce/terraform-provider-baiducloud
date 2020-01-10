provider "baiducloud" {}

resource "baiducloud_bos_bucket" "default" {
  bucket = "terraform-bucket-20191106"
  acl    = "public-read-write"
}

data "baiducloud_zones" "default" {}

resource "baiducloud_vpc" "default" {
  name = "terraform-vpc"
  cidr = "192.168.0.0/16"
}

resource "baiducloud_subnet" "default" {
  name      = "terraform-subnet"
  zone_name = "${data.baiducloud_zones.default.zones.0.zone_name}"
  cidr      = "192.168.1.0/24"
  vpc_id    = "${baiducloud_vpc.default.id}"
}

resource "baiducloud_security_group" "default" {
  name   = "terraform-sg"
  vpc_id = "${baiducloud_vpc.default.id}"
}

resource "baiducloud_cfc_function" "default" {
  function_name = "${var.name}"
  description   = "${var.description}"
  environment = {
    "aaa" : "bbb"
    "ccc" : "ddd"
  }
  handler        = "index.handler"
  memory_size    = 256
  runtime        = "nodejs8.5"
  time_out       = 20
  code_file_name = "../../baiducloud/testFiles/cfcTestCode.zip"
  //code_file_dir                  = "../../baiducloud/testFiles/cfcTestCode"
  //code_bos_bucket                = "testBucket"
  //code_bos_object                = "cfcTestCode.zip"
  reserved_concurrent_executions = 20
  vpc_config {
    subnet_ids         = ["${baiducloud_subnet.default.id}"]
    security_group_ids = ["${baiducloud_security_group.default.id}"]
  }
  log_type    = "bos"
  log_bos_dir = "${baiducloud_bos_bucket.default.bucket}"
}

resource "baiducloud_cfc_version" "default" {
  function_name       = "${baiducloud_cfc_function.default.function_name}"
  version_description = "terraformVersion"
  code_sha256         = "${baiducloud_cfc_function.default.code_sha256}"
  log_type            = "none"
}

resource "baiducloud_cfc_alias" "default" {
  function_name    = "${baiducloud_cfc_version.default.function_name}"
  function_version = "${baiducloud_cfc_version.default.version}"
  alias_name       = "terraformAlias"
  description      = "terraform create alias"
}

resource "baiducloud_cfc_trigger" "http-trigger" {
  source_type   = "http"
  target        = "${baiducloud_cfc_version.default.function_brn}"
  resource_path = "/aaabbs"
  method        = ["GET", "PUT"]
  auth_type     = "iam"
}

resource "baiducloud_cfc_trigger" "bos-trigger" {
  source_type    = "bos"
  bucket         = "${baiducloud_bos_bucket.default.bucket}"
  target         = "${baiducloud_cfc_version.default.function_brn}"
  name           = "hehehehe"
  status         = "enabled"
  bos_event_type = ["PutObject", "PostObject"]
  resource       = "/undefined"
}

resource "baiducloud_cfc_trigger" "crontab-trigger" {
  source_type         = "crontab"
  target              = "${baiducloud_cfc_version.default.function_brn}"
  name                = "hahahaha"
  enabled             = "Enabled"
  schedule_expression = "cron(* * * * *)"
}

resource "baiducloud_cfc_trigger" "dueros-trigger" {
  source_type = "dueros"
  target      = "${baiducloud_cfc_version.default.function_brn}"
}

resource "baiducloud_cfc_trigger" "duedge-trigger" {
  source_type = "duedge"
  target      = "${baiducloud_cfc_version.default.function_brn}"
}

resource "baiducloud_cfc_trigger" "cdn-trigger" {
  source_type    = "cdn"
  target         = "${baiducloud_cfc_version.default.function_brn}"
  cdn_event_type = "CachedObjectsBlocked"
  status         = "enabled"
}

data "baiducloud_cfc_function" "default" {
  function_name = "${baiducloud_cfc_function.default.function_name}"
  qualifier     = "${baiducloud_cfc_version.default.version}"
}