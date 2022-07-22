package baiducloud

import (
	"fmt"
	"testing"

	"github.com/baidubce/bce-sdk-go/services/cfc"
	"github.com/baidubce/bce-sdk-go/services/cfc/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccCFCTriggerResourceType = "baiducloud_cfc_trigger"
	testAccCFCTriggerResourceName = testAccCFCTriggerResourceType + "." + BaiduCloudTestResourceName
)

func TestAccBaiduCloudCFCTrigger_HttpTrigger(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCFCTriggerDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccCfcHttpTriggerConfig(BaiduCloudTestResourceTypeNameCfcTrigger),
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
				Config: testAccCfcHttpTriggerConfigUpdate(BaiduCloudTestResourceTypeNameCfcTrigger),
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

func TestAccBaiduCloudCFCTrigger_CDNTrigger(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCFCTriggerDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccCfcCDNTriggerConfig(BaiduCloudTestResourceTypeNameCfcTrigger),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "source_type", "cdn"),
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "cdn_event_type", "CachedObjectsBlocked"),
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "status", "disabled"),
					resource.TestCheckResourceAttrSet(testAccCFCTriggerResourceName, "relation_id"),
					resource.TestCheckResourceAttrSet(testAccCFCTriggerResourceName, "target"),
				),
			},
			{
				Config: testAccCfcCDNTriggerConfigUpdate(BaiduCloudTestResourceTypeNameCfcTrigger),
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

func TestAccBaiduCloudCFCTrigger_BOSTrigger(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCFCTriggerDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccCfcBOSTriggerConfig(BaiduCloudTestResourceTypeNameCfcTrigger),
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
				Config: testAccCfcBOSTriggerConfigUpdate(BaiduCloudTestResourceTypeNameCfcTrigger),
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

func TestAccBaiduCloudCFCTrigger_DuerOSTrigger(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCFCTriggerDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccCfcDuerOSTriggerConfig(BaiduCloudTestResourceTypeNameCfcTrigger),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCFCTriggerResourceName, "source_type", "dueros"),
					resource.TestCheckResourceAttrSet(testAccCFCTriggerResourceName, "relation_id"),
					resource.TestCheckResourceAttrSet(testAccCFCTriggerResourceName, "target"),
				),
			},
		},
	})
}

func TestAccBaiduCloudCFCTrigger_CrontabTrigger(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCFCTriggerDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccCfcCrontabTriggerConfig(BaiduCloudTestResourceTypeNameCfcTrigger),
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
				Config: testAccCfcCrontabTriggerConfigUpdate(BaiduCloudTestResourceTypeNameCfcTrigger),
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

func testAccCfcHttpTriggerConfig(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_cfc_function" "default" {
  function_name  = "%s"
  description    = "created by terraform"
  handler        = "index.handler"
  memory_size    = 128
  runtime        = "nodejs12"
  time_out       = 3
  code_file_name = "testFiles/cfcTestCode.zip"
}

resource "baiducloud_cfc_trigger" "default" {
  source_type   = "http"
  target        = baiducloud_cfc_function.default.function_brn
  resource_path = "/test"
  method        = ["GET","PUT"]
  auth_type     = "iam"
}
`, name)
}

func testAccCfcHttpTriggerConfigUpdate(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_cfc_function" "default" {
  function_name  = "%s"
  description    = "created by terraform"
  handler        = "index.handler"
  memory_size    = 128
  runtime        = "nodejs12"
  time_out       = 3
  code_file_name = "testFiles/cfcTestCode.zip"
}

resource "baiducloud_cfc_trigger" "default" {
  source_type   = "http"
  target        = baiducloud_cfc_function.default.function_brn
  resource_path = "/test2"
  method        = ["GET","PUT","POST"]
  auth_type     = "iam"
}
`, name)
}

func testAccCfcCDNTriggerConfig(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_cfc_function" "default" {
  function_name  = "%s"
  description    = "created by terraform"
  handler        = "index.handler"
  memory_size    = 128
  runtime        = "nodejs12"
  time_out       = 3
  code_file_name = "testFiles/cfcTestCode.zip"
}

