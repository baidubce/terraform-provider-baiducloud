package baiducloud

import (
	"fmt"
	"testing"

	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccSecurityGroupResourceType = "baiducloud_security_group"
	testAccSecurityGroupResourceName = testAccSecurityGroupResourceType + "." + BaiduCloudTestResourceName
)

func TestAccBaiduCloudSecurityGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSecurityGroupDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccSecurityGroupResourceName),
					resource.TestCheckResourceAttr(testAccSecurityGroupResourceName, "name", BaiduCloudTestResourceAttrNamePrefix+"SecurityGroup"),
					resource.TestCheckResourceAttr(testAccSecurityGroupResourceName, "tags.#", "1"),
					resource.TestCheckResourceAttrSet(testAccSecurityGroupResourceName, "vpc_id"),
				),
			},
			{
				ResourceName:      testAccSecurityGroupResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSecurityGroupDestory(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	bccService := BccService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccSnapshotResourceType {
			continue
		}

		listArgs := &api.ListSecurityGroupArgs{
			VpcId: rs.Primary.Attributes["vpc_id"],
		}
		sgList, err := bccService.ListAllSecurityGroups(listArgs)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		for _, sg := range sgList {
			if sg.Id == rs.Primary.ID {
				return WrapError(Error("SecurityGroup still exist"))
			}
		}
	}

	return nil
}

func testAccSecurityGroupConfig() string {
	return fmt.Sprintf(`
resource "%s" "%s" {
  name        = "%s"
  description = "Baidu acceptance test"
  tags {
    tag_key   = "testKey"
    tag_value = "testValue"
  }
}
`, testAccSecurityGroupResourceType, BaiduCloudTestResourceName, BaiduCloudTestResourceAttrNamePrefix+"SecurityGroup")
}
