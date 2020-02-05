package baiducloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccCdsAttachmentResourceType = "baiducloud_cds_attachment"
	testAccCdsAttachmentResourceName = testAccCdsAttachmentResourceType + "." + BaiduCloudTestResourceName
)

func init() {
	resource.AddTestSweepers(testAccCdsAttachmentResourceType, &resource.Sweeper{
		Name: testAccCdsAttachmentResourceType,
		F:    testSweepCdsAttachment,
	})
}

func testSweepCdsAttachment(region string) error {
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

		if c.Status == api.VolumeStatusINUSE || c.Attachments[0].InstanceId != "" {
			log.Printf("[INFO] Detach CDS: %s (%s)", c.Name, c.Id)
			instanceId := c.Attachments[0].InstanceId
			err := bccService.DetachCDSVolume(c.Id, instanceId)
			if err != nil {
				log.Printf("[ERROR] Failed to Detach CDS %s (%s) from instance %s", c.Name, c.Id, instanceId)
			}
		}
	}

	return nil
}

//lintignore:AT003
func TestAccBaiduCloudCdsAttachment(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCdsAttachmentDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccCdsAttachmentConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccCdsAttachmentResourceName),
					resource.TestCheckResourceAttrSet(testAccCdsAttachmentResourceName, "cds_id"),
					resource.TestCheckResourceAttrSet(testAccCdsAttachmentResourceName, "instance_id"),
					resource.TestCheckResourceAttrSet(testAccCdsAttachmentResourceName, "attachment_device"),
					resource.TestCheckResourceAttrSet(testAccCdsAttachmentResourceName, "attachment_serial"),
				),
			},
			{
				ResourceName:      testAccCdsAttachmentResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCdsAttachmentDestory(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	bccService := BccService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccCdsAttachmentResourceType {
			continue
		}

		volume, err := bccService.GetCDSVolumeDetail(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		if volume.Status == api.VolumeStatusINUSE {
			return WrapError(Error("CDS attachment still exist"))
		}
	}

	return nil
}

func testAccCdsAttachmentConfig() string {
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

resource "baiducloud_cds" "default" {
  depends_on      = [baiducloud_instance.default]
  name            = "%s"
  disk_size_in_gb = 5
  payment_timing  = "Postpaid"
}

resource "%s" "%s" {
  cds_id      = baiducloud_cds.default.id
  instance_id = baiducloud_instance.default.id
}
`, BaiduCloudTestResourceAttrNamePrefix+"BCC",
		BaiduCloudTestResourceAttrNamePrefix+"CDS",
		testAccCdsAttachmentResourceType, BaiduCloudTestResourceName)
}
