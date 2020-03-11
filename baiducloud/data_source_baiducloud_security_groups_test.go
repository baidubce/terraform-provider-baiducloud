package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
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
				Config: testAccSecurityGroupsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccSecurityGroupsDataSourceName),
					resource.TestCheckResourceAttr(testAccSecurityGroupsDataSourceName, "security_groups.#", "1"),
				),
			},
		},
	})
}

func testAccSecurityGroupsDataSourceConfig() string {
	return fmt.Sprintf(`
resource "baiducloud_vpc" "default" {
  name        = "%s"
  description = "test"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_security_group" "default" {
  name        = "%s"
  description = "Baidu acceptance test"
  vpc_id      = baiducloud_vpc.default.id

  tags = {
    "testKey" = "testValue"
  }
}

data "baiducloud_security_groups" "default" {
  vpc_id = baiducloud_security_group.default.vpc_id

  filter {
    name = "name"
    values = ["test-BaiduAcc*"]
  }
}
`, BaiduCloudTestResourceAttrNamePrefix+"VPC",
		BaiduCloudTestResourceAttrNamePrefix+"SecurityGroup")
}
