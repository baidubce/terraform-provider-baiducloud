/*
Provide a resource to manage an IAM group.

Example Usage

```hcl
resource "baiducloud_iam_group" "my-group" {
  name = "my_group_name"
  description = "group description"
  force_destroy    = true
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/iam"
	"github.com/baidubce/bce-sdk-go/services/iam/api"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"time"
)

func resourceBaiduCloudIamGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudIamGroupCreate,
		Read:   resourceBaiduCloudIamGroupRead,
		Update: resourceBaiduCloudIamGroupUpdate,
		Delete: resourceBaiduCloudIamGroupDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"unique_id": {
				Type:        schema.TypeString,
				Description: "Unique ID of group.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of group.",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the group.",
				Optional:    true,
			},
			"force_destroy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Delete group and its related user memberships and policy attachments.",
			},
		},
	}
}

func resourceBaiduCloudIamGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	action := "Create Group " + name

	group, err := client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
		return iamClient.CreateGroup(&api.CreateGroupArgs{
			Name:        name,
			Description: description,
		})
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_iam_group", action, BCESDKGoERROR)
	}
	addDebug(action, group)

	d.SetId(name)
	return resourceBaiduCloudIamGroupRead(d, meta)
}

func resourceBaiduCloudIamGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("name") || d.HasChange("description") {
		client := meta.(*connectivity.BaiduClient)
		on, nn := d.GetChange("name")
		name := on.(string)
		newName := nn.(string)
		description := d.Get("description").(string)
		action := "Update Group " + name

		group, err := client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
			return iamClient.UpdateGroup(name, &api.UpdateGroupArgs{
				Name:        newName,
				Description: description,
			})
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_iam_group", action, BCESDKGoERROR)
		}
		addDebug(action, group)

		d.SetId(newName)
		return resourceBaiduCloudIamGroupRead(d, meta)
	}
	return nil
}

func resourceBaiduCloudIamGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	name := d.Id()
	action := "Query Group " + name

	raw, err := client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
		return iamClient.GetGroup(name)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_iam_group", action, BCESDKGoERROR)
	}
	addDebug(action, raw)

	group, _ := raw.(*api.GetGroupResult)
	d.Set("unique_id", group.Id)
	d.Set("name", group.Name)
	d.Set("description", group.Description)
	return nil
}

func resourceBaiduCloudIamGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	iamService := IamService{client}

	name := d.Id()
	action := "Delete Group " + name

	if d.Get("force_destroy").(bool) {
		if err := iamService.ClearGroupAttachedPolicy(name); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_iam_group", action, BCESDKGoERROR)
		}
		if err := iamService.ClearUserFromGroup(name); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_iam_group", action, BCESDKGoERROR)
		}
	}
	_, err := client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
		return nil, iamClient.DeleteGroup(name)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_iam_group", action, BCESDKGoERROR)
	}
	addDebug(action, name)
	return nil
}
