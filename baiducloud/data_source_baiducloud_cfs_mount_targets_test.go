package baiducloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const (
	testAccCfsMountTargetDataSourceName          = "data.baiducloud_cfs_mount_targets.default"
	testAccCfsMountTargetDataSourceAttrKeyPrefix = "mount_targets.0."
)

func TestAccBaiduCloudCfsMountTargetsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccCfsMountTargetsConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccCfsMountTargetDataSourceName),
					resource.TestCheckResourceAttrSet(testAccCfsMountTargetDataSourceName, testAccCfsMountTargetDataSourceAttrKeyPrefix+"domain"),
				),
			},
		},
	})
}

func testAccCfsMountTargetsConfig() string {
	return fmt.Sprintf(`
resource "baiducloud_vpc" "vpc" {
  name = "terraform_vpc"
  cidr = "172.16.0.0/16"
}
resource "baiducloud_subnet" "subnet" {
  name        = "terraform_subnet"
  zone_name   = "cn-bj-c"
  cidr        = "172.16.128.0/24"
  vpc_id      = baiducloud_vpc.vpc.id
  description = "terraform test subnet"
}

resource "baiducloud_cfs" "default" {
  name = "terraform_test"
  zone = "zoneD"
}

resource "baiducloud_cfs_mount_target" "default" {
  fs_id = baiducloud_cfs.default.id
  subnet_id = baiducloud_subnet.subnet.id
  vpc_id = baiducloud_vpc.vpc.id
}

data "baiducloud_cfs_mount_targets" "default" {
  fs_id = baiducloud_cfs.default.id
  filter{
    name = "mount_id"
    values = [baiducloud_cfs_mount_target.default.id]
  }
}
`)
}
