package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccBLBSecurityGroupsDataSourceName = "data.baiducloud_blb_securitygroups.default"
)

//lintignore:AT003
func TestAccBaiduCloudBLBSecurityGroupsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupDataSourceConfig(BaiduCloudTestResourceTypeNameblbSecurityGroup),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBLBSecurityGroupsDataSourceName),
					resource.TestCheckResourceAttrSet(testAccBLBSecurityGroupsDataSourceName, "blb_id"),
				),
			},
		},
	})
}

func testAccSecurityGroupDataSourceConfig(name string) string {
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

data "baiducloud_blb_securitygroups" "default" {
    depends_on  = [baiducloud_blb_securitygroup.default]
    blb_id = baiducloud_blb.default.id
}


`, name)
}
