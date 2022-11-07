package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccSMSTemplateDataSourceName          = "data.baiducloud_sms_template.default"
	testAccSMSTemplateDataSourceAttrKeyPrefix = "template_info."
)

//lintignore:AT003
func TestAccBaiduCloudSMSTemplateDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSMSTemplateDataSourceConfig(BaiduCloudTestResourceTypeNameSMSTemplate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccSMSTemplateDataSourceName),
					resource.TestCheckResourceAttrSet(testAccSMSTemplateDataSourceName, testAccSMSTemplateDataSourceAttrKeyPrefix+"name"),
					resource.TestCheckResourceAttrSet(testAccSMSTemplateDataSourceName, testAccSMSTemplateDataSourceAttrKeyPrefix+"content"),
					resource.TestCheckResourceAttrSet(testAccSMSTemplateDataSourceName, testAccSMSTemplateDataSourceAttrKeyPrefix+"status"),
					resource.TestCheckResourceAttr(testAccSMSTemplateDataSourceName, testAccSMSTemplateDataSourceAttrKeyPrefix+"country_type", "GLOBAL"),
				),
			},
		},
	})
}

func testAccSMSTemplateDataSourceConfig(name string) string {
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

data "baiducloud_sms_template" "default" {
	template_id = "${baiducloud_sms_template.default.id}"
}
`, name)
}
