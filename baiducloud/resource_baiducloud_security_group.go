/*
Provide a resource to create a security group.

Example Usage

```hcl
resource "baiducloud_security_group" "default" {
  name        = "testSecurityGroup"
  description = "default"
  tags {
    tag_key   = "testKey"
    tag_value = "testValue"
  }
}
```

Import

Bcc SecurityGroup can be imported, e.g.

```hcl
$ terraform import baiducloud_security_group.default security_group_id
```
*/
package baiducloud

import (
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudSecurityGroupCreate,
		Read:   resourceBaiduCloudSecurityGroupRead,
		Delete: resourceBaiduCloudSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "SecurityGroup name",
				Required:    true,
				ForceNew:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "SecurityGroup description",
				Optional:    true,
				ForceNew:    true,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "SecurityGroup binded VPC id",
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceBaiduCloudSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	createSecurityGroupArgs := buildBaiduCloudSecurityGroupArgs(d, meta)

	action := "Create SecurityGroup " + createSecurityGroupArgs.Name
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithBccClient(func(bccClient *bcc.Client) (interface{}, error) {
			return bccClient.CreateSecurityGroup(createSecurityGroupArgs)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)

		response, _ := raw.(*api.CreateSecurityGroupResult)
		d.SetId(response.SecurityGroupId)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_security_group", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudSecurityGroupRead(d, meta)
}

func resourceBaiduCloudSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	securityGroupID := d.Id()

	listSecurityGroupArgs := &api.ListSecurityGroupArgs{}
	if v, ok := d.GetOk("vpc_id"); ok && v.(string) != "" {
		listSecurityGroupArgs.VpcId = v.(string)
	}
	action := "Query SecurityGroup " + securityGroupID

	isTruncated := true
	for isTruncated {
		raw, err := client.WithBccClient(func(bccClient *bcc.Client) (interface{}, error) {
			return bccClient.ListSecurityGroup(listSecurityGroupArgs)
		})
		addDebug(action, raw)

		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_security_group", action, BCESDKGoERROR)
		}
		response, _ := raw.(*api.ListSecurityGroupResult)

		for _, sg := range response.SecurityGroups {
			if sg.Id == securityGroupID {
				d.Set("name", sg.Name)
				d.Set("description", sg.Desc)
				d.Set("vpc_id", sg.VpcId)
				d.Set("tags", flattenTagsToMap(sg.Tags))

				return nil
			}
		}

		listSecurityGroupArgs.Marker = response.Marker
		listSecurityGroupArgs.MaxKeys = response.MaxKeys
		isTruncated = response.IsTruncated
	}

	// no found securityGroup
	d.SetId("")
	return nil
}

func resourceBaiduCloudSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	securityGroupID := d.Id()
	action := "Delete SecurityGroup " + securityGroupID

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithBccClient(func(bccClient *bcc.Client) (interface{}, error) {
			return securityGroupID, bccClient.DeleteSecurityGroup(securityGroupID)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR, SECURITYGROUP_INUSE_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)
		return nil
	})

	if err != nil {
		if IsExceptedErrors(err, ObjectNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_security_group", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudSecurityGroupArgs(d *schema.ResourceData, meta interface{}) *api.CreateSecurityGroupArgs {
	request := &api.CreateSecurityGroupArgs{
		ClientToken: buildClientToken(),
	}

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		request.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.Desc = v.(string)
	}

	if v, ok := d.GetOk("vpc_id"); ok && v.(string) != "" {
		request.VpcId = v.(string)
	}

	if v, ok := d.GetOk("tags"); ok {
		request.Tags = tranceTagMapToModel(v.(*schema.Set).List())
	}

	return request
}
