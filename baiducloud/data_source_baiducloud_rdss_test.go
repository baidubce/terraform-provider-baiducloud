package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const (
	testAccRdssDataSourceName          = "data.baiducloud_rdss.default"
	testAccRdssDataSourceAttrKeyPrefix = "rdss.0."
)

func TestAccBaiduCloudRdssDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRdssDataSourceConfig(BaiduCloudTestResourceTypeNameRdsInstance),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccRdssDataSourceName),
					resource.TestCheckResourceAttr(testAccRdssDataSourceName, testAccRdssDataSourceAttrKeyPrefix+"memory_capacity", "1"),
					resource.TestCheckResourceAttrSet(testAccRdssDataSourceName, testAccRdssDataSourceAttrKeyPrefix+"engine_version"),
					resource.TestCheckResourceAttrSet(testAccRdssDataSourceName, testAccRdssDataSourceAttrKeyPrefix+"engine"),
					resource.TestCheckResourceAttrSet(testAccRdssDataSourceName, testAccRdssDataSourceAttrKeyPrefix+"cpu_count"),
					resource.TestCheckResourceAttrSet(testAccRdssDataSourceName, testAccRdssDataSourceAttrKeyPrefix+"volume_capacity"),
					resource.TestCheckResourceAttrSet(testAccRdssDataSourceName, testAccRdssDataSourceAttrKeyPrefix+"payment_timing"),
				),
			},
		},
	})
}

func testAccRdssDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_rds_instance" "default" {
    instance_name             = "%s"
    billing = {
        payment_timing        = "Postpaid"
    }
    engine_version            = "5.6"
    engine                    = "MySQL"
    cpu_count                 = 1
    memory_capacity           = 1
    volume_capacity           = 5
}

data "baiducloud_rdss" "default" {
    name_regex            = "tf-test-acc*"
    filter {
        name              = "memory_capacity"
        values            = [baiducloud_rds_instance.default.memory_capacity]
    }
}
`, name)
}
