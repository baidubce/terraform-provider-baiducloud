package baiducloud

import (
	"fmt"
	"testing"

	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/hashicorp/terraform/helper/resource"
)

const (
	testAccCDSsDataSourceName          = "data.baiducloud_cdss.default"
	testAccCDSsDataSourceAttrKeyPrefix = "cdss.0."
)

func TestAccBaiduCloudCdsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCdsDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccCdsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccCDSsDataSourceName),
					resource.TestCheckResourceAttr(testAccCDSsDataSourceName, testAccCDSsDataSourceAttrKeyPrefix+"name", BaiduCloudTestResourceAttrNamePrefix+"CDS"),
					resource.TestCheckResourceAttr(testAccCDSsDataSourceName, testAccCDSsDataSourceAttrKeyPrefix+"disk_size_in_gb", "5"),
					resource.TestCheckResourceAttr(testAccCDSsDataSourceName, testAccCDSsDataSourceAttrKeyPrefix+"payment_timing", "Postpaid"),
					resource.TestCheckResourceAttr(testAccCDSsDataSourceName, testAccCDSsDataSourceAttrKeyPrefix+"is_system_volume", "false"),
					resource.TestCheckResourceAttr(testAccCDSsDataSourceName, testAccCDSsDataSourceAttrKeyPrefix+"status", string(api.VolumeStatusINUSE)),
					resource.TestCheckResourceAttr(testAccCDSsDataSourceName, testAccCDSsDataSourceAttrKeyPrefix+"attachments.#", "1"),
					resource.TestCheckResourceAttrSet(testAccCDSsDataSourceName, testAccCDSsDataSourceAttrKeyPrefix+"attachments.0.instance_id"),
					resource.TestCheckResourceAttrSet(testAccCDSsDataSourceName, testAccCDSsDataSourceAttrKeyPrefix+"cds_id"),
					resource.TestCheckResourceAttrSet(testAccCDSsDataSourceName, testAccCDSsDataSourceAttrKeyPrefix+"storage_type"),
				),
			},
		},
	})
}

func testAccCdsDataSourceConfig() string {
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
  disk_size_in_gb        = 5
  payment_timing         = "Postpaid"
}

resource "baiducloud_cds_attachment" "default" {
  cds_id      = "${baiducloud_cds.default.id}"
  instance_id = "${baiducloud_instance.default.id}"
}

data "baiducloud_cdss" "default" {
  instance_id = "${baiducloud_cds_attachment.default.instance_id}"
  zone_name   = "${baiducloud_cds.default.zone_name}"
}
`, BaiduCloudTestResourceAttrNamePrefix+"BCC",
		BaiduCloudTestResourceAttrNamePrefix+"CDS")
}
