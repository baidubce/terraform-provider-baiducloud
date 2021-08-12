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
		if !strings.HasPrefix(asp.Name, BaiduCloudTestResourceTypeName) {
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

//lintignore:AT003
func TestAccBaiduCloudAutoSnapshotPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccAutoSnapshotPolicyDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccAutoSnapshotPolicyConfig(BaiduCloudTestResourceTypeNameAutoSnapshotPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAutoSnapshotPolicyResourceName),
					resource.TestCheckResourceAttr(testAccAutoSnapshotPolicyResourceName, "name", BaiduCloudTestResourceTypeNameAutoSnapshotPolicy),
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
				Config: testAccAutoSnapshotPolicyConfigUpdate(BaiduCloudTestResourceTypeNameAutoSnapshotPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAutoSnapshotPolicyResourceName),
					resource.TestCheckResourceAttr(testAccAutoSnapshotPolicyResourceName, "name", BaiduCloudTestResourceTypeNameAutoSnapshotPolicy),
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
				Config: testAccAutoSnapshotPolicyConfig(BaiduCloudTestResourceTypeNameAutoSnapshotPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAutoSnapshotPolicyResourceName),
					resource.TestCheckResourceAttr(testAccAutoSnapshotPolicyResourceName, "name", BaiduCloudTestResourceTypeNameAutoSnapshotPolicy),
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

func testAccAutoSnapshotPolicyConfig(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_auto_snapshot_policy" "default" {
  name            = "%s"
  time_points     = [0, 22]
  repeat_weekdays = [0, 3]
  retention_days  = -1
}
`, name)
}

func testAccAutoSnapshotPolicyConfigUpdate(name string) string {
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

resource "baiducloud_cds" "default" {
  depends_on      = [baiducloud_instance.default]
  name            = "${var.name}"
  description     = "created by terraform"
  disk_size_in_gb = 5
  payment_timing  = "Postpaid"
  zone_name       = data.baiducloud_zones.default.zones.0.zone_name
}

resource "baiducloud_cds_attachment" "default" {
  cds_id      = baiducloud_cds.default.id
  instance_id = baiducloud_instance.default.id
}

resource "baiducloud_auto_snapshot_policy" "default" {
  name            = "${var.name}"
  time_points     = [0, 20, 22]
  repeat_weekdays = [0]
  retention_days  = 2
  volume_ids      = [baiducloud_cds_attachment.default.id]
}
`, name)
}
