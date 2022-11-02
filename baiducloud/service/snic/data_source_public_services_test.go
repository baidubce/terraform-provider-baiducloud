package snic_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/acctest"
	"testing"
)

func TestAccSNICPublicServices(t *testing.T) {
	dataSourceName := "data.baiducloud_snic_public_services" + ".test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acctest.PreCheck(t) },
		Providers: acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccSNICPublicServicesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "services.#"),
				),
			},
		},
	})
}

func testAccSNICPublicServicesConfig() string {
	return fmt.Sprintf(`
data "baiducloud_snic_public_services" "test" { 
}`)
}
