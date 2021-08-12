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
	testAccIamUserResourceType = "baiducloud_iam_user"
	testAccIamUserResourceName = testAccIamUserResourceType + "." + BaiduCloudTestResourceName
)

func init() {
	resource.AddTestSweepers(testAccIamUserResourceType, &resource.Sweeper{
		Name: testAccIamUserResourceType,
		F:    testSweepIamUsers,
	})
}

func testSweepIamUsers(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("get BaiduCloud client error: %s", err)
	}

	client := rawClient.(*connectivity.BaiduClient)

	raw, err := client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
		return iamClient.ListUser()
	})
	if err != nil {
		return fmt.Errorf("list users error: %v", err)
	}

	result, _ := raw.(*api.ListUserResult)
	for _, user := range result.Users {
		if !strings.HasPrefix(user.Name, BaiduCloudTestResourceTypeNameUnderLine) {
			continue
		}
		log.Printf("[INFO] Deleting user: %s", user.Name)
		_, err := client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
			return nil, iamClient.DeleteUser(user.Name)
		})
		if err != nil {
			return fmt.Errorf("delete user error: %v", err)
		}
	}
	return nil
}

func TestAccBaiduCloudIamUser(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccIamUserDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccIamUserConfig(BaiduCloudTestResourceTypeNameUnderLine),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccIamUserResourceName),
					resource.TestCheckResourceAttr(testAccIamUserResourceName, "name", BaiduCloudTestResourceTypeNameUnderLine),
					resource.TestCheckResourceAttr(testAccIamUserResourceName, "description", "created by terraform"),
					resource.TestCheckResourceAttrSet(testAccIamUserResourceName, "unique_id"),
				),
			},
			{
				ResourceName:            testAccIamUserResourceName,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccIamUserConfig(BaiduCloudTestResourceTypeNameUnderLine + "_update"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccIamUserResourceName),
					resource.TestCheckResourceAttr(testAccIamUserResourceName, "name", BaiduCloudTestResourceTypeNameUnderLine+"_update"),
					resource.TestCheckResourceAttr(testAccIamUserResourceName, "description", "created by terraform"),
					resource.TestCheckResourceAttrSet(testAccIamUserResourceName, "unique_id"),
				),
			},
		},
	})
}

func testAccIamUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccIamUserResourceType {
			continue
		}

		_, err := client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
			return iamClient.GetUser(rs.Primary.ID)
		})
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		return WrapError(Error("Iam User still exist"))
	}

	return nil
}

func testAccIamUserConfig(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_iam_user" "default" {
  name = "%s"
  description = "created by terraform"
  force_destroy    = true
}`, name)
}
