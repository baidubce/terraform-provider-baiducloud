package baiducloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const (
	testAccLocalDnsRecordsDataSourceName = "data.baiducloud_local_dns_records.default"
	testAccLocalDnsRecordsAttrKeyPrefix  = "records.0."
)

//lintignore:AT003
func TestAccBaiduCloudLocalDnsRecordsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLocalDnsRecordsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccLocalDnsRecordsDataSourceName),
					resource.TestCheckResourceAttr(testAccLocalDnsRecordsDataSourceName, testAccLocalDnsRecordsAttrKeyPrefix+"rr", "www"),
					resource.TestCheckResourceAttr(testAccLocalDnsRecordsDataSourceName, testAccLocalDnsRecordsAttrKeyPrefix+"value", "1.1.1.1"),
					resource.TestCheckResourceAttr(testAccLocalDnsRecordsDataSourceName, testAccLocalDnsRecordsAttrKeyPrefix+"type", "A"),
					resource.TestCheckResourceAttr(testAccLocalDnsRecordsDataSourceName, testAccLocalDnsRecordsAttrKeyPrefix+"description", "terraform-test"),
					resource.TestCheckResourceAttr(testAccLocalDnsRecordsDataSourceName, testAccLocalDnsRecordsAttrKeyPrefix+"priority", "0"),
					resource.TestCheckResourceAttr(testAccLocalDnsRecordsDataSourceName, testAccLocalDnsRecordsAttrKeyPrefix+"status", "pause"),
				),
			},
		},
	})
}

const testAccLocalDnsRecordsDataSourceConfig = `
resource "baiducloud_localdns_record" "local-dns-test" {
  zone_id     = "zone-1mytixsfqpku"
  rr          = "www"
  value       = "1.1.1.1"
  type        = "A"
  description = "terraform-test"
  ttl         = ""
  priority    = ""
  status      = "Enable"
}
data "baiducloud_localdns_records" "local-dns-data" {
  zone_id = "zone-1mytixsfqpku"
  filter {
    name = "description"
    values = ["terraform_test"]
  }
}
`
