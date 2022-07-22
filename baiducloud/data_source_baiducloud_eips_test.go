package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccEipsDataSourceName          = "data.baiducloud_eips.default"
	testAccEipsDataSourceAttrKeyPrefix = "eips.0."
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
					resource.TestCheckResourceAttrSet(testAccEipsDataSourceName, testAccEipsDataSourceAttrKeyPrefix+"eip"),
					resource.TestCheckResourceAttr(testAccEipsDataSourceName, testAccEipsDataSourceAttrKeyPrefix+"bandwidth_in_mbps", "100"),
					resource.TestCheckResourceAttr(testAccEipsDataSourceName, testAccEipsDataSourceAttrKeyPrefix+"tags.testKey", "testValue"),
				),
			},
		},
	})
}

func testAccEipsDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_eip" "my-eip" {
  name              = "%s"
  bandwidth_in_mbps = 100
  payment_timing    = "Postpaid"
  billing_method    = "ByTraffic"

  tags = {
    "testKey" = "testValue"
  }
}

data "baiducloud_eips" "default" {
  eip = baiducloud_eip.my-eip.id

  filter {
    name = "name"
    values = ["tf-test-acc*"]
  }
}
`, name)
}
