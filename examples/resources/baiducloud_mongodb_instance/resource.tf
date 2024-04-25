resource "baiducloud_mongodb_instance" "example" {

  cpu_count = 2
  memory_capacity = 4
  storage = 20
  engine_version = "3.6"

  voting_member_num = 3
  readonly_node_num = 2

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