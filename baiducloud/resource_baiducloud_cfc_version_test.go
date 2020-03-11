package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccCFCVersionResourceType = "baiducloud_cfc_version"
	testAccCFCVersionResourceName = testAccCFCVersionResourceType + "." + BaiduCloudTestResourceName
)

//lintignore:AT003
func TestAccBaiduCloudCFCVersion(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCFCVersionDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccCfcVersionConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCFCVersionResourceName, "version", "1"),
					resource.TestCheckResourceAttr(testAccCFCVersionResourceName, "version_description", "test-BaiduAccVersion"),
					resource.TestCheckResourceAttr(testAccCFCVersionResourceName, "function_name", "test-BaiduAccCFC"),
					resource.TestCheckResourceAttr(testAccCFCVersionResourceName, "description", "terraform create"),
					resource.TestCheckResourceAttr(testAccCFCVersionResourceName, "memory_size", "128"),
					resource.TestCheckResourceAttr(testAccCFCVersionResourceName, "handler", "index.handler"),
					resource.TestCheckResourceAttr(testAccCFCVersionResourceName, "runtime", "nodejs12"),
					resource.TestCheckResourceAttr(testAccCFCVersionResourceName, "time_out", "3"),
					resource.TestCheckResourceAttr(testAccCFCVersionResourceName, "log_type", "none"),
					resource.TestCheckResourceAttrSet(testAccCFCVersionResourceName, "update_time"),
					resource.TestCheckResourceAttrSet(testAccCFCVersionResourceName, "last_modified"),
					resource.TestCheckResourceAttrSet(testAccCFCVersionResourceName, "code_sha256"),
					resource.TestCheckResourceAttrSet(testAccCFCVersionResourceName, "function_brn"),
					resource.TestCheckResourceAttrSet(testAccCFCVersionResourceName, "function_arn"),
					resource.TestCheckResourceAttrSet(testAccCFCVersionResourceName, "commit_id"),
					resource.TestCheckResourceAttrSet(testAccCFCVersionResourceName, "uid"),
					resource.TestCheckResourceAttrSet(testAccCFCVersionResourceName, "region"),
				),
			},
		},
	})
}

func testAccCFCVersionDestory(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	cfcService := CFCService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccCFCVersionResourceType {
			continue
		}

		functionName := rs.Primary.Attributes["function_name"]
		functionVersion := rs.Primary.Attributes["version"]
		_, err := cfcService.CFCGetVersionsByFunction(functionName, functionVersion)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		return WrapError(Error("CFC Function version still exist"))
	}

	return nil
}

func testAccCfcVersionConfig() string {
	return fmt.Sprintf(`
data "baiducloud_zones" "default" {}

resource "baiducloud_vpc" "default" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "baiducloud_subnet" "default" {
  name        = "%s"
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.3.0/24"
  vpc_id      = baiducloud_vpc.default.id
  subnet_type = "BCC"
}

resource "baiducloud_security_group" "default" {
  name   = "%s"
  vpc_id = baiducloud_vpc.default.id
}

resource "baiducloud_cfc_function" "default" {
  function_name  = "%s"
  description    = "terraform create"
  handler        = "index.handler"
  memory_size    = 128
  runtime        = "nodejs12"
  time_out       = 3
  code_file_name = "testFiles/cfcTestCode.zip"
  vpc_config {
  	subnet_ids         = [baiducloud_subnet.default.id]
  	security_group_ids = [baiducloud_security_group.default.id]
  }
}

resource "%s" "%s" {
  function_name       = baiducloud_cfc_function.default.function_name
  version_description = "%s"
  code_sha256         = baiducloud_cfc_function.default.code_sha256
  log_type            = "none"
}
`, BaiduCloudTestResourceAttrNamePrefix+"VPC",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceAttrNamePrefix+"SecurityGroup",
		BaiduCloudTestResourceAttrNamePrefix+"CFC",
		testAccCFCVersionResourceType, BaiduCloudTestResourceName, BaiduCloudTestResourceAttrNamePrefix+"Version")
}
