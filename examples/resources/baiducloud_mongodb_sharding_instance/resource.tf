resource "baiducloud_mongodb_sharding_instance" "example" {

  engine_version = "3.6"
  mongos_count = 2
  mongos_cpu_count = 2
  mongos_memory_capacity = 4

  shard_count = 2
  shard_cpu_count = 2
  shard_memory_capacity = 4
  shard_storage = 20
  shard_storage_type = "CDS_ENHANCED_SSD"

  name = "mongo_example"
  vpc_id = "vpc-abc123"
  subnets {
    subnet_id = "sbn-abc123"
    zone_name = "cn-bj-a"
  }
  tags = {
    TagA = "valueA"
    TagB = "valueB"
  }

}