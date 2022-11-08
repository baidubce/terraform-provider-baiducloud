package baiducloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const (
	testAccCfsResourceType = "baiducloud_cfs"
	testAccCfsResourceName = testAccCfsResourceType + "." + BaiduCloudTestResourceName
)

func TestAccBaiduCloudCfs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccCfsConfig(BaiduCloudTestResourceTypeNameCfs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccCfsResourceName),
					resource.TestCheckResourceAttr(testAccCfsResourceName, "name", "tf-test-acc-cfs"),
				),
			},
		},
	})
}

func testAccCfsConfig(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_cfs" "default" {
    name     = "%s"
    zone     = "zoneD"
}
`, name)
}
