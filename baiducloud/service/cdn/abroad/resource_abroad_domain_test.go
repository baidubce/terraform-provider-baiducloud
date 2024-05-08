package abroad_test

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
	testAccAbroadCdnDomainResourceType = "baiducloud_abroad_cdn_domain"
)

func TestAbroadAccDomain(t *testing.T) {
	domain := "acc.test.com"
	resourceName := testAccAbroadCdnDomainResourceType + ".test"

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
					resource.TestCheckResourceAttr(resourceName, "origin.0.type", "IP"),
					resource.TestCheckResourceAttr(resourceName, "origin.0.addr", "1.2.3.4"),
					resource.TestCheckResourceAttr(resourceName, "origin.0.backup", "false"),
				),
			},
		},
	})
}

func testAccCheckDomainDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*connectivity.BaiduClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccAbroadCdnDomainResourceType {
			continue
		}

		_, err := cdn.FindDomainStatusByName(conn, rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*resource.NotFoundError); ok {
				continue
			}
			return err
		}
		return fmt.Errorf("Abroad CDN Domain %s still exist", rs.Primary.ID)
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
	resource "baiducloud_abroad_cdn_domain" "default" {
	  domain = "%s"
	  origin {
		backup = false
		type   = "IP"
		addr   = "1.2.3.4"
	  }
	}`, domain)
}