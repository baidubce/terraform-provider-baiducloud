provider "baiducloud" {
  # option config, you can use assume role as the operation account
  #assume_role {
  #  account_id = "your-account-id"
  #  role_name = "your-role-name"
  #}
}

resource "baiducloud_bos_bucket" "default" {
  bucket = var.bucket
}

resource "baiducloud_bos_bucket_object" "default" {
  bucket              = baiducloud_bos_bucket.default.bucket
  key                 = var.key
  content             = var.content
  acl                 = var.acl
  cache_control       = "no-cache"
  content_disposition = "inline"
  storage_class       = "COLD"
  user_meta = {
    metaA = "metaA"
    metaB = "metaB"
  }
}

data "baiducloud_bos_bucket_objects" "default" {
  bucket     = baiducloud_bos_bucket.default.bucket
  depends_on = [baiducloud_bos_bucket_object.default]
}