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
	testAccIamGroupPolicyAttachmentResourceType = "baiducloud_iam_group_policy_attachment"
	testAccIamGroupPolicyAttachmentResourceName = testAccIamGroupPolicyAttachmentResourceType + "." + BaiduCloudTestResourceName
)

func init() {
	resource.AddTestSweepers(testAccIamGroupPolicyAttachmentResourceType, &resource.Sweeper{
		Name: testAccIamGroupPolicyAttachmentResourceType,
		F:    testSweepIamGroupPolicyAttachments,
	})
}

func testSweepIamGroupPolicyAttachments(region string) error {
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
		if err := iamService.ClearGroupAttachedPolicy(group.Name); err != nil {
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
		return iamClient.ListPolicy("", api.POLICY_TYPE_CUSTOM)
	})
	if err != nil {
		return fmt.Errorf("list policy error: %v", err)
	}

	policies, _ := raw.(*api.ListPolicyResult)
	for _, policy := range policies.Policies {
		if !strings.HasPrefix(policy.Name, testAccIamPolicyPrefix) {
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

func TestAccBaiduCloudIamGroupPolicyAttachment(t *testing.T) {
	groupName := strings.ReplaceAll(acctest.RandomWithPrefix(testAccIamGroupPrefix), "-", "_")
	policyName := strings.ReplaceAll(acctest.RandomWithPrefix(testAccIamPolicyPrefix), "-", "_")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccIamGroupPolicyAttachmentDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccIamGroupPolicyAttachmentConfig(groupName, policyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccIamGroupPolicyAttachmentResourceName),
					resource.TestCheckResourceAttr(testAccIamGroupPolicyAttachmentResourceName, "group", groupName),
					resource.TestCheckResourceAttr(testAccIamGroupPolicyAttachmentResourceName, "policy", policyName),
					resource.TestCheckResourceAttr(testAccIamGroupPolicyAttachmentResourceName, "policy_type", api.POLICY_TYPE_CUSTOM),
				),
			},
		},
	})
}

func testAccIamGroupPolicyAttachmentConfig(groupName, policyName string) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
  name = "%s"
  force_destroy = true
}
resource "%s" "%s" {
  name = "%s"
  document = <<EOF
  {"accessControlList": [{"region":"bj","service":"bcc","resource":["*"],"permission":["*"],"effect":"Allow"}]}
  EOF
}
resource "%s" "%s" {
  group = "${%s}"
  policy = "${%s}"
}
`,
		testAccIamGroupResourceType, BaiduCloudTestResourceName, groupName,
		testAccIamPolicyResourceType, BaiduCloudTestResourceName, policyName,
		testAccIamGroupPolicyAttachmentResourceType, BaiduCloudTestResourceName, testAccIamGroupResourceName+".name",
		testAccIamPolicyResourceName+".name")
}

func testAccIamGroupPolicyAttachmentDestroy(s *terraform.State) error {
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
