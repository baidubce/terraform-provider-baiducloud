provider "baiducloud" {
  # option config, you can use assume role as the operation account
  #assume_role {
  #  account_id = "your-account-id"
  #  role_name = "your-role-name"
  #}
}

resource "baiducloud_bos_bucket" "my-bucket" {
  bucket = var.my_bucket
  acl    = var.acl
}

resource "baiducloud_bos_bucket" "default" {
  bucket = var.bucket
  acl    = var.acl

  replication_configuration {
    id       = var.rc_id
    status   = "enabled"
    resource = [var.resource]
    destination {
      bucket = baiducloud_bos_bucket.my-bucket.bucket
    }
    replicate_deletes = "disabled"
  }

  force_destroy = true

  logging {
    target_bucket = var.bucket
    target_prefix = var.logging_prefix
  }

  lifecycle_rule {
    id       = var.lr_id
    status   = "enabled"
    resource = [var.resource]
    condition {
      time {
        date_greater_than = var.date_greater_than_date
      }
    }
    action {
      name = var.action
    }
  }

  storage_class               = var.storage_class
  server_side_encryption_rule = var.server_side_encryption_rule

  website {
    index_document = var.index
    error_document = var.err
  }

  cors_rule {
    allowed_origins = [var.allowed_origins]
    allowed_methods = [var.allowed_methods]
    max_age_seconds = 1800
  }

  copyright_protection {
    resource = [var.resource]
  }
}

data "baiducloud_bos_buckets" "default" {
  bucket = baiducloud_bos_bucket.default.bucket
}