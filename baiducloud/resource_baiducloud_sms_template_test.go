package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccSMSTemplateResourceType = "baiducloud_sms_template"
	testAccSMSTemplateResourceName = testAccSMSTemplateResourceType + "." + BaiduCloudTestResourceName
)

func TestAccBaiduCloudSMSTemplate_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSMSTemplateDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccSMSTemplateConfig(BaiduCloudTestResourceTypeNameSMSTemplate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccSMSTemplateResourceName),
					resource.TestCheckResourceAttrSet(testAccSMSTemplateResourceName, "user_id"),
					resource.TestCheckResourceAttr(testAccSMSTemplateResourceName, "status", "SUBMITTED"),
				),
			},
		},
	})
}

func testAccSMSTemplateDestory(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	smsService := SMSService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccSMSTemplateResourceType {
			continue
		}

		_, err := smsService.GetSMSTemplateDetail(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}

	return nil
}

func testAccSMSTemplateConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_sms_template" "default" {
  name	         = "My test template"
  content        = "Test content"
  sms_type       = "CommonNotice"
  country_type   = "GLOBAL"
  description    = "this is a test sms template"
}
`, name)
}
