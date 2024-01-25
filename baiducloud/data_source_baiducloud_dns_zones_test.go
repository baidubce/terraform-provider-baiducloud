package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccDnszonesDataSourceName          = "data.baiducloud_dns_zones.default"
	testAccDnszonesDataSourceAttrKeyPrefix = "zones.0."
)

//lintignore:AT003
func TestAccBaiduCloudDnszonesDataSourceSimple(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnszoneByNameDataSourceSimpleConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccDnszonesDataSourceName),
					resource.TestCheckResourceAttrSet(testAccDnszonesDataSourceName, testAccDnszonesDataSourceAttrKeyPrefix+"name"),
				),
			},
		},
	})
}

func TestAccBaiduCloudDnszonesByNameDataSourceSimple(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnszoneDataSourceSimpleConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccDnszonesDataSourceName),
					resource.TestCheckResourceAttrSet(testAccDnszonesDataSourceName, testAccDnszonesDataSourceAttrKeyPrefix+"name"),
				),
			},
		},
	})
}

func TestAccBaiduCloudDnszonesDataSourceFull(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnszoneDataSourceFullConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccDnszonesDataSourceName),
					resource.TestCheckResourceAttrSet(testAccDnszonesDataSourceName, testAccDnszonesDataSourceAttrKeyPrefix+"name"),
				),
			},
		},
	})
}

func TestAccBaiduCloudDnszonesByNameDataSourceFull(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnszoneByNameDataSourceFullConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccDnszonesDataSourceName),
					resource.TestCheckResourceAttrSet(testAccDnszonesDataSourceName, testAccDnszonesDataSourceAttrKeyPrefix+"name"),
				),
			},
		},
	})
}

func testAccDnszoneByNameDataSourceSimpleConfig() string {
	return fmt.Sprintf(`
data "baiducloud_dns_zones" "default" { 
}
`)
}

func testAccDnszoneDataSourceSimpleConfig() string {
	return fmt.Sprintf(`
data "baiducloud_dns_zones" "default" {
    name = "terraform.com"    
}
`)
}

func testAccDnszoneDataSourceFullConfig() string {
	return fmt.Sprintf(`
resource "baiducloud_dns_zone" "default" {
  name         = "terraform.com"
}

data "baiducloud_dns_zones" "default" {
}
`)
}

func testAccDnszoneByNameDataSourceFullConfig() string {
	return fmt.Sprintf(`
resource "baiducloud_dns_zone" "default" {
      name = "terraform.com"
}

data "baiducloud_dns_zones" "default" {
	  name = "terraform.com"
}
`)
}
