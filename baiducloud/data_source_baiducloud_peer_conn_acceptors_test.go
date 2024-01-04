package baiducloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccPeerConnAcceptorsDataSourceName          = "data.baiducloud_peer_conn_acceptors.default"
	testAccPeerConnAcceptorsDataSourceAttrKeyPrefix = "peer_conn_acceptors.0."
)

//lintignore:AT003
func TestAccBaiduCloudPeerConnAcceptorsDataSource(t *testing.T) {
	region := os.Getenv("BAIDUCLOUD_REGION")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccPeerConnsDataSourceConfigForPeerconnAcceptors(BaiduCloudTestResourceTypeNamePeerConnAcceptor, region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccPeerConnAcceptorsDataSourceName),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"peer_conn_id"),
					resource.TestCheckResourceAttr(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"bandwidth_in_mbps", "20"),
					resource.TestCheckResourceAttr(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"description", "created by terraform"),
					resource.TestCheckResourceAttr(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"local_if_name", "local-interface"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"local_vpc_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"peer_vpc_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"peer_region"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"peer_account_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"created_time"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"dns_status"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"role"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"status"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"local_if_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"local_region"),
					resource.TestCheckResourceAttr(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"payment_timing", "Postpaid"),
				),
			},
			{
				Config: testAccPeerConnAcceptorsDataSourceConfigForAll(BaiduCloudTestResourceTypeNamePeerConnAcceptor, region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccPeerConnAcceptorsDataSourceName),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"peer_conn_id"),
					resource.TestCheckResourceAttr(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"bandwidth_in_mbps", "20"),
					resource.TestCheckResourceAttr(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"description", "created by terraform"),
					resource.TestCheckResourceAttr(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"local_if_name", "local-interface"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"local_vpc_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"peer_vpc_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"peer_region"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"peer_account_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"created_time"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"dns_status"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"role"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"status"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"local_if_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"local_region"),
					resource.TestCheckResourceAttr(testAccPeerConnAcceptorsDataSourceName, testAccPeerConnAcceptorsDataSourceAttrKeyPrefix+"payment_timing", "Postpaid"),
				),
			},
		},
	})
}

func testAccPeerConnsDataSourceConfigForPeerconnAcceptors(name, region string) string {
	return fmt.Sprintf(`
resource "baiducloud_vpc" "local-vpc" {
  name = "%s"
  cidr = "172.17.0.0/16"
}

resource "baiducloud_vpc" "peer-vpc" {
  name = "%s"
  cidr = "172.18.0.0/16"
}

resource "baiducloud_peer_conn_acceptors" "default" {
  bandwidth_in_mbps = 20
  local_vpc_id      = baiducloud_vpc.local-vpc.id
  peer_vpc_id       = baiducloud_vpc.peer-vpc.id
  peer_region       = "%s"
  description       = "created by terraform"
  local_if_name     = "local-interface"
  billing = {
    payment_timing = "Postpaid"
  }
}

data "baiducloud_peer_conn_acceptors" "default" {
  peer_conn_id = baiducloud_peer_conn.default.id

  filter {
    name = "bandwidth_in_mbps"
    values = ["20"]
  }
}
`, name+"-local", name+"-peer", region)
}

func testAccPeerConnAcceptorsDataSourceConfigForAll(name, region string) string {
	return fmt.Sprintf(`
resource "baiducloud_vpc" "local-vpc" {
  name = "%s"
  cidr = "172.17.0.0/16"
}

resource "baiducloud_vpc" "peer-vpc" {
  name = "%s"
  cidr = "172.18.0.0/16"
}

resource "baiducloud_peer_conn_acceptor" "default" {
  bandwidth_in_mbps = 20
  local_vpc_id      = baiducloud_vpc.local-vpc.id
  peer_vpc_id       = baiducloud_vpc.peer-vpc.id
  peer_region       = "%s"
  description       = "created by terraform"
  local_if_name     = "local-interface"
  billing = {
    payment_timing = "Postpaid"
  }
}

data "baiducloud_peer_conn_acceptors" "default" {
  vpc_id = baiducloud_vpc.local-vpc.id

  filter {
    name = "bandwidth_in_mbps"
    values = ["20"]
  }
}
`, name+"-local", name+"-peer", region)
}
