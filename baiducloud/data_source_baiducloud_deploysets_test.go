package baiducloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const (
	testAccDeploySetDataSourceName          = "data.baiducloud_deploysets.default"
	testAccDeploySetDataSourceAttrKeyPrefix = "deploy_sets.0."
)

func TestAccBaiduCloudDeploySetsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccDeploySetsConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccDeploySetDataSourceName),
					resource.TestCheckResourceAttr(testAccDeploySetDataSourceName, testAccDeploySetDataSourceAttrKeyPrefix+"strategy", "HOST_HA"),
				),
			},
		},
	})
}

func testAccDeploySetsConfig() string {
	return fmt.Sprintf(`
data "baiducloud_deploysets" "default"{

}
`)
}
