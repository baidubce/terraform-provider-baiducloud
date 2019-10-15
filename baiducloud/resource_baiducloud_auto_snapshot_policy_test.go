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
	testAccAutoSnapshotPolicyResourceType = "baiducloud_auto_snapshot_policy"
	testAccAutoSnapshotPolicyResourceName = testAccAutoSnapshotPolicyResourceType + "." + BaiduCloudTestResourceName
)

func init() {
	resource.AddTestSweepers(testAccAutoSnapshotPolicyResourceType, &resource.Sweeper{
		Name: testAccAutoSnapshotPolicyResourceType,
		F:    testSweepAutoSnapshotPolicys,
	})
}

func testSweepAutoSnapshotPolicys(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("get BaiduCloud client error: %s", err)
	}
	client := rawClient.(*connectivity.BaiduClient)
	bccService := BccService{client}

	listArgs := &api.ListASPArgs{}
	aspList, err := bccService.ListAllAutoSnapshotPolicies(listArgs)
	if err != nil {
		return fmt.Errorf("get AutoSnapshotPolicies error: %s", err)
	}

	for _, asp := range aspList {
		if !strings.HasPrefix(asp.Name, BaiduCloudTestResourceAttrNamePrefix) {
			log.Printf("[INFO] Skipping AutoSnapshotPolicy: %s (%s)", asp.Name, asp.Id)
			continue
		}

		log.Printf("[INFO] Deleting AutoSnapshotPolicy: %s (%s)", asp.Name, asp.Id)
		_, err := client.WithBccClient(func(client *bcc.Client) (i interface{}, e error) {
			return nil, client.DeleteAutoSnapshotPolicy(asp.Id)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete AutoSnapshotPolicy %s (%s)", asp.Name, asp.Id)
		}
	}

	return nil
}

func TestAccBaiduCloudAutoSnapshotPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccAutoSnapshotPolicyDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccAutoSnapshotPolicyConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAutoSnapshotPolicyResourceName),
					resource.TestCheckResourceAttr(testAccAutoSnapshotPolicyResourceName, "name", BaiduCloudTestResourceAttrNamePrefix+"ASP"),
					resource.TestCheckResourceAttr(testAccAutoSnapshotPolicyResourceName, "time_points.#", "2"),
					resource.TestCheckResourceAttr(testAccAutoSnapshotPolicyResourceName, "repeat_weekdays.#", "2"),
					resource.TestCheckResourceAttr(testAccAutoSnapshotPolicyResourceName, "retention_days", "-1"),
					resource.TestCheckResourceAttrSet(testAccAutoSnapshotPolicyResourceName, "status"),
					resource.TestCheckResourceAttrSet(testAccAutoSnapshotPolicyResourceName, "created_time"),
				),
			},
			{
				ResourceName:      testAccAutoSnapshotPolicyResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAutoSnapshotPolicyConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAutoSnapshotPolicyResourceName),
					resource.TestCheckResourceAttr(testAccAutoSnapshotPolicyResourceName, "name", BaiduCloudTestResourceAttrNamePrefix+"ASP"),
					resource.TestCheckResourceAttr(testAccAutoSnapshotPolicyResourceName, "time_points.#", "3"),
					resource.TestCheckResourceAttr(testAccAutoSnapshotPolicyResourceName, "repeat_weekdays.#", "1"),
					resource.TestCheckResourceAttr(testAccAutoSnapshotPolicyResourceName, "retention_days", "2"),
					resource.TestCheckResourceAttr(testAccAutoSnapshotPolicyResourceName, "volume_ids.#", "1"),
					resource.TestCheckResourceAttr(testAccAutoSnapshotPolicyResourceName, "volume_count", "1"),
					resource.TestCheckResourceAttrSet(testAccAutoSnapshotPolicyResourceName, "status"),
					resource.TestCheckResourceAttrSet(testAccAutoSnapshotPolicyResourceName, "created_time"),
				),
			},
			{
				Config: testAccAutoSnapshotPolicyConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAutoSnapshotPolicyResourceName),
					resource.TestCheckResourceAttr(testAccAutoSnapshotPolicyResourceName, "name", BaiduCloudTestResourceAttrNamePrefix+"ASP"),
					resource.TestCheckResourceAttr(testAccAutoSnapshotPolicyResourceName, "time_points.#", "2"),
					resource.TestCheckResourceAttr(testAccAutoSnapshotPolicyResourceName, "repeat_weekdays.#", "2"),
					resource.TestCheckResourceAttr(testAccAutoSnapshotPolicyResourceName, "retention_days", "-1"),
					resource.TestCheckResourceAttr(testAccAutoSnapshotPolicyResourceName, "volume_ids.#", "0"),
					resource.TestCheckResourceAttrSet(testAccAutoSnapshotPolicyResourceName, "status"),
					resource.TestCheckResourceAttrSet(testAccAutoSnapshotPolicyResourceName, "created_time"),
				),
			},
		},
	})
}

func testAccAutoSnapshotPolicyDestory(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccAutoSnapshotPolicyResourceType {
			continue
		}

		_, err := client.WithBccClient(func(client *bcc.Client) (i interface{}, e error) {
			return client.GetAutoSnapshotPolicy(rs.Primary.ID)
		})
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		return WrapError(Error("AutoSnapshotPolicy still exist"))
	}

	return nil
}

func testAccAutoSnapshotPolicyConfig() string {
	return fmt.Sprintf(`
resource "%s" "%s" {
  name            = "%s"
  time_points     = [0, 22]
  repeat_weekdays = [0, 3]
  retention_days  = -1
}
`, testAccAutoSnapshotPolicyResourceType, BaiduCloudTestResourceName, BaiduCloudTestResourceAttrNamePrefix+"ASP")
}

func testAccAutoSnapshotPolicyConfigUpdate() string {
	return fmt.Sprintf(`
data "baiducloud_specs" "default" {}

data "baiducloud_zones" "default" {}

data "baiducloud_images" "default" {
  image_type = "System"
}

resource "baiducloud_instance" "default" {
  name                  = "%s"
  image_id              = "${data.baiducloud_images.default.images.0.id}"
  availability_zone     = "${data.baiducloud_zones.default.zones.0.zone_name}"
  cpu_count             = "${data.baiducloud_specs.default.specs.0.cpu_count}"
  memory_capacity_in_gb = "${data.baiducloud_specs.default.specs.0.memory_size_in_gb}"
  billing = {
    payment_timing = "Postpaid"
  }
}

resource "baiducloud_cds" "default" {
  name                   = "%s"
  description            = ""
  disk_size_in_gb        = 5
  payment_timing         = "Postpaid"
}

resource "baiducloud_cds_attachment" "default" {
  cds_id      = "${baiducloud_cds.default.id}"
  instance_id = "${baiducloud_instance.default.id}"
}

resource "%s" "%s" {
  name            = "%s"
  time_points     = [0, 20, 22]
  repeat_weekdays = [0]
  retention_days  = 2
  volume_ids      = ["${baiducloud_cds_attachment.default.id}"]
}
`, BaiduCloudTestResourceAttrNamePrefix+"BCC", BaiduCloudTestResourceAttrNamePrefix+"CDS",
		testAccAutoSnapshotPolicyResourceType, BaiduCloudTestResourceName, BaiduCloudTestResourceAttrNamePrefix+"ASP")
}
