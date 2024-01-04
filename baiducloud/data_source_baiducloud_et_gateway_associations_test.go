package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccEtGatewayAssociationsDataSourceName          = "data.baiducloud_et_gateway_associations.default"
	testAccEtGatewayAssociationsDataSourceAttrKeyPrefix = "gateway_associations.0."
)

//lintignore:AT003
func TestAccBaiduCloudEtGatewayAssociationsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEtGatewayAssociationsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccEtGatewayDataSourceName),
					resource.TestCheckResourceAttrSet(testAccEtGatewayDataSourceName, testAccEtGatewayDataSourceAttrKeyPrefix+"name"),
				),
			},
		},
	})
}

func testAccEtGatewayAssociationsDataSourceConfig() string {
	return fmt.Sprintf(`
data "baiducloud_et_gateway_associations" "default" {
    et_gateway_id = "xxxxx"   
}
`)
}
