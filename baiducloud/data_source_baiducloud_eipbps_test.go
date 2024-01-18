package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccEipbpsDataSourceName          = "data.baiducloud_eipbps.default"
	testAccEipbpsDataSourceAttrKeyPrefix = "eip_bps.0."
)

//lintignore:AT003
func TestAccBaiduCloudEipbpsByNameDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEipbpsByNameDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccEipbpsDataSourceName),
					resource.TestCheckResourceAttrSet(testAccEipbpsDataSourceName, testAccEipbpsDataSourceAttrKeyPrefix+"name"),
				),
			},
		},
	})
}

func TestAccBaiduCloudEipbpsByIdDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEipbpsByIdDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccEipbpsDataSourceName),
					resource.TestCheckResourceAttrSet(testAccEipbpsDataSourceName, testAccEipbpsDataSourceAttrKeyPrefix+"id"),
				),
			},
		},
	})
}

func TestAccBaiduCloudEipbpsByBindTypeDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEipbpsByBindTypeDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccEipbpsDataSourceName),
					resource.TestCheckResourceAttrSet(testAccEipbpsDataSourceName, testAccEipbpsDataSourceAttrKeyPrefix+"bind_type"),
				),
			},
		},
	})
}

func TestAccBaiduCloudEipbpsByTypeDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEipbpsByTypeDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccEipbpsDataSourceName),
					resource.TestCheckResourceAttrSet(testAccEipbpsDataSourceName, testAccEipbpsDataSourceAttrKeyPrefix+"type"),
				),
			},
		},
	})
}

func testAccEipbpsByNameDataSourceConfig() string {
	return fmt.Sprintf(`
data "baiducloud_eipbps" "default" {
    name = "xxxx"    
}
`)
}

func testAccEipbpsByIdDataSourceConfig() string {
	return fmt.Sprintf(`
data "baiducloud_eipbps" "default" {
    id = "xxxx"    
}
`)
}

func testAccEipbpsByBindTypeDataSourceConfig() string {
	return fmt.Sprintf(`
data "baiducloud_eipbps" "default" {
    bind_type = "xxxx"    
}
`)
}

func testAccEipbpsByTypeDataSourceConfig() string {
	return fmt.Sprintf(`
data "baiducloud_eipbps" "default" {
    type = "xxxx"    
}
`)
}
