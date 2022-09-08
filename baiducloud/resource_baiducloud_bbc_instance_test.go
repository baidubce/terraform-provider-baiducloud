package baiducloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"testing"
)

const (
	testAccBbcInstanceResourceType = "baiducloud_bbc_instance"
	testAccBbcInstanceResourceName = testAccBbcInstanceResourceType + "." + BaiduCloudTestResourceName
)

func TestAccBaiduCloudBbcInstance(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccBbcInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBbcInstanceConfig(BaiduCloudTestResourceTypeNameBbcInstance),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBbcInstanceResourceName),
					resource.TestCheckResourceAttr(testAccBbcInstanceResourceName, "name", BaiduCloudTestResourceTypeNameBbcInstance),
					resource.TestCheckResourceAttr(testAccBbcInstanceResourceName, "description", "created by terraform"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "image_id"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "root_disk_size_in_gb"),
					resource.TestCheckResourceAttr(testAccBbcInstanceResourceName, "billing.payment_timing", "Postpaid"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "subnet_id"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "status"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "create_time"),
					resource.TestCheckResourceAttrSet(testAccBbcInstanceResourceName, "internal_ip"),
				),
			},
			{
				Config: testAccBbcInstanceConfigUpdate(BaiduCloudTestResourceTypeNameBbcInstance),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBbcInstanceResourceName),
					resource.TestCheckResourceAttr(testAccBbcInstanceResourceName, "description", "created by terraform-Update"),
				),
			},
		},
	})
}

func testAccBbcInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	bbcService := BbcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccBbcInstanceResourceType {
			continue
		}

		bbcInstance, _ := bbcService.GetBbcInstanceDetail(rs.Primary.ID)
		if bbcInstance != nil {
			return WrapError(Error("bbc instance still exist"))
		}
	}

	return nil
}
func testAccBbcInstanceConfig(name string) string {
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
resource "baiducloud_bbc_image" "test-image" {
  image_name = "terrform-bbc-image-test"
  instance_id = "i-qwIq4vKi"
}
data "baiducloud_bbc_instances" "data_bbc_instance" {
  internal_ip = "172.16.16.4"
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
    "业务"  = "terraform_test"
  }
  description = "created by terraform"
}
`, name)
}
func testAccBbcInstanceConfigUpdate(name string) string {
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
resource "baiducloud_bbc_image" "test-image" {
  image_name = "terrform-bbc-image-test"
  instance_id = "i-qwIq4vKi"
}
data "baiducloud_bbc_instances" "data_bbc_instance" {
  internal_ip = "172.16.16.4"
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
    "业务"  = "terraform_test"
  }
  description = "created by terraform-Update"
}
`, name)
}
