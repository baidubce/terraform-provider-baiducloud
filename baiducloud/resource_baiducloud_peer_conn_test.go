package baiducloud

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/baidubce/bce-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccPeerConnResourceType = "baiducloud_peer_conn"
	testAccPeerConnResourceName = testAccPeerConnResourceType + "." + BaiduCloudTestResourceName
)

func init() {
	resource.AddTestSweepers(testAccPeerConnResourceType, &resource.Sweeper{
		Name: testAccPeerConnResourceType,
		F:    testSweepPeerConns,
	})
}

func testSweepPeerConns(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("get BaiduCloud client error: %s", err)
	}

	client := rawClient.(*connectivity.BaiduClient)
	vpcService := &VpcService{client}

	peerconnList, err := vpcService.ListAllPeerConns("")
	if err != nil {
		return fmt.Errorf("get peer connections error: %v", err)
	}

	for _, peerconn := range peerconnList {
		if !strings.HasPrefix(peerconn.Description, BaiduCloudTestResourceAttrNamePrefix) {
			log.Printf("[INFO] Skipping PeerConn: %s", peerconn.PeerConnId)
			continue
		}

		log.Printf("[INFO] Deleting PeerConn: %s", peerconn.PeerConnId)
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return nil, vpcClient.DeletePeerConn(peerconn.PeerConnId, buildClientToken())
		})
		if err != nil {
			if IsExceptedErrors(err, PeerConnNotFound) {
				continue
			}
			log.Printf("[ERROR] Failed to delete PeerConn %s", peerconn.PeerConnId)
			return err
		}
	}

	return nil
}

//lintignore:AT003
func TestAccBaiduCloudPeerConn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccPeerConnDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccPeerConnConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccPeerConnResourceName),
					resource.TestCheckResourceAttr(testAccPeerConnResourceName, "bandwidth_in_mbps", "20"),
					resource.TestCheckResourceAttr(testAccPeerConnResourceName, "description", BaiduCloudTestResourceAttrNamePrefix+"PeerConn"),
					resource.TestCheckResourceAttr(testAccPeerConnResourceName, "local_if_name", "local-interface"),
					resource.TestCheckResourceAttrSet(testAccPeerConnResourceName, "local_if_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnResourceName, "local_vpc_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnResourceName, "peer_vpc_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnResourceName, "peer_region"),
					resource.TestCheckResourceAttrSet(testAccPeerConnResourceName, "peer_account_id"),
					resource.TestCheckResourceAttr(testAccPeerConnResourceName, "peer_if_name", "peer-interface"),
					resource.TestCheckResourceAttrSet(testAccPeerConnResourceName, "created_time"),
					resource.TestCheckResourceAttr(testAccPeerConnResourceName, "dns_sync", "true"),
					resource.TestCheckResourceAttr(testAccPeerConnResourceName, "dns_status", "open"),
					resource.TestCheckResourceAttr(testAccPeerConnResourceName, "role", string(vpc.PEERCONN_ROLE_INITIATOR)),
					resource.TestCheckResourceAttrSet(testAccPeerConnResourceName, "status"),
					resource.TestCheckResourceAttrSet(testAccPeerConnResourceName, "created_time"),
				),
			},
			{
				ResourceName:            testAccPeerConnResourceName,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"dns_sync"},
			},
			{
				Config: testAccPeerConnConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccPeerConnResourceName),
					resource.TestCheckResourceAttr(testAccPeerConnResourceName, "bandwidth_in_mbps", "30"),
					resource.TestCheckResourceAttr(testAccPeerConnResourceName, "description", BaiduCloudTestResourceAttrNamePrefix+"PeerConnUpdate"),
					resource.TestCheckResourceAttr(testAccPeerConnResourceName, "local_if_name", "local-interface-update"),
					resource.TestCheckResourceAttrSet(testAccPeerConnResourceName, "local_if_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnResourceName, "local_vpc_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnResourceName, "peer_vpc_id"),
					resource.TestCheckResourceAttrSet(testAccPeerConnResourceName, "peer_region"),
					resource.TestCheckResourceAttrSet(testAccPeerConnResourceName, "peer_account_id"),
					resource.TestCheckResourceAttr(testAccPeerConnResourceName, "peer_if_name", "peer-interface"),
					resource.TestCheckResourceAttrSet(testAccPeerConnResourceName, "created_time"),
					resource.TestCheckResourceAttr(testAccPeerConnResourceName, "dns_sync", "false"),
					resource.TestCheckResourceAttr(testAccPeerConnResourceName, "dns_status", "close"),
					resource.TestCheckResourceAttr(testAccPeerConnResourceName, "role", string(vpc.PEERCONN_ROLE_INITIATOR)),
					resource.TestCheckResourceAttrSet(testAccPeerConnResourceName, "status"),
					resource.TestCheckResourceAttrSet(testAccPeerConnResourceName, "created_time"),
				),
			},
		},
	})
}

func testAccPeerConnDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	vpcService := &VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccPeerConnResourceType {
			continue
		}

		_, err := vpcService.GetPeerConnDetail(rs.Primary.ID, vpc.PEERCONN_ROLE_INITIATOR)
		if err != nil {
			if IsExceptedErrors(err, PeerConnNotFound) {
				continue
			}
			return WrapError(err)
		}
		return WrapError(Error("PeerConn still exist"))
	}

	return nil
}

func testAccPeerConnConfig() string {
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
  peer_if_name      = "peer-interface"
  description       = "%s"
  local_if_name     = "local-interface"
  dns_sync = true
  billing = {
    payment_timing = "Postpaid"
  }
}
`, BaiduCloudTestResourceAttrNamePrefix+"VPC-local",
		BaiduCloudTestResourceAttrNamePrefix+"VPC-peer", region,
		BaiduCloudTestResourceAttrNamePrefix+"PeerConn")
}

func testAccPeerConnConfigUpdate() string {
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
  bandwidth_in_mbps = 30
  local_vpc_id      = baiducloud_vpc.local-vpc.id
  peer_vpc_id       = baiducloud_vpc.peer-vpc.id
  peer_region       = "%s"
  peer_if_name      = "peer-interface"
  description       = "%s"
  local_if_name     = "local-interface-update"
  dns_sync          = false
  billing = {
    payment_timing = "Postpaid"
  }
}
`, BaiduCloudTestResourceAttrNamePrefix+"VPC-local",
		BaiduCloudTestResourceAttrNamePrefix+"VPC-peer", region,
		BaiduCloudTestResourceAttrNamePrefix+"PeerConnUpdate")
}
