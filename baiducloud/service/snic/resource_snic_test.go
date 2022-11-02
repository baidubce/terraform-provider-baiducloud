package snic_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/acctest"
	"testing"
)

func TestAccSNIC(t *testing.T) {
	resourceName := "baiducloud_snic" + ".test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acctest.PreCheck(t) },
		Providers: acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccSNICConfig_create(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "name_for_test"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_id"),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "service"),
					resource.TestCheckResourceAttr(resourceName, "description", "description_for_test"),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSNICConfig_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "name_for_test_update"),
					resource.TestCheckResourceAttr(resourceName, "description", "description_for_test_update"),
				),
			},
		},
	})
}

func testAccSNICConfig_create() string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnet(), testAccSNICPublicServicesConfig(), fmt.Sprintf(`
resource "baiducloud_snic" "test" {
    name = "name_for_test"
    vpc_id = baiducloud_vpc.test.id
    subnet_id = baiducloud_subnet.test.id
	service = data.baiducloud_snic_public_services.test.services.0
	description = "description_for_test"
}`))
}

func testAccSNICConfig_update() string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnet(), testAccSNICPublicServicesConfig(), fmt.Sprintf(`
resource "baiducloud_snic" "test" {
    name = "name_for_test_update"
    vpc_id = baiducloud_vpc.test.id
    subnet_id = baiducloud_subnet.test.id
	service = data.baiducloud_snic_public_services.test.services.0
	description = "description_for_test_update"
}`))
}
