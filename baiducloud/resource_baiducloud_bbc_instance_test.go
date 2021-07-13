package baiducloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/baidubce/bce-sdk-go/services/bbc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccBbcInstanceResourceType = "baiducloud_bbc_instance"
	testAccBbcInstanceResourceName = testAccBbcInstanceResourceType + "." + BaiduCloudTestResourceName
)

func init() {
	resource.AddTestSweepers(testAccBbcInstanceResourceType, &resource.Sweeper{
		Name: testAccBbcInstanceResourceType,
		F:    testSweepBbcInstances,
		Dependencies: []string{
			testAccEipAssociationResourceType,
		},
	})
}

func testSweepBbcInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("get BaiduCloud client error: %s", err)
	}

	client := rawClient.(*connectivity.BaiduClient)
	bbcService := &BbcService{client}

	args := &bbc.ListInstancesArgs{}
	instList, err := bbcService.ListAllInstance(args)
	if err != nil {
		return fmt.Errorf("get BBC instances error: %v", err)
	}

	for _, inst := range instList {
		if !strings.HasPrefix(inst.Name, BaiduCloudTestResourceAttrNamePrefix) {
			log.Printf("[INFO] Skipping BBC instance: %s (%s)", inst.Id, inst.Name)
			continue
		}

		log.Printf("[INFO] Deleting BBC instance: %s (%s)", inst.Id, inst.Name)
		_, err := client.WithBbcClient(func(bbcClient *bbc.Client) (i interface{}, e error) {
			return nil, bbcClient.DeleteInstance(inst.Id)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete BBC instance %s (%s)", inst.Id, inst.Name)
		}
	}

	return nil
}

//lintignore:AT003
func TestAccBaiduCloudBbcInstance(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccBbcInstanceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccBbcInstanceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBbcInstanceResourceName),
					resource.TestCheckResourceAttr(testAccBbcInstanceResourceName, "name", BaiduCloudTestResourceAttrNamePrefix+"BBC"),
					resource.TestCheckResourceAttr(testAccBbcInstanceResourceName, "description", "terraform test instance"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "image_id"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "flavor_id"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "availability_zone"),
					resource.TestCheckResourceAttr(testAccBbcInstanceResourceName, "billing.payment_timing", "Postpaid"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "subnet_id"),
					resource.TestCheckResourceAttr(testAccBbcInstanceResourceName, "security_group", "1"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "status"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "create_time"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "internal_ip"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "placement_policy"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "vpc_id"),
					resource.TestCheckResourceAttr(testAccBbcInstanceResourceName, "tags.%", "1"),
				),
			},
			{
				ResourceName:            testAccBbcInstanceResourceName,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"auto_renew_time_length"},
			},
			{
				Config: testAccBbcInstanceConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBbcInstanceResourceName),
					resource.TestCheckResourceAttr(testAccBbcInstanceResourceName, "name", BaiduCloudTestResourceAttrNamePrefix+"BBC-update"),
					resource.TestCheckResourceAttr(testAccBbcInstanceResourceName, "description", "terraform test update instance"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "image_id"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "flavor_id"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "raid_id"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "availability_zone"),
					resource.TestCheckResourceAttr(testAccBbcInstanceResourceName, "billing.payment_timing", "Postpaid"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "subnet_id"),
					resource.TestCheckResourceAttr(testAccBbcInstanceResourceName, "security_group", "1"),
					resource.TestCheckResourceAttr(testAccBbcInstanceResourceName, "status", "Running"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "create_time"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "internal_ip"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "placement_policy"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "vpc_id"),
					resource.TestCheckResourceAttr(testAccBbcInstanceResourceName, "tags.%", "1"),
				),
			},
			{
				Config: testAccBbcInstanceActionUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBbcInstanceResourceName),
					resource.TestCheckResourceAttr(testAccBbcInstanceResourceName, "name", BaiduCloudTestResourceAttrNamePrefix+"BBC-update"),
					resource.TestCheckResourceAttr(testAccBbcInstanceResourceName, "description", "terraform test update instance"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "image_id"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "flavor_id"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "raid_id"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "availability_zone"),
					resource.TestCheckResourceAttr(testAccBbcInstanceResourceName, "billing.payment_timing", "Postpaid"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "subnet_id"),
					resource.TestCheckResourceAttr(testAccBbcInstanceResourceName, "security_group", "1"),
					resource.TestCheckResourceAttr(testAccBbcInstanceResourceName, "status", "Stopped"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "create_time"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "internal_ip"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "placement_policy"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "vpc_id"),
					resource.TestCheckResourceAttr(testAccBbcInstanceResourceName, "tags.%", "1"),
				),
			},
		},
	})
}

func testAccBbcInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	bccService := &BccService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccBbcInstanceResourceType {
			continue
		}

		_, err := bccService.GetInstanceDetail(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		return WrapError(Error("BBC instance still exist"))
	}

	return nil
}

func testAccBbcInstanceConfig() string {
	return fmt.Sprintf(`
data "baiducloud_bbc_flavors" "default" {

}

data "baiducloud_zones" "default" {}

data "baiducloud_bbc_images" "default" {}

data "baiducloud_bbc_raids" "default" {
	flavor_id = data.baiducloud_bbc_flavors.default.flavors.0.flavor_id
}

resource "baiducloud_vpc" "default" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "baiducloud_subnet" "default" {
  name      = "%s"
  zone_name = data.baiducloud_zones.default.zones.0.zone_name
  cidr      = "192.168.1.0/24"
  vpc_id    = baiducloud_vpc.default.id
}

resource "baiducloud_security_group" "default" {
  name        = "%s"
  description = "security group created by terraform"
  vpc_id      = baiducloud_vpc.default.id
}

resource "baiducloud_bbc_instance" "default" {
  image_id              = data.baiducloud_bbc_images.default.images.0.id
  flavor_id 			= data.baiducloud_bbc_flavors.default.flavors.0.flavor_id
  raid_id 				= data.baiducloud_bbc_raids.default.raids.0.raid_id
  name                  = "%s"
  description           = "terraform test instance"
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  billing = {
    payment_timing = "Postpaid"
  }

  subnet_id       = baiducloud_subnet.default.id
  security_group = baiducloud_security_group.default.id

  tags = {
    "testKey" = "testValue"
  }
}

data "baiducloud_bbc_instances" "default" {
	vpc_id    = baiducloud_vpc.default.id
}
`, BaiduCloudTestResourceAttrNamePrefix+"VPC",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceAttrNamePrefix+"SG",
		BaiduCloudTestResourceAttrNamePrefix+"BBC")
}

