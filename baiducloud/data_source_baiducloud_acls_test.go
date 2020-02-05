package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
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
				Config: testAccACLsDataSourceConfigBySubnet(),
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
				Config: testAccACLsDataSourceConfigByVPC(),
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

func testAccACLsDataSourceConfigBySubnet() string {
	return fmt.Sprintf(`
data "baiducloud_zones" "default" {}

resource "baiducloud_vpc" "default" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "baiducloud_subnet" "default" {
  name = "%s"
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
}
`, BaiduCloudTestResourceAttrNamePrefix+"VPC", BaiduCloudTestResourceAttrNamePrefix+"Subnet")
}

func testAccACLsDataSourceConfigByVPC() string {
	return fmt.Sprintf(`
data "baiducloud_zones" "default" {}

resource "baiducloud_vpc" "default" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "baiducloud_subnet" "default" {
  name = "%s"
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
}
`, BaiduCloudTestResourceAttrNamePrefix+"VPC", BaiduCloudTestResourceAttrNamePrefix+"Subnet")
}
