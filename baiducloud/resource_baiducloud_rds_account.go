/*
Use this resource to get information about a RDS Account.

Example Usage

```hcl
resource "baiducloud_rds_account" "default" {
}
```

Import

RDS Account can be imported, e.g.

```hcl
$ terraform import baiducloud_rds_account.default id
```
*/
package baiducloud

import (
	"fmt"
	"strings"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/rds"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudRdsAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudRdsAccountCreate,
		Read:   resourceBaiduCloudRdsAccountRead,
		Delete: resourceBaiduCloudRdsAccountDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Description: "ID of the rds instance.",
				Required:    true,
				ForceNew:    true,
			},
			"account_name": {
				Type:        schema.TypeString,
				Description: "Account name.",
				Required:    true,
				ForceNew:    true,
			},
			"account_type": {
				Type:         schema.TypeString,
				Description:  "Type of the Account, Available values are Common、Super. The default is Common",
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Common", "Super"}, false),
			},
			"password": {
				Type:        schema.TypeString,
				Description: "Operation password.",
				Required:    true,
				ForceNew:    true,
				Sensitive:   true,
			},
			"desc": {
				Type:        schema.TypeString,
				Description: "description.",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Status of the Account.",
				Computed:    true,
			},
		},
	}
}

func resourceBaiduCloudRdsAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	args := &rds.CreateAccountArgs{
		ClientToken: buildClientToken(),
		AccountName: d.Get("account_name").(string),
		Password:    d.Get("password").(string),
	}

	instanceID := d.Get("instance_id").(string)

	if accountType, ok := d.GetOk("account_type"); ok {
		args.AccountType = accountType.(string)
	}

	if desc, ok := d.GetOk("desc"); ok {
		args.Desc = desc.(string)
	}

	action := "Create RDS Account " + args.AccountName
	addDebug(action, args)

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return instanceID, rdsClient.CreateAccount(instanceID, args)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_rds_account", action, BCESDKGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", instanceID, COLON_SEPARATED, args.AccountName))

	return resourceBaiduCloudRdsAccountRead(d, meta)
}

func resourceBaiduCloudRdsAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	items := strings.Split(d.Id(), COLON_SEPARATED)
	instanceID := items[0]
	accountName := items[1]

	action := "Query RDS Account " + accountName

	raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.GetAccount(instanceID, accountName)
	})

	addDebug(action, raw)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_rds_account", action, BCESDKGoERROR)
	}

	result, _ := raw.(*rds.Account)

	d.Set("account_name", result.AccountName)
	d.Set("account_type", result.AccountType)
	d.Set("status", result.Status)
	d.Set("desc", result.Desc)
	return nil
}

func resourceBaiduCloudRdsAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	items := strings.Split(d.Id(), COLON_SEPARATED)
	instanceID := items[0]
	accountName := items[1]

	action := "Delete RDS Account " + accountName

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return accountName, rdsClient.DeleteAccount(instanceID, accountName)
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
		if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_rds_account", action, BCESDKGoERROR)
	}

	return nil
}
