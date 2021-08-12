package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccSubnetResourceType = "baiducloud_subnet"
	testAccSubnetResourceName = testAccSubnetResourceType + "." + BaiduCloudTestResourceName
)

//lintignore:AT003
func TestAccBaiduCloudSubnet(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSubnetDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccSubnetConfig(BaiduCloudTestResourceTypeNameSubnet),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccSubnetResourceName),
					resource.TestCheckResourceAttr(testAccSubnetResourceName, "name", BaiduCloudTestResourceTypeNameSubnet),
					resource.TestCheckResourceAttr(testAccSubnetResourceName, "cidr", "192.168.3.0/24"),
					resource.TestCheckResourceAttr(testAccSubnetResourceName, "description", "created by terraform"),
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
				Config: testAccSubnetConfigUpdate(BaiduCloudTestResourceTypeNameSubnet),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccSubnetResourceName),
					resource.TestCheckResourceAttr(testAccSubnetResourceName, "name", BaiduCloudTestResourceTypeNameSubnet+"-update"),
					resource.TestCheckResourceAttr(testAccSubnetResourceName, "cidr", "192.168.3.0/24"),
					resource.TestCheckResourceAttr(testAccSubnetResourceName, "description", "created by terraform"),
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

func testAccSubnetConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "baiducloud_zones" "default" {
  name_regex = ".*e$"
}

resource "baiducloud_vpc" "default" {
  name = var.name
  cidr = "192.168.0.0/16"
}

resource "baiducloud_subnet" "default" {
  name        = var.name
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.3.0/24"
  description = "created by terraform"
  vpc_id      = baiducloud_vpc.default.id
  subnet_type = "BCC"
  tags = {
    "tagKey" = "tagValue"
  }
}
`, name)
}

func testAccSubnetConfigUpdate(name string) string {
	return fmt.Sprintf(`
data "baiducloud_zones" "default" {
  name_regex = ".*e$"
}

resource "baiducloud_vpc" "default" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "baiducloud_subnet" "default" {
  name        = "%s"
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.3.0/24"
  description = "created by terraform"
  vpc_id      = baiducloud_vpc.default.id
  subnet_type = "BCC"
  tags = {
    "tagKey" = "tagValue"
  }
}
`, name, name+"-update")
}
