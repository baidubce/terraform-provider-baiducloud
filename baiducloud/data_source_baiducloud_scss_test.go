package baiducloud

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

const (
	testAccScssDataSourceName          = "data.baiducloud_scss.default"
	testAccScssDataSourceAttrKeyPrefix = "scss.0."
)

func TestAccBaiduCloudScssDataSource(t *testing.T) {
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	name := BaiduCloudTestResourceAttrNamePrefix + "Scss-" + timeStamp
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScssDataSourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccScssDataSourceName),
					resource.TestCheckResourceAttrSet(testAccScssDataSourceName, testAccScssDataSourceAttrKeyPrefix+"engine_version"),
					resource.TestCheckResourceAttrSet(testAccScssDataSourceName, testAccScssDataSourceAttrKeyPrefix+"instance_id"),
					resource.TestCheckResourceAttrSet(testAccScssDataSourceName, testAccScssDataSourceAttrKeyPrefix+"port"),
					resource.TestCheckResourceAttrSet(testAccScssDataSourceName, testAccScssDataSourceAttrKeyPrefix+"capacity"),
					resource.TestCheckResourceAttrSet(testAccScssDataSourceName, testAccScssDataSourceAttrKeyPrefix+"payment_timing"),
					resource.TestCheckResourceAttrSet(testAccScssDataSourceName, testAccScssDataSourceAttrKeyPrefix+"create_time"),
				),
			},
		},
	})
}

func testAccScssDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_scs" "default" {
    instance_name           = "%s"
    billing = {
   		payment_timing 		= "Postpaid"
    }
    purchase_count 			= 1
 	port 					= 6379
	engine_version 			= "3.2"
	node_type 				= "cache.n1.micro"
	cluster_type 			= "master_slave"
	replication_num 		= 1
	shard_num 				= 1
	proxy_num 				= 0
}
data "baiducloud_scss" "default" {
    name_regex        = "test-BaiduAccScss-*"
	filter {
		name = "cluster_type"
	 	values = [baiducloud_scs.default.cluster_type]
	}
}
`, name)
}
