package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccLocalDnsVPCResourceType = "baiducloud_localdns_vpc"
	testAccLocalDnsVPCResourceName = testAccLocalDnsVPCResourceType + "." + BaiduCloudTestResourceName
)

//lintignore:AT003
func TestAccBaiduCloudLocalDnsVPC(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccLocalDnsVPCDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccLocalDnsVPCConfig(BaiduCloudTestResourceTypeNameLocalDnsVPC),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccLocalDnsVPCResourceName),
					resource.TestCheckResourceAttr(testAccLocalDnsVPCResourceName, "region", "bj"),
				),
			},
		},
	})
}

func testAccLocalDnsVPCDestory(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	localDnsService := LocalDnsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccLocalDnsVPCResourceType {
			continue
		}

		zone, _ := localDnsService.GetPrivateZoneDetail(rs.Primary.ID)
		if zone != nil {
			return WrapError(Error("private zone still exist"))
		}
	}

	return nil
}

func testAccLocalDnsVPCConfig(name string) string {
	return fmt.Sprintf(`

resource "baiducloud_localdns_privatezone" "default" {
    zone_name = "%s.com"
}

resource "baiducloud_vpc" "default" {
    name = "terra-test-vpc1"
    description = "baiducloud vpc created by terraform"
    cidr = "192.168.0.0/16"
    tags = {
    "terraform" = "terraform-test"
    }
}

resource "baiducloud_localdns_vpc" "default" {
   zone_id = "${baiducloud_localdns_privatezone.default.id}"
   vpc_ids = ["${baiducloud_vpc.default.id}"]
   region = "bj"

}

`, name)
}
