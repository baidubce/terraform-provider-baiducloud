/*
Provide a resource to attach an IAM Policy to IAM User.

Example Usage

```hcl
resource "baiducloud_iam_user" "my-user" {
  name = "my_user_name"
  force_destroy    = true
}
resource "baiducloud_iam_policy" "my-policy" {
  name = "my_policy"
  document = <<EOF
{"accessControlList": [{"region":"bj","service":"bcc","resource":["*"],"permission":["*"],"effect":"Allow"}]}
  EOF
}
resource "baiducloud_iam_user_policy_attachment" "my-user-policy-attachment" {
  user = "${baiducloud_iam_user.my-user.name}"
  policy = "${baiducloud_iam_policy.my-policy.name}"
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/iam"
	"github.com/baidubce/bce-sdk-go/services/iam/api"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"log"
	"strings"
	"time"
)

func resourceBaiduCloudIamUserPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudIamUserPolicyAttachmentCreate,
		Read:   resourceBaiduCloudIamUserPolicyAttachmentRead,
		Delete: resourceBaiduCloudIamUserPolicyAttachmentDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"user": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of user.",
			},
			"policy": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of policy.",
			},
			"policy_type": {
				Type:        schema.TypeString,
				Default:     api.POLICY_TYPE_CUSTOM,
				Optional:    true,
				ForceNew:    true,
				Description: "Type of policy, valid values are Custom/System.",
			},
		},
	}
}

func resourceBaiduCloudIamUserPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	name := d.Get("user").(string)
	policy := d.Get("policy").(string)
	policyType := d.Get("policy_type").(string)
	action := "Create User Policy Attachment for user " + name + " with " + policyType + " policy " + policy

	_, err := client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
		return nil, iamClient.AttachPolicyToUser(&api.AttachPolicyToUserArgs{
			UserName:   name,
			PolicyName: policy,
			PolicyType: policyType,
		})
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_iam_user_policy_attachment",
			action, BCESDKGoERROR)
	}
	addDebug(action, nil)

	d.SetId(getUserPolicyAttachmentResourceId(name, policy, policyType))
	return resourceBaiduCloudIamUserPolicyAttachmentRead(d, meta)
}

func resourceBaiduCloudIamUserPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	user, policy, policyType, err := parseUserPolicyAttachmentResourceId(d.Id())
	if err != nil {
		return err
	}
	action := "List User Policy Attachment for user " + user

	raw, err := client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
		return iamClient.ListUserAttachedPolicies(user)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_iam_user_policy_attachment",
			action, BCESDKGoERROR)
	}
	policies, _ := raw.(*api.ListPolicyResult)
	addDebug(action, policies)

	var found bool
	for _, p := range policies.Policies {
		if p.Name == policy && p.Type == policyType {
			found = true
		}
	}
	if !found {
		log.Printf("[WARN] Unable to find Policy Attachment for user %s with policy %s", user, policy)
		d.SetId("")
	}
	return nil
}

func resourceBaiduCloudIamUserPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	iamService := IamService{client}
	user, policy, policyType, err := parseUserPolicyAttachmentResourceId(d.Id())
	if err != nil {
		return err
	}
	action := "Delete User Policy Attachment for user " + user + " with policy " + policy

	err = iamService.DetachPolicyFromUser(user, policy, policyType)
	if err != nil {
		return err
	}
	addDebug(action, nil)
	return nil
}

func parseUserPolicyAttachmentResourceId(id string) (string, string, string, error) {
	parts := strings.Split(id, ":")
	if len(parts) != 4 {
		return "", "", "", WrapErrorf(nil, DefaultErrorMsg, "baiducloud_iam_user_policy_attachment",
			"parse attachment resource id", BCESDKGoERROR)
	}
	return parts[1], parts[3], parts[2], nil
}

func getUserPolicyAttachmentResourceId(user string, policy string, policyType string) string {
	return strings.Join([]string{"user", user, policyType, policy}, ":")
}
