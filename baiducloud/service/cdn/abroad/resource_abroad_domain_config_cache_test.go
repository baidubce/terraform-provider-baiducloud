package abroad_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/acctest"
	"testing"
)

func TestAccDomainConfigCache(t *testing.T) {
	domain := "acc.test.com"
	resourceName := "baiducloud_abroad_cdn_domain_config_cache.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		Providers:    acctest.Providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainConfigCacheCreate(domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cache_ttl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cache_full_url", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDomainConfigCacheUpdate(domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cache_ttl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "cache_full_url", "false"),
				),
			},
		},
	})
}

func testAccDomainConfigCacheCreate(domain string) string {
	return fmt.Sprintf(`
resource "baiducloud_abroad_cdn_domain" "test" {
    domain = "%s"
    origin {
        backup = false
        host   = "acc1.test.com"
        peer   = "http://2.3.4.5:80"
        weight = 30
    }
}
resource "baiducloud_abroad_cdn_domain_cache" "test" {
    domain = baiducloud_abroad_cdn_domain.test.id
	cache_ttl {
        type = "suffix"
        value = ".jpg"
        ttl = 36000
        weight = 30
    }
}`, domain)
}

func testAccDomainConfigCacheUpdate(domain string) string {
	return fmt.Sprintf(`
resource "baiducloud_abroad_cdn_domain" "test" {
    domain = "%s"
    origin {
        backup = false
        host   = "acc1.test.com"
        peer   = "http://2.3.4.5:80"
        weight = 30
    }
}
resource "baiducloud_abroad_cdn_domain_cache" "test" {
    domain = baiducloud_abroad_cdn_domain.test.id
	cache_ttl {
        type = "suffix"
        value = ".jpg"
        ttl = 36000
        weight = 30
    }
	cache_ttl {
		type   = "path"
		value  = "/to/my/file"
		ttl    = 1800
		weight = 5
  	}
    cache_full_url=false
}`, domain)
}
