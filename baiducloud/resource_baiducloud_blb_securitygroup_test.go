package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccBLBSecurityGroupResourceType = "baiducloud_blb_securitygroup"
	testAccBLBSecurityGroupResourceName = testAccBLBSecurityGroupResourceType + "." + BaiduCloudTestResourceName
)

func TestAccBaiduCloudBLBSecurityGroup_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccBLBSecurityGroupDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccBLBSecurityGroupConfig(BaiduCloudTestResourceTypeNameblbSecurityGroup),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBLBSecurityGroupResourceName),
					resource.TestCheckResourceAttrSet(testAccBLBSecurityGroupResourceName, "blb_id"),
					resource.TestCheckResourceAttr(testAccBLBSecurityGroupResourceName, "bind_security_groups.#", "2"),
				),
			},
		},
	})
}

func testAccBLBSecurityGroupDestory(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	blbService := BLBService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccBLBSecurityGroupResourceType {
			continue
		}

		raw, err := blbService.GetBlbSecurityGroup(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		if len(raw.BlbSecurityGroups) != 0 {
			return WrapError(Error("BLB security group still exist"))
		}
	}

	return nil
}

func testAccBLBSecurityGroupConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_vpc" "default" {
    name = "terra-test-vpc"
    description = "baiducloud vpc created by terraform"
    cidr = "192.168.0.0/16"
    tags = {
    "terraform" = "terraform-test"
    }
}

resource "baiducloud_subnet" "default" {
  name = "terra-subnet"
  zone_name = "cn-bj-a"
  cidr = "192.168.3.0/24"
  vpc_id = "${baiducloud_vpc.default.id}"
}

resource "baiducloud_security_group" "default1" {
  name        = "terra-security-group-1"
  description = "created by terraform"
  vpc_id      = "${baiducloud_vpc.default.id}"
}

resource "baiducloud_security_group" "default2" {
  name        = "terra-security-group-2"
  description = "created by terraform"
  vpc_id      = "${baiducloud_vpc.default.id}"
}

resource "baiducloud_blb" "default" {
  name        = "terratestLoadBalance"
  description = "this is a test LoadBalance instance"
  vpc_id      = "${baiducloud_vpc.default.id}"
  subnet_id   = "${baiducloud_subnet.default.id}"
}

resource "baiducloud_blb_securitygroup" "default" {
  blb_id      = "${baiducloud_blb.default.id}"
  security_group_ids = ["${baiducloud_security_group.default1.id}","${baiducloud_security_group.default2.id}"] 
}
`, name)
}
