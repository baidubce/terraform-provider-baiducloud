package baiducloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const (
	testAccRdsSecurityIpResourceType = "baiducloud_rds_security_ip"
	testAccRdsSecurityIpResourceName = testAccRdsSecurityIpResourceType + "." + BaiduCloudTestResourceName
)

func TestAccBaiduCloudRdsSecurityIp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccRdsSecurityIpConfig(BaiduCloudTestResourceTypeNameRdsSecurityIp),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccRdsSecurityIpResourceName),
					resource.TestCheckResourceAttr(testAccRdsSecurityIpResourceName, "instance_id", "rds-BIFDrIl9"),
				),
			},
		},
	})
}

func testAccRdsSecurityIpConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "baiducloud_rds_security_ip" "default" {
    instance_id = "rds-BIFDrIl9"
    security_ips = ["192.168.3.5"]
}
`, name)
}
