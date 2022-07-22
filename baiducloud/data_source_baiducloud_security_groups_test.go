package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccSecurityGroupsDataSourceName = "data.baiducloud_security_groups.default"
)

//lintignore:AT003
func TestAccBaiduCloudSecurityGroupsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupsDataSourceConfig(BaiduCloudTestResourceTypeNameSecurityGroup),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccSecurityGroupsDataSourceName),
					resource.TestCheckResourceAttr(testAccSecurityGroupsDataSourceName, "security_groups.#", "1"),
				),
			},
		},
	})
}

func testAccSecurityGroupsDataSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_vpc" "default" {
  name        = var.name
  description = "created by terraform"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_security_group" "default" {
  name        = var.name
  description = "created by terraform"
  vpc_id      = baiducloud_vpc.default.id

  tags = {
    "testKey" = "testValue"
  }
}

data "baiducloud_security_groups" "default" {
  vpc_id = baiducloud_security_group.default.vpc_id

  filter {
    name = "name"
    values = ["tf-test-acc*"]
  }
}
`, name)
}
