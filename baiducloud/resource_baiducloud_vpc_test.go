package baiducloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/baidubce/bce-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccVPCResourceType     = "baiducloud_vpc"
	testAccVPCResourceName     = testAccVPCResourceType + "." + BaiduCloudTestResourceName
	testAccVPCResourceAttrName = BaiduCloudTestResourceAttrNamePrefix + "VPC"
)

func init() {
	resource.AddTestSweepers(testAccVPCResourceType, &resource.Sweeper{
		Name: testAccVPCResourceType,
		F:    testSweepVPCs,
		Dependencies: []string{
			testAccInstanceResourceType,
			testAccAppBLBResourceType,
			testAccPeerConnResourceType,
		},
	})
}

func testSweepVPCs(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("get BaiduCloud client error: %s", err)
	}

	client := rawClient.(*connectivity.BaiduClient)
	vpcService := &VpcService{client}

	vpcList, err := vpcService.ListAllVpcs()
	if err != nil {
		return fmt.Errorf("get VPCs error: %v", err)
	}

	for _, v := range vpcList {
		if !strings.HasPrefix(v.Name, BaiduCloudTestResourceAttrNamePrefix) {
			log.Printf("[INFO] Skipping VPC: %s (%s)", v.VPCID, v.Name)
			continue
		}

		// if nat gateways exist, sweep them first
		args := &vpc.ListNatGatewayArgs{VpcId: v.VPCID}
		natList, err := vpcService.ListAllNatGateways(args)
		if err != nil {
			return fmt.Errorf("get NatGateways error: %v", err)
		}

		for _, nat := range natList {
			log.Printf("[INFO] Deleting Nat Gateway: %s (%s)", nat.Id, nat.Name)
			_, err = client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
				return nil, vpcClient.DeleteNatGateway(nat.Id, "")
			})
			if err != nil {
				if IsExceptedErrors(err, NatGatewayNotFound) {
					continue
				}
				log.Printf("[ERROR] Failed to delete Nat Gateway %s (%s)", nat.Id, nat.Name)
				return err
			}
		}

		log.Printf("[INFO] Deleting VPC: %s (%s)", v.VPCID, v.Name)
		_, err = client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return nil, vpcClient.DeleteVPC(v.VPCID, "")
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete VPC %s (%s)", v.VPCID, v.Name)
		}
	}

	return nil
}

func TestAccBaiduCloudVPC(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccVPCDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccVPCConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccVPCResourceName),
					resource.TestCheckResourceAttr(testAccVPCResourceName, "name", testAccVPCResourceAttrName),
					resource.TestCheckResourceAttr(testAccVPCResourceName, "description", "vpc create"),
					resource.TestCheckResourceAttr(testAccVPCResourceName, "cidr", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(testAccVPCResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttrSet(testAccVPCResourceName, "route_table_id"),
					resource.TestCheckResourceAttr(testAccVPCResourceName, "secondary_cidrs.#", "0"),
				),
			},
			{
				ResourceName:      testAccVPCResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccVPCConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccVPCResourceName),
					resource.TestCheckResourceAttr(testAccVPCResourceName, "name", testAccVPCResourceAttrName+"Update"),
					resource.TestCheckResourceAttr(testAccVPCResourceName, "description", "vpc update"),
					resource.TestCheckResourceAttr(testAccVPCResourceName, "cidr", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(testAccVPCResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttrSet(testAccVPCResourceName, "route_table_id"),
					resource.TestCheckResourceAttr(testAccVPCResourceName, "secondary_cidrs.#", "0"),
				),
			},
		},
	})
}

func testAccVPCDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	vpcService := &VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccVPCResourceType {
			continue
		}

		_, err := vpcService.GetVPCDetail(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		return WrapError(Error("VPC still exist"))
	}

	return nil
}

func testAccVPCConfig() string {
	return fmt.Sprintf(`
resource "%s" "%s" {
  name        = "%s"
  description = "vpc create"
  cidr        = "192.168.0.0/24"
  tags = {
	"tagKey" = "tagValue"
  }
}`, testAccVPCResourceType, BaiduCloudTestResourceName, testAccVPCResourceAttrName)
}

func testAccVPCConfigUpdate() string {
	return fmt.Sprintf(`
resource "%s" "%s" {
  name        = "%s"
  description = "vpc update"
  cidr        = "192.168.0.0/24"
  tags = {
	"tagKey" = "tagValue"
  }
}`, testAccVPCResourceType, BaiduCloudTestResourceName, testAccVPCResourceAttrName+"Update")
}
