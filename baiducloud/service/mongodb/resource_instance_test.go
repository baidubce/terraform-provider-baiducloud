package mongodb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/acctest"
)

func TestAccInstance(t *testing.T) {
	resourceName := "baiducloud_mongodb_instance.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acctest.PreCheck(t) },
		Providers: acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfig_create(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "storage_type"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_string"),
					resource.TestCheckResourceAttrSet(resourceName, "port"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),

					resource.TestCheckResourceAttr(resourceName, "payment_timing", "Postpaid"),
					resource.TestCheckResourceAttr(resourceName, "readonly_node_num", "0"),
					resource.TestCheckResourceAttr(resourceName, "status", "RUNNING"),
					resource.TestCheckResourceAttr(resourceName, "storage_engine", "WiredTiger"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccInstanceConfig_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "terraform-test"),
					resource.TestCheckResourceAttr(resourceName, "storage", "10"),
				),
			},
		},
	})
}

func testAccInstanceConfig_create() string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnet(), fmt.Sprintf(`
resource "baiducloud_mongodb_instance" "test" {
  cpu_count = 1
  memory_capacity = 2
  storage = 10
  engine_version = "3.6"
  voting_member_num = 1

  vpc_id = baiducloud_vpc.test.id
  subnets {
    subnet_id = baiducloud_subnet.test.id
    zone_name = baiducloud_subnet.test.zone_name
  }
  tags = {
	Usage = "test"
    CreatedBy = "terraform"
  }
}
`))
}

func testAccInstanceConfig_update() string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnet(), fmt.Sprintf(`
resource "baiducloud_mongodb_instance" "test" {
  cpu_count = 1
  memory_capacity = 2
  storage = 5
  engine_version = "3.6"
  voting_member_num = 1

  vpc_id = baiducloud_vpc.test.id
  subnets {
    subnet_id = baiducloud_subnet.test.id
    zone_name = baiducloud_subnet.test.zone_name
  }
  tags = {
	Usage = "test"
    CreatedBy = "terraform"
  }
  name = "terraform-test"
  account_password = "1a2b3c!4d5e6f"
}
`))
}
