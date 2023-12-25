package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccScsSecurityIpDataSourceName          = "data.baiducloud_scs_security_ips.default"
	testAccScsSecurityIpDataSourceAttrKeyPrefix = "security_ips.0."
)

//lintignore:AT003
func TestAccBaiduCloudScsSecurityIpDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScsSecurityIpDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccScsSecurityIpDataSourceName),
					resource.TestCheckResourceAttrSet(testAccScsSecurityIpDataSourceName, testAccScsSecurityIpDataSourceAttrKeyPrefix+"ip"),
				),
			},
		},
	})
}

func testAccScsSecurityIpDataSourceConfig() string {
	return fmt.Sprintf(`
data "baiducloud_scs_security_ips" "default" {
    instance_id = "scs-bj-hzsywuljybfy"    
}
`)
}

func testAccScsSecurityIpFullConfig() string {
	return fmt.Sprintf(`

resource "baiducloud_scs" "default" {
    instance_name           = "scs-test"
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

data "baiducloud_scs_security_ips" "default" {
  instance_id = baiducloud_scs.default.id

}

`)
}
