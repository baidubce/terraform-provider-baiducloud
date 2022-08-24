package cdn_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/acctest"
	"testing"
)

func TestAccDomainsDataSource(t *testing.T) {
	dataSourceName := "data.baiducloud_cdn_domains.test"
	resourceName := "baiducloud_cdn_domain.test"
	domain := "acc.test.com"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		Providers:    acctest.Providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainsDataSource(domain),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "domains.#", "1"),
					resource.TestCheckResourceAttrPair(dataSourceName, "domains.0.domain", resourceName, "domain"),
					resource.TestCheckResourceAttrPair(dataSourceName, "domains.0.status", resourceName, "status"),
				),
			},
		},
	})
}

func testAccDomainsDataSource(domain string) string {
	return fmt.Sprintf(`
resource "baiducloud_cdn_domain" "test" {
    domain = "%s"
    origin {
        backup = false
        host   = "acc1.test.com"
        peer   = "http://2.3.4.5:80"
    }
}

data "baiducloud_cdn_domains" "test" {
	rule = baiducloud_cdn_domain.test.id
}`, domain)
}
