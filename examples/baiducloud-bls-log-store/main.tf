provider "baiducloud" {
  # option config, you can use assume role as the operation account
  #assume_role {
  #  account_id = "your-account-id"
  #  role_name = "your-role-name"
  #}
}


resource "baiducloud_bls_log_store" "default" {
  log_store_name   = "MyTest"
  retention        = 10

}

data "baiducloud_bls_log_stores" "default" {
  name_pattern = "My"
  order = "asc"
  order_by = "retention"
  page_no = 1
  page_size = 10
}