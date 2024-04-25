package mongodb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/acctest"
)

func TestAccShardingInstance(t *testing.T) {
	resourceName := "baiducloud_mongodb_sharding_instance.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acctest.PreCheck(t) },
		Providers: acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccShardingInstanceConfig_create(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_string"),
					resource.TestCheckResourceAttrSet(resourceName, "engine_version"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),

					resource.TestCheckResourceAttr(resourceName, "payment_timing", "Postpaid"),
					resource.TestCheckResourceAttr(resourceName, "status", "RUNNING"),
					resource.TestCheckResourceAttr(resourceName, "storage_engine", "WiredTiger"),
					resource.TestCheckResourceAttr(resourceName, "mongos_list.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "mongos_list.0.node_id"),
					resource.TestCheckResourceAttr(resourceName, "mongos_list.0.cpu_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "mongos_list.0.memory_capacity", "2"),
					resource.TestCheckResourceAttr(resourceName, "mongos_list.0.status", "RUNNING"),

					resource.TestCheckResourceAttr(resourceName, "shard_list.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "shard_list.0.node_id"),
					resource.TestCheckResourceAttr(resourceName, "shard_list.0.cpu_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "shard_list.0.memory_capacity", "2"),
					resource.TestCheckResourceAttr(resourceName, "shard_list.0.storage", "5"),
					resource.TestCheckResourceAttr(resourceName, "shard_list.0.status", "RUNNING"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"mongos_cpu_count", "mongos_memory_capacity", "shard_cpu_count",
					"shard_memory_capacity", "shard_storage", "shard_storage_type"},
			},
			{
				Config: testAccShardingInstanceConfig_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "terraform-test"),
				),
			},
		},
	})
}

func testAccShardingInstanceConfig_create() string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnet(), fmt.Sprintf(`
resource "baiducloud_mongodb_sharding_instance" "test" {
  vpc_id = baiducloud_vpc.test.id
  subnets {
    subnet_id = baiducloud_subnet.test.id
    zone_name = baiducloud_subnet.test.zone_name
  }

  mongos_count = 2
  mongos_cpu_count = 1
  mongos_memory_capacity = 2

  shard_count = 2
  shard_cpu_count = 1
  shard_memory_capacity = 2
  shard_storage = 5
  shard_storage_type = "CDS_ENHANCED_SSD"

  tags = {
	Usage = "test"
    CreatedBy = "terraform"
  }
}
`))
}

func testAccShardingInstanceConfig_update() string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnet(), fmt.Sprintf(`
resource "baiducloud_mongodb_sharding_instance" "test" {
  vpc_id = baiducloud_vpc.test.id
  subnets {
    subnet_id = baiducloud_subnet.test.id
    zone_name = baiducloud_subnet.test.zone_name
  }

  mongos_count = 2
  mongos_cpu_count = 1
  mongos_memory_capacity = 2

  shard_count = 2
  shard_cpu_count = 1
  shard_memory_capacity = 2
  shard_storage = 5
  shard_storage_type = "CDS_ENHANCED_SSD"

  tags = {
	Usage = "test"
    CreatedBy = "terraform"
  }
  name = "terraform-test"
  account_password = "1a2b3c!4d5e6f"
}
`))
}
