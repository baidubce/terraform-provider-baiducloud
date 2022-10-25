package cdn_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/acctest"
	"testing"
)

func TestAccDomainConfigAdvanced(t *testing.T) {
	domain := "acc.test.com"
	resourceName := "baiducloud_cdn_domain_config_advanced.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		Providers:    acctest.Providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainConfigAdvancedCreate(domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ipv6_dispatch.0.enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "http_header.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "media_drag.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "media_drag.0.mp4.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "media_drag.0.flv.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "seo_switch.0.directly_origin", "OFF"),
					resource.TestCheckResourceAttr(resourceName, "file_trim", "false"),
					resource.TestCheckResourceAttr(resourceName, "compress.0.allow", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDomainConfigAdvancedUpdate(domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ipv6_dispatch.0.enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "http_header.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "media_drag.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "media_drag.0.mp4.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "media_drag.0.mp4.0.start_arg_name", "abcd"),
					resource.TestCheckResourceAttr(resourceName, "media_drag.0.flv.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "media_drag.0.flv.0.start_arg_name", "start"),
					resource.TestCheckResourceAttr(resourceName, "seo_switch.0.directly_origin", "ON"),
					resource.TestCheckResourceAttr(resourceName, "file_trim", "true"),
					resource.TestCheckResourceAttr(resourceName, "compress.0.allow", "true"),
					resource.TestCheckResourceAttr(resourceName, "compress.0.type", "all"),
				),
			},
		},
	})
}

func testAccDomainConfigAdvancedCreate(domain string) string {
	return fmt.Sprintf(`
resource "baiducloud_cdn_domain" "test" {
   domain = "%s"
   origin {
	   backup = false
	   host   = "acc1.test.com"
	   peer   = "http://2.3.4.5:80"
   }
}
resource "baiducloud_cdn_domain_config_advanced" "test" {
  domain = baiducloud_cdn_domain.test.id

  ipv6_dispatch {
  }
  media_drag {
  }
  seo_switch {
  }
  compress {
  }
}
`, domain)
}

func testAccDomainConfigAdvancedUpdate(domain string) string {
	return fmt.Sprintf(`
resource "baiducloud_cdn_domain" "test" {
   domain = "%s"
   origin {
	   backup = false
	   host   = "acc1.test.com"
	   peer   = "http://2.3.4.5:80"
   }
}
resource "baiducloud_cdn_domain_config_advanced" "test" {
  domain = baiducloud_cdn_domain.test.id

  ipv6_dispatch {
    enable = true
  }
  http_header {
    type     = "response"
    header   = "Cache-Control"
    value    = "allowFull"
    action   = "add"
    describe = "Specifies the caching mechanism."
  }
  http_header {
    type     = "origin"
    header   = "Cache-Control"
    action   = "remove"
    describe = "Specifies the caching mechanism."
  }
  media_drag {
    mp4 {
      file_suffix = ["mp4"]
      start_arg_name = "abcd"
      end_arg_name = "1234"
      drag_mode = "second"
    }
    flv {
      file_suffix = ["flv"]
      drag_mode = "byteAV"
    }
  }
  seo_switch {
    directly_origin = "ON"
  }
  file_trim = true
  compress {
    allow = true
    type = "all"
  }
}
`, domain)
}