func testAccBbcInstanceConfigUpdate() string {
	return fmt.Sprintf(`
data "baiducloud_bbc_flavors" "default" {

}

data "baiducloud_zones" "default" {}

data "baiducloud_bbc_images" "default" {}

data "baiducloud_bbc_raids" "default" {
	flavor_id = data.baiducloud_bbc_flavors.default.flavors.0.flavor_id
}

resource "baiducloud_eip" "default" {
  name              = "%s"
  bandwidth_in_mbps = 1
  payment_timing    = "Postpaid"
  billing_method    = "ByTraffic"
}

resource "baiducloud_vpc" "default" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "baiducloud_subnet" "default" {
  name      = "%s"
  zone_name = data.baiducloud_zones.default.zones.0.zone_name
  cidr      = "192.168.1.0/24"
  vpc_id    = baiducloud_vpc.default.id
}

resource "baiducloud_security_group" "default" {
  name        = "%s"
  description = "security group created by terraform"
  vpc_id      = baiducloud_vpc.default.id
}

resource "baiducloud_subnet" "default02" {
  name      = "%s"
  zone_name = data.baiducloud_zones.default.zones.0.zone_name
  cidr      = "192.168.2.0/24"
  vpc_id    = baiducloud_vpc.default.id
}

resource "baiducloud_security_group" "default02" {
  name        = "%s"
  description = "security group created by terraform"
  vpc_id      = baiducloud_vpc.default.id
}

resource "baiducloud_bbc_instance" "default" {
	image_id              = data.baiducloud_bbc_images.default.images.0.id
	flavor_id 			= data.baiducloud_bbc_flavors.default.flavors.0.flavor_id
	raid_id 				= data.baiducloud_bbc_raids.default.raids.0.raid_id
  name                  = "%s"
  description           = "terraform test update instance"
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  billing = {
    payment_timing = "Postpaid"
  }
  admin_pass = "terraform@123"

  subnet_id       = baiducloud_subnet.default02.id
  security_group = baiducloud_security_group.default02.id

  tags = {
    "testKey" = "testValue"
  }
}

resource "baiducloud_eip_association" "default" {
  eip           = baiducloud_eip.default.id
  instance_type = "BBC"
  instance_id   = baiducloud_bbc_instance.default.id
}
`, BaiduCloudTestResourceAttrNamePrefix+"EIP",
		BaiduCloudTestResourceAttrNamePrefix+"VPC",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceAttrNamePrefix+"SG",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet02",
		BaiduCloudTestResourceAttrNamePrefix+"SG02",
		BaiduCloudTestResourceAttrNamePrefix+"BBC-update")
}

func testAccBbcInstanceActionUpdate() string {
	return fmt.Sprintf(`
data "baiducloud_bbc_flavors" "default" {

}

data "baiducloud_zones" "default" {}

data "baiducloud_bbc_images" "default" {}

data "baiducloud_bbc_raids" "default" {
	flavor_id = data.baiducloud_bbc_flavors.default.flavors.0.flavor_id
}

resource "baiducloud_eip" "default" {
  name              = "%s"
  bandwidth_in_mbps = 1
  payment_timing    = "Postpaid"
  billing_method    = "ByTraffic"
}

resource "baiducloud_vpc" "default" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "baiducloud_subnet" "default" {
  name      = "%s"
  zone_name = data.baiducloud_zones.default.zones.0.zone_name
  cidr      = "192.168.1.0/24"
  vpc_id    = baiducloud_vpc.default.id
}

resource "baiducloud_security_group" "default" {
  name        = "%s"
  description = "security group created by terraform"
  vpc_id      = baiducloud_vpc.default.id
}

resource "baiducloud_subnet" "default02" {
  name      = "%s"
  zone_name = data.baiducloud_zones.default.zones.0.zone_name
  cidr      = "192.168.2.0/24"
  vpc_id    = baiducloud_vpc.default.id
}

resource "baiducloud_security_group" "default02" {
  name        = "%s"
  description = "security group created by terraform"
  vpc_id      = baiducloud_vpc.default.id
}

resource "baiducloud_bbc_instance" "default" {
	image_id              = data.baiducloud_bbc_images.default.images.0.id
	flavor_id 			= data.baiducloud_bbc_flavors.default.flavors.0.flavor_id
	raid_id 				= data.baiducloud_bbc_raids.default.raids.0.raid_id
  name                  = "%s"
  description           = "terraform test update instance"
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  billing = {
    payment_timing = "Postpaid"
  }
  admin_pass = "terraform@123"

  subnet_id       = baiducloud_subnet.default02.id
  security_group = baiducloud_security_group.default02.id

  tags = {
    "testKey1" = "testValue"
  }

}

resource "baiducloud_eip_association" "default" {
  eip           = baiducloud_eip.default.id
  instance_type = "BCC"
  instance_id   = baiducloud_bbc_instance.default.id
}
`, BaiduCloudTestResourceAttrNamePrefix+"EIP",
		BaiduCloudTestResourceAttrNamePrefix+"VPC",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceAttrNamePrefix+"SG",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet02",
		BaiduCloudTestResourceAttrNamePrefix+"SG02",
		BaiduCloudTestResourceAttrNamePrefix+"BBC-update")
}
