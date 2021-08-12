package baiducloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const (
	testAccSubnetsDataSourceName          = "data.baiducloud_subnets.default"
	testAccSubnetsDataSourceAttrKeyPrefix = "subnets.0."
)

//lintignore:AT003
func TestAccBaiduCloudSubnetsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccSubnetsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccSubnetsDataSourceName),
					resource.TestCheckResourceAttrSet(testAccSubnetsDataSourceName, testAccSubnetsDataSourceAttrKeyPrefix+"subnet_id"),
					resource.TestCheckResourceAttr(testAccSubnetsDataSourceName, testAccSubnetsDataSourceAttrKeyPrefix+"cidr", "192.168.1.0/24"),
					resource.TestCheckResourceAttr(testAccSubnetsDataSourceName, testAccSubnetsDataSourceAttrKeyPrefix+"name", BaiduCloudTestResourceTypeNameSubnet),
					resource.TestCheckResourceAttrSet(testAccSubnetsDataSourceName, testAccSubnetsDataSourceAttrKeyPrefix+"zone_name"),
					resource.TestCheckResourceAttrSet(testAccSubnetsDataSourceName, testAccSubnetsDataSourceAttrKeyPrefix+"vpc_id"),
					resource.TestCheckResourceAttrSet(testAccSubnetsDataSourceName, testAccSubnetsDataSourceAttrKeyPrefix+"subnet_type"),
					resource.TestCheckResourceAttr(testAccSubnetsDataSourceName, testAccSubnetsDataSourceAttrKeyPrefix+"description", "created by terraform"),
					resource.TestCheckResourceAttrSet(testAccSubnetsDataSourceName, testAccSubnetsDataSourceAttrKeyPrefix+"available_ip"),
					resource.TestCheckResourceAttr(testAccSubnetsDataSourceName, testAccSubnetsDataSourceAttrKeyPrefix+"tags.%", "1"),
					resource.TestCheckResourceAttr(testAccSubnetsDataSourceName, testAccSubnetsDataSourceAttrKeyPrefix+"tags.testKey", "testValue"),
				),
			},
		},
	})
}

const testAccSubnetsDataSourceConfig = `
data "baiducloud_zones" "default" {
  name_regex = ".*e$"
}

resource "baiducloud_vpc" "default" {
  name        = "tf-test-acc"
  description = "created by terraform"
  cidr        = "192.168.0.0/16"
}

resource "baiducloud_subnet" "default" {
  name        = "tf-test-acc-subnet"
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.1.0/24"
  description = "created by terraform"
  vpc_id      = baiducloud_vpc.default.id
  tags = {
    "testKey" = "testValue"
  }
}

data "baiducloud_subnets" "default" {
  subnet_id = baiducloud_subnet.default.id

  filter {
    name = "name"
    values = ["test-filter", "tf-test-acc*"]
  }

  filter {
    name = "cidr"
    values = ["192.168.1.0/24"]
  }
}
`
