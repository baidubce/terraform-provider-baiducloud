package baiducloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const (
	testAccEniResourceType = "baiducloud_eni"
	testAccEniResourceName = testAccEniResourceType + "." + BaiduCloudTestResourceTypeNameEni
)

func TestAccBaiduCloudEni(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccEniConfig(BaiduCloudTestResourceTypeNameEni),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccEniResourceName),
					resource.TestCheckResourceAttr(testAccEniResourceName, "description", "terraform test"),
				),
			},
		},
	})
}

func testAccEniConfig(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_vpc" "vpc" {
  name = "terraform_vpc"
  cidr = "172.16.0.0/20"
}
resource "baiducloud_subnet" "subnet" {
  name        = "terraform_subnet"
  zone_name   = "cn-bj-d"
  cidr        = "172.16.0.0/24"
  vpc_id      = baiducloud_vpc.vpc.id
  description = "terraform test subnet"
}
resource "baiducloud_security_group" "sg" {
  name        = "terraform-sg"
  description = "security group created by terraform"
  vpc_id      = baiducloud_vpc.vpc.id
}
resource "baiducloud_security_group_rule" "sgr1_in" {
  security_group_id = baiducloud_security_group.sg.id
  remark            = "remark"
  protocol          = "icmp"
  port_range        = ""
  direction         = "ingress"
  source_ip         = "all"
}
resource "baiducloud_security_group_rule" "sgr1_out" {
  security_group_id = baiducloud_security_group.sg.id
  remark            = "remark"
  protocol          = "icmp"
  port_range        = ""
  direction         = "egress"
  dest_ip           = "all"
}
resource "baiducloud_eni" "%s" {
  name      = "terraform-eni"
  subnet_id = baiducloud_subnet.subnet.id

  description        = "terraform test"
  security_group_ids = [
    baiducloud_security_group.sg.id
  ]
  private_ip {
    primary            = true
    private_ip_address = ""
  }
}
`, name)
}
