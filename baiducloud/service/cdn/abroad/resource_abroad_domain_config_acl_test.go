package abroad

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/acctest"
	"testing"
)

func TestAccAbroadDomainConfigACL(t *testing.T) {
	domain := "acc.test.com"
	resourceName := "baiducloud_abroad_cdn_domain_config_acl.test"

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
					resource.TestCheckResourceAttr(resourceName, "allow_empty", "true"),
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
					resource.TestCheckResourceAttr(resourceName, "allow_empty", "false"),
				),
			},
		},
	})
}

func testAccDomainConfigACLCreate(domain string) string {
	return fmt.Sprintf(`
resource "baiducloud_abroad_cdn_domain" "test" {
   domain = "%s"
   origin {
	   backup = false
	   host   = "acc1.test.com"
	   peer   = "http://2.3.4.5:80"
   }
}
resource "baiducloud_cdn_domain_config_acl" "test" {
    domain = baiducloud_cdn_domain.test.id
	allow_empty = true
  	referer_acl {
    	black_list = ["www.xxx.com", "*.abcde.com"]
  	}
  	ip_acl {
    	black_list = ["1.2.3.4", "2.3.4.5"]
  	}
}`, domain)
}

func testAccDomainConfigACLUpdate(domain string) string {
	return fmt.Sprintf(`
resource "baiducloud_abroad_cdn_domain" "test" {
   domain = "%s"
   origin {
	   backup = false
	   host   = "acc1.test.com"
	   peer   = "http://2.3.4.5:80"
   }
}
resource "baiducloud_cdn_domain_config_acl" "test" {
    domain = baiducloud_cdn_domain.test.id
	allow_empty = false
  	referer_acl {
    	white_list = ["www.xxx.com", "*.abcde.com"]
  	}
  	ip_acl {
    	white_list = ["1.2.3.4", "2.3.4.5"]
  	}
}`, domain)
}
