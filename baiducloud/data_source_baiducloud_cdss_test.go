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

//lintignore:AT003
func TestAccBaiduCloudCdsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCdsDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccCdsDataSourceConfig(BaiduCloudTestResourceTypeNameCds),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccCDSsDataSourceName),
					resource.TestCheckResourceAttr(testAccCDSsDataSourceName, testAccCDSsDataSourceAttrKeyPrefix+"name", BaiduCloudTestResourceTypeNameCds),
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

func testAccCdsDataSourceConfig(name string) string {
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
  disk_size_in_gb = 5
  payment_timing  = "Postpaid"
  zone_name       = data.baiducloud_zones.default.zones.0.zone_name
}

resource "baiducloud_cds_attachment" "default" {
  cds_id      = baiducloud_cds.default.id
  instance_id = baiducloud_instance.default.id
}

data "baiducloud_cdss" "default" {
  instance_id = baiducloud_cds_attachment.default.instance_id
  zone_name   = baiducloud_cds.default.zone_name

  filter {
    name = "name"
    values = ["tf-test-acc*"]
  }
}
`, name)
}
