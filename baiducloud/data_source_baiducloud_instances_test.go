package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const (
	testAccInstancesDataSourceName          = "data.baiducloud_instances.default"
	testAccInstancesDataSourceAttrKeyPrefix = "instances.0."
)

func TestAccBaiduCloudInstancesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstancesDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccInstancesDataSourceName),
					resource.TestCheckResourceAttr(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"name", BaiduCloudTestResourceAttrNamePrefix+"BCC"),
					resource.TestCheckResourceAttr(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"tags.%", "1"),
					resource.TestCheckResourceAttr(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"tags.testKey", "testValue"),
					resource.TestCheckResourceAttrSet(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"image_id"),
					resource.TestCheckResourceAttrSet(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"zone_name"),
					resource.TestCheckResourceAttrSet(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"cpu_count"),
					resource.TestCheckResourceAttrSet(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"memory_capacity_in_gb"),
					resource.TestCheckResourceAttrSet(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"payment_timing"),
					resource.TestCheckResourceAttrSet(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"subnet_id"),
					resource.TestCheckResourceAttrSet(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"status"),
					resource.TestCheckResourceAttrSet(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"create_time"),
					resource.TestCheckResourceAttrSet(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"internal_ip"),
					resource.TestCheckResourceAttrSet(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"placement_policy"),
					resource.TestCheckResourceAttrSet(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"vpc_id"),
					resource.TestCheckResourceAttrSet(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"root_disk_size_in_gb"),
					resource.TestCheckResourceAttrSet(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"root_disk_storage_type"),
					resource.TestCheckResourceAttrSet(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"auto_renew"),
				),
			},
		},
	})
}

func testAccInstancesDataSourceConfig() string {
	return fmt.Sprintf(`
data "baiducloud_specs" "default" {}

data "baiducloud_zones" "default" {}

data "baiducloud_images" "default" {}

resource "baiducloud_instance" "default" {
  image_id              = data.baiducloud_images.default.images.0.id
  name                  = "%s"
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  billing = {
    payment_timing = "Postpaid"
  }

  tags = {
    "testKey" = "testValue"
  }
}

data "baiducloud_instances" "default" {
  internal_ip = baiducloud_instance.default.internal_ip
  zone_name   = baiducloud_instance.default.availability_zone
}
`, BaiduCloudTestResourceAttrNamePrefix+"BCC")
}
