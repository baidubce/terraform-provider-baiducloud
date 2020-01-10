package baiducloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccInstanceResourceType = "baiducloud_instance"
	testAccInstanceResourceName = testAccInstanceResourceType + "." + BaiduCloudTestResourceName
)

func init() {
	resource.AddTestSweepers(testAccInstanceResourceType, &resource.Sweeper{
		Name: testAccInstanceResourceType,
		F:    testSweepInstances,
		Dependencies: []string{
			testAccEipAssociationResourceType,
		},
	})
}

func testSweepInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("get BaiduCloud client error: %s", err)
	}

	client := rawClient.(*connectivity.BaiduClient)
	bccService := &BccService{client}

	args := &api.ListInstanceArgs{}
	instList, err := bccService.ListAllInstance(args)
	if err != nil {
		return fmt.Errorf("get BCC instances error: %v", err)
	}

	for _, inst := range instList {
		if !strings.HasPrefix(inst.InstanceName, BaiduCloudTestResourceAttrNamePrefix) {
			log.Printf("[INFO] Skipping BCC instance: %s (%s)", inst.InstanceId, inst.InstanceName)
			continue
		}

		log.Printf("[INFO] Deleting BCC instance: %s (%s)", inst.InstanceId, inst.InstanceName)
		_, err := client.WithBccClient(func(bccClient *bcc.Client) (i interface{}, e error) {
			return nil, bccClient.DeleteInstance(inst.InstanceId)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete BCC instance %s (%s)", inst.InstanceId, inst.InstanceName)
		}
	}

	return nil
}

func TestAccBaiduCloudInstance(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccInstanceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccInstanceResourceName),
					resource.TestCheckResourceAttr(testAccInstanceResourceName, "name", BaiduCloudTestResourceAttrNamePrefix+"BCC"),
					resource.TestCheckResourceAttr(testAccInstanceResourceName, "description", "terraform test instance"),
					resource.TestCheckResourceAttrSet(testAccInstanceResourceName, "image_id"),
					resource.TestCheckResourceAttrSet(testAccInstanceResourceName, "availability_zone"),
					resource.TestCheckResourceAttrSet(testAccInstanceResourceName, "cpu_count"),
					resource.TestCheckResourceAttrSet(testAccInstanceResourceName, "memory_capacity_in_gb"),
					resource.TestCheckResourceAttrSet(testAccInstanceResourceName, "root_disk_size_in_gb"),
					resource.TestCheckResourceAttrSet(testAccInstanceResourceName, "root_disk_storage_type"),
					resource.TestCheckResourceAttr(testAccInstanceResourceName, "ephemeral_disks.#", "0"),
					resource.TestCheckResourceAttr(testAccInstanceResourceName, "billing.payment_timing", "Postpaid"),
					resource.TestCheckResourceAttrSet(testAccInstanceResourceName, "subnet_id"),
					resource.TestCheckResourceAttr(testAccInstanceResourceName, "security_groups.#", "1"),
					resource.TestCheckResourceAttrSet(testAccInstanceResourceName, "status"),
					resource.TestCheckResourceAttrSet(testAccInstanceResourceName, "create_time"),
					resource.TestCheckResourceAttrSet(testAccInstanceResourceName, "internal_ip"),
					resource.TestCheckResourceAttrSet(testAccInstanceResourceName, "placement_policy"),
					resource.TestCheckResourceAttrSet(testAccInstanceResourceName, "vpc_id"),
					resource.TestCheckResourceAttr(testAccInstanceResourceName, "cds_disks.#", "1"),
					resource.TestCheckResourceAttr(testAccInstanceResourceName, "cds_disks.0.cds_size_in_gb", "50"),
					resource.TestCheckResourceAttr(testAccInstanceResourceName, "cds_disks.0.storage_type", "cloud_hp1"),
					resource.TestCheckResourceAttr(testAccInstanceResourceName, "tags.%", "1"),
				),
			},
			{
				ResourceName:            testAccInstanceResourceName,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"auto_renew_time_length", "cds_auto_renew", "delete_cds_snapshot_flag", "related_release_flag"},
			},
			{
				Config: testAccInstanceConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccInstanceResourceName),
					resource.TestCheckResourceAttr(testAccInstanceResourceName, "name", BaiduCloudTestResourceAttrNamePrefix+"BCC-update"),
					resource.TestCheckResourceAttr(testAccInstanceResourceName, "description", "terraform test update instance"),
					resource.TestCheckResourceAttrSet(testAccInstanceResourceName, "image_id"),
					resource.TestCheckResourceAttrSet(testAccInstanceResourceName, "availability_zone"),
					resource.TestCheckResourceAttrSet(testAccInstanceResourceName, "cpu_count"),
					resource.TestCheckResourceAttrSet(testAccInstanceResourceName, "memory_capacity_in_gb"),
					resource.TestCheckResourceAttrSet(testAccInstanceResourceName, "root_disk_size_in_gb"),
					resource.TestCheckResourceAttrSet(testAccInstanceResourceName, "root_disk_storage_type"),
					resource.TestCheckResourceAttr(testAccInstanceResourceName, "billing.payment_timing", "Postpaid"),
					resource.TestCheckResourceAttrSet(testAccInstanceResourceName, "subnet_id"),
					resource.TestCheckResourceAttr(testAccInstanceResourceName, "security_groups.#", "1"),
					resource.TestCheckResourceAttrSet(testAccInstanceResourceName, "status"),
					resource.TestCheckResourceAttrSet(testAccInstanceResourceName, "create_time"),
					resource.TestCheckResourceAttrSet(testAccInstanceResourceName, "internal_ip"),
					resource.TestCheckResourceAttrSet(testAccInstanceResourceName, "placement_policy"),
					resource.TestCheckResourceAttrSet(testAccInstanceResourceName, "vpc_id"),
					resource.TestCheckResourceAttr(testAccInstanceResourceName, "cds_disks.0.cds_size_in_gb", "50"),
					resource.TestCheckResourceAttr(testAccInstanceResourceName, "cds_disks.0.storage_type", "cloud_hp1"),
					resource.TestCheckResourceAttr(testAccInstanceResourceName, "tags.%", "1"),
				),
			},
		},
	})
}

func testAccInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	bccService := &BccService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccInstanceResourceType {
			continue
		}

		_, err := bccService.GetInstanceDetail(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		return WrapError(Error("BCC instance still exist"))
	}

	return nil
}

func testAccInstanceConfig() string {
	return fmt.Sprintf(`
data "baiducloud_specs" "default" {}

data "baiducloud_zones" "default" {}

data "baiducloud_images" "default" {}

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

resource "baiducloud_instance" "default" {
  image_id              = data.baiducloud_images.default.images.0.id
  name                  = "%s"
  description           = "terraform test instance"
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  billing = {
    payment_timing = "Postpaid"
  }

  subnet_id       = baiducloud_subnet.default.id
  security_groups = [baiducloud_security_group.default.id]

  related_release_flag     = true
  delete_cds_snapshot_flag = true

  cds_disks {
    cds_size_in_gb = 50
    storage_type   = "cloud_hp1"
  }

  tags = {
    "testKey" = "testValue"
  }
}
`, BaiduCloudTestResourceAttrNamePrefix+"VPC",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceAttrNamePrefix+"SG",
		BaiduCloudTestResourceAttrNamePrefix+"BCC")
}

func testAccInstanceConfigUpdate() string {
	return fmt.Sprintf(`
data "baiducloud_specs" "default" {}

data "baiducloud_zones" "default" {}

data "baiducloud_images" "default" {}

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

resource "baiducloud_instance" "default" {
  image_id              = data.baiducloud_images.default.images.0.id
  name                  = "%s"
  description           = "terraform test update instance"
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  billing = {
    payment_timing = "Postpaid"
  }
  admin_pass = "terraform@123"

  subnet_id       = baiducloud_subnet.default02.id
  security_groups = [baiducloud_security_group.default02.id]

  related_release_flag     = true
  delete_cds_snapshot_flag = true

  cds_disks {
    cds_size_in_gb = 50
    storage_type   = "cloud_hp1"
  }

  tags = {
    "testKey" = "testValue"
  }
}

resource "baiducloud_eip_association" "default" {
  eip           = baiducloud_eip.default.id
  instance_type = "BCC"
  instance_id   = baiducloud_instance.default.id
}
`, BaiduCloudTestResourceAttrNamePrefix+"EIP",
		BaiduCloudTestResourceAttrNamePrefix+"VPC",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceAttrNamePrefix+"SG",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet02",
		BaiduCloudTestResourceAttrNamePrefix+"SG02",
		BaiduCloudTestResourceAttrNamePrefix+"BCC-update")
}
