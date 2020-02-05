package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const (
	testAccAutoSnapshotPoliciesDataSourceName          = "data.baiducloud_auto_snapshot_policies.default"
	testAccAutoSnapshotPoliciesDataSourceAttrKeyPrefix = "auto_snapshot_policies.0."
)

//lintignore:AT003
func TestAccBaiduCloudAutoSnapshotPoliciesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAutoSnapshotPoliciesDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccAutoSnapshotPoliciesDataSourceName),
					resource.TestCheckResourceAttr(testAccAutoSnapshotPoliciesDataSourceName, testAccAutoSnapshotPoliciesDataSourceAttrKeyPrefix+"name", BaiduCloudTestResourceAttrNamePrefix+"ASP"),
					resource.TestCheckResourceAttr(testAccAutoSnapshotPoliciesDataSourceName, testAccAutoSnapshotPoliciesDataSourceAttrKeyPrefix+"time_points.#", "2"),
					resource.TestCheckResourceAttr(testAccAutoSnapshotPoliciesDataSourceName, testAccAutoSnapshotPoliciesDataSourceAttrKeyPrefix+"time_points.0", "0"),
					resource.TestCheckResourceAttr(testAccAutoSnapshotPoliciesDataSourceName, testAccAutoSnapshotPoliciesDataSourceAttrKeyPrefix+"time_points.1", "22"),
					resource.TestCheckResourceAttr(testAccAutoSnapshotPoliciesDataSourceName, testAccAutoSnapshotPoliciesDataSourceAttrKeyPrefix+"repeat_weekdays.#", "2"),
					resource.TestCheckResourceAttr(testAccAutoSnapshotPoliciesDataSourceName, testAccAutoSnapshotPoliciesDataSourceAttrKeyPrefix+"repeat_weekdays.0", "0"),
					resource.TestCheckResourceAttr(testAccAutoSnapshotPoliciesDataSourceName, testAccAutoSnapshotPoliciesDataSourceAttrKeyPrefix+"repeat_weekdays.1", "3"),
					resource.TestCheckResourceAttr(testAccAutoSnapshotPoliciesDataSourceName, testAccAutoSnapshotPoliciesDataSourceAttrKeyPrefix+"retention_days", "-1"),
					resource.TestCheckResourceAttr(testAccAutoSnapshotPoliciesDataSourceName, testAccAutoSnapshotPoliciesDataSourceAttrKeyPrefix+"volume_count", "1"),
					resource.TestCheckResourceAttrSet(testAccAutoSnapshotPoliciesDataSourceName, testAccAutoSnapshotPoliciesDataSourceAttrKeyPrefix+"status"),
					resource.TestCheckResourceAttrSet(testAccAutoSnapshotPoliciesDataSourceName, testAccAutoSnapshotPoliciesDataSourceAttrKeyPrefix+"created_time"),
				),
			},
		},
	})
}

func testAccAutoSnapshotPoliciesDataSourceConfig() string {
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

resource "baiducloud_cds_attachment" "default" {
  cds_id      = baiducloud_cds.default.id
  instance_id = baiducloud_instance.default.id
}

resource "baiducloud_auto_snapshot_policy" "default" {
  name            = "%s"
  time_points     = [0, 22]
  repeat_weekdays = [0, 3]
  retention_days  = -1
  volume_ids      = [baiducloud_cds_attachment.default.cds_id]
}

data "baiducloud_auto_snapshot_policies" "default" {
   asp_name    = baiducloud_auto_snapshot_policy.default.name
   volume_name = baiducloud_cds.default.name
}
`, BaiduCloudTestResourceAttrNamePrefix+"BCC",
		BaiduCloudTestResourceAttrNamePrefix+"CDS",
		BaiduCloudTestResourceAttrNamePrefix+"ASP")
}
