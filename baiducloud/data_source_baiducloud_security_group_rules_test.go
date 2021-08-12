package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const (
	testAccSecurityGroupRulesDataSourceName          = "data.baiducloud_security_group_rules.default"
	testAccSecurityGroupRulesDataSourceAttrKeyPrefix = "rules.0."
)

//lintignore:AT003
func TestAccBaiduCloudSecurityGroupRulesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRulesDataSourceConfig(BaiduCloudTestResourceTypeNameSecurityGroupRule),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccSecurityGroupRulesDataSourceName),
					resource.TestCheckResourceAttr(testAccSecurityGroupRulesDataSourceName, "rules.#", "1"),
					resource.TestCheckResourceAttr(testAccSecurityGroupRulesDataSourceName, testAccSecurityGroupRulesDataSourceAttrKeyPrefix+"direction", "ingress"),
					resource.TestCheckResourceAttr(testAccSecurityGroupRulesDataSourceName, testAccSecurityGroupRulesDataSourceAttrKeyPrefix+"protocol", "udp"),
					resource.TestCheckResourceAttr(testAccSecurityGroupRulesDataSourceName, testAccSecurityGroupRulesDataSourceAttrKeyPrefix+"port_range", "1-65523"),
					resource.TestCheckResourceAttr(testAccSecurityGroupRulesDataSourceName, testAccSecurityGroupRulesDataSourceAttrKeyPrefix+"remark", "remark"),
					resource.TestCheckResourceAttr(testAccSecurityGroupRulesDataSourceName, testAccSecurityGroupRulesDataSourceAttrKeyPrefix+"ether_type", "IPv4"),
					resource.TestCheckResourceAttr(testAccSecurityGroupRulesDataSourceName, testAccSecurityGroupRulesDataSourceAttrKeyPrefix+"source_ip", "all"),
				),
			},
		},
	})
}

func testAccSecurityGroupRulesDataSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_vpc" "default" {
  name = var.name
  description = "created by terraform"
  cidr = "192.168.0.0/24"
}

resource "baiducloud_security_group" "default" {
  name        = var.name
  description = "created by terraform"
  vpc_id      = baiducloud_vpc.default.id
}

resource "baiducloud_security_group_rule" "default" {
  security_group_id = baiducloud_security_group.default.id
  remark            = "remark"
  protocol          = "udp"
  port_range        = "1-65523"
  direction         = "ingress"
}

data "baiducloud_security_group_rules" "default" {
  security_group_id = baiducloud_security_group_rule.default.security_group_id
  vpc_id            = baiducloud_security_group.default.vpc_id

  filter {
    name = "protocol"
    values = ["tcp", "udp"]
  }
}
`, name)
}
