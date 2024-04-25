package mongodb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/acctest"
)

func TestAccDataSourceInstances(t *testing.T) {
	dataSourceName := "data.baiducloud_mongodb_instances.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acctest.PreCheck(t) },
		Providers: acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccInstancesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "instance_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instance_list.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instance_list.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instance_list.0.payment_timing"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instance_list.0.engine_version"),
				),
			},
		},
	})
}

func testAccInstancesConfig() string {
	return acctest.ConfigCompose(testAccInstanceConfig_create(),
		fmt.Sprintf(`
	data "baiducloud_mongodb_instances" "test" {
		type   		   = "replica"
		storage_engine = baiducloud_mongodb_instance.test.storage_engine
	}
	`))
}