resource "baiducloud_cfc_trigger" "default" {
  source_type    = "cdn"
  target         = baiducloud_cfc_function.default.function_brn
}
`, name)
}

//TODO cfc 版本适应
//cdn_event_type = "CachedObjectsBlocked"
//status         = "disabled"

func testAccCfcCDNTriggerConfigUpdate(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_cfc_function" "default" {
  function_name  = "%s"
  description    = "created by terraform"
  handler        = "index.handler"
  memory_size    = 128
  runtime        = "nodejs12"
  time_out       = 3
  code_file_name = "testFiles/cfcTestCode.zip"
}

resource "baiducloud_cfc_trigger" "default" {
  source_type    = "cdn"
  target         = baiducloud_cfc_function.default.function_brn
  cdn_event_type = "CachedObjectsPushed"
  status         = "enabled"
}
`, name)
}

func testAccCfcBOSTriggerConfig(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_bos_bucket" "default" {
  bucket = "%s"
  acl    = "public-read-write"
}

resource "baiducloud_cfc_function" "default" {
  function_name  = "%s"
  description    = "created by terraform"
  handler        = "index.handler"
  memory_size    = 128
  runtime        = "nodejs12"
  time_out       = 3
  code_file_name = "testFiles/cfcTestCode.zip"
}

resource "baiducloud_cfc_trigger" "default" {
  source_type    = "bos"
  bucket         = baiducloud_bos_bucket.default.bucket
  target         = baiducloud_cfc_function.default.function_brn
  name           = "hehehehe"
  status         = "disabled"
  bos_event_type = ["PutObject", "PostObject"]
  resource       = "/undefined"
}
`, name+"-bucket-new", name)
}

func testAccCfcBOSTriggerConfigUpdate(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_bos_bucket" "default" {
  bucket = "%s"
  acl    = "public-read-write"
}

resource "baiducloud_cfc_function" "default" {
  function_name  = "%s"
  description    = "created by terraform"
  handler        = "index.handler"
  memory_size    = 128
  runtime        = "nodejs12"
  time_out       = 3
  code_file_name = "testFiles/cfcTestCode.zip"
}

resource "baiducloud_cfc_trigger" "default" {
  source_type    = "bos"
  bucket         = baiducloud_bos_bucket.default.bucket
  target         = baiducloud_cfc_function.default.function_brn
  name           = "hehehehe"
  status         = "enabled"
  bos_event_type = ["PostObject"]
  resource       = "/undefined"
}
`, name+"-bucket-new", name)
}

func testAccCfcDuerOSTriggerConfig(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_cfc_function" "default" {
  function_name  = "%s"
  description    = "created by terraform"
  handler        = "index.handler"
  memory_size    = 128
  runtime        = "nodejs12"
  time_out       = 3
  code_file_name = "testFiles/cfcTestCode.zip"
}

resource "baiducloud_cfc_trigger" "default" {
  source_type = "dueros"
  target      = baiducloud_cfc_function.default.function_brn
}
`, name)
}

func testAccCfcCrontabTriggerConfig(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_cfc_function" "default" {
  function_name  = "%s"
  description    = "created by terraform"
  handler        = "index.handler"
  memory_size    = 128
  runtime        = "nodejs12"
  time_out       = 3
  code_file_name = "testFiles/cfcTestCode.zip"
}

resource "baiducloud_cfc_trigger" "default" {
  source_type         = "crontab"
  target              = baiducloud_cfc_function.default.function_brn
  name                = "hahahaha"
  enabled             = "Disabled"
  schedule_expression = "cron(* * * * *)"
}
`, name)
}

func testAccCfcCrontabTriggerConfigUpdate(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_cfc_function" "default" {
  function_name  = "%s"
  description    = "created by terraform"
  handler        = "index.handler"
  memory_size    = 128
  runtime        = "nodejs12"
  time_out       = 3
  code_file_name = "testFiles/cfcTestCode.zip"
}

resource "baiducloud_cfc_trigger" "default" {
  source_type         = "crontab"
  target              = baiducloud_cfc_function.default.function_brn
  name                = "hahahaha"
  enabled             = "Enabled"
  schedule_expression = "cron(0 10 * * ?)"
}
`, name)
}
