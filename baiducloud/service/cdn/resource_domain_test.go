package cdn_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/acctest"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/service/cdn"
	"testing"
)

const (
	testAccCdnDomainResourceType = "baiducloud_cdn_domain"
)

func TestAccDomain(t *testing.T) {
	domain := "acc.test.com"
	resourceName := testAccCdnDomainResourceType + ".test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainCreate(domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "domain", domain),
					resource.TestCheckResourceAttr(resourceName, "form", "default"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "cname"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
					resource.TestCheckResourceAttrSet(resourceName, "last_modify_time"),
					resource.TestCheckResourceAttrSet(resourceName, "is_ban"),
					resource.TestCheckResourceAttr(resourceName, "origin.0.weight", "20"),
					resource.TestCheckResourceAttr(resourceName, "origin.0.isp", "un"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"status"},
			},
			{
				Config: testAccDomainOriginUpdate(domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "origin.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "origin.1.backup", "true"),
				),
			},
		},
	})
}

func testAccCheckDomainDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*connectivity.BaiduClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccCdnDomainResourceType {
			continue
		}

		_, err := cdn.FindDomainStatusByName(conn, rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*resource.NotFoundError); ok {
				continue
			}
			return err
		}
		return fmt.Errorf("CDN Domain %s still exist", rs.Primary.ID)
	}

	return nil
}

func testAccCheckDomainExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, err := acctest.CheckResource(resourceName, state)
		if err != nil {
			return err
		}

		conn := acctest.Provider.Meta().(*connectivity.BaiduClient)
		_, err = cdn.FindDomainStatusByName(conn, rs.Primary.ID)

		if err != nil {
			return err
		}
		return nil
	}
}

func testAccDomainCreate(domain string) string {
	return fmt.Sprintf(`
resource "baiducloud_cdn_domain" "test" {
    domain = "%s"
    origin {
        backup = false
        host   = "acc.test.com"
        peer   = "http://2.3.4.5:80"
       	weight = 20
       	isp    = "un"
    }
}`, domain)
}

func testAccDomainOriginUpdate(domain string) string {
	return fmt.Sprintf(`
resource "baiducloud_cdn_domain" "test" {
    domain = "%s"
	default_host = "%s"
    origin {
        backup = false
        host   = "acc.test.com"
        peer   = "http://2.3.4.5:80"
       	weight = 20
    }
	origin {
        backup = true
        host   = "acc2.test.com"
        peer   = "https://2.3.4.5:443"
       	weight = 30
    }
}`, domain, domain)
}
