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
	testAccCFCTriggerResourceType = "baiducloud_cfc_trigger"
	testAccCFCTriggerResourceName = testAccCFCTriggerResourceType + "." + BaiduCloudTestResourceName
)

func TestAccBaiduCloudCFCHttpTrigger(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCFCTriggerDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccCfcHttpTriggerConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "source_type", "http"),
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "resource_path", "/test"),
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "auth_type", "iam"),
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "method.#", "2"),
					resource.TestCheckResourceAttrSet(testAccCFCTriggerResourceName, "relation_id"),
					resource.TestCheckResourceAttrSet(testAccCFCTriggerResourceName, "target"),
				),
			},
			{
				Config: testAccCfcHttpTriggerConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "source_type", "http"),
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "resource_path", "/test2"),
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "auth_type", "iam"),
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "method.#", "3"),
					resource.TestCheckResourceAttrSet(testAccCFCTriggerResourceName, "relation_id"),
					resource.TestCheckResourceAttrSet(testAccCFCTriggerResourceName, "target"),
				),
			},
		},
	})
}

func TestAccBaiduCloudCFCCDNTrigger(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCFCTriggerDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccCfcCDNTriggerConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "source_type", "cdn"),
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "cdn_event_type", "CachedObjectsBlocked"),
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "status", "disabled"),
					resource.TestCheckResourceAttrSet(testAccCFCTriggerResourceName, "relation_id"),
					resource.TestCheckResourceAttrSet(testAccCFCTriggerResourceName, "target"),
				),
			},
			{
				Config: testAccCfcCDNTriggerConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "source_type", "cdn"),
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "cdn_event_type", "CachedObjectsPushed"),
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "status", "enabled"),
					resource.TestCheckResourceAttrSet(testAccCFCTriggerResourceName, "relation_id"),
					resource.TestCheckResourceAttrSet(testAccCFCTriggerResourceName, "target"),
				),
			},
		},
	})
}

func TestAccBaiduCloudCFCBOSTrigger(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCFCTriggerDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccCfcBOSTriggerConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "source_type", "bos"),
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "name", "hehehehe"),
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "status", "disabled"),
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "bos_event_type.#", "2"),
					resource.TestCheckResourceAttrSet(testAccCFCTriggerResourceName, "relation_id"),
					resource.TestCheckResourceAttrSet(testAccCFCTriggerResourceName, "target"),
				),
			},
			{
				Config: testAccCfcBOSTriggerConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "source_type", "bos"),
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "name", "hehehehe"),
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "status", "enabled"),
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "bos_event_type.#", "1"),
					resource.TestCheckResourceAttrSet(testAccCFCTriggerResourceName, "relation_id"),
					resource.TestCheckResourceAttrSet(testAccCFCTriggerResourceName, "target"),
				),
			},
		},
	})
}

func TestAccBaiduCloudCFCDuerOSTrigger(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCFCTriggerDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccCfcDuerOSTriggerConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "source_type", "dueros"),
					resource.TestCheckResourceAttrSet(testAccCFCTriggerResourceName, "relation_id"),
					resource.TestCheckResourceAttrSet(testAccCFCTriggerResourceName, "target"),
				),
			},
		},
	})
}

func TestAccBaiduCloudCFCCrontabTrigger(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCFCTriggerDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccCfcCrontabTriggerConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "source_type", "crontab"),
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "name", "hahahaha"),
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "enabled", "Disabled"),
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "schedule_expression", "cron(* * * * *)"),
					resource.TestCheckResourceAttrSet(testAccCFCTriggerResourceName, "relation_id"),
					resource.TestCheckResourceAttrSet(testAccCFCTriggerResourceName, "target"),
				),
			},
			{
				Config: testAccCfcCrontabTriggerConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "source_type", "crontab"),
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "name", "hahahaha"),
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "enabled", "Enabled"),
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "schedule_expression", "cron(0 10 * * ?)"),
					resource.TestCheckResourceAttrSet(testAccCFCTriggerResourceName, "relation_id"),
					resource.TestCheckResourceAttrSet(testAccCFCTriggerResourceName, "target"),
				),
			},
		},
	})
}

