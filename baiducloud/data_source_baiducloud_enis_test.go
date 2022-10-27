package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccEnisDataSourceName          = "data.baiducloud_enis.default"
	testAccEnisDataSourceAttrKeyPrefix = "enis.0."
)

//lintignore:AT003
func TestAccBaiduCloudEnisDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEnisConfig(BaiduCloudTestResourceTypeNameEni),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccEnisDataSourceName),
					resource.TestCheckResourceAttr(testAccEnisDataSourceName, testAccEnisDataSourceAttrKeyPrefix+"name", "tf-test-acc-eni"),
				),
			},
		},
	})
}

func testAccEnisConfig(name string) string {
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
resource "baiducloud_security_group" "sg" {
  name        = "terraform-sg"
  description = "security group created by terraform"
  vpc_id      = baiducloud_vpc.vpc.id
}
resource "baiducloud_eni" "default" {
  name      = "%s"
  subnet_id = baiducloud_subnet.subnet.id

  description        = "terraform test"
  security_group_ids = [
    baiducloud_security_group.sg.id
  ]
  private_ip {
    primary            = true
    private_ip_address = "172.16.0.13"
  }
}
data "baiducloud_enis" "default" {
  vpc_id = baiducloud_vpc.vpc.id
  filter{
    name = "eni_id"
    values = [baiducloud_eni.default.id]
  }
}
`, name)
}
