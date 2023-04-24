package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccInstancesDataSourceName          = "data.baiducloud_instances.default"
	testAccInstancesDataSourceAttrKeyPrefix = "instances.0."
)

//lintignore:AT003
func TestAccBaiduCloudInstancesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstancesDataSourceConfig(BaiduCloudTestResourceTypeNameInstance),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccInstancesDataSourceName),
					resource.TestCheckResourceAttr(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"name", BaiduCloudTestResourceTypeNameInstance),
					resource.TestCheckResourceAttr(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"tags.%", "1"),
					resource.TestCheckResourceAttr(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"tags.testKey", "testValue"),
					resource.TestCheckResourceAttrSet(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"image_id"),
					resource.TestCheckResourceAttrSet(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"zone_name"),
					resource.TestCheckResourceAttrSet(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"cpu_count"),
					resource.TestCheckResourceAttrSet(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"memory_capacity_in_gb"),
					resource.TestCheckResourceAttrSet(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"payment_timing"),
					resource.TestCheckResourceAttrSet(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"subnet_id"),
					resource.TestCheckResourceAttrSet(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"status"),
					resource.TestCheckResourceAttrSet(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"create_time"),
					resource.TestCheckResourceAttrSet(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"internal_ip"),
					resource.TestCheckResourceAttrSet(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"placement_policy"),
					resource.TestCheckResourceAttrSet(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"vpc_id"),
					resource.TestCheckResourceAttrSet(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"root_disk_size_in_gb"),
					resource.TestCheckResourceAttrSet(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"root_disk_storage_type"),
					resource.TestCheckResourceAttrSet(testAccInstancesDataSourceName, testAccInstancesDataSourceAttrKeyPrefix+"auto_renew"),
				),
			},
		},
	})
}

func testAccInstancesDataSourceConfig(name string) string {
	return fmt.Sprintf(`
data "baiducloud_images" "default" {}

resource "baiducloud_vpc" "test" {
  name        = "vpc_terraform_test"
  description = "created by terraform for test"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "test" {
  name        = "subnet_terraform_test"
  zone_name   = "cn-bj-d"
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.test.id
  description = "created by terraform for test"
}

resource "baiducloud_security_group" "test" {
  name        = "security_group_terraform_test"
  description = "created by terraform for test"
  vpc_id      = baiducloud_vpc.test.id
}

resource "baiducloud_deployset" "test" {
  name     = "deployset_terraform_test"
  desc     = "created by terraform for test"
  strategy = "HOST_HA"
}

resource "baiducloud_instance" "test" {
  image_id              = data.baiducloud_images.default.images.0.id
  name                  = "%s"
  availability_zone     = "cn-bj-d"
  cpu_count             = 1
  memory_capacity_in_gb = 1
  instance_type = "N5"
  payment_timing = "Postpaid"

  tags = {
    "testKey" = "testValue"
  }

  subnet_id = baiducloud_subnet.test.id
  security_groups = [baiducloud_security_group.test.id]
  deploy_set_ids = [baiducloud_deployset.test.id]
}

data "baiducloud_instances" "default" {
  internal_ip = baiducloud_instance.test.internal_ip
  zone_name   = baiducloud_instance.test.availability_zone
  instance_ids = baiducloud_instance.test.id 
  instance_names = "%s" 
  deploy_set_ids = baiducloud_deployset.test.id
  security_group_ids = baiducloud_security_group.test.id 
  payment_timing = "Postpaid" 
  tags = "testKey:testValue" 
  vpc_id = baiducloud_vpc.test.id
  private_ips = baiducloud_instance.test.internal_ip
  status = "Running"

  filter {
    name = "name"
    values = ["tf-test-acc*"]
  }
}
`, name, name)
}
