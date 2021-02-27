package baiducloud

import (
	"fmt"
	"github.com/baidubce/bce-sdk-go/services/iam"
	"github.com/baidubce/bce-sdk-go/services/iam/api"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"log"
	"strings"
	"testing"
)

const (
	testAccIamGroupResourceType = "baiducloud_iam_group"
	testAccIamGroupPrefix       = "test_BaiduAcc"
	testAccIamGroupResourceName = testAccIamGroupResourceType + "." + BaiduCloudTestResourceName
	testAccIamGroupDescription  = "test_description"
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
		if !strings.HasPrefix(group.Name, testAccIamGroupPrefix) {
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
	name := strings.ReplaceAll(acctest.RandomWithPrefix(testAccIamGroupPrefix), "-", "_")
	nameUpdate := strings.ReplaceAll(acctest.RandomWithPrefix(testAccIamGroupPrefix), "-", "_")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccIamGroupDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccIamGroupConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccIamGroupResourceName),
					resource.TestCheckResourceAttr(testAccIamGroupResourceName, "name", name),
					resource.TestCheckResourceAttr(testAccIamGroupResourceName, "description", testAccIamGroupDescription),
					resource.TestCheckResourceAttrSet(testAccIamGroupResourceName, "unique_id"),
				),
			},
			{
				ResourceName:            testAccIamGroupResourceName,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccIamGroupConfig(nameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccIamGroupResourceName),
					resource.TestCheckResourceAttr(testAccIamGroupResourceName, "name", nameUpdate),
					resource.TestCheckResourceAttr(testAccIamGroupResourceName, "description", testAccIamGroupDescription),
					resource.TestCheckResourceAttrSet(testAccIamGroupResourceName, "unique_id"),
				),
			},
		},
	})
}

func testAccIamGroupConfig(name string) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
  name = "%s"
  description = "%s"
  force_destroy    = true
}`, testAccIamGroupResourceType, BaiduCloudTestResourceName, name, testAccIamGroupDescription)
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
