package baiducloud

import (
	"fmt"
	"github.com/baidubce/bce-sdk-go/services/iam"
	"github.com/baidubce/bce-sdk-go/services/iam/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"log"
	"strings"
	"testing"
)

const (
	testAccIamGroupResourceType = "baiducloud_iam_group"
	testAccIamGroupResourceName = testAccIamGroupResourceType + "." + BaiduCloudTestResourceName
)

func init() {
	resource.AddTestSweepers(testAccIamGroupResourceType, &resource.Sweeper{
		Name: testAccIamGroupResourceType,
		F:    testSweepIamGroups,
	})
}

func testSweepIamGroups(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("get BaiduCloud client error: %s", err)
	}

	client := rawClient.(*connectivity.BaiduClient)

	raw, err := client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
		return iamClient.ListGroup()
	})
	if err != nil {
		return fmt.Errorf("list groups error: %v", err)
	}

	result, _ := raw.(*api.ListGroupResult)
	for _, group := range result.Groups {
		if !strings.HasPrefix(group.Name, BaiduCloudTestResourceTypeNameUnderLine) {
			continue
		}
		log.Printf("[INFO] Deleting group: %s", group.Name)
		_, err := client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
			return nil, iamClient.DeleteGroup(group.Name)
		})
		if err != nil {
			return fmt.Errorf("delete group error: %v", err)
		}
	}
	return nil
}
func TestAccBaiduCloudIamGroup(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccIamGroupDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccIamGroupConfig(BaiduCloudTestResourceTypeNameUnderLine),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccIamGroupResourceName),
					resource.TestCheckResourceAttr(testAccIamGroupResourceName, "name", BaiduCloudTestResourceTypeNameUnderLine),
					resource.TestCheckResourceAttr(testAccIamGroupResourceName, "description", "created by terraform"),
					resource.TestCheckResourceAttrSet(testAccIamGroupResourceName, "unique_id"),
				),
			},
			{
				ResourceName:            testAccIamGroupResourceName,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccIamGroupConfig(BaiduCloudTestResourceTypeNameUnderLine),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccIamGroupResourceName),
					resource.TestCheckResourceAttr(testAccIamGroupResourceName, "name", BaiduCloudTestResourceTypeNameUnderLine),
					resource.TestCheckResourceAttr(testAccIamGroupResourceName, "description", "created by terraform"),
					resource.TestCheckResourceAttrSet(testAccIamGroupResourceName, "unique_id"),
				),
			},
		},
	})
}

func testAccIamGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccIamGroupResourceType {
			continue
		}

		_, err := client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
			return iamClient.GetGroup(rs.Primary.ID)
		})
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		return WrapError(Error("Iam Group still exist"))
	}

	return nil
}

func testAccIamGroupConfig(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_iam_group" "default" {
  name = "%s"
  description = "created by terraform"
  force_destroy    = true
}`, name)
}
