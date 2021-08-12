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
				Config: testAccRdsAccountConfig(BaiduCloudTestResourceTypeNameRdsAccount),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccRdsAccountResourceName),
					resource.TestCheckResourceAttr(testAccRdsAccountResourceName, "account_name", "mysqlaccount"),
				),
			},
		},
	})
}

func testAccRdsAccountConfig(name string) string {
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

resource "baiducloud_rds_account" "default" {
    instance_id         = baiducloud_rds_instance.default.instance_id
    account_name        = "mysqlaccount"
    password            = "password12"
    account_type        = "Super"
    desc                = "test"
}
`, name+"-rds-account")
}
