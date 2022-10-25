package cdn_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/acctest"
	"testing"
)

func TestAccDomainConfigHttps(t *testing.T) {
	domain := "acc.test.com"
	resourceName := "baiducloud_cdn_domain_config_https.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		Providers:    acctest.Providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainConfigHttpsCreate(domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "https.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "https.0.http_redirect", "true"),
					resource.TestCheckResourceAttr(resourceName, "https.0.http_redirect_code", "301"),
					resource.TestCheckResourceAttr(resourceName, "https.0.https_redirect", "false"),
					resource.TestCheckResourceAttr(resourceName, "https.0.http2_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "https.0.verify_client", "false"),
					resource.TestCheckResourceAttr(resourceName, "https.0.ssl_protocols.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "ocsp", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDomainConfigHttpsUpdate(domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "https.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "https.0.http_redirect", "false"),
					resource.TestCheckResourceAttr(resourceName, "https.0.https_redirect", "true"),
					resource.TestCheckResourceAttr(resourceName, "https.0.https_redirect_code", "302"),
					resource.TestCheckResourceAttr(resourceName, "https.0.http2_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "https.0.verify_client", "true"),
					resource.TestCheckResourceAttr(resourceName, "https.0.ssl_protocols.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ocsp", "false"),
				),
			},
		},
	})

}

func testAccDomainConfigHttpsCreate(domain string) string {
	return fmt.Sprintf(`
data "baiducloud_cdn_domain_certificate" "test" {
   domain = "%s"
}
resource "baiducloud_cdn_domain_config_https" "test" {
  domain = "%s"
  https {
    enabled             = true
    cert_id             = "${data.baiducloud_cdn_domain_certificate.test.certificate.0.cert_id}"
    http_redirect       = true
    http_redirect_code  = 301
    ssl_protocols       = ["TLSv1.1", "TLSv1.2", "TLSv1.3"]
  }
  ocsp = true
}
`, domain, domain)
}

func testAccDomainConfigHttpsUpdate(domain string) string {
	return fmt.Sprintf(`
data "baiducloud_cdn_domain_certificate" "test" {
   domain = "%s"
}
resource "baiducloud_cdn_domain_config_https" "test" {
  domain = "%s"
  https {
    enabled             = true
    cert_id             = "${data.baiducloud_cdn_domain_certificate.test.certificate.0.cert_id}"
    https_redirect      = true
    http2_enabled       = true
    verify_client       = true
    ssl_protocols       = ["TLSv1.2", "TLSv1.3"]
  }
}
`, domain, domain)
}
