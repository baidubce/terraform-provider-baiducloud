package baiducloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const (
	testAccCceClusterIpCidrDataSourceName = "data.baiducloud_ccev2_clusterip_cidr.default"
)

func TestAccBaiduCloudCCEv2ClusterIPCidrDataSource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCceClusterIpCidrDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccCceClusterIpCidrDataSourceName),
					resource.TestCheckResourceAttrSet(testAccCceClusterIpCidrDataSourceName, "is_success"),
					resource.TestCheckResourceAttrSet(testAccCceClusterIpCidrDataSourceName, "request_id"),
					resource.TestCheckResourceAttrSet(testAccCceClusterIpCidrDataSourceName, "recommended_clusterip_cidrs.#"),
				),
				Destroy: true,
			},
		},
	})
}

const testAccCceClusterIpCidrDataSourceConfig = `
variable "vpc_cidr" {
  default = "192.168.0.0/16"
}
variable "container_cidr" {
  default = "172.28.0.0/16"
}
data "baiducloud_ccev2_clusterip_cidr" "default" {
  vpc_cidr = var.vpc_cidr
  container_cidr = var.container_cidr
  cluster_max_service_num = 32
  private_net_cidrs = ["172.16.0.0/12",]
  ip_version = "ipv4"
}
`
