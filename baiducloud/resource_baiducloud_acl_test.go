package baiducloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccACLResourceType = "baiducloud_acl"
	testAccACLResourceName = testAccACLResourceType + "." + BaiduCloudTestResourceName
)

func init() {
	resource.AddTestSweepers(testAccACLResourceType, &resource.Sweeper{
		Name: testAccACLResourceType,
		F:    testSweepACLs,
		Dependencies: []string{
			testAccInstanceResourceType,
			testAccAppBLBResourceType,
			testAccPeerConnResourceType,
			testAccCcev2ClusterResourceType,
			testAccVPCResourceType,
		},
	})
}

func testSweepACLs(region string) error {
	log.Printf("[INFO] Skipping ACL,Nothing to do)")
	return nil
}

//lintignore:AT003
func TestAccBaiduCloudACL(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccACLDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccACLConfig(BaiduCloudTestResourceTypeNameAcl),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccACLResourceName),
					resource.TestCheckResourceAttrSet(testAccACLResourceName, "subnet_id"),
					resource.TestCheckResourceAttr(testAccACLResourceName, "protocol", "tcp"),
					resource.TestCheckResourceAttr(testAccACLResourceName, "source_ip_address", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(testAccACLResourceName, "destination_ip_address", "192.168.1.0/24"),
					resource.TestCheckResourceAttr(testAccACLResourceName, "source_port", "8888"),
					resource.TestCheckResourceAttr(testAccACLResourceName, "destination_port", "9999"),
					resource.TestCheckResourceAttr(testAccACLResourceName, "position", "20"),
					resource.TestCheckResourceAttr(testAccACLResourceName, "direction", "ingress"),
					resource.TestCheckResourceAttr(testAccACLResourceName, "action", "allow"),
					resource.TestCheckResourceAttr(testAccACLResourceName, "description", "created by terraform"),
				),
			},
			{
				Config: testAccACLConfigUpdate(BaiduCloudTestResourceTypeNameAcl),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccACLResourceName),
					resource.TestCheckResourceAttrSet(testAccACLResourceName, "subnet_id"),
					resource.TestCheckResourceAttr(testAccACLResourceName, "protocol", "udp"),
					resource.TestCheckResourceAttr(testAccACLResourceName, "source_ip_address", "192.168.2.0/24"),
					resource.TestCheckResourceAttr(testAccACLResourceName, "destination_ip_address", "192.168.1.0/24"),
					resource.TestCheckResourceAttr(testAccACLResourceName, "source_port", "6666"),
					resource.TestCheckResourceAttr(testAccACLResourceName, "destination_port", "7777"),
					resource.TestCheckResourceAttr(testAccACLResourceName, "position", "30"),
					resource.TestCheckResourceAttr(testAccACLResourceName, "direction", "ingress"),
					resource.TestCheckResourceAttr(testAccACLResourceName, "action", "allow"),
					resource.TestCheckResourceAttr(testAccACLResourceName, "description", "updated by terraform"),
				),
			},
		},
	})
}

func testAccACLDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	vpcService := &VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccACLResourceType {
			continue
		}

		subnetID := rs.Primary.Attributes["subnet_id"]
		_, err := vpcService.GetSubnetDetail(subnetID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		result, err := vpcService.ListAllAclRulesWithSubnetID(subnetID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		for _, acl := range result {
			if acl.Id == rs.Primary.ID {
				return WrapError(Error("ACL still exist"))
			}
		}
	}

	return nil
}

func testAccACLConfig(name string) string {
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
  name      = "${var.name}"
  zone_name = data.baiducloud_zones.default.zones.0.zone_name
  cidr      = "192.168.1.0/24"
  vpc_id    = baiducloud_vpc.default.id
}

resource "baiducloud_acl" "default" {
  subnet_id              = baiducloud_subnet.default.id
  protocol               = "tcp"
  source_ip_address      = "192.168.0.0/24"
  destination_ip_address = "192.168.1.0/24"
  source_port            = "8888"
  destination_port       = "9999"
  position               = 20
  direction              = "ingress"
  action                 = "allow"
  description            = "created by terraform"
}
`, name)
}

func testAccACLConfigUpdate(name string) string {
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
  zone_name = data.baiducloud_zones.default.zones.0.zone_name
  cidr      = "192.168.1.0/24"
  vpc_id    = baiducloud_vpc.default.id
}

resource "baiducloud_acl" "default" {
  subnet_id              = baiducloud_subnet.default.id
  protocol               = "udp"
  source_ip_address      = "192.168.2.0/24"
  destination_ip_address = "192.168.1.0/24"
  source_port            = "6666"
  destination_port       = "7777"
  position               = 30
  direction              = "ingress"
  action                 = "allow"
  description            = "updated by terraform"
}
`, name)
}
