package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccEipgroupsDataSourceName          = "data.baiducloud_eipgroups.default"
	testAccEipgroupsDataSourceAttrKeyPrefix = "eips.0."
)

//lintignore:AT003
func TestAccBaiduCloudEipsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEipsDataSourceConfig(BaiduCloudTestResourceTypeNameEip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccEipsDataSourceName),
					resource.TestCheckResourceAttrSet(testAccEipsDataSourceName, testAccEipsDataSourceAttrKeyPrefix+"name"),
				),
			},
		},
	})
}

func testAccEtGatewayDataSourceConfig() string {
	return fmt.Sprintf(`
data "baiducloud_eipgroups" "default" {
    name = "xxxx"    
}
`)
}
