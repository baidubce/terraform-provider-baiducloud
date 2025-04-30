package abroad

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/acctest"
	"testing"
)

func TestAccAbroadDomainConfigHttps(t *testing.T) {
	domain := "acc.test.com"
	resourceName := "baiducloud_abroad_cdn_domain_config_https.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		Providers:    acctest.Providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainConfigHttpsCreate(domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "http2_enabled", "false"),
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
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "http2_enabled", "true"),
				),
			},
		},
	})
}

func testAccDomainConfigHttpsCreate(domain string) string {
	return fmt.Sprintf(`
resource "baiducloud_abroad_cdn_domain" "test" {
   domain = "%s"
   origin {
	   backup = false
	   host   = "acc1.test.com"
	   peer   = "http://2.3.4.5:80"
   }
}
resource "baiducloud_abroad_cdn_domain_config_https" "example" {
  domain = baiducloud_abroad_cdn_domain.test.domain
  cert_id = "cert-xxxxxxx"
  enabled = true
  http_redirect = false
  http2_enabled = false
}`, domain)
}

func testAccDomainConfigHttpsUpdate(domain string) string {
	return fmt.Sprintf(`
resource "baiducloud_abroad_cdn_domain" "test" {
   domain = "%s"
   origin {
	   backup = false
	   host   = "acc1.test.com"
	   peer   = "http://2.3.4.5:80"
   }
}
resource "baiducloud_abroad_cdn_domain_config_https" "example" {
  domain = baiducloud_abroad_cdn_domain.test.domain
  cert_id = "cert-xxxxxxx"
  enabled = true
  http_redirect = false
  http2_enabled = true
}`, domain)
}
