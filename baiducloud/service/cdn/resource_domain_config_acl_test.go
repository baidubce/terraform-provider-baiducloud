package cdn_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/acctest"
	"testing"
)

func TestAccDomainConfigACL(t *testing.T) {
	domain := "acc.test.com"
	resourceName := "baiducloud_cdn_domain_config_acl.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		Providers:    acctest.Providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainConfigACLCreate(domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "referer_acl.0.black_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ip_acl.0.black_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ua_acl.0.black_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "cors.0.origin_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "access_limit.0.limit", "500"),
					resource.TestCheckResourceAttr(resourceName, "traffic_limit.0.limit_rate", "500"),
					resource.TestCheckResourceAttr(resourceName, "request_auth.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDomainConfigACLUpdate(domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "referer_acl.0.white_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ip_acl.0.white_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ua_acl.0.white_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "cors.0.allow", "off"),
					resource.TestCheckResourceAttr(resourceName, "access_limit.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "traffic_limit.0.enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "request_auth.#", "1"),
				),
			},
		},
	})
}

func testAccDomainConfigACLCreate(domain string) string {
	return fmt.Sprintf(`
resource "baiducloud_cdn_domain" "test" {
   domain = "%s"
   origin {
	   backup = false
	   host   = "acc1.test.com"
	   peer   = "http://2.3.4.5:80"
   }
}
resource "baiducloud_cdn_domain_config_acl" "test" {
    domain = baiducloud_cdn_domain.test.id
  	referer_acl {
    	allow_empty = false
    	black_list = ["www.xxx.com", "*.abcde.com"]
  	}
  	ip_acl {
    	black_list = ["1.2.3.4", "2.3.4.5"]
  	}
  	ua_acl {
    	black_list = ["MQQBrowser/5.3/Mozilla/5.0", "Mozilla/5.0 (Linux; Android 7.0"]
  	}
  	cors {
    	allow = "on"
    	origin_list = ["https://www.baidu.com", "http://*.bce.com"]
  	}
  	access_limit {
    	enabled = true
    	limit   = 500
  	}
  	traffic_limit {
    	enable           = true
    	limit_start_hour = 1
    	limit_end_hour   = 23
    	limit_rate       = 500
  	}
}`, domain)
}

func testAccDomainConfigACLUpdate(domain string) string {
	return fmt.Sprintf(`
resource "baiducloud_cdn_domain" "test" {
   domain = "%s"
   origin {
	   backup = false
	   host   = "acc1.test.com"
	   peer   = "http://2.3.4.5:80"
   }
}
resource "baiducloud_cdn_domain_config_acl" "test" {
    domain = baiducloud_cdn_domain.test.id
  	referer_acl {
    	allow_empty = false
    	white_list = ["www.xxx.com", "*.abcde.com"]
  	}
  	ip_acl {
    	white_list = ["1.2.3.4", "2.3.4.5"]
  	}
  	ua_acl {
    	white_list = ["MQQBrowser/5.3/Mozilla/5.0", "Mozilla/5.0 (Linux; Android 7.0"]
  	}
  	cors {
    	allow = "off"
    	origin_list = ["https://www.baidu.com", "http://*.bce.com"]
  	}
  	access_limit {
    	enabled = false
    	limit   = 500
  	}
  	traffic_limit {
    	enable           = false
    	limit_start_hour = 1
    	limit_end_hour   = 23
    	limit_rate       = 500
  	}
    request_auth {
        type = "b"
        key1 = "1234abcd1"
        key2 = "5678abcd2"
        timeout = 1802
        timestamp_metric = 10
    }
}`, domain)
}
