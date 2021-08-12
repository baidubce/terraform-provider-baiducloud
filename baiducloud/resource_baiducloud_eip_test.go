package baiducloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/baidubce/bce-sdk-go/services/eip"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccEipResourceType = "baiducloud_eip"
	testAccEipResourceName = testAccEipResourceType + "." + BaiduCloudTestResourceName
)

func init() {
	resource.AddTestSweepers(testAccEipResourceType, &resource.Sweeper{
		Name: testAccEipResourceType,
		F:    testSweepEips,
		Dependencies: []string{
			testAccEipAssociationResourceType,
			testAccInstanceResourceType,
			testAccAppBLBResourceType,
			testAccVPCResourceType,
		},
	})
}

func testSweepEips(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("get BaiduCloud client error: %s", err)
	}
	client := rawClient.(*connectivity.BaiduClient)
	eipService := EipService{client}

	listArgs := &eip.ListEipArgs{}
	eipList, err := eipService.ListAllEips(listArgs)
	if err != nil {
		return fmt.Errorf("get EIPs error: %s", err)
	}

	for _, e := range eipList {
		name := e.Name
		ip := e.Eip
		if !strings.HasPrefix(e.Name, BaiduCloudTestResourceTypeName) {
			log.Printf("[INFO] Skipping EIP: %s (%s)", name, ip)
			continue
		}

		log.Printf("[INFO] Deleting EIP: %s (%s)", e.Name, e.Eip)
		_, err := client.WithEipClient(func(client *eip.Client) (i interface{}, e error) {
			return nil, client.DeleteEip(ip, buildClientToken())
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete EIP %s (%s)", name, ip)
		}
	}

	return nil
}

//lintignore:AT003
func TestAccBaiduCloudEip(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccEIPDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccEipConfig(BaiduCloudTestResourceTypeNameEip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccEipResourceName),
					resource.TestCheckResourceAttr(testAccEipResourceName, "name", BaiduCloudTestResourceTypeNameEip),
					resource.TestCheckResourceAttr(testAccEipResourceName, "bandwidth_in_mbps", "1"),
					resource.TestCheckResourceAttr(testAccEipResourceName, "billing_method", "ByTraffic"),
					resource.TestCheckResourceAttr(testAccEipResourceName, "payment_timing", "Postpaid"),
					resource.TestCheckResourceAttr(testAccEipResourceName, "tags.%", "1"),
					resource.TestCheckNoResourceAttr(testAccEipResourceName, "reservation_length"),
				),
			},
			{
				ResourceName:      testAccEipResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEipConfigUpdate(BaiduCloudTestResourceTypeNameEip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccEipResourceName),
					resource.TestCheckResourceAttr(testAccEipResourceName, "name", BaiduCloudTestResourceTypeNameEip),
					resource.TestCheckResourceAttr(testAccEipResourceName, "bandwidth_in_mbps", "2"),
					resource.TestCheckResourceAttr(testAccEipResourceName, "billing_method", "ByTraffic"),
					resource.TestCheckResourceAttr(testAccEipResourceName, "payment_timing", "Postpaid"),
					resource.TestCheckResourceAttr(testAccEipResourceName, "tags.%", "1"),
					resource.TestCheckNoResourceAttr(testAccEipResourceName, "reservation_length"),
				),
			},
		},
	})
}

func testAccEIPDestory(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	eipService := EipService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccEipResourceType {
			continue
		}

		_, err := eipService.EipGetDetail(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		return WrapError(Error("EIP still exist"))
	}

	return nil
}

func testAccEipConfig(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_eip" "default" {
  name               = "%s"
  bandwidth_in_mbps  = 1
  payment_timing     = "Postpaid"
  billing_method     = "ByTraffic"
  reservation_length = 1

  tags = {
    "testKey" = "testValue"
  }
}
`, name)
}

func testAccEipConfigUpdate(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_eip" "default" {
  name              = "%s"
  bandwidth_in_mbps = 2
  payment_timing    = "Postpaid"
  billing_method    = "ByTraffic"

  tags = {
    "testKey" = "testValue"
  }
}
`, name)
}
