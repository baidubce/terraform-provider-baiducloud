package baiducloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/baidubce/bce-sdk-go/services/appblb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccAppBLBResourceType = "baiducloud_appblb"
	testAccAppBLBResourceName = testAccAppBLBResourceType + "." + BaiduCloudTestResourceName
)

func init() {
	resource.AddTestSweepers(testAccAppBLBResourceType, &resource.Sweeper{
		Name: testAccAppBLBResourceType,
		F:    testSweepAppBLBs,
		Dependencies: []string{
			testAccInstanceResourceType,
			testAccVPCResourceType,
		},
	})
}

func testSweepAppBLBs(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("get BaiduCloud client error: %s", err)
	}
	client := rawClient.(*connectivity.BaiduClient)
	appblbService := APPBLBService{client}

	listArgs := &appblb.DescribeLoadBalancersArgs{}
	appblbList, _, err := appblbService.ListAllAppBLB(listArgs)
	if err != nil {
		return fmt.Errorf("get APPBLBs error: %s", err)
	}

	for _, blb := range appblbList {
		name := blb.Name
		blbId := blb.BlbId
		if !strings.HasPrefix(name, BaiduCloudTestResourceTypeName) {
			log.Printf("[INFO] Skipping APPBLB: %s (%s)", name, blbId)
			continue
		}

		log.Printf("[INFO] Deleting APPBLB: %s (%s)", name, blbId)
		_, err := client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
			return nil, client.DeleteLoadBalancer(blbId)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete APPBLB %s (%s)", name, blbId)
		}
	}

	return nil
}

//lintignore:AT003
func TestAccBaiduCloudAppBLB(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccAppBLBDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccAppBLBConfig(BaiduCloudTestResourceTypeNameAppblb),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBResourceName),
					resource.TestCheckResourceAttr(testAccAppBLBResourceName, "name", BaiduCloudTestResourceTypeNameAppblb),
					resource.TestCheckResourceAttr(testAccAppBLBResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(testAccAppBLBResourceName, "cidr", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(testAccAppBLBResourceName, "subnet_cidr", "192.168.0.0/24"),
					resource.TestCheckResourceAttrSet(testAccAppBLBResourceName, "create_time"),
					resource.TestCheckResourceAttrSet(testAccAppBLBResourceName, "vpc_id"),
					resource.TestCheckResourceAttrSet(testAccAppBLBResourceName, "subnet_id"),
					resource.TestCheckResourceAttrSet(testAccAppBLBResourceName, "vpc_name"),
					resource.TestCheckResourceAttrSet(testAccAppBLBResourceName, "subnet_name"),
				),
			},
			{
				ResourceName:      testAccAppBLBResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAppBLBConfigUpdate(BaiduCloudTestResourceTypeNameAppblb),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBResourceName),
					resource.TestCheckResourceAttr(testAccAppBLBResourceName, "name", BaiduCloudTestResourceTypeNameAppblb+"-update"),
					resource.TestCheckResourceAttr(testAccAppBLBResourceName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(testAccAppBLBResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(testAccAppBLBResourceName, "cidr", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(testAccAppBLBResourceName, "subnet_cidr", "192.168.0.0/24"),
					resource.TestCheckResourceAttrSet(testAccAppBLBResourceName, "create_time"),
					resource.TestCheckResourceAttrSet(testAccAppBLBResourceName, "vpc_id"),
					resource.TestCheckResourceAttrSet(testAccAppBLBResourceName, "subnet_id"),
					resource.TestCheckResourceAttrSet(testAccAppBLBResourceName, "vpc_name"),
					resource.TestCheckResourceAttrSet(testAccAppBLBResourceName, "subnet_name"),
				),
			},
		},
	})
}

func testAccAppBLBDestory(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	appblbService := APPBLBService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccAppBLBResourceType {
			continue
		}

		_, _, err := appblbService.GetAppBLBDetail(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		return WrapError(Error("APPBLB still exist"))
	}

	return nil
}

func testAccAppBLBConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "baiducloud_specs" "default" {}

data "baiducloud_zones" "default" {
  name_regex = ".*e$"
}

data "baiducloud_images" "default" {
  image_type = "System"
}

resource "baiducloud_instance" "default" {
  name                  = "${var.name}"
  image_id              = data.baiducloud_images.default.images.0.id
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  billing = {
    payment_timing = "Postpaid"
  }
}

resource "baiducloud_vpc" "default" {
  name        = "${var.name}"
  description = "created by terraform"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = "${var.name}"
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created by terraform"
}

resource "baiducloud_appblb" "default" {
  depends_on  = [baiducloud_instance.default]
  name        = "${var.name}"
  description = "created by terraform"
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id

  tags = {
    "testKey" = "testValue"
  }
}
`, name)
}

func testAccAppBLBConfigUpdate(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "baiducloud_specs" "default" {}

data "baiducloud_zones" "default" {
  name_regex = ".*e$"
}

data "baiducloud_images" "default" {
  image_type = "System"
}

resource "baiducloud_instance" "default" {
  name                  = "${var.name}"
  image_id              = data.baiducloud_images.default.images.0.id
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  billing = {
    payment_timing = "Postpaid"
  }
}

resource "baiducloud_vpc" "default" {
  name        = "${var.name}"
  description = "created by terraform"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = "${var.name}"
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created by terraform"
}

resource "baiducloud_eip" "default" {
  name              = "${var.name}"
  bandwidth_in_mbps = 1
  payment_timing    = "Postpaid"
  billing_method    = "ByTraffic"
}

resource "baiducloud_eip_association" "default" {
  eip           = baiducloud_eip.default.id
  instance_type = "BLB"
  instance_id   = baiducloud_appblb.default.id
}

resource "baiducloud_appblb" "default" {
  depends_on  = [baiducloud_instance.default]
  name        = "%s"
  description = "created by terraform"
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id

  tags = {
    "testKey" = "testValue"
  }
}
`, name, name+"-update")
}
