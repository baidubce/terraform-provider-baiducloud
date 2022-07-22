/*
Provide a resource to create a CFC Function Alias.

Example Usage

```hcl
resource "baiducloud_cfc_alias" "default" {
  function_name    = "terraform-cfc"
  function_version = "$LATEST"
  alias_name       = "terraformAlias"
  description      = "terraform create alias"
}
```

```
*/
package baiducloud

import (
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/cfc"
	"github.com/baidubce/bce-sdk-go/services/cfc/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudCFCAlias() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudCFCAliasCreate,
		Read:   resourceBaiduCloudCFCAliasRead,
		Update: resourceBaiduCloudCFCAliasUpdate,
		Delete: resourceBaiduCloudCFCAliasDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"alias_name": {
				Type:        schema.TypeString,
				Description: "CFC Function alias name",
				Required:    true,
				ForceNew:    true,
			},
			"function_name": {
				Type:        schema.TypeString,
				Description: "CFC Function name",
				Required:    true,
				ForceNew:    true,
			},
			"function_version": {
				Type:        schema.TypeString,
				Description: "CFC Function version this alias binded",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "CFC Function alias description",
				Optional:    true,
			},
			"alias_brn": {
				Type:        schema.TypeString,
				Description: "CFC Function alias brn",
				Computed:    true,
			},
			"alias_arn": {
				Type:        schema.TypeString,
				Description: "CFC Function alias arn",
				Computed:    true,
			},
			"uid": {
				Type:        schema.TypeString,
				Description: "CFC Function uid",
				Computed:    true,
			},
			"update_time": {
				Type:        schema.TypeString,
				Description: "CFC Function alias update time",
				Computed:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "CFC Function alias create time",
				Computed:    true,
			},
		},
	}
}

func resourceBaiduCloudCFCAliasCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	createArgs := &api.CreateAliasArgs{
		Name:            d.Get("alias_name").(string),
		FunctionName:    d.Get("function_name").(string),
		FunctionVersion: d.Get("function_version").(string),
	}
	if value, ok := d.GetOk("description"); ok {
		createArgs.Description = value.(string)
	}

	action := "Create CFC function " + createArgs.FunctionName + " alias " + createArgs.Name
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithCFCClient(func(client *cfc.Client) (i interface{}, e error) {
			return client.CreateAlias(createArgs)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		response, _ := raw.(*api.CreateAliasResult)

		addDebug(action, raw)
		d.SetId(response.FunctionName + "-" + response.Name)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfc_alias", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudCFCAliasRead(d, meta)
}

func resourceBaiduCloudCFCAliasRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	getArgs := &api.GetAliasArgs{
		FunctionName: d.Get("function_name").(string),
		AliasName:    d.Get("alias_name").(string),
	}
	action := "Query Function " + getArgs.FunctionName + " with alias " + getArgs.AliasName

	raw, err := client.WithCFCClient(func(client *cfc.Client) (i interface{}, e error) {
		return client.GetAlias(getArgs)
	})

	if err != nil {
		d.SetId("")
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfc_alias", action, BCESDKGoERROR)
	}

	addDebug(action, raw)

	response, _ := raw.(*api.GetAliasResult)
	d.Set("alias_name", response.Name)
	d.Set("function_name", response.FunctionName)
	d.Set("function_version", response.FunctionVersion)
	d.Set("description", response.Description)
	d.Set("alias_brn", response.AliasBrn)
	d.Set("alias_arn", response.AliasArn)
	d.Set("uid", response.Uid)
	d.Set("update_time", response.UpdatedAt.String())
	d.Set("create_time", response.CreatedAt.String())

	return nil
}

func resourceBaiduCloudCFCAliasUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	updateAliasArgs := &api.UpdateAliasArgs{
		FunctionName:    d.Get("function_name").(string),
		AliasName:       d.Get("alias_name").(string),
		FunctionVersion: d.Get("function_version").(string),
	}

	update := d.HasChange("function_version")
	if d.HasChange("description") {
		update = true
		updateAliasArgs.Description = d.Get("description").(string)
	}

	if update {
		action := "Update CFC Function " + updateAliasArgs.FunctionName + " alias " + updateAliasArgs.AliasName
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			raw, err := client.WithCFCClient(func(client *cfc.Client) (i interface{}, e error) {
				return client.UpdateAlias(updateAliasArgs)
			})

			if err != nil {
				if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			addDebug(action, raw)
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfc_alias", action, BCESDKGoERROR)
		}
	}

	return resourceBaiduCloudCFCAliasRead(d, meta)
}

func resourceBaiduCloudCFCAliasDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	deleteArgs := &api.DeleteAliasArgs{
		FunctionName: d.Get("function_name").(string),
		AliasName:    d.Get("alias_name").(string),
	}

	action := "Delete CFC Function " + deleteArgs.FunctionName + " alias " + deleteArgs.AliasName
	err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		_, err := client.WithCFCClient(func(client *cfc.Client) (i interface{}, e error) {
			return nil, client.DeleteAlias(deleteArgs)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		addDebug(action, deleteArgs)
		return nil
	})

	if err != nil {
		if IsExceptedErrors(err, ObjectNotFound) {
			return nil
		}

		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfc_alias", action, BCESDKGoERROR)
	}

	return nil
}
