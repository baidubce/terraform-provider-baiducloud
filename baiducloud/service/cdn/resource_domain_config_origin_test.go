package cdn_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/acctest"
	"testing"
)

func TestAccDomainConfigOrigin(t *testing.T) {
	domain := "acc.test.com"
	resourceName := "baiducloud_cdn_domain_config_origin.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		Providers:    acctest.Providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainConfigOriginCreate(domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "range_switch", "on"),
					resource.TestCheckResourceAttr(resourceName, "origin_protocol.0.value", "*"),
					resource.TestCheckResourceAttr(resourceName, "offline_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "client_ip.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "client_ip.0.name", "True-Client-Ip"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDomainConfigOriginUpdate(domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "range_switch", "off"),
					resource.TestCheckResourceAttr(resourceName, "origin_protocol.0.value", "http"),
					resource.TestCheckResourceAttr(resourceName, "offline_mode", "false"),
					resource.TestCheckResourceAttr(resourceName, "client_ip.0.enabled", "false"),
				),
			},
		},
	})

}

func testAccDomainConfigOriginCreate(domain string) string {
	return fmt.Sprintf(`
resource "baiducloud_cdn_domain" "test" {
   domain = "%s"
   origin {
	   backup = false
	   host   = "acc1.test.com"
	   peer   = "http://2.3.4.5:80"
   }
}
resource "baiducloud_cdn_domain_config_origin" "test" {
  domain = baiducloud_cdn_domain.test.id
  range_switch = "on"
  origin_protocol {
    value = "*"
  }
  offline_mode = true
  client_ip {
    enabled = true
    name    = "True-Client-Ip"
  }
}
`, domain)
}

func testAccDomainConfigOriginUpdate(domain string) string {
	return fmt.Sprintf(`
resource "baiducloud_cdn_domain" "test" {
   domain = "%s"
   origin {
	   backup = false
	   host   = "acc1.test.com"
	   peer   = "http://2.3.4.5:80"
   }
}
resource "baiducloud_cdn_domain_config_origin" "test" {
  domain = baiducloud_cdn_domain.test.id
  range_switch = "off"
  origin_protocol {
    value = "http"
  }
  offline_mode = false
  client_ip {
    enabled = false
  }
}
`, domain)
}
