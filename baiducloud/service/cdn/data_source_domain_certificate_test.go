package cdn_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/acctest"
	"testing"
)

func TestAccDomainCertificateDataSource(t *testing.T) {
	domain := "acc.test.com"
	dataSourceName := "data.baiducloud_cdn_domain_certificate.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		Providers:    acctest.Providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainCertificate(domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "certificate.0.cert_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "certificate.0.cert_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "certificate.0.cert_common_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "certificate.0.cert_dns_names"),
					resource.TestCheckResourceAttrSet(dataSourceName, "certificate.0.cert_start_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "certificate.0.cert_stop_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "certificate.0.cert_create_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "certificate.0.cert_update_time"),
				),
			},
		},
	})
}

func testAccDomainCertificate(domain string) string {
	return fmt.Sprintf(`
data "baiducloud_cdn_domain_certificate" "test" {
  domain = "%s"
}
`, domain)
}
