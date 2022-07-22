package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccPeerConnAcceptorResourceType = "baiducloud_peer_conn_acceptor"
	testAccPeerConnAcceptorResourceName = testAccPeerConnAcceptorResourceType + "." + BaiduCloudTestResourceName
)

//lintignore:AT003
func TestAccBaiduCloudPeerConnAcceptor(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccPeerConnAcceptorConfig(BaiduCloudTestResourceTypeNamePeerConnAcceptor),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccPeerConnAcceptorResourceName),
					resource.TestCheckResourceAttr(testAccPeerConnAcceptorResourceName, "bandwidth_in_mbps", "20"),
					resource.TestCheckResourceAttr(testAccPeerConnAcceptorResourceName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(testAccPeerConnAcceptorResourceName, "local_if_name", "local-interface"),
					resource.TestCheckResourceAttrSet(testAccPeerConnResourceName, "local_if_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorResourceName, "local_vpc_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorResourceName, "peer_vpc_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorResourceName, "peer_region"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorResourceName, "peer_account_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnResourceName, "status"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorResourceName, "created_time"),
					resource.TestCheckResourceAttr(testAccPeerConnAcceptorResourceName, "dns_status", "close"),
				),
			},
			{
				ResourceName:            testAccPeerConnAcceptorResourceName,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"dns_sync"},
			},
			{
				Config: testAccPeerConnAcceptorConfigUpdate(BaiduCloudTestResourceTypeNamePeerConnAcceptor),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccPeerConnAcceptorResourceName),
					resource.TestCheckResourceAttr(testAccPeerConnAcceptorResourceName, "bandwidth_in_mbps", "20"),
					resource.TestCheckResourceAttr(testAccPeerConnAcceptorResourceName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(testAccPeerConnAcceptorResourceName, "local_if_name", "local-interface"),
					resource.TestCheckResourceAttrSet(testAccPeerConnResourceName, "local_if_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorResourceName, "local_vpc_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorResourceName, "peer_vpc_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorResourceName, "peer_region"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorResourceName, "peer_account_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnAcceptorResourceName, "created_time"),
					resource.TestCheckResourceAttr(testAccPeerConnAcceptorResourceName, "dns_status", "open"),
					resource.TestCheckResourceAttrSet(testAccPeerConnResourceName, "status"),
					resource.TestCheckResourceAttrSet(testAccPeerConnResourceName, "created_time"),
				),
			},
		},
	})
}

func testAccPeerConnAcceptorConfig(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_vpc" "local-vpc" {
  name     = "%s"
  cidr     = "172.17.0.0/16"
}

resource "baiducloud_vpc" "peer-vpc" {
  name     = "%s"
  cidr     = "172.18.0.0/16"
}

resource "baiducloud_peer_conn" "default" {
  bandwidth_in_mbps = 20
  local_vpc_id      = baiducloud_vpc.local-vpc.id
  peer_vpc_id       = baiducloud_vpc.peer-vpc.id
  peer_region       = "bj"
  description       = "created by terraform"
  local_if_name     = "local-interface"
  billing = {
    payment_timing = "Postpaid"
  }
}

resource "baiducloud_peer_conn_acceptor" "default" {
  peer_conn_id = baiducloud_peer_conn.default.id
  auto_accept  = true
  dns_sync     = false
}
`, name+"-local", name+"-peer")
}

func testAccPeerConnAcceptorConfigUpdate(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_vpc" "local-vpc" {
  name     = "%s"
  cidr     = "172.17.0.0/16"
}

resource "baiducloud_vpc" "peer-vpc" {
  name     = "%s"
  cidr     = "172.18.0.0/16"
}

resource "baiducloud_peer_conn" "default" {
  bandwidth_in_mbps = 20
  local_vpc_id      = baiducloud_vpc.local-vpc.id
  peer_vpc_id       = baiducloud_vpc.peer-vpc.id
  peer_region       = "bj"
  description       = "created by terraform"
  local_if_name     = "local-interface"
  billing = {
    payment_timing = "Postpaid"
  }
}

resource "baiducloud_peer_conn_acceptor" "default" {
  peer_conn_id = baiducloud_peer_conn.default.id
  auto_accept  = true
  dns_sync     = true
}
`, name+"-local", name+"-peer")
}
