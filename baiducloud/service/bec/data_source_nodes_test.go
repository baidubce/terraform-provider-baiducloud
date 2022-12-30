package bec_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/acctest"
)

func TestAccNodes(t *testing.T) {
	dataSourceName := "data.baiducloud_bec_nodes" + ".test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acctest.PreCheck(t) },
		Providers: acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccNodesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "region_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "region_list.0.region"),
					resource.TestCheckResourceAttrSet(dataSourceName, "region_list.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "region_list.0.country"),
					resource.TestCheckResourceAttrSet(dataSourceName, "region_list.0.country_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "region_list.0.city_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "region_list.0.city_list.0.city"),
					resource.TestCheckResourceAttrSet(dataSourceName, "region_list.0.city_list.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "region_list.0.city_list.0.service_provider_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "region_list.0.city_list.0.service_provider_list.0.service_provider"),
					resource.TestCheckResourceAttrSet(dataSourceName, "region_list.0.city_list.0.service_provider_list.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "region_list.0.city_list.0.service_provider_list.0.region_id"),
				),
			},
		},
	})
}

func testAccNodesConfig() string {
	return fmt.Sprintf(`
data "baiducloud_bec_nodes" "test" { 
	type = "vm"
}`)
}
