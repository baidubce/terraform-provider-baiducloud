package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccBbcInstancesDataSourceName          = "data.baiducloud_bbc_instances.default"
	testAccBbcInstancesDataSourceAttrKeyPrefix = "instances.0."
)

func TestAccBaiduCloudBbcInstancesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBbcInstancesDataSourceConfig(BaiduCloudTestResourceTypeNameBbcInstance),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBbcInstancesDataSourceName),
					resource.TestCheckResourceAttr(testAccBbcInstancesDataSourceName, testAccBbcInstancesDataSourceAttrKeyPrefix+"name", BaiduCloudTestResourceTypeNameInstance),
					resource.TestCheckResourceAttr(testAccBbcInstancesDataSourceName, testAccBbcInstancesDataSourceAttrKeyPrefix+"tags.%", "1"),
					resource.TestCheckResourceAttr(testAccBbcInstancesDataSourceName, testAccBbcInstancesDataSourceAttrKeyPrefix+"tags.testKey", "testValue"),
				),
			},
		},
	})
}

func testAccBbcInstancesDataSourceConfig(name string) string {
	return fmt.Sprintf(`
data "baiducloud_bbc_images" "bbc_images" {
  image_type = "BbcSystem"
  os_name    = "CentOS"
  filter {
    name   = "id"
    values = ["m-i2aoqIlx"]
  }
}
data "baiducloud_security_groups" "sg" {
  filter {
    name   = "name"
    values = ["默认安全组"]
  }
  filter {
    name   = "id"
    values = ["g-mwi8hx6dy4qb"]
  }
}
data "baiducloud_subnets" "subnets" {
  filter {
    name   = "zone_name"
    values = ["cn-bd-b"]
  }
  filter {
    name   = "name"
    values = ["系统预定义子网B"]
  }
}

data "baiducloud_bbc_flavors" "bbc_flavors" {
  filter {
    name   = "flavor_id"
    values = ["BBC-I4-HC04S"]
  }
}
resource "baiducloud_bbc_instance" "bbc_instance" {
  name = "%s"
  flavor_id = "${data.baiducloud_bbc_flavors.bbc_flavors.flavors.0.flavor_id}"
  image_id = "${data.baiducloud_bbc_images.bbc_images.images.0.id}"
  raid = "NoRaid"
  root_disk_size_in_gb = 40
  purchase_count = 1
  zone_name = "cn-bd-b"
  subnet_id = "${data.baiducloud_subnets.subnets.subnets.0.subnet_id}"
  security_groups = [
    "${data.baiducloud_security_groups.sg.security_groups.0.id}",
  ]
  billing = {
    payment_timing = "Postpaid"
  }
  tags = {
    "testKey"  = "testValue"
  }
  description = "terraform-test"
}
data "baiducloud_bbc_instances" "default" {
  internal_ip = baiducloud_bbc_instance.bbc_instance.internal_ip
}
`, name)
}
