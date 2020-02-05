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
	testAccCdsResourceType = "baiducloud_cds"
	testAccCdsResourceName = testAccCdsResourceType + "." + BaiduCloudTestResourceName
)

func init() {
	resource.AddTestSweepers(testAccCdsResourceType, &resource.Sweeper{
		Name:         testAccCdsResourceType,
		F:            testSweepCds,
		Dependencies: []string{testAccCdsAttachmentResourceType},
	})
}

func testSweepCds(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("get BaiduCloud client error: %s", err)
	}
	client := rawClient.(*connectivity.BaiduClient)
	bccService := BccService{client}

	listArgs := &api.ListCDSVolumeArgs{}
	cdsList, err := bccService.ListAllCDSVolumeDetail(listArgs)
	if err != nil {
		return fmt.Errorf("get CDS list error: %s", err)
	}

	for _, c := range cdsList {
		if !strings.HasPrefix(c.Name, BaiduCloudTestResourceAttrNamePrefix) {
			log.Printf("[INFO] Skipping CDS: %s (%s)", c.Name, c.Id)
			continue
		}

		if c.Status == api.VolumeStatusINUSE {
			instanceId := c.Attachments[0].InstanceId
			err := bccService.DetachCDSVolume(c.Id, instanceId)
			if err != nil {
				log.Printf("[ERROR] Failed to Detach CDS %s (%s) from instance %s", c.Name, c.Id, instanceId)
			}
		}

		log.Printf("[INFO] Deleting CDS: %s (%s)", c.Name, c.Id)
		deleteArgs := &api.DeleteCDSVolumeArgs{
			AutoSnapshot:   "on",
			ManualSnapshot: "on",
		}
		_, err := client.WithBccClient(func(client *bcc.Client) (i interface{}, e error) {
			return nil, client.DeleteCDSVolumeNew(c.Id, deleteArgs)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete CDS %s (%s)", c.Name, c.Id)
		}
	}

	return nil
}

//lintignore:AT003
func TestAccBaiduCloudCds(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCdsDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccCdsConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccCdsResourceName),
					resource.TestCheckResourceAttr(testAccCdsResourceName, "name", BaiduCloudTestResourceAttrNamePrefix+"CDS"),
					resource.TestCheckResourceAttr(testAccCdsResourceName, "disk_size_in_gb", "5"),
					resource.TestCheckResourceAttr(testAccCdsResourceName, "payment_timing", "Postpaid"),
					resource.TestCheckResourceAttr(testAccCdsResourceName, "status", string(api.VolumeStatusAVAILABLE)),
					resource.TestCheckResourceAttrSet(testAccCdsResourceName, "storage_type"),
					resource.TestCheckNoResourceAttr(testAccCdsResourceName, "description"),
				),
			},
			{
				ResourceName:      testAccCdsResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCdsConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccCdsResourceName),
					resource.TestCheckResourceAttr(testAccCdsResourceName, "name", BaiduCloudTestResourceAttrNamePrefix+"CDSUpdate"),
					resource.TestCheckResourceAttr(testAccCdsResourceName, "description", "test update"),
					resource.TestCheckResourceAttr(testAccCdsResourceName, "disk_size_in_gb", "10"),
					resource.TestCheckResourceAttr(testAccCdsResourceName, "payment_timing", "Postpaid"),
					resource.TestCheckResourceAttr(testAccCdsResourceName, "status", string(api.VolumeStatusAVAILABLE)),
					resource.TestCheckResourceAttrSet(testAccCdsResourceName, "storage_type"),
					resource.TestCheckResourceAttrSet(testAccCdsResourceName, "description"),
				),
			},
		},
	})
}

func testAccCdsDestory(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	bccService := BccService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccCdsResourceType {
			continue
		}

		_, err := bccService.GetCDSVolumeDetail(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		return WrapError(Error("CDS still exist"))
	}

	return nil
}

func testAccCdsConfig() string {
	return fmt.Sprintf(`
data "baiducloud_specs" "default" {}

data "baiducloud_zones" "default" {}

data "baiducloud_images" "default" {
  image_type = "System"
}

resource "baiducloud_instance" "default" {
  name                  = "%s"
  image_id              = data.baiducloud_images.default.images.0.id
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  billing = {
    payment_timing = "Postpaid"
  }
}

resource "%s" "%s" {
  depends_on      = [baiducloud_instance.default]
  name            = "%s"
  disk_size_in_gb = 5
  payment_timing  = "Postpaid"
}
`, BaiduCloudTestResourceAttrNamePrefix+"BCC",
		testAccCdsResourceType, BaiduCloudTestResourceName, BaiduCloudTestResourceAttrNamePrefix+"CDS")
}

func testAccCdsConfigUpdate() string {
	return fmt.Sprintf(`
data "baiducloud_specs" "default" {}

data "baiducloud_zones" "default" {}

data "baiducloud_images" "default" {
  image_type = "System"
}

resource "baiducloud_instance" "default" {
  name                  = "%s"
  image_id              = data.baiducloud_images.default.images.0.id
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  billing = {
    payment_timing = "Postpaid"
  }
}

resource "%s" "%s" {
  depends_on      = [baiducloud_instance.default]
  name            = "%s"
  description     = "test update"
  disk_size_in_gb = 10
  payment_timing  = "Postpaid"
}
`, BaiduCloudTestResourceAttrNamePrefix+"BCC",
		testAccCdsResourceType, BaiduCloudTestResourceName, BaiduCloudTestResourceAttrNamePrefix+"CDSUpdate")
}
