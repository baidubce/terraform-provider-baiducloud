package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRdsSecurityIpDataSourceName          = "data.baiducloud_rds_security_ips.default"
	testAccRdsSecurityIpDataSourceAttrKeyPrefix = "security_ips.0."
)

//lintignore:AT003
func TestAccBaiduCloudRdsSecurityIpDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsSecurityIpDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccRdsSecurityIpDataSourceName),
					resource.TestCheckResourceAttrSet(testAccRdsSecurityIpDataSourceName, testAccRdsSecurityIpDataSourceAttrKeyPrefix+"ip"),
				),
			},
		},
	})
}

func testAccRdsSecurityIpDataSourceConfig() string {
	return fmt.Sprintf(`
data "baiducloud_rds_security_ips" "default" {
    instance_id = "rds-BIFDrIl9"    
}
`)
}
