/*
Provide a resource to manage IAM Group membership for IAM Users.

Example Usage

```hcl
resource "baiducloud_iam_group" "my-group" {
  name = "my_group_name"
  force_destroy = true
}
resource "baiducloud_iam_user" "my-user" {
  name = "my_user_name"
  force_destroy = true
}
resource "baiducloud_iam_group_membership" "my-group-membership" {
  group = "${baiducloud_iam_group.my-group.name}"
  users = ["${baiducloud_iam_user.my-user.name}"]
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

func resourceBaiduCloudIamGroupMembership() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudIamGroupMembershipCreate,
		Read:   resourceBaiduCloudIamGroupMembershipRead,
		Update: resourceBaiduCloudIamGroupMembershipUpdate,
		Delete: resourceBaiduCloudIamGroupMembershipDelete,

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
			"users": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "Names of users to add into group.",
			},
		},
	}
}

func resourceBaiduCloudIamGroupMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	name := d.Get("group").(string)
	userList := expandStringSet(d.Get("users").(*schema.Set))
	action := "Create Group Membership for group " + name

	if err := addUsersToGroup(client, userList, name); err != nil {
		return err
	}
	addDebug(action, name)

	d.SetId(name)
	return resourceBaiduCloudIamGroupMembershipRead(d, meta)
}

func resourceBaiduCloudIamGroupMembershipRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	name := d.Id()
	action := "Read Group Membership for group " + name
	raw, err := client.WithIamClient(func(iamClient *iam.Client) (i interface{}, e error) {
		return iamClient.ListUsersInGroup(name)
	})
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_iam_group_membership", "list users in group",
			BCESDKGoERROR)
	}
	addDebug(action, name)

	listUsersResult := raw.(*api.ListUsersInGroupResult)
	var users []string
	for _, user := range listUsersResult.Users {
		users = append(users, user.Name)
	}
	d.Set("group", d.Id())
	if err := d.Set("users", users); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_iam_group_membership", "setting users in group",
			BCESDKGoERROR)
	}
	return nil
}

func resourceBaiduCloudIamGroupMembershipUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	d.Partial(true)
	if d.HasChange("users") {
		d.SetPartial("users")
		group := d.Get("group").(string)
		action := "Update Group Membership for group " + group
		o, n := d.GetChange("users")
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}

		oldSet := o.(*schema.Set)
		newSet := n.(*schema.Set)

		remove := expandStringSet(oldSet.Difference(newSet))
		add := expandStringSet(newSet.Difference(oldSet))
		addDebug("adding user to group", add)
		addDebug("removing user from group", remove)

		if err := removeUsersFromGroup(client, remove, group); err != nil {
			return err
		}
		if err := addUsersToGroup(client, add, group); err != nil {
			return err
		}
		addDebug(action, group)
	}
	d.Partial(false)
	return resourceBaiduCloudIamGroupMembershipRead(d, meta)
}

func resourceBaiduCloudIamGroupMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	name := d.Get("group").(string)
	userList := expandStringSet(d.Get("users").(*schema.Set))
	action := "Delete Group Membership for group " + name

	if err := removeUsersFromGroup(client, userList, name); err != nil {
		return err
	}
	addDebug(action, name)
	return nil
}

func removeUsersFromGroup(client *connectivity.BaiduClient, users []string, group string) error {
	iamService := IamService{client}
	for _, user := range users {
		err := iamService.DeleteUserFromGroup(user, group)
		if err != nil {
			return err
		}
	}
	return nil
}

func addUsersToGroup(client *connectivity.BaiduClient, users []string, group string) error {
	iamService := IamService{client}
	for _, user := range users {
		err := iamService.AddUserToGroup(user, group)
		if err != nil {
			return err
		}
	}
	return nil
}
