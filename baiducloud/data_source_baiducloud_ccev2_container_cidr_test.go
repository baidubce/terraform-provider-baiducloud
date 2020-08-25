package baiducloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const (
	testAccCceContainerCidrDataSourceName = "data.baiducloud_ccev2_container_cidr.default"
)

func TestAccBaiduCloudCCEv2ContainerCidrDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCceContainerCidrDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccCceContainerCidrDataSourceName),
					resource.TestCheckResourceAttrSet(testAccCceContainerCidrDataSourceName, "is_success"),
					resource.TestCheckResourceAttrSet(testAccCceContainerCidrDataSourceName, "request_id"),
					resource.TestCheckResourceAttrSet(testAccCceContainerCidrDataSourceName, "recommended_container_cidrs.#"),
				),
			},
		},
	})
}

const testAccCceContainerCidrDataSourceConfig = `
resource "baiducloud_vpc" "default" {
  name        = "test-vpc-tf-auto"
  description = "test-vpc-tf-auto"
  cidr        = "192.168.0.0/16"
}
data "baiducloud_ccev2_container_cidr" "default" {
  vpc_id = baiducloud_vpc.default.id
  vpc_cidr = baiducloud_vpc.default.cidr
  cluster_max_node_num = 16
  max_pods_per_node = 32
  private_net_cidrs = ["192.168.0.0/16",]
  k8s_version = "1.16.8"
  ip_version = "ipv4"
}
`