func testAccCFCTriggerDestory(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	cfcService := CFCService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccCFCTriggerResourceType {
			continue
		}

		functionBrn := rs.Primary.Attributes["target"]
		_, err := client.WithCFCClient(func(client *cfc.Client) (i interface{}, e error) {
			return client.GetFunction(&api.GetFunctionArgs{FunctionName: functionBrn})
		})
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		relationId := rs.Primary.Attributes["relation_id"]
		_, err = cfcService.CFCGetTriggerByFunction(functionBrn, relationId)
		if err != nil {
			if NotFoundError(err) {
				continue
			}

			return WrapError(err)
		}

		return WrapError(Error("CFC Function Trigger still exist"))
	}

	return nil
}

func testAccCfcHttpTriggerConfig() string {
	return fmt.Sprintf(`
resource "baiducloud_cfc_function" "default" {
  function_name  = "%s"
  description    = "terraform create"
  handler        = "index.handler"
  memory_size    = 128
  runtime        = "nodejs8.5"
  time_out       = 3
  code_file_name = "testFiles/cfcTestCode.zip"
}

resource "%s" "%s" {
  source_type   = "http"
  target        = baiducloud_cfc_function.default.function_brn
  resource_path = "/test"
  method        = ["GET","PUT"]
  auth_type     = "iam"
}
`, BaiduCloudTestResourceAttrNamePrefix+"CFC",
		testAccCFCTriggerResourceType, BaiduCloudTestResourceName)
}

func testAccCfcHttpTriggerConfigUpdate() string {
	return fmt.Sprintf(`
resource "baiducloud_cfc_function" "default" {
  function_name  = "%s"
  description    = "terraform create"
  handler        = "index.handler"
  memory_size    = 128
  runtime        = "nodejs8.5"
  time_out       = 3
  code_file_name = "testFiles/cfcTestCode.zip"
}

resource "%s" "%s" {
  source_type   = "http"
  target        = baiducloud_cfc_function.default.function_brn
  resource_path = "/test2"
  method        = ["GET","PUT","POST"]
  auth_type     = "iam"
}
`, BaiduCloudTestResourceAttrNamePrefix+"CFC",
		testAccCFCTriggerResourceType, BaiduCloudTestResourceName)
}

func testAccCfcCDNTriggerConfig() string {
	return fmt.Sprintf(`
resource "baiducloud_cfc_function" "default" {
  function_name  = "%s"
  description    = "terraform create"
  handler        = "index.handler"
  memory_size    = 128
  runtime        = "nodejs8.5"
  time_out       = 3
  code_file_name = "testFiles/cfcTestCode.zip"
}

resource "%s" "%s" {
  source_type    = "cdn"
  target         = baiducloud_cfc_function.default.function_brn
  cdn_event_type = "CachedObjectsBlocked"
  status         = "disabled"
}
`, BaiduCloudTestResourceAttrNamePrefix+"CFC",
		testAccCFCTriggerResourceType, BaiduCloudTestResourceName)
}

func testAccCfcCDNTriggerConfigUpdate() string {
	return fmt.Sprintf(`
resource "baiducloud_cfc_function" "default" {
  function_name  = "%s"
  description    = "terraform create"
  handler        = "index.handler"
  memory_size    = 128
  runtime        = "nodejs8.5"
  time_out       = 3
  code_file_name = "testFiles/cfcTestCode.zip"
}

resource "%s" "%s" {
  source_type    = "cdn"
  target         = baiducloud_cfc_function.default.function_brn
  cdn_event_type = "CachedObjectsPushed"
  status         = "enabled"
}
`, BaiduCloudTestResourceAttrNamePrefix+"CFC",
		testAccCFCTriggerResourceType, BaiduCloudTestResourceName)
}

