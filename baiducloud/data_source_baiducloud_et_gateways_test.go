package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccEtGatewayDataSourceName          = "data.baiducloud_et_gateways.default"
	testAccEtGatewayDataSourceAttrKeyPrefix = "gateways.0."
)

//lintignore:AT003
func TestAccBaiduCloudEtGatewayDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEtGatewayDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccEtGatewayDataSourceName),
					resource.TestCheckResourceAttrSet(testAccEtGatewayDataSourceName, testAccEtGatewayDataSourceAttrKeyPrefix+"name"),
				),
			},
		},
	})
}

func testAccEtGatewayDataSourceConfig() string {
	return fmt.Sprintf(`
data "baiducloud_et_gateways" "default" {
    vpc_id = "vpc-ud3hmp5ziuvm"    
}
`)
}
