package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccSMSSignatureGroupResourceType = "baiducloud_sms_signature"
	testAccSMSSignatureGroupResourceName = testAccSMSSignatureGroupResourceType + "." + BaiduCloudTestResourceName
)

func TestAccBaiduCloudSMSSignature_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccSMSSignatureDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccSMSSignatureConfig(BaiduCloudTestResourceTypeNameSMSSignature),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccSMSSignatureGroupResourceName),
					resource.TestCheckResourceAttrSet(testAccSMSSignatureGroupResourceName, "user_id"),
					resource.TestCheckResourceAttr(testAccSMSSignatureGroupResourceName, "status", "SUBMITTED"),
				),
			},
		},
	})
}

func testAccSMSSignatureDestory(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	smsService := SMSService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccSMSSignatureGroupResourceType {
			continue
		}

		_, err := smsService.GetSMSSignatureDetail(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}

	return nil
}

func testAccSMSSignatureConfig(name string) string {
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
`, name)
}
