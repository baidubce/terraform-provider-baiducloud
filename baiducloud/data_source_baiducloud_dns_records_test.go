package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccDnsrecordsDataSourceName          = "data.baiducloud_dns_records.default"
	testAccDnsrecordsDataSourceAttrKeyPrefix = "records.0."
)

//lintignore:AT003
func TestAccBaiduCloudDnsrecordsDataSourceSimple(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsrecordByNameDataSourceSimpleConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccDnsrecordsDataSourceName),
					resource.TestCheckResourceAttrSet(testAccDnsrecordsDataSourceAttrKeyPrefix, testAccDnsrecordsDataSourceAttrKeyPrefix+"zone_name"),
				),
			},
		},
	})
}

func TestAccBaiduCloudDnsrecordsDataSourceFull(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsrecordDataSourceFullConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccDnsrecordsDataSourceName),
					resource.TestCheckResourceAttrSet(testAccDnsrecordsDataSourceName, testAccDnsrecordsDataSourceAttrKeyPrefix+"zone_name"),
				),
			},
		},
	})
}

func testAccDnsrecordByNameDataSourceSimpleConfig() string {
	return fmt.Sprintf(`
data "baiducloud_dns_records" "default" { 
	zone_name = "testname"
}
`)
}

func testAccDnsrecordDataSourceFullConfig() string {
	return fmt.Sprintf(`
resource "baiducloud_dns_record" "default" {
  name              = "testname"
  rr                     = "test"
  type                   = "test"
  value                  = "test"
}

data "baiducloud_dns_records" "default" {
  name              = "testname"
}
`)
}
