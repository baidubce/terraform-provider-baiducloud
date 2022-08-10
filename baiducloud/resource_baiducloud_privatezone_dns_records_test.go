package baiducloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

const (
	testAccLocalDnsRecordResourceType = "baiducloud_local_dns_record"
	testAccLocalDnsRecordResourceName = testAccLocalDnsRecordResourceType + "." + BaiduCloudTestResourceName
)

func TestAccLocalDnsRecord(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccLocalDnsRecordDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccLocalDnsRecordConfig(BaiduCloudTestResourceTypeNameNatSnatRule),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccLocalDnsRecordResourceName),
					resource.TestCheckResourceAttr(testAccLocalDnsRecordResourceName, "rr", "www"),
					resource.TestCheckResourceAttr(testAccLocalDnsRecordResourceName, "value", "1.1.1.1"),
					resource.TestCheckResourceAttr(testAccLocalDnsRecordResourceName, "type", "A"),
					resource.TestCheckResourceAttr(testAccLocalDnsRecordResourceName, "description", "terraform-test"),
				),
			},
			{
				Config: testAccLocalDnsRecordUpdate(BaiduCloudTestResourceTypeNameNatSnatRule),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccLocalDnsRecordResourceName),
					resource.TestCheckResourceAttr(testAccLocalDnsRecordResourceName, "rr", "aaa"),
					resource.TestCheckResourceAttr(testAccLocalDnsRecordResourceName, "value", "2.2.2.2"),
					resource.TestCheckResourceAttr(testAccLocalDnsRecordResourceName, "type", "A"),
					resource.TestCheckResourceAttr(testAccLocalDnsRecordResourceName, "description", "terraform-test"),
				),
			},
		},
	})
}

func testAccLocalDnsRecordDestroy(s *terraform.State) error {
	return nil
}

func testAccLocalDnsRecordConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "baiducloud_privatezone_dns_record" "local-dns-test" {
  zone_id     = "zone-1mytixsfqpku"
  rr          = "www"
  value       = "1.1.1.1"
  type        = "A"
  description = "terraform-test"
  ttl         = ""
  priority    = ""
}
`, name)
}

func testAccLocalDnsRecordUpdate(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "baiducloud_privatezone_dns_record" "local-dns-test" {
  zone_id     = "zone-1mytixsfqpku"
  rr          = "aaa"
  value       = "2.2.2.2"
  type        = "A"
  description = "terraform-test"
  ttl         = ""
  priority    = ""
}
`, name)
}
