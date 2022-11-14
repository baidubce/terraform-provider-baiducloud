package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccBLSLogStoresDataSourceName          = "data.baiducloud_bls_log_stores.default"
	testAccBLSLogStoresDataSourceAttrKeyPrefix = "log_stores.0."
)

//lintignore:AT003
func TestAccBaiduCloudBLSLogStoresDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBLSLogStoreDataSourceConfig(BaiduCloudTestResourceTypeNameblsLogStore),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBLSLogStoresDataSourceName),
					resource.TestCheckResourceAttrSet(testAccBLSLogStoresDataSourceName, "name_pattern"),
				),
			},
		},
	})
}

func testAccBLSLogStoreDataSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "baiducloud_bls_log_store" "default" {
  log_store_name   = "MyTest"
  retention        = 10

}

data "baiducloud_bls_log_stores" "default" {
  name_pattern = "My"
}

`, name)
}
