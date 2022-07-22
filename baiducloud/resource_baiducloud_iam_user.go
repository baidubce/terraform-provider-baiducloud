/*
Provide a resource to manage an IAM user.

Example Usage

```hcl
resource "baiducloud_iam_user" "my-user" {
  name = "my_user_name"
  description = "user description"
  force_destroy    = true
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/iam"
	"github.com/baidubce/bce-sdk-go/services/iam/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"time"
)

func resourceBaiduCloudIamUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudIamUserCreate,
		Read:   resourceBaiduCloudIamUserRead,
		Update: resourceBaiduCloudIamUserUpdate,
		Delete: resourceBaiduCloudIamUserDelete,

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
				Description: "Unique ID of user.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of user.",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the user.",
				Optional:    true,
			},
			"force_destroy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Delete user and its related access keys, group memberships and policy attachments.",
			},
		},
	}
}

func resourceBaiduCloudIamUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	action := "Create User " + name

	user, err := client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
		return iamClient.CreateUser(&api.CreateUserArgs{
			Name:        name,
			Description: description,
		})
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_iam_user", action, BCESDKGoERROR)
	}
	addDebug(action, user)

	d.SetId(name)
	return resourceBaiduCloudIamUserRead(d, meta)
}

func resourceBaiduCloudIamUserUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("name") || d.HasChange("description") {
		client := meta.(*connectivity.BaiduClient)
		on, nn := d.GetChange("name")
		name := on.(string)
		newName := nn.(string)
		description := d.Get("description").(string)
		action := "Update User " + name

		user, err := client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
			return iamClient.UpdateUser(name, &api.UpdateUserArgs{
				Name:        newName,
				Description: description,
			})
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_iam_user", action, BCESDKGoERROR)
		}
		addDebug(action, user)

		d.SetId(newName)
		return resourceBaiduCloudIamUserRead(d, meta)
	}
	return nil
}

func resourceBaiduCloudIamUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	name := d.Id()
	action := "Query User " + name

	raw, err := client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
		return iamClient.GetUser(name)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_iam_user", action, BCESDKGoERROR)
	}
	addDebug(action, raw)

	user, _ := raw.(*api.GetUserResult)
	d.Set("unique_id", user.Id)
	d.Set("name", user.Name)
	d.Set("description", user.Description)
	return nil
}

func resourceBaiduCloudIamUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	iamService := IamService{client}

	name := d.Id()
	action := "Delete User " + name

	if d.Get("force_destroy").(bool) {
		if err := iamService.ClearUserAttachedPolicy(name); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_iam_user", action, BCESDKGoERROR)
		}
		if err := iamService.ClearUserGroupMembership(name); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_iam_user", action, BCESDKGoERROR)
		}
	}
	_, err := client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
		return nil, iamClient.DeleteUser(name)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_iam_user", action, BCESDKGoERROR)
	}
	addDebug(action, name)
	return nil
}
