package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccBLSLogStoreResourceType = "baiducloud_bls_log_store"
	testAccBLSLogStoreResourceName = testAccBLSLogStoreResourceType + "." + BaiduCloudTestResourceName
)

//lintignore:AT003
func TestAccBaiduCloudBLSLogStore(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccBLSLogStoreDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccBLSLogStoreConfig(BaiduCloudTestResourceTypeNameblsLogStore),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBLSLogStoreResourceName),
					resource.TestCheckResourceAttr(testAccBLSLogStoreResourceName, "retention", "10"),
				),
			},
			{
				Config: testAccBLSLogStoreUpdate(BaiduCloudTestResourceTypeNameblsLogStore),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBLSLogStoreResourceName),
					resource.TestCheckResourceAttr(testAccBLSLogStoreResourceName, "retention", "5"),
				),
			},
		},
	})
}

func testAccBLSLogStoreDestory(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	blsService := BLSService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccBLSLogStoreResourceType {
			continue
		}

		_, err := blsService.GetBLSLogStoreDetail("MyTest")
		if err != nil {
			continue
		}
		return WrapError(Error("BLS still exist"))
	}

	return nil
}

func testAccBLSLogStoreConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "baiducloud_bls_log_store" "default" {
  log_store_name   = "MyTest"
  retention        = 10

}
`, name)
}

func testAccBLSLogStoreUpdate(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_bls_log_store" "default" {
  log_store_name   = "MyTest"
  retention        = 5

}
`, name+"-update")
}
