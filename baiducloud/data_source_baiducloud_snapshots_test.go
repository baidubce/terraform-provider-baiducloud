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

//lintignore:AT003
func TestAccBaiduCloudSnapshotsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotsDataSourceConfig(BaiduCloudTestResourceTypeNameSnapshot),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccSnapshotsDataSourceName),
					resource.TestCheckResourceAttr(testAccSnapshotsDataSourceName, testAccSnapshotsDataSourceAttrKeyPrefix+"name", BaiduCloudTestResourceTypeNameSnapshot),
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

func testAccSnapshotsDataSourceConfig(name string) string {
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

resource "baiducloud_cds" "default" {
  depends_on      = [baiducloud_instance.default]
  name            = var.name
  description     = "created by terraform"
  disk_size_in_gb = 5
  payment_timing  = "Postpaid"
  zone_name     = data.baiducloud_zones.default.zones.0.zone_name
}

resource "baiducloud_snapshot" "default" {
  name        = var.name
  description = "created by terraform"
  volume_id   = baiducloud_cds.default.id
}

data "baiducloud_snapshots" "default" {
  volume_id = baiducloud_snapshot.default.volume_id

  filter {
    name = "name"
    values = ["tf-test-acc*"]
  }
}
`, name)
}
