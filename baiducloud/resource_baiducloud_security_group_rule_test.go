package baiducloud

import (
	"fmt"
	"testing"

	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccSecurityGroupRuleResourceType = "baiducloud_security_group_rule"
	testAccSecurityGroupRuleResourceName = testAccSecurityGroupRuleResourceType + "." + BaiduCloudTestResourceName
)

//lintignore:AT003
func TestAccBaiduCloudSecurityRuleGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSecurityGroupRuleDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleConfig(BaiduCloudTestResourceTypeNameSecurityGroup),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccSecurityGroupRuleResourceName),
					resource.TestCheckResourceAttr(testAccSecurityGroupRuleResourceName, "protocol", "udp"),
					resource.TestCheckResourceAttr(testAccSecurityGroupRuleResourceName, "ether_type", "IPv4"),
					resource.TestCheckResourceAttr(testAccSecurityGroupRuleResourceName, "direction", "ingress"),
					resource.TestCheckResourceAttr(testAccSecurityGroupRuleResourceName, "source_ip", "all"),
				),
			},
		},
	})
}

func testAccSecurityGroupRuleDestory(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	bccService := BccService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccSecurityGroupRuleResourceType {
			continue
		}

		listArgs := &api.ListSecurityGroupArgs{
			VpcId: rs.Primary.Attributes["vpc_id"],
		}
		sgList, err := bccService.ListAllSecurityGroups(listArgs)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		for _, sg := range sgList {
			if sg.Id == rs.Primary.Attributes["security_group_id"] {
				for _, rule := range sg.Rules {
					ruleInfo, err := bccService.parseSecurityGroupRuleId(rs.Primary.ID)
					if err != nil {
						return WrapError(err)
					}

					if compareSecurityGroupRule(&rule, ruleInfo) {
						return fmt.Errorf("security Group Rule still exist")
					}
				}
			}
		}
	}

	return nil
}

func testAccSecurityGroupRuleConfig(name string) string {
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

`, name)
}
