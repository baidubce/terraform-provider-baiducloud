data "baiducloud_mongodb_instances" "example" {

  type = "sharding"
  engine_version = "3.6"
  storage_engine = "WiredTiger"

}