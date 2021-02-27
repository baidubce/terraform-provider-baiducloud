/*
Provide a resource to attach an IAM Policy to IAM Group.

Example Usage

```hcl
resource "baiducloud_iam_group" "my-group" {
  name = "my_group_name"
  force_destroy    = true
}
resource "baiducloud_iam_policy" "my-policy" {
   name = "my_policy"
  document = <<EOF
{"accessControlList": [{"region":"bj","service":"bcc","resource":["*"],"permission":["*"],"effect":"Allow"}]}
  EOF
}
resource "baiducloud_iam_group_policy_attachment" "my-group-policy-attachment" {
  group = "${baiducloud_iam_group.my-group.name}"
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

func resourceBaiduCloudIamGroupPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudIamGroupPolicyAttachmentCreate,
		Read:   resourceBaiduCloudIamGroupPolicyAttachmentRead,
		Delete: resourceBaiduCloudIamGroupPolicyAttachmentDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"group": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of group.",
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

func resourceBaiduCloudIamGroupPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	name := d.Get("group").(string)
	policy := d.Get("policy").(string)
	policyType := d.Get("policy_type").(string)
	action := "Create Group Policy Attachment for group " + name + " with " + policyType + " policy " + policy

	_, err := client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
		return nil, iamClient.AttachPolicyToGroup(&api.AttachPolicyToGroupArgs{
			GroupName:  name,
			PolicyName: policy,
			PolicyType: policyType,
		})
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_iam_group_policy_attachment",
			action, BCESDKGoERROR)
	}
	addDebug(action, nil)

	d.SetId(getGroupPolicyAttachmentResourceId(name, policy, policyType))
	return resourceBaiduCloudIamGroupPolicyAttachmentRead(d, meta)
}

func resourceBaiduCloudIamGroupPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	group, policy, policyType, err := parseGroupPolicyAttachmentResourceId(d.Id())
	if err != nil {
		return err
	}
	action := "List Group Policy Attachment for group " + group

	raw, err := client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
		return iamClient.ListGroupAttachedPolicies(group)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_iam_group_policy_attachment",
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
		log.Printf("[WARN] Unable to find Policy Attachment for group %s with policy %s", group, policy)
		d.SetId("")
	}
	return nil
}

func resourceBaiduCloudIamGroupPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	iamService := IamService{client}
	group, policy, policyType, err := parseGroupPolicyAttachmentResourceId(d.Id())
	if err != nil {
		return err
	}
	action := "Delete Group Policy Attachment for group " + group + " with policy " + policy

	err = iamService.DetachPolicyFromGroup(group, policy, policyType)
	if err != nil {
		return err
	}
	addDebug(action, nil)
	return nil
}

func parseGroupPolicyAttachmentResourceId(id string) (string, string, string, error) {
	parts := strings.Split(id, ":")
	if len(parts) != 4 {
		return "", "", "", WrapErrorf(nil, DefaultErrorMsg, "baiducloud_iam_group_policy_attachment",
			"parse attachment resource id", BCESDKGoERROR)
	}
	return parts[1], parts[3], parts[2], nil
}

func getGroupPolicyAttachmentResourceId(group string, policy string, policyType string) string {
	return strings.Join([]string{"group", group, policyType, policy}, ":")
}
