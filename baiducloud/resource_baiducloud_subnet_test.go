package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccSubnetResourceType     = "baiducloud_subnet"
	testAccSubnetResourceName     = testAccSubnetResourceType + "." + BaiduCloudTestResourceName
	testAccSubnetResourceAttrName = BaiduCloudTestResourceAttrNamePrefix + "Subnet"
)

func TestAccBaiduCloudSubnet(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSubnetDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccSubnetConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccSubnetResourceName),
					resource.TestCheckResourceAttr(testAccSubnetResourceName, "name", testAccSubnetResourceAttrName),
					resource.TestCheckResourceAttr(testAccSubnetResourceName, "cidr", "192.168.3.0/24"),
					resource.TestCheckResourceAttr(testAccSubnetResourceName, "description", "test"),
					resource.TestCheckResourceAttr(testAccSubnetResourceName, "subnet_type", "BCC"),
					resource.TestCheckResourceAttrSet(testAccSubnetResourceName, "vpc_id"),
					resource.TestCheckResourceAttrSet(testAccSubnetResourceName, "zone_name"),
					resource.TestCheckResourceAttr(testAccSubnetResourceName, "tags.%", "1"),
				),
			},
			{
				ResourceName:      testAccSubnetResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSubnetConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccSubnetResourceName),
					resource.TestCheckResourceAttr(testAccSubnetResourceName, "name", testAccSubnetResourceAttrName+"Update"),
					resource.TestCheckResourceAttr(testAccSubnetResourceName, "cidr", "192.168.3.0/24"),
					resource.TestCheckResourceAttr(testAccSubnetResourceName, "description", "test update"),
					resource.TestCheckResourceAttr(testAccSubnetResourceName, "subnet_type", "BCC"),
					resource.TestCheckResourceAttrSet(testAccSubnetResourceName, "vpc_id"),
					resource.TestCheckResourceAttrSet(testAccSubnetResourceName, "zone_name"),
					resource.TestCheckResourceAttr(testAccSubnetResourceName, "tags.%", "1"),
				),
			},
		},
	})
}

func testAccSubnetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	vpcService := &VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccSubnetResourceType {
			continue
		}

		_, err := vpcService.GetSubnetDetail(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		return WrapError(Error("Subnet still exist"))
	}

	return nil
}

func testAccSubnetConfig() string {
	return fmt.Sprintf(`
data "baiducloud_zones" "default" {}

resource "baiducloud_vpc" "default" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "%s" "%s" {
  name        = "%s"
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.3.0/24"
  description = "test"
  vpc_id      = baiducloud_vpc.default.id
  subnet_type = "BCC"
  tags = {
    "tagKey" = "tagValue"
  }
}
`, BaiduCloudTestResourceAttrNamePrefix+"VPC", testAccSubnetResourceType,
		BaiduCloudTestResourceName, testAccSubnetResourceAttrName)
}

func testAccSubnetConfigUpdate() string {
	return fmt.Sprintf(`
data "baiducloud_zones" "default" {}

resource "baiducloud_vpc" "default" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "%s" "%s" {
  name        = "%s"
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.3.0/24"
  description = "test update"
  vpc_id      = baiducloud_vpc.default.id
  subnet_type = "BCC"
  tags = {
    "tagKey" = "tagValue"
  }
}
`, BaiduCloudTestResourceAttrNamePrefix+"VPC", testAccSubnetResourceType,
		BaiduCloudTestResourceName, testAccSubnetResourceAttrName+"Update")
}
