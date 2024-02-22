package baiducloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const (
	testAccDnsrecordResourceType = "baiducloud_dns_record"
	testAccDnsrecordResourceName = testAccDnsrecordResourceType + "." + BaiduCloudTestResourceName
)

func TestAccBaiduCloudDnsrecordSimple(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccDnsrecordConfig("tf-test-acc-dns_record"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccDnsrecordResourceName),
					resource.TestCheckResourceAttrSet(testAccDnsrecordResourceName, "zone_name"),
				),
			},
		},
	})
}

func TestAccBaiduCloudDnsrecordUpdaterr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccBLBListenerDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccDnsrecordConfig("tf-test-acc-dns_record"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccDnsrecordResourceName),
					resource.TestCheckResourceAttrSet(testAccDnsrecordResourceName, "zone_name"),
					resource.TestCheckResourceAttrSet(testAccDnsrecordResourceName, "rr"),
				),
			},
			{
				Config: testAccDnsrecordConfigUpdaterr("tf-test-acc-dns_record"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccDnsrecordResourceName),
					resource.TestCheckResourceAttrSet(testAccDnsrecordResourceName, "zone_name"),
					resource.TestCheckResourceAttrSet(testAccDnsrecordResourceName, "rr"),
				),
			},
		},
	})
}

func TestAccBaiduCloudDnsrecordUpdatetype(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccBLBListenerDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccDnsrecordConfig("tf-test-acc-dns_record"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccDnsrecordResourceName),
					resource.TestCheckResourceAttrSet(testAccDnsrecordResourceName, "name"),
					resource.TestCheckResourceAttrSet(testAccDnsrecordResourceName, "type"),
				),
			},
			{
				Config: testAccDnsrecordConfigUpdatetype("tf-test-acc-dns_record"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccDnsrecordResourceName),
					resource.TestCheckResourceAttrSet(testAccDnsrecordResourceName, "name"),
					resource.TestCheckResourceAttrSet(testAccDnsrecordResourceName, "type"),
				),
			},
		},
	})
}

func TestAccBaiduCloudDnsrecordUpdatevalue(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccBLBListenerDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccDnsrecordConfig("tf-test-acc-dns_record"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccDnsrecordResourceName),
					resource.TestCheckResourceAttrSet(testAccDnsrecordResourceName, "name"),
					resource.TestCheckResourceAttrSet(testAccDnsrecordResourceName, "value"),
				),
			},
			{
				Config: testAccDnsrecordConfigUpdatevalue("tf-test-acc-dns_record"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccDnsrecordResourceName),
					resource.TestCheckResourceAttrSet(testAccDnsrecordResourceName, "name"),
					resource.TestCheckResourceAttrSet(testAccDnsrecordResourceName, "value"),
				),
			},
		},
	})
}

func TestAccBaiduCloudDnsrecordUpdateenable(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccBLBListenerDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccDnsrecordConfig("tf-test-acc-dns_record"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccDnsrecordResourceName),
					resource.TestCheckResourceAttrSet(testAccDnsrecordResourceName, "name"),
					resource.TestCheckResourceAttrSet(testAccDnsrecordResourceName, "value"),
				),
			},
			{
				Config: testAccDnsrecordConfigUpdateenable("tf-test-acc-dns_record"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccDnsrecordResourceName),
					resource.TestCheckResourceAttrSet(testAccDnsrecordResourceName, "name"),
					resource.TestCheckResourceAttrSet(testAccDnsrecordResourceName, "value"),
				),
			},
		},
	})
}

func TestAccBaiduCloudDnsrecordUpdatedisable(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccBLBListenerDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccDnsrecordConfig("tf-test-acc-dns_record"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccDnsrecordResourceName),
					resource.TestCheckResourceAttrSet(testAccDnsrecordResourceName, "name"),
					resource.TestCheckResourceAttrSet(testAccDnsrecordResourceName, "value"),
				),
			},
			{
				Config: testAccDnsrecordConfigUpdatedisable("tf-test-acc-dns_record"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccDnsrecordResourceName),
					resource.TestCheckResourceAttrSet(testAccDnsrecordResourceName, "name"),
					resource.TestCheckResourceAttrSet(testAccDnsrecordResourceName, "value"),
				),
			},
		},
	})
}

func testAccDnsrecordConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_dns_record" "default" {
  zone_name              = "testZoneName"
  rr                     = "rr"
  type                   = "type"
  value                  = "value"
}
`, name)
}

func testAccDnsrecordConfigUpdaterr(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_dns_record" "default" {
  zone_name              = "testZoneName"
  rr                     = "rr1"
  type                   = "type"
  value                  = "value"
}
`, name)
}

func testAccDnsrecordConfigUpdatetype(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_dns_record" "default" {
  zone_name              = "testZoneName"
  rr                     = "rr"
  type                   = "type1"
  value                  = "value"
}
`, name)
}

func testAccDnsrecordConfigUpdatevalue(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_dns_record" "default" {
  zone_name              = "testZoneName"
  rr                     = "rr"
  type                   = "type"
  value                  = "value1"
}
`, name)
}

func testAccDnsrecordConfigUpdatedisable(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_dns_record" "default" {
  zone_name              = "testZoneName"
  rr                     = "rr"
  type                   = "type"
  value                  = "value"
  action                 = "disable"
}
`, name)
}

func testAccDnsrecordConfigUpdateenable(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_dns_record" "default" {
  zone_name              = "testZoneName"
  rr                     = "rr"
  type                   = "type"
  value                  = "value"
  action                 = "enable"
}
`, name)
}
