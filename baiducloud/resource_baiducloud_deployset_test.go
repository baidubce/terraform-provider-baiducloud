package baiducloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const (
	testAccDeploySetResourceType = "baiducloud_deployset"
	testAccDeploySetResourceName = testAccDeploySetResourceType + "." + BaiduCloudTestResourceName
)

func TestAccBaiduCloudDeploySet(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccDeploySetConfig(BaiduCloudTestResourceTypeNameDeploySet),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccDeploySetResourceName),
					resource.TestCheckResourceAttr(testAccDeploySetResourceName, "desc", "test desc0"),
				),
			},
		},
	})
}

func testAccDeploySetConfig(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_deployset" "default" {
  name     = "%s"
  desc     = "test desc0"
  strategy = "HOST_HA"
}
`, name)
}
