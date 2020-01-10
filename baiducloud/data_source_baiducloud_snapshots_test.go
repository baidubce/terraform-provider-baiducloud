package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const (
	testAccSnapshotsDataSourceName          = "data.baiducloud_snapshots.default"
	testAccSnapshotsDataSourceAttrKeyPrefix = "snapshots.0."
)

func TestAccBaiduCloudSnapshotsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccSnapshotsDataSourceName),
					resource.TestCheckResourceAttr(testAccSnapshotsDataSourceName, testAccSnapshotsDataSourceAttrKeyPrefix+"name", BaiduCloudTestResourceAttrNamePrefix+"Snapshot"),
					resource.TestCheckResourceAttr(testAccSnapshotsDataSourceName, testAccSnapshotsDataSourceAttrKeyPrefix+"size_in_gb", "5"),
					resource.TestCheckResourceAttr(testAccSnapshotsDataSourceName, testAccSnapshotsDataSourceAttrKeyPrefix+"status", "Available"),
					resource.TestCheckResourceAttrSet(testAccSnapshotsDataSourceName, testAccSnapshotsDataSourceAttrKeyPrefix+"create_method"),
					resource.TestCheckResourceAttrSet(testAccSnapshotsDataSourceName, testAccSnapshotsDataSourceAttrKeyPrefix+"create_time"),
					resource.TestCheckResourceAttrSet(testAccSnapshotsDataSourceName, testAccSnapshotsDataSourceAttrKeyPrefix+"volume_id"),
				),
			},
		},
	})
}

func testAccSnapshotsDataSourceConfig() string {
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

resource "baiducloud_snapshot" "default" {
  name        = "%s"
  description = "Baidu acceptance test"
  volume_id   = baiducloud_cds.default.id
}

data "baiducloud_snapshots" "default" {
  volume_id = baiducloud_snapshot.default.volume_id
}
`, BaiduCloudTestResourceAttrNamePrefix+"BCC",
		BaiduCloudTestResourceAttrNamePrefix+"CDS",
		BaiduCloudTestResourceAttrNamePrefix+"Snapshot")
}
