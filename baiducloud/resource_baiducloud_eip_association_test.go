package baiducloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/baidubce/bce-sdk-go/services/eip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccEipAssociationResourceType = "baiducloud_eip_association"
	testAccEipAssociationResourceName = testAccEipAssociationResourceType + "." + BaiduCloudTestResourceName
)

func init() {
	resource.AddTestSweepers(testAccEipAssociationResourceType, &resource.Sweeper{
		Name: testAccEipAssociationResourceType,
		F:    testSweepEipsAssociate,
	})
}

func testSweepEipsAssociate(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("get BaiduCloud client error: %s", err)
	}
	client := rawClient.(*connectivity.BaiduClient)
	eipService := EipService{client}

	listArgs := &eip.ListEipArgs{}
	eipList, err := eipService.ListAllEips(listArgs)
	if err != nil {
		return fmt.Errorf("get EIPs error: %s", err)
	}

	for _, e := range eipList {
		name := e.Name
		ip := e.Eip
		if !strings.HasPrefix(e.Name, BaiduCloudTestResourceTypeName) || e.Status != EIPStatusBinded {
			log.Printf("[INFO] Skipping EIP: %s (%s)", name, ip)
			continue
		}

		log.Printf("[INFO] Unbind EIP: %s (%s)", e.Name, e.Eip)
		err = eipService.EipUnBind(e.Eip)
		if err != nil {
			log.Printf("[ERROR] Unbind to delete EIP %s (%s)", name, ip)
		}
	}

	return nil
}

//lintignore:AT003
func TestAccBaiduCloudEipAssociate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccEIPAssociateDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccEipAssociateConfig(BaiduCloudTestResourceTypeNameEipAssociation),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccEipAssociationResourceName),
					resource.TestCheckResourceAttrSet(testAccEipAssociationResourceName, "eip"),
					resource.TestCheckResourceAttrSet(testAccEipAssociationResourceName, "instance_id"),
					resource.TestCheckResourceAttrSet(testAccEipAssociationResourceName, "instance_type"),
				),
			},
			{
				ResourceName:      testAccEipAssociationResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEIPAssociateDestory(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	eipService := EipService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccEipAssociationResourceType {
			continue
		}

		result, err := eipService.EipGetDetail(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		if result.Status == EIPStatusBinded {
			return WrapError(Error("EIP association still exist"))
		}
	}

	return nil
}

func testAccEipAssociateConfig(name string) string {
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
  name                  = var.name
  image_id              = data.baiducloud_images.default.images.0.id
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  billing = {
    payment_timing = "Postpaid"
  }
}

resource "baiducloud_vpc" "default" {
  name        = var.name
  description = "created by terraform"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = var.name
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "created by terraform"
}

resource "baiducloud_eip" "default" {
  name              = var.name
  bandwidth_in_mbps = 1
  payment_timing    = "Postpaid"
  billing_method    = "ByTraffic"
}

resource "baiducloud_appblb" "default" {
  depends_on  = [baiducloud_instance.default]
  name        = var.name
  description = "created by terraform"
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}

resource "baiducloud_eip_association" "default" {
  eip           = baiducloud_eip.default.id
  instance_type = "BLB"
  instance_id   = baiducloud_appblb.default.id
}
`, name)
}
