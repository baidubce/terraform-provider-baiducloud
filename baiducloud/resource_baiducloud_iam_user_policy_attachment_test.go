package baiducloud

import (
	"fmt"
	"github.com/baidubce/bce-sdk-go/services/iam"
	"github.com/baidubce/bce-sdk-go/services/iam/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"log"
	"strings"
	"testing"
)

const (
	testAccIamUserPolicyAttachmentResourceType = "baiducloud_iam_user_policy_attachment"
	testAccIamUserPolicyAttachmentResourceName = testAccIamUserPolicyAttachmentResourceType + "." + BaiduCloudTestResourceName
)

func init() {
	resource.AddTestSweepers(testAccIamUserPolicyAttachmentResourceType, &resource.Sweeper{
		Name: testAccIamUserPolicyAttachmentResourceType,
		F:    testSweepIamUserPolicyAttachments,
	})
}

func testSweepIamUserPolicyAttachments(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("get BaiduCloud client error: %s", err)
	}

	client := rawClient.(*connectivity.BaiduClient)
	iamService := IamService{client}

	raw, err := client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
		return iamClient.ListUser()
	})
	if err != nil {
		return fmt.Errorf("list groups error: %v", err)
	}

	result, _ := raw.(*api.ListUserResult)
	for _, user := range result.Users {
		if !strings.HasPrefix(user.Name, BaiduCloudTestResourceTypeNameUnderLine) {
			continue
		}
		log.Printf("[INFO] Deleting user: %s", user.Name)
		if err := iamService.ClearUserAttachedPolicy(user.Name); err != nil {
			return err
		}
		_, err := client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
			return nil, iamClient.DeleteUser(user.Name)
		})
		if err != nil {
			return fmt.Errorf("delete user error: %v", err)
		}
	}

	raw, err = client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
		return iamClient.ListPolicy("", api.POLICY_TYPE_CUSTOM)
	})
	if err != nil {
		return fmt.Errorf("list policy error: %v", err)
	}

	policies, _ := raw.(*api.ListPolicyResult)
	for _, policy := range policies.Policies {
		if !strings.HasPrefix(policy.Name, BaiduCloudTestResourceTypeNameUnderLine) {
			continue
		}
		_, err := client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
			return nil, iamClient.DeletePolicy(policy.Name)
		})
		if err != nil {
			return fmt.Errorf("delete policy error: %v", err)
		}
	}
	return nil
}

func TestAccBaiduCloudIamUserPolicyAttachment(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccIamUserPolicyAttachmentDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccIamUserPolicyAttachmentConfig(BaiduCloudTestResourceTypeNameUnderLine),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccIamUserPolicyAttachmentResourceName),
					resource.TestCheckResourceAttr(testAccIamUserPolicyAttachmentResourceName, "user", BaiduCloudTestResourceTypeNameUnderLine),
					resource.TestCheckResourceAttr(testAccIamUserPolicyAttachmentResourceName, "policy", BaiduCloudTestResourceTypeNameUnderLine),
					resource.TestCheckResourceAttr(testAccIamUserPolicyAttachmentResourceName, "policy_type", api.POLICY_TYPE_CUSTOM),
				),
			},
		},
	})
}

func testAccIamUserPolicyAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)

	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
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
		case testAccIamPolicyResourceType:
			_, err := client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
				return iamClient.GetPolicy(rs.Primary.ID, api.POLICY_TYPE_CUSTOM)
			})
			if err != nil {
				if NotFoundError(err) {
					continue
				}
				return WrapError(err)
			}
			return WrapError(Error("Iam Policy still exist"))
		default:
			continue
		}
	}
	return nil
}

func testAccIamUserPolicyAttachmentConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "baiducloud_iam_user" "default" {
  name = var.name
  force_destroy = true
}
resource "baiducloud_iam_policy" "default" {
  name = var.name
  document = <<EOF
  {"accessControlList": [{"region":"bj","service":"bcc","resource":["*"],"permission":["*"],"effect":"Allow"}]}
  EOF
}
resource "baiducloud_iam_user_policy_attachment" "default" {
  user = baiducloud_iam_user.default.name
  policy = baiducloud_iam_policy.default.name
}
`, name)
}
