provider "baiducloud" {
  # option config, you can use assume role as the operation account
  #assume_role {
  #  account_id = "your-account-id"
  #  role_name = "your-role-name"
  #}
}
resource "baiducloud_dts" "default" {
  product_type = "postpay"
  type = "migration"
  standard = "Large"
  source_instance_type = "public"
  target_instance_type = "public"
  cross_region_tag = 0
  task_name = var.dts_name
  data_type = ["schema", "base"]
  src_connection = {
    region = "public"
    db_type = "mysql"
    db_user = "your-username"
    db_pass = "your-password"
    db_port = 3306
    db_host = "192.168.0.1"
    instance_id = "your-instanceId"
    instance_type = "public"
  }
  dst_connection = {
    region = "public"
    db_type = "mysql"
    db_user = "your-username"
    db_pass = "your-password"
    db_port = 3306
    db_host = "192.168.0.1"
    instance_id = "your-instanceId"
    instance_type = "public"
  }
  schema_mapping {
    type = "db"
    src = "db1"
    dst = "db2"
    where = ""
  }
}
data "baiducloud_dtss" "default" {
  dts_name = baiducloud_dts.default.task_name
  type = "migration"
}