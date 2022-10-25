package cdn_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/acctest"
	"testing"
)

func TestAccDomainConfigCache(t *testing.T) {
	domain := "acc.test.com"
	resourceName := "baiducloud_cdn_domain_config_cache.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		Providers:    acctest.Providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainConfigCacheCreate(domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cache_ttl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cache_url_args.0.cache_full_url", "true"),
					resource.TestCheckResourceAttr(resourceName, "error_page.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cache_share.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "mobile_access.0.distinguish_client", "false"),
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
					resource.TestCheckResourceAttr(resourceName, "cache_url_args.0.cache_full_url", "false"),
					resource.TestCheckResourceAttr(resourceName, "cache_url_args.0.cache_url_args.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "error_page.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "cache_share.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "mobile_access.0.distinguish_client", "true"),
				),
			},
		},
	})
}

func testAccDomainConfigCacheCreate(domain string) string {
	return fmt.Sprintf(`
resource "baiducloud_cdn_domain" "test" {
    domain = "%s"
    origin {
        backup = false
        host   = "acc1.test.com"
        peer   = "http://2.3.4.5:80"
        weight = 30
    }
}
resource "baiducloud_cdn_domain_config_cache" "test" {
    domain = baiducloud_cdn_domain.test.id
	cache_ttl {
        type = "suffix"
        value = ".jpg"
        ttl = 36000
        weight = 30
    }
    error_page {
		code = 403
		url = "403.html"
	}
}`, domain)
}

func testAccDomainConfigCacheUpdate(domain string) string {
	return fmt.Sprintf(`
data "baiducloud_cdn_domains" "test" {
}
resource "baiducloud_cdn_domain" "test" {
    domain = "%s"
    origin {
        backup = false
        host   = "acc1.test.com"
        peer   = "http://2.3.4.5:80"
        weight = 30
    }
}
resource "baiducloud_cdn_domain_config_cache" "test" {
    domain = baiducloud_cdn_domain.test.id
	cache_ttl {
        type = "suffix"
        value = ".jpg"
        ttl = 36000
        weight = 30
    }
    cache_ttl {
        type = "suffix"
        value = ".mp4"
        ttl = 36000
        weight = 30
    }
	cache_url_args {
		cache_full_url = false
  		cache_url_args = ["test1", "test2", "test3"]
	}
    error_page {
		code = 403
		url = "403.html"
	}
    error_page {
		code = 404
		url = "404.html"
    }
    cache_share {
        enabled = true
		domain = "${data.baiducloud_cdn_domains.test.domains.1.domain}"
    }
    mobile_access {
        distinguish_client = true
    }
}`, domain)
}
