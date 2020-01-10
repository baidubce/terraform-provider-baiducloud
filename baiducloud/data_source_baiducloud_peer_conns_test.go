package baiducloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const (
	testAccPeerConnsDataSourceName          = "data.baiducloud_peer_conns.default"
	testAccPeerConnsDataSourceAttrKeyPrefix = "peer_conns.0."
)

func TestAccBaiduCloudPeerConnsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccPeerConnsDataSourceConfigForPeerconn(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccPeerConnsDataSourceName),
					resource.TestCheckResourceAttrSet(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"peer_conn_id"),
					resource.TestCheckResourceAttr(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"bandwidth_in_mbps", "20"),
					resource.TestCheckResourceAttr(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"description", "test peer conn"),
					resource.TestCheckResourceAttr(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"local_if_name", "local-interface"),
					resource.TestCheckResourceAttrSet(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"local_vpc_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"peer_vpc_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"peer_region"),
					resource.TestCheckResourceAttrSet(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"peer_account_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"created_time"),
					resource.TestCheckResourceAttrSet(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"dns_status"),
					resource.TestCheckResourceAttrSet(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"role"),
					resource.TestCheckResourceAttrSet(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"status"),
					resource.TestCheckResourceAttrSet(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"local_if_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"local_region"),
					resource.TestCheckResourceAttr(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"payment_timing", "Postpaid"),
				),
			},
			{
				Config: testAccPeerConnsDataSourceConfigForAll(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccPeerConnsDataSourceName),
					resource.TestCheckResourceAttrSet(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"peer_conn_id"),
					resource.TestCheckResourceAttr(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"bandwidth_in_mbps", "20"),
					resource.TestCheckResourceAttr(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"description", "test peer conn"),
					resource.TestCheckResourceAttr(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"local_if_name", "local-interface"),
					resource.TestCheckResourceAttrSet(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"local_vpc_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"peer_vpc_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"peer_region"),
					resource.TestCheckResourceAttrSet(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"peer_account_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"created_time"),
					resource.TestCheckResourceAttrSet(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"dns_status"),
					resource.TestCheckResourceAttrSet(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"role"),
					resource.TestCheckResourceAttrSet(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"status"),
					resource.TestCheckResourceAttrSet(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"local_if_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"local_region"),
					resource.TestCheckResourceAttr(testAccPeerConnsDataSourceName, testAccPeerConnsDataSourceAttrKeyPrefix+"payment_timing", "Postpaid"),
				),
			},
		},
	})
}

func testAccPeerConnsDataSourceConfigForPeerconn() string {
	region := os.Getenv("BAIDUCLOUD_REGION")
	return fmt.Sprintf(`
resource "baiducloud_vpc" "local-vpc" {
  name = "%s"
  cidr = "172.17.0.0/16"
}

resource "baiducloud_vpc" "peer-vpc" {
  name = "%s"
  cidr = "172.18.0.0/16"
}

resource "baiducloud_peer_conn" "default" {
  bandwidth_in_mbps = 20
  local_vpc_id      = baiducloud_vpc.local-vpc.id
  peer_vpc_id       = baiducloud_vpc.peer-vpc.id
  peer_region       = "%s"
  description       = "test peer conn"
  local_if_name     = "local-interface"
  billing = {
    payment_timing = "Postpaid"
  }
}

data "baiducloud_peer_conns" "default" {
  peer_conn_id = baiducloud_peer_conn.default.id
}
`, BaiduCloudTestResourceAttrNamePrefix+"VPC", BaiduCloudTestResourceAttrNamePrefix+"VPC", region)
}

func testAccPeerConnsDataSourceConfigForAll() string {
	region := os.Getenv("BAIDUCLOUD_REGION")
	return fmt.Sprintf(`
resource "baiducloud_vpc" "local-vpc" {
  name = "%s"
  cidr = "172.17.0.0/16"
}

resource "baiducloud_vpc" "peer-vpc" {
  name = "%s"
  cidr = "172.18.0.0/16"
}

resource "baiducloud_peer_conn" "default" {
  bandwidth_in_mbps = 20
  local_vpc_id      = baiducloud_vpc.local-vpc.id
  peer_vpc_id       = baiducloud_vpc.peer-vpc.id
  peer_region       = "%s"
  description       = "test peer conn"
  local_if_name     = "local-interface"
  billing = {
    payment_timing = "Postpaid"
  }
}

data "baiducloud_peer_conns" "default" {
  vpc_id = baiducloud_vpc.local-vpc.id
}
`, BaiduCloudTestResourceAttrNamePrefix+"VPC", BaiduCloudTestResourceAttrNamePrefix+"VPC", region)
}
