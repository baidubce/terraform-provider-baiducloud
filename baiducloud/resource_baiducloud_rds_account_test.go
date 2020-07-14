package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const (
	testAccRdsAccountResourceType = "baiducloud_rds_account"
	testAccRdsAccountResourceName = testAccRdsAccountResourceType + "." + BaiduCloudTestResourceName
)

func TestAccBaiduCloudRdsAccount(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccRdsAccountConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccRdsAccountResourceName),
					resource.TestCheckResourceAttr(testAccRdsAccountResourceName, "account_name", "mysqlaccount"),
				),
			},
		},
	})
}

func testAccRdsAccountConfig() string {
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

resource "%s" "%s" {
    instance_id         = baiducloud_rds_instance.default.instance_id
    account_name        = "mysqlaccount"
    password            = "password12"
    account_type        = "Super"
    desc                = "test"
}
`, BaiduCloudTestResourceAttrNamePrefix+"Rds_Account", testAccRdsAccountResourceType, BaiduCloudTestResourceName)
}
