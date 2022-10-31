package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccSMSSignautureDataSourceName         = "data.baiducloud_sms_signature.default"
	testAccSMSSignatureDataSourceAttrKeyPrefix = "signature_info."
)

//lintignore:AT003
func TestAccBaiduCloudSMSSignatureDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSMSSignatureDataSourceConfig(BaiduCloudTestResourceTypeNameSMSSignature),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccSMSSignautureDataSourceName),
					resource.TestCheckResourceAttrSet(testAccSMSSignautureDataSourceName, testAccSMSSignatureDataSourceAttrKeyPrefix+"content"),
					resource.TestCheckResourceAttrSet(testAccSMSSignautureDataSourceName, testAccSMSSignatureDataSourceAttrKeyPrefix+"content_type"),
					resource.TestCheckResourceAttrSet(testAccSMSSignautureDataSourceName, testAccSMSSignatureDataSourceAttrKeyPrefix+"status"),
					resource.TestCheckResourceAttr(testAccSMSSignautureDataSourceName, testAccSMSSignatureDataSourceAttrKeyPrefix+"content", "baidu"),
				),
			},
		},
	})
}

func testAccSMSSignatureDataSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "baiducloud_sms_signature" "default" {
  content      = "baidu"
  content_type = "Enterprise"
  description  = "terraform test"
  country_type = "DOMESTIC"
    
}

data "baiducloud_sms_signature" "default" {
	signature_id = "${baiducloud_sms_signature.default.id}"
}
`, name)
}
