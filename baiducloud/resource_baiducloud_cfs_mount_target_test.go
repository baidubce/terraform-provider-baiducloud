package baiducloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const (
	testAccCfsMountTargetResourceType = "baiducloud_cfs_mount_target"
	testAccCfsMountTargetResourceName = testAccCfsMountTargetResourceType + "." + BaiduCloudTestResourceTypeNameCfsMountTarget
)

func TestAccBaiduCloudCfsMountTarget(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccCfsMountTargetConfig(BaiduCloudTestResourceTypeNameCfsMountTarget),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccCfsMountTargetResourceName),
					resource.TestCheckResourceAttrSet(testAccCfsMountTargetResourceName, "domain"),
				),
			},
		},
	})
}

func testAccCfsMountTargetConfig(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_vpc" "vpc" {
  name = "terraform_vpc"
  cidr = "172.16.0.0/20"
}
resource "baiducloud_subnet" "subnet" {
  name        = "terraform_subnet"
  zone_name   = "cn-bj-d"
  cidr        = "172.16.0.0/24"
  vpc_id      = baiducloud_vpc.vpc.id
  description = "terraform test subnet"
}

resource "baiducloud_cfs" "default" {
  name = "terraform_test"
  zone = "zoneD"
}

resource "baiducloud_cfs_mount_target" "%s" {
  fs_id = baiducloud_cfs.default.id
  subnet_id = baiducloud_subnet.subnet.id
  vpc_id = baiducloud_vpc.vpc.id
}
`, name)
}
