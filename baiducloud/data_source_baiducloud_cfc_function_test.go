package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

//lintignore:AT003
func TestAccBaiduCloudCFCFunctionDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCFCFunctionDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccCfcDataSourceConfig(BaiduCloudTestResourceTypeNameCfcFunction),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCFCResourceName, "function_name", BaiduCloudTestResourceTypeNameCfcFunction),
					resource.TestCheckResourceAttr(testAccCFCResourceName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(testAccCFCResourceName, "memory_size", "128"),
					resource.TestCheckResourceAttr(testAccCFCResourceName, "handler", "index.handler"),
					resource.TestCheckResourceAttr(testAccCFCResourceName, "runtime", "nodejs12"),
					resource.TestCheckResourceAttr(testAccCFCResourceName, "time_out", "3"),
					resource.TestCheckResourceAttr(testAccCFCResourceName, "environment.%", "2"),
					resource.TestCheckResourceAttr(testAccCFCResourceName, "vpc_config.#", "1"),
					resource.TestCheckResourceAttrSet(testAccCFCResourceName, "update_time"),
					resource.TestCheckResourceAttrSet(testAccCFCResourceName, "last_modified"),
					resource.TestCheckResourceAttrSet(testAccCFCResourceName, "code_sha256"),
					resource.TestCheckResourceAttrSet(testAccCFCResourceName, "function_brn"),
					resource.TestCheckResourceAttrSet(testAccCFCResourceName, "function_arn"),
					resource.TestCheckResourceAttrSet(testAccCFCResourceName, "commit_id"),
					resource.TestCheckResourceAttrSet(testAccCFCResourceName, "uid"),
					resource.TestCheckResourceAttrSet(testAccCFCResourceName, "region"),
				),
			},
		},
	})
}

func testAccCfcDataSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "baiducloud_zones" "default" {
  name_regex = ".*e$"
}

resource "baiducloud_vpc" "default" {
  name = var.name
  cidr = "192.168.0.0/16"
}

resource "baiducloud_subnet" "default" {
  name        = var.name
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.3.0/24"
  vpc_id      = baiducloud_vpc.default.id
  subnet_type = "BCC"
}

resource "baiducloud_security_group" "default" {
  name   = var.name
  vpc_id = baiducloud_vpc.default.id
}

resource "baiducloud_cfc_function" "default" {
  function_name = var.name
  description   = "created by terraform"
  environment = {
    "aaa": "bbb"
    "ccc": "ddd"
  }
  handler                        = "index.handler"
  memory_size                    = 128
  runtime                        = "nodejs12"
  time_out                       = 3
  code_file_name                 = "testFiles/cfcTestCode.zip"
  reserved_concurrent_executions = 10
  vpc_config {
    subnet_ids         = [baiducloud_subnet.default.id]
    security_group_ids = [baiducloud_security_group.default.id]
  }
}

data "baiducloud_cfc_function" "default" {
  function_name = baiducloud_cfc_function.default.function_name
}
`, name)
}
