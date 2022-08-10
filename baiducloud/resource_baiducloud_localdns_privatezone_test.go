package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccLocalDnsPrivateZoneResourceType = "baiducloud_localdns_privatezone"
	testAccLocalDnsPrivateZoneResourceName = testAccLocalDnsPrivateZoneResourceType + "." + BaiduCloudTestResourceName
)

//lintignore:AT003
func TestAccBaiduCloudLocalDnsPrivateZone(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccLocalDnsPrivateZoneDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccLocalDnsPrivateZoneConfig(BaiduCloudTestResourceTypeNameLocalDnsPrivatezone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccLocalDnsPrivateZoneResourceName),
					resource.TestCheckResourceAttr(testAccLocalDnsPrivateZoneResourceName, "record_count", "2"),
				),
			},
		},
	})
}

func testAccLocalDnsPrivateZoneDestory(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	localDnsService := LocalDnsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccLocalDnsPrivateZoneResourceType {
			continue
		}

		zone, _ := localDnsService.GetPrivateZoneDetail(rs.Primary.ID)
		if zone != nil {
			return WrapError(Error("private zone still exist"))
		}
	}

	return nil
}

func testAccLocalDnsPrivateZoneConfig(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_localdns_privatezone" "default" {
  zone_name               = "%s.com"
}
`, name)
}
