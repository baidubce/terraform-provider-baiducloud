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
	testAccIamGroupMembershipResourceType = "baiducloud_iam_group_membership"
	testAccIamGroupMembershipResourceName = testAccIamGroupMembershipResourceType + "." + BaiduCloudTestResourceName
)

func init() {
	resource.AddTestSweepers(testAccIamGroupMembershipResourceType, &resource.Sweeper{
		Name: testAccIamGroupMembershipResourceType,
		F:    testSweepIamGroupMemberships,
	})
}

func testSweepIamGroupMemberships(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("get BaiduCloud client error: %s", err)
	}

	client := rawClient.(*connectivity.BaiduClient)
	iamService := IamService{client}

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
		if err := iamService.ClearUserFromGroup(group.Name); err != nil {
			return err
		}
		_, err := client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
			return nil, iamClient.DeleteGroup(group.Name)
		})
		if err != nil {
			return fmt.Errorf("delete group error: %v", err)
		}
	}

	raw, err = client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
		return iamClient.ListUser()
	})
	if err != nil {
		return fmt.Errorf("list user error: %v", err)
	}

	users, _ := raw.(*api.ListUserResult)
	for _, user := range users.Users {
		if !strings.HasPrefix(user.Name, testAccIamUserPrefix) {
			continue
		}
		if err := iamService.ClearUserGroupMembership(user.Name); err != nil {
			return err
		}
		_, err := client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
			return nil, iamClient.DeleteUser(user.Name)
		})
		if err != nil {
			return fmt.Errorf("delete user error: %v", err)
		}
	}
	return nil
}

func TestAccBaiduCloudIamGroupMembership(t *testing.T) {
	groupName := strings.ReplaceAll(acctest.RandomWithPrefix(testAccIamGroupPrefix), "-", "_")
	userName1 := acctest.RandomWithPrefix(testAccIamUserPrefix) + "_1"
	userName2 := acctest.RandomWithPrefix(testAccIamUserPrefix) + "_2"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: testAccIamGroupMembershipResourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccIamGroupMembershipDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccIamGroupMembership1UserConfig(groupName, userName1, userName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccIamGroupMembershipResourceName),
					resource.TestCheckResourceAttr(testAccIamGroupMembershipResourceName, "group", groupName),
					resource.TestCheckResourceAttr(testAccIamGroupMembershipResourceName, "users.#", "1"),
				),
			},
			{
				Config: testAccIamGroupMembership2UserConfig(groupName, userName1, userName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccIamGroupMembershipResourceName),
					resource.TestCheckResourceAttr(testAccIamGroupMembershipResourceName, "group", groupName),
					resource.TestCheckResourceAttr(testAccIamGroupMembershipResourceName, "users.#", "2"),
				),
			},
			{
				Config: testAccIamGroupMembership1UserConfig(groupName, userName2, userName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccIamGroupMembershipResourceName),
					resource.TestCheckResourceAttr(testAccIamGroupMembershipResourceName, "group", groupName),
					resource.TestCheckResourceAttr(testAccIamGroupMembershipResourceName, "users.#", "1"),
				),
			},
		},
	})
}

func testAccIamGroupMembership1UserConfig(groupName, userName1, userName2 string) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
  name = "%s"
  force_destroy = true
}
resource "%s" "%s" {
  name = "%s"
  force_destroy = true
}
resource "%s" "%s" {
  name = "%s"
  force_destroy = true
}
resource "%s" "%s" {
  group = "${%s}"
  users = ["${%s}"]
}
`,
		testAccIamGroupResourceType, BaiduCloudTestResourceName, groupName,
		testAccIamUserResourceType, BaiduCloudTestResourceName+"_"+userName1, userName1,
		testAccIamUserResourceType, BaiduCloudTestResourceName+"_"+userName2, userName2,
		testAccIamGroupMembershipResourceType, BaiduCloudTestResourceName, testAccIamGroupResourceName+".name",
		getTestAccIamUserResourceNameAttr(userName1))
}

func testAccIamGroupMembership2UserConfig(groupName, userName1, userName2 string) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
  name = "%s"
  force_destroy = true
}
resource "%s" "%s" {
  name = "%s"
  force_destroy = true
}
resource "%s" "%s" {
  name = "%s"
  force_destroy = true
}
resource "%s" "%s" {
  group = "${%s}"
  users = ["${%s}", "${%s}"]
}
`,
		testAccIamGroupResourceType, BaiduCloudTestResourceName, groupName,
		testAccIamUserResourceType, BaiduCloudTestResourceName+"_"+userName1, userName1,
		testAccIamUserResourceType, BaiduCloudTestResourceName+"_"+userName2, userName2,
		testAccIamGroupMembershipResourceType, BaiduCloudTestResourceName, testAccIamGroupResourceName+".name",
		getTestAccIamUserResourceNameAttr(userName1), getTestAccIamUserResourceNameAttr(userName2))
}

func getTestAccIamUserResourceNameAttr(name string) string {
	return testAccIamUserResourceName + "_" + name + ".name"
}

func testAccIamGroupMembershipDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)

	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case testAccIamGroupResourceType:
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
		case testAccIamUserResourceType:
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
		default:
			continue
		}
	}
	return nil
}
