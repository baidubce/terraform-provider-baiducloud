package baiducloud

import (
	"fmt"
	"testing"

	"github.com/baidubce/bce-sdk-go/services/cfc"
	"github.com/baidubce/bce-sdk-go/services/cfc/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccCFCAliasResourceType = "baiducloud_cfc_alias"
	testAccCFCAliasResourceName = testAccCFCAliasResourceType + "." + BaiduCloudTestResourceName
)

func TestAccBaiduCloudCFCAlias(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCFCAliasDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccCfcAliasConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCFCAliasResourceName, "alias_name", "test-BaiduAccAlias"),
					resource.TestCheckResourceAttr(testAccCFCAliasResourceName, "function_name", "test-BaiduAccCFC"),
					resource.TestCheckResourceAttr(testAccCFCAliasResourceName, "function_version", "$LATEST"),
					resource.TestCheckResourceAttr(testAccCFCAliasResourceName, "description", "terraform create alias"),
					resource.TestCheckResourceAttrSet(testAccCFCAliasResourceName, "update_time"),
					resource.TestCheckResourceAttrSet(testAccCFCAliasResourceName, "create_time"),
					resource.TestCheckResourceAttrSet(testAccCFCAliasResourceName, "uid"),
					resource.TestCheckResourceAttrSet(testAccCFCAliasResourceName, "alias_brn"),
					resource.TestCheckResourceAttrSet(testAccCFCAliasResourceName, "alias_arn"),
				),
			},
			{
				Config: testAccCfcAliasConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCFCAliasResourceName, "alias_name", "test-BaiduAccAlias"),
					resource.TestCheckResourceAttr(testAccCFCAliasResourceName, "function_name", "test-BaiduAccCFC"),
					resource.TestCheckResourceAttr(testAccCFCAliasResourceName, "function_version", "$LATEST"),
					resource.TestCheckResourceAttr(testAccCFCAliasResourceName, "description", "terraform update alias"),
					resource.TestCheckResourceAttrSet(testAccCFCAliasResourceName, "update_time"),
					resource.TestCheckResourceAttrSet(testAccCFCAliasResourceName, "create_time"),
					resource.TestCheckResourceAttrSet(testAccCFCAliasResourceName, "uid"),
					resource.TestCheckResourceAttrSet(testAccCFCAliasResourceName, "alias_brn"),
					resource.TestCheckResourceAttrSet(testAccCFCAliasResourceName, "alias_arn"),
				),
			},
		},
	})
}

func testAccCFCAliasDestory(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccCFCAliasResourceType {
			continue
		}

		args := &api.GetAliasArgs{
			FunctionName: rs.Primary.Attributes["function_name"],
			AliasName:    rs.Primary.Attributes["alias_name"],
		}
		_, err := client.WithCFCClient(func(client *cfc.Client) (i interface{}, e error) {
			return client.GetAlias(args)
		})
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		return WrapError(Error("CFC Function Alias still exist"))
	}

	return nil
}

func testAccCfcAliasConfig() string {
	return fmt.Sprintf(`
resource "baiducloud_cfc_function" "default" {
  function_name     = "%s"
  description       = "terraform create"
  environment = {
    "aaa": "bbb"
    "ccc": "ddd"
  }
  handler        = "index.handler"
  memory_size    = 128
  runtime        = "nodejs8.5"
  time_out       = 3
  code_file_name = "testFiles/cfcTestCode.zip"
}

resource "%s" "%s" {
  function_name    = baiducloud_cfc_function.default.function_name
  function_version = baiducloud_cfc_function.default.version
  alias_name       = "%s"
  description      = "terraform create alias"
}
`, BaiduCloudTestResourceAttrNamePrefix+"CFC",
		testAccCFCAliasResourceType, BaiduCloudTestResourceName, BaiduCloudTestResourceAttrNamePrefix+"Alias")
}

func testAccCfcAliasConfigUpdate() string {
	return fmt.Sprintf(`
resource "baiducloud_cfc_function" "default" {
  function_name = "%s"
  description   = "terraform create"
  environment = {
    "aaa": "bbb"
    "ccc": "ddd"
  }
  handler        = "index.handler"
  memory_size    = 128
  runtime        = "nodejs8.5"
  time_out       = 3
  code_file_name = "testFiles/cfcTestCode.zip"
}

resource "%s" "%s" {
  function_name    = baiducloud_cfc_function.default.function_name
  function_version = baiducloud_cfc_function.default.version
  alias_name       = "%s"
  description      = "terraform update alias"
}
`, BaiduCloudTestResourceAttrNamePrefix+"CFC",
		testAccCFCAliasResourceType, BaiduCloudTestResourceName, BaiduCloudTestResourceAttrNamePrefix+"Alias")
}
