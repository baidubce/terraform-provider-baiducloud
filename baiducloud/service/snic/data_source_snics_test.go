package snic_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/acctest"
	"testing"
)

func TestAccSNICs(t *testing.T) {
	resourceName := "baiducloud_snic" + ".test"
	dataSourceName := "data.baiducloud_snics" + ".test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acctest.PreCheck(t) },
		Providers: acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccSNICsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "snics.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "snics.0.name"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", dataSourceName, "snics.0.vpc_id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id", dataSourceName, "snics.0.subnet_id"),
					resource.TestCheckResourceAttrPair(resourceName, "ip_address", dataSourceName, "snics.0.ip_address"),
					resource.TestCheckResourceAttrPair(resourceName, "service", dataSourceName, "snics.0.service"),
					resource.TestCheckResourceAttrPair(resourceName, "description", dataSourceName, "snics.0.description"),
					resource.TestCheckResourceAttrPair(resourceName, "status", dataSourceName, "snics.0.status"),
				),
			},
		},
	})
}

func testAccSNICsConfig() string {
	return acctest.ConfigCompose(testAccSNICConfig_create(), fmt.Sprintf(`
data "baiducloud_snics" "test" {
	vpc_id = baiducloud_snic.test.vpc_id
}`))
}
