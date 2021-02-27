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
	testAccIamUserResourceType = "baiducloud_iam_user"
	testAccIamUserPrefix       = "test_BaiduAcc"
	testAccIamUserResourceName = testAccIamUserResourceType + "." + BaiduCloudTestResourceName
	testAccIamUserDescription  = "test_description"
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
		if !strings.HasPrefix(user.Name, testAccIamUserPrefix) {
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
	name := acctest.RandomWithPrefix(testAccIamUserPrefix)
	nameUpdate := acctest.RandomWithPrefix(testAccIamUserPrefix)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccIamUserDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccIamUserConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccIamUserResourceName),
					resource.TestCheckResourceAttr(testAccIamUserResourceName, "name", name),
					resource.TestCheckResourceAttr(testAccIamUserResourceName, "description", testAccIamUserDescription),
					resource.TestCheckResourceAttrSet(testAccIamUserResourceName, "unique_id"),
				),
			},
			{
				ResourceName:            testAccIamUserResourceName,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccIamUserConfig(nameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccIamUserResourceName),
					resource.TestCheckResourceAttr(testAccIamUserResourceName, "name", nameUpdate),
					resource.TestCheckResourceAttr(testAccIamUserResourceName, "description", testAccIamUserDescription),
					resource.TestCheckResourceAttrSet(testAccIamUserResourceName, "unique_id"),
				),
			},
		},
	})
}

func testAccIamUserConfig(name string) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
  name = "%s"
  description = "%s"
  force_destroy    = true
}`, testAccIamUserResourceType, BaiduCloudTestResourceName, name, testAccIamUserDescription)
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
