package baiducloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/baidubce/bce-sdk-go/services/cfc"
	"github.com/baidubce/bce-sdk-go/services/cfc/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccCFCResourceType = "baiducloud_cfc_function"
	testAccCFCResourceName = testAccCFCResourceType + "." + BaiduCloudTestResourceName
)

func init() {
	resource.AddTestSweepers(testAccCFCResourceType, &resource.Sweeper{
		Name: testAccCFCResourceType,
		F:    testSweepCFCFunctions,
	})
}

func testSweepCFCFunctions(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("get BaiduCloud client error: %s", err)
	}
	client := rawClient.(*connectivity.BaiduClient)
	cfcService := CFCService{client}

	functions, err := cfcService.ListAllFunctions()
	if err != nil {
		return err
	}

	for _, f := range functions {
		name := f.FunctionName
		if !strings.HasPrefix(f.FunctionName, BaiduCloudTestResourceAttrNamePrefix) {
			log.Printf("[INFO] Skipping CFC Function: %s ", name)
			continue
		}

		log.Printf("[INFO] Deleting CFC Function: %s ", name)
		_, err := client.WithCFCClient(func(client *cfc.Client) (i interface{}, e error) {
			return nil, client.DeleteFunction(&api.DeleteFunctionArgs{
				FunctionName: name,
			})
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete CFC Function %s", name)
		}
	}

	return nil
}

func TestAccBaiduCloudCFCFunction(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCFCFunctionDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccCfcConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCFCResourceName, "function_name", "test-BaiduAccCFC"),
					resource.TestCheckResourceAttr(testAccCFCResourceName, "description", "terraform create"),
					resource.TestCheckResourceAttr(testAccCFCResourceName, "memory_size", "128"),
					resource.TestCheckResourceAttr(testAccCFCResourceName, "handler", "index.handler"),
					resource.TestCheckResourceAttr(testAccCFCResourceName, "runtime", "nodejs8.5"),
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
					resource.TestCheckResourceAttrSet(testAccCFCResourceName, "code_id"),
				),
			},
			{
				ResourceName:            testAccCFCResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"code_file_name", "code_bos_bucket", "code_bos_object", "code_file_dir", "reserved_concurrent_executions", "code_storage.location"},
			},
			{
				Config: testAccCfcConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCFCResourceName, "function_name", "test-BaiduAccCFC"),
					resource.TestCheckResourceAttr(testAccCFCResourceName, "description", "terraform update"),
					resource.TestCheckResourceAttr(testAccCFCResourceName, "memory_size", "256"),
					resource.TestCheckResourceAttr(testAccCFCResourceName, "handler", "index.handler2"),
					resource.TestCheckResourceAttr(testAccCFCResourceName, "runtime", "python2"),
					resource.TestCheckResourceAttr(testAccCFCResourceName, "time_out", "5"),
					resource.TestCheckResourceAttr(testAccCFCResourceName, "environment.%", "1"),
					resource.TestCheckResourceAttr(testAccCFCResourceName, "vpc_config.#", "0"),
					resource.TestCheckResourceAttrSet(testAccCFCResourceName, "update_time"),
					resource.TestCheckResourceAttrSet(testAccCFCResourceName, "last_modified"),
					resource.TestCheckResourceAttrSet(testAccCFCResourceName, "code_sha256"),
					resource.TestCheckResourceAttrSet(testAccCFCResourceName, "function_brn"),
					resource.TestCheckResourceAttrSet(testAccCFCResourceName, "function_arn"),
					resource.TestCheckResourceAttrSet(testAccCFCResourceName, "commit_id"),
					resource.TestCheckResourceAttrSet(testAccCFCResourceName, "uid"),
					resource.TestCheckResourceAttrSet(testAccCFCResourceName, "region"),
					resource.TestCheckResourceAttrSet(testAccCFCResourceName, "code_id"),
				),
			},
		},
	})
}

func testAccCFCFunctionDestory(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccEipResourceType {
			continue
		}

		_, err := client.WithCFCClient(func(client *cfc.Client) (i interface{}, e error) {
			return client.GetFunction(&api.GetFunctionArgs{FunctionName: rs.Primary.ID})
		})
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		return WrapError(Error("CFC Function still exist"))
	}

	return nil
}

func testAccCfcConfig() string {
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

resource "%s" "%s" {
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
  reserved_concurrent_executions = 10
  vpc_config {
    subnet_ids         = [baiducloud_subnet.default.id]
    security_group_ids = [baiducloud_security_group.default.id]
  }
}
`, BaiduCloudTestResourceAttrNamePrefix+"VPC",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceAttrNamePrefix+"SecurityGroup",
		testAccCFCResourceType, BaiduCloudTestResourceName, BaiduCloudTestResourceAttrNamePrefix+"CFC")
}

func testAccCfcConfigUpdate() string {
	return fmt.Sprintf(`
resource "%s" "%s" {
  function_name     = "%s"
  description       = "terraform update"
  environment = {
    "aaa": "bbb"
  }
  handler        = "index.handler2"
  memory_size    = 256
  runtime        = "python2"
  time_out       = 5
  code_file_dir  = "testFiles/cfcTestCode"
}
`, testAccCFCResourceType, BaiduCloudTestResourceName, BaiduCloudTestResourceAttrNamePrefix+"CFC")
}
