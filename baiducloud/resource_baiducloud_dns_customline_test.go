package baiducloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const (
	testAccDnscustomlineResourceType = "baiducloud_dns_customline"
	testAccDnscustomlineResourceName = testAccDnscustomlineResourceType + "." + BaiduCloudTestResourceName
)

func TestAccBaiduCloudDnscustomlineSimple(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccDnscustomlineConfig("tf-test-acc-dns_customline"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccDnscustomlineResourceName),
					resource.TestCheckResourceAttrSet(testAccDnscustomlineResourceName, "name"),
				),
			},
		},
	})
}

func TestAccBaiduCloudDnscustomlineUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccBLBListenerDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccDnscustomlineConfig("tf-test-acc-dns_customline"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccDnscustomlineResourceName),
					resource.TestCheckResourceAttrSet(testAccDnscustomlineResourceName, "name"),
					resource.TestCheckResourceAttrSet(testAccDnscustomlineResourceName, "lines"),
				),
			},
			{
				Config: testAccDnscustomlineConfig("tf-test-acc-dns_customline"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccDnscustomlineResourceName),
					resource.TestCheckResourceAttrSet(testAccDnscustomlineResourceName, "name"),
					resource.TestCheckResourceAttrSet(testAccDnscustomlineResourceName, "lines"),
				),
			},
		},
	})
}

func TestAccBaiduCloudDnscustomlineUpdateName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccBLBListenerDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccDnscustomlineConfig("tf-test-acc-dns_customline"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccDnscustomlineResourceName),
					resource.TestCheckResourceAttrSet(testAccDnscustomlineResourceName, "name"),
					resource.TestCheckResourceAttrSet(testAccDnscustomlineResourceName, "lines"),
				),
			},
			{
				Config: testAccDnscustomlineConfigUpdateName("tf-test-acc-dns_customline"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccDnscustomlineResourceName),
					resource.TestCheckResourceAttrSet(testAccDnscustomlineResourceName, "name"),
					resource.TestCheckResourceAttrSet(testAccDnscustomlineResourceName, "lines"),
				),
			},
		},
	})
}

func TestAccBaiduCloudDnscustomlineUpdateLine(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccBLBListenerDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccDnscustomlineConfig("tf-test-acc-dns_customline"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccDnscustomlineResourceName),
					resource.TestCheckResourceAttrSet(testAccDnscustomlineResourceName, "name"),
					resource.TestCheckResourceAttrSet(testAccDnscustomlineResourceName, "lines"),
				),
			},
			{
				Config: testAccDnscustomlineConfigUpdateLines("tf-test-acc-dns_customline"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccDnscustomlineResourceName),
					resource.TestCheckResourceAttrSet(testAccDnscustomlineResourceName, "name"),
					resource.TestCheckResourceAttrSet(testAccDnscustomlineResourceName, "lines"),
				),
			},
		},
	})
}

func testAccDnscustomlineConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_dns_customline" "default" {
  name              = "testname"
  lines             = ["zhejiang.ct", "shanxi.ct"]
}
`, name)
}

func testAccDnscustomlineConfigUpdateName(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_dns_customline" "default" {
  name              = "testname1"
  lines             = ["zhejiang.ct", "shanxi.ct"]
}
`, name)
}

func testAccDnscustomlineConfigUpdateLines(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_dns_customline" "default" {
  name              = "testname"
  lines             = ["zhejiang.ct"]
}
`, name)
}
