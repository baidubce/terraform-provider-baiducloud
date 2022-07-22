package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccACLsDataSourceName          = "data.baiducloud_acls.default"
	testAccACLsDataSourceAttrKeyPrefix = "acls.0."
)

//lintignore:AT003
func TestAccBaiduCloudACLsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccACLsDataSourceConfigBySubnet(BaiduCloudTestResourceTypeNameAcl),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccACLsDataSourceName),
					resource.TestCheckResourceAttrSet(testAccACLsDataSourceName, testAccACLsDataSourceAttrKeyPrefix+"acl_id"),
					resource.TestCheckResourceAttrSet(testAccACLsDataSourceName, testAccACLsDataSourceAttrKeyPrefix+"subnet_id"),
					resource.TestCheckResourceAttr(testAccACLsDataSourceName, testAccACLsDataSourceAttrKeyPrefix+"description", "created by terraform"),
					resource.TestCheckResourceAttr(testAccACLsDataSourceName, testAccACLsDataSourceAttrKeyPrefix+"protocol", "tcp"),
					resource.TestCheckResourceAttr(testAccACLsDataSourceName, testAccACLsDataSourceAttrKeyPrefix+"source_ip_address", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(testAccACLsDataSourceName, testAccACLsDataSourceAttrKeyPrefix+"destination_ip_address", "192.168.1.0/24"),
					resource.TestCheckResourceAttr(testAccACLsDataSourceName, testAccACLsDataSourceAttrKeyPrefix+"source_port", "8888"),
					resource.TestCheckResourceAttr(testAccACLsDataSourceName, testAccACLsDataSourceAttrKeyPrefix+"destination_port", "9999"),
					resource.TestCheckResourceAttr(testAccACLsDataSourceName, testAccACLsDataSourceAttrKeyPrefix+"position", "20"),
					resource.TestCheckResourceAttr(testAccACLsDataSourceName, testAccACLsDataSourceAttrKeyPrefix+"direction", "ingress"),
					resource.TestCheckResourceAttr(testAccACLsDataSourceName, testAccACLsDataSourceAttrKeyPrefix+"action", "allow"),
				),
			},
			{
				Config: testAccACLsDataSourceConfigByVPC(BaiduCloudTestResourceTypeNameAcl),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccACLsDataSourceName),
					resource.TestCheckResourceAttrSet(testAccACLsDataSourceName, testAccACLsDataSourceAttrKeyPrefix+"acl_id"),
					resource.TestCheckResourceAttrSet(testAccACLsDataSourceName, testAccACLsDataSourceAttrKeyPrefix+"subnet_id"),
					resource.TestCheckResourceAttr(testAccACLsDataSourceName, testAccACLsDataSourceAttrKeyPrefix+"description", "created by terraform"),
					resource.TestCheckResourceAttr(testAccACLsDataSourceName, testAccACLsDataSourceAttrKeyPrefix+"protocol", "tcp"),
					resource.TestCheckResourceAttr(testAccACLsDataSourceName, testAccACLsDataSourceAttrKeyPrefix+"source_ip_address", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(testAccACLsDataSourceName, testAccACLsDataSourceAttrKeyPrefix+"destination_ip_address", "192.168.1.0/24"),
					resource.TestCheckResourceAttr(testAccACLsDataSourceName, testAccACLsDataSourceAttrKeyPrefix+"source_port", "8888"),
					resource.TestCheckResourceAttr(testAccACLsDataSourceName, testAccACLsDataSourceAttrKeyPrefix+"destination_port", "9999"),
					resource.TestCheckResourceAttr(testAccACLsDataSourceName, testAccACLsDataSourceAttrKeyPrefix+"position", "20"),
					resource.TestCheckResourceAttr(testAccACLsDataSourceName, testAccACLsDataSourceAttrKeyPrefix+"direction", "ingress"),
					resource.TestCheckResourceAttr(testAccACLsDataSourceName, testAccACLsDataSourceAttrKeyPrefix+"action", "allow"),
				),
			},
		},
	})
}

func testAccACLsDataSourceConfigBySubnet(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "baiducloud_zones" "default" {
  name_regex = ".*e$"
}

resource "baiducloud_vpc" "default" {
  name = "${var.name}"
  cidr = "192.168.0.0/16"
}

resource "baiducloud_subnet" "default" {
  name = "${var.name}"
  zone_name = "${data.baiducloud_zones.default.zones.0.zone_name}"
  cidr = "192.168.1.0/24"
  vpc_id = "${baiducloud_vpc.default.id}"
}

resource "baiducloud_acl" "default" {
  subnet_id = "${baiducloud_subnet.default.id}"
  protocol = "tcp"
  source_ip_address = "192.168.0.0/24"
  destination_ip_address = "192.168.1.0/24"
  source_port = "8888"
  destination_port = "9999"
  position = 20
  direction = "ingress"
  action = "allow"
  description = "created by terraform"
}

data "baiducloud_acls" "default" {
  acl_id = "${baiducloud_acl.default.id}"
  subnet_id = "${baiducloud_subnet.default.id}"

  filter {
    name = "direction"
    values = ["ingress"]
  }
}
`, name)
}

func testAccACLsDataSourceConfigByVPC(name string) string {
	return fmt.Sprintf(`

variable "name" {
  default = "%s"
}

data "baiducloud_zones" "default" {
  name_regex = ".*e$"
}

resource "baiducloud_vpc" "default" {
  name = "${var.name}"
  cidr = "192.168.0.0/16"
}

resource "baiducloud_subnet" "default" {
  name = "${var.name}"
  zone_name = "${data.baiducloud_zones.default.zones.0.zone_name}"
  cidr = "192.168.1.0/24"
  vpc_id = "${baiducloud_vpc.default.id}"
}

resource "baiducloud_acl" "default" {
  subnet_id = "${baiducloud_subnet.default.id}"
  protocol = "tcp"
  source_ip_address = "192.168.0.0/24"
  destination_ip_address = "192.168.1.0/24"
  source_port = "8888"
  destination_port = "9999"
  position = 20
  direction = "ingress"
  action = "allow"
  description = "created by terraform"
}

data "baiducloud_acls" "default" {
  acl_id = "${baiducloud_acl.default.id}"
  vpc_id = "${baiducloud_vpc.default.id}"

  filter {
    name = "action"
    values = ["allow"]
  }
}
`, name)
}