func testAccCfcBOSTriggerConfig() string {
	return fmt.Sprintf(`
resource "baiducloud_bos_bucket" "default" {
  bucket = "%s"
  acl    = "public-read-write"
}

resource "baiducloud_cfc_function" "default" {
  function_name  = "%s"
  description    = "terraform create"
  handler        = "index.handler"
  memory_size    = 128
  runtime        = "nodejs8.5"
  time_out       = 3
  code_file_name = "testFiles/cfcTestCode.zip"
}

resource "%s" "%s" {
  source_type    = "bos"
  bucket         = baiducloud_bos_bucket.default.bucket
  target         = baiducloud_cfc_function.default.function_brn
  name           = "hehehehe"
  status         = "disabled"
  bos_event_type = ["PutObject", "PostObject"]
  resource       = "/undefined"
}
`, BaiduCloudTestBucketResourceAttrNamePrefix+"bossss",
		BaiduCloudTestResourceAttrNamePrefix+"CFC",
		testAccCFCTriggerResourceType, BaiduCloudTestResourceName)
}

func testAccCfcBOSTriggerConfigUpdate() string {
	return fmt.Sprintf(`
resource "baiducloud_bos_bucket" "default" {
  bucket = "%s"
  acl    = "public-read-write"
}

resource "baiducloud_cfc_function" "default" {
  function_name  = "%s"
  description    = "terraform create"
  handler        = "index.handler"
  memory_size    = 128
  runtime        = "nodejs8.5"
  time_out       = 3
  code_file_name = "testFiles/cfcTestCode.zip"
}

resource "%s" "%s" {
  source_type    = "bos"
  bucket         = baiducloud_bos_bucket.default.bucket
  target         = baiducloud_cfc_function.default.function_brn
  name           = "hehehehe"
  status         = "enabled"
  bos_event_type = ["PostObject"]
  resource       = "/undefined"
}
`, BaiduCloudTestBucketResourceAttrNamePrefix+"bossss",
		BaiduCloudTestResourceAttrNamePrefix+"CFC",
		testAccCFCTriggerResourceType, BaiduCloudTestResourceName)
}

func testAccCfcDuerOSTriggerConfig() string {
	return fmt.Sprintf(`
resource "baiducloud_cfc_function" "default" {
  function_name  = "%s"
  description    = "terraform create"
  handler        = "index.handler"
  memory_size    = 128
  runtime        = "nodejs8.5"
  time_out       = 3
  code_file_name = "testFiles/cfcTestCode.zip"
}

resource "%s" "%s" {
  source_type = "dueros"
  target      = baiducloud_cfc_function.default.function_brn
}
`, BaiduCloudTestResourceAttrNamePrefix+"CFC",
		testAccCFCTriggerResourceType, BaiduCloudTestResourceName)
}

func testAccCfcCrontabTriggerConfig() string {
	return fmt.Sprintf(`
resource "baiducloud_cfc_function" "default" {
  function_name  = "%s"
  description    = "terraform create"
  handler        = "index.handler"
  memory_size    = 128
  runtime        = "nodejs8.5"
  time_out       = 3
  code_file_name = "testFiles/cfcTestCode.zip"
}

resource "%s" "%s" {
  source_type         = "crontab"
  target              = baiducloud_cfc_function.default.function_brn
  name                = "hahahaha"
  enabled             = "Disabled"
  schedule_expression = "cron(* * * * *)"
}
`, BaiduCloudTestResourceAttrNamePrefix+"CFC",
		testAccCFCTriggerResourceType, BaiduCloudTestResourceName)
}

func testAccCfcCrontabTriggerConfigUpdate() string {
	return fmt.Sprintf(`
resource "baiducloud_cfc_function" "default" {
  function_name  = "%s"
  description    = "terraform create"
  handler        = "index.handler"
  memory_size    = 128
  runtime        = "nodejs8.5"
  time_out       = 3
  code_file_name = "testFiles/cfcTestCode.zip"
}

resource "%s" "%s" {
  source_type         = "crontab"
  target              = baiducloud_cfc_function.default.function_brn
  name                = "hahahaha"
  enabled             = "Enabled"
  schedule_expression = "cron(0 10 * * ?)"
}
`, BaiduCloudTestResourceAttrNamePrefix+"CFC",
		testAccCFCTriggerResourceType, BaiduCloudTestResourceName)
}
