package baiducloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"testing"
)

const (
	testAccBbcImageResourceType = "baiducloud_bbc_image"
	testAccBbcImageResourceName = testAccBbcImageResourceType + "." + BaiduCloudTestResourceName
)

func TestAccBaiduCloudBbcImage(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccBbcImageDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccBbcImageConfig(BaiduCloudTestResourceTypeNameBbcImage),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBbcImageResourceName),
				),
			},
		},
	})

}
func testAccBbcImageDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	bbcService := BbcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccBbcImageResourceType {
			continue
		}

		bbcImage, _ := bbcService.GetBbcImageDetails(rs.Primary.ID)
		if bbcImage != nil {
			return WrapError(Error("bbc image still exist"))
		}
	}

	return nil
}

func testAccBbcImageConfig(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_bbc_image" "test-image" {
  image_name = "%s"
  instance_id = "i-qwIq4vKi"
}
`, name)
}
