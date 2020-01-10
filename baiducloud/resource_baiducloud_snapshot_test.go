package baiducloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccSnapshotResourceType = "baiducloud_snapshot"
	testAccSnapshotResourceName = testAccSnapshotResourceType + "." + BaiduCloudTestResourceName
)

func init() {
	resource.AddTestSweepers(testAccSnapshotResourceType, &resource.Sweeper{
		Name: testAccSnapshotResourceType,
		F:    testSweepSnapshots,
	})
}

func testSweepSnapshots(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("get BaiduCloud client error: %s", err)
	}
	client := rawClient.(*connectivity.BaiduClient)
	bccService := BccService{client}

	spList, err := bccService.ListAllSnapshots("")
	if err != nil {
		return fmt.Errorf("get Snapshots error: %s", err)
	}

	for _, sp := range spList {
		if !strings.HasPrefix(sp.Name, BaiduCloudTestResourceAttrNamePrefix) {
			log.Printf("[INFO] Skipping Snapshot: %s (%s)", sp.Name, sp.Id)
			continue
		}

		log.Printf("[INFO] Deleting Snapshot: %s (%s)", sp.Name, sp.Id)
		_, err := client.WithBccClient(func(client *bcc.Client) (i interface{}, e error) {
			return nil, client.DeleteSnapshot(sp.Id)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Snapshot %s (%s)", sp.Name, sp.Id)
		}
	}

	return nil
}

func TestAccBaiduCloudSnapshot(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSnapshotDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccSnapshotResourceName),
					resource.TestCheckResourceAttr(testAccSnapshotResourceName, "name", BaiduCloudTestResourceAttrNamePrefix+"Snapshot"),
					resource.TestCheckResourceAttr(testAccSnapshotResourceName, "size_in_gb", "5"),
					resource.TestCheckResourceAttr(testAccSnapshotResourceName, "status", "Available"),
					resource.TestCheckResourceAttrSet(testAccSnapshotResourceName, "create_method"),
					resource.TestCheckResourceAttrSet(testAccSnapshotResourceName, "create_time"),
					resource.TestCheckResourceAttrSet(testAccSnapshotResourceName, "volume_id"),
				),
			},
			{
				ResourceName:      testAccSnapshotResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSnapshotDestory(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	bccService := BccService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccSnapshotResourceType {
			continue
		}

		instanceId := rs.Primary.Attributes["volume_id"]
		spList, err := bccService.ListAllSnapshots(instanceId)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		for _, sp := range spList {
			if sp.Id == rs.Primary.ID {
				return WrapError(Error("Snapshot still exist"))
			}
		}
	}

	return nil
}

func testAccSnapshotConfig() string {
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
  description     = ""
  disk_size_in_gb = 5
  payment_timing  = "Postpaid"
}

resource "%s" "%s" {
  name        = "%s"
  description = "Baidu acceptance test"
  volume_id   = baiducloud_cds.default.id
}
`, BaiduCloudTestResourceAttrNamePrefix+"BCC",
		BaiduCloudTestResourceAttrNamePrefix+"CDS",
		testAccSnapshotResourceType,
		BaiduCloudTestResourceName,
		BaiduCloudTestResourceAttrNamePrefix+"Snapshot")
}
