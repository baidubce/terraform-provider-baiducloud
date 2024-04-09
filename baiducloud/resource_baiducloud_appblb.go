/*
Provide a resource to create an APPBLB.

Example Usage

```hcl
resource "baiducloud_appblb" "default" {
  name        = "testLoadBalance"
  description = "this is a test LoadBalance instance"
  vpc_id      = "vpc-gxaava4knqr1"
  subnet_id   = "sbn-m4x3f2i6c901"

  tags = {
    "tagAKey" = "tagAValue"
    "tagBKey" = "tagBValue"
  }
}
```

Import

APPBLB can be imported, e.g.

```hcl
$ terraform import baiducloud_appblb.default id
```
*/
package baiducloud

import (
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/appblb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudAppBLB() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudAppBLBCreate,
		Read:   resourceBaiduCloudAppBLBRead,
		Update: resourceBaiduCloudAppBLBUpdate,
		Delete: resourceBaiduCloudAppBLBDelete,

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
				Type:         schema.TypeString,
				Description:  "LoadBalance instance's name, length must be between 1 and 65 bytes, and will be automatically generated if not set",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 65),
			},
			"eip": {
				Type:        schema.TypeString,
				Description: "eip of the LoadBalance",
				Optional:    true,
			},
			"auto_renew_length": {
				Type:         schema.TypeInt,
				Description:  "The automatic renewal time is 1-9 per month and 1-3 per year.",
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 9),
			},
			"auto_renew_time_unit": {
				Type:         schema.TypeString,
				Description:  "Monthly payment or annual payment, month is month and year is year.",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"month", "year"}, false),
			},
			"security_groups": {
				Type:        schema.TypeSet,
				Description: "security group ids of the APPBLB.",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enterprise_security_groups": {
				Type:        schema.TypeSet,
				Description: "enterprise security group ids of the APPBLB",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"performance_level": {
				Type:        schema.TypeString,
				Description: "performance level, available values are small1, small2, medium1, medium2, large1, large2, large3",
				Optional:    true,
				ForceNew:    true,
				ValidateFunc: validation.StringInSlice([]string{"small1", "small2", "medium1", "medium2", "large1", "large2", "large3",}, false),
			},
			"description": {
				Type:         schema.TypeString,
				Description:  "LoadBalance's description, length must be between 0 and 450 bytes, and support Chinese",
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 450),
			},
			"status": {
				Type:        schema.TypeString,
				Description: "LoadBalance instance's status, see https://cloud.baidu.com/doc/BLB/s/Pjwvxnxdm/#blbstatus for detail",
				Computed:    true,
			},
			"address": {
				Type:        schema.TypeString,
				Description: "LoadBalance instance's service IP, instance can be accessed through this IP",
				Optional:    true,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "The VPC short ID to which the LoadBalance instance belongs",
				Required:    true,
				ForceNew:    true,
			},
			"vpc_name": {
				Type:        schema.TypeString,
				Description: "The VPC name to which the LoadBalance instance belongs",
				Computed:    true,
			},
			"public_ip": {
				Type:        schema.TypeString,
				Description: "LoadBalance instance's public ip",
				Computed:    true,
			},
			"cidr": {
				Type:        schema.TypeString,
				Description: "Cidr of the network where the LoadBalance instance reside",
				Computed:    true,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Description: "The subnet ID to which the LoadBalance instance belongs",
				Required:    true,
				ForceNew:    true,
			},
			"subnet_name": {
				Type:        schema.TypeString,
				Description: "The subnet name to which the LoadBalance instance belongs",
				Computed:    true,
			},
			"subnet_cidr": {
				Type:        schema.TypeString,
				Description: "Cidr of the subnet which the LoadBalance instance belongs",
				Computed:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "LoadBalance instance's create time",
				Computed:    true,
			},
			"release_time": {
				Type:        schema.TypeString,
				Description: "LoadBalance instance's auto release time",
				Computed:    true,
			},
			"listener": {
				Type:        schema.TypeSet,
				Description: "List of listeners mounted under the instance",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:        schema.TypeInt,
							Description: "Listening port",
							Computed:    true,
						},
						"type": {
							Type:        schema.TypeString,
							Description: "Listening protocol type",
							Computed:    true,
						},
					},
				},
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceBaiduCloudAppBLBCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	appblbService := APPBLBService{client}

	createArgs := buildBaiduCloudCreateAppBlbArgs(d)
	action := "Create APPBLB " + createArgs.Name

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
			return client.CreateLoadBalancer(createArgs)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		addDebug(action, raw)
		response, _ := raw.(*appblb.CreateLoadBalanceResult)
		d.SetId(response.BlbId)
		d.Set("address", response.Address)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appblb", action, BCESDKGoERROR)
	}

	stateConf := buildStateConf(
		APPBLBProcessingStatus,
		APPBLBAvailableStatus,
		d.Timeout(schema.TimeoutCreate),
		appblbService.APPBLBStateRefreshFunc(d.Id(), APPBLBFailedStatus))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	if _, ok := d.GetOk("security_groups"); ok {
		err := appblbService.updateAppBlbSecurityGroups(d, meta)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appblb", action, BCESDKGoERROR)
		}
	}

	if _, ok := d.GetOk("enterprise_security_groups"); ok {
		err := appblbService.updateAppBlbEnterpriseSecurityGroups(d, meta)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appblb", action, BCESDKGoERROR)
		}
	}
	return resourceBaiduCloudAppBLBRead(d, meta)
}
func resourceBaiduCloudAppBLBRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	appblbService := APPBLBService{client}

	blbId := d.Id()
	action := "Query APPBLB " + blbId

	blbModel, blbDetail, err := appblbService.GetAppBLBDetail(blbId)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appblb", action, BCESDKGoERROR)
	}

	d.Set("name", blbModel.Name)
	d.Set("status", blbDetail.Status)
	d.Set("address", blbDetail.Address)
	d.Set("description", blbDetail.Description)
	d.Set("vpc_id", blbModel.VpcId)
	d.Set("vpc_name", blbDetail.VpcName)
	d.Set("subnet_id", blbModel.SubnetId)
	d.Set("subnet_name", blbDetail.SubnetName)
	d.Set("cidr", blbDetail.Cidr)
	d.Set("public_ip", blbDetail.PublicIp)
	d.Set("subnet_cidr", blbDetail.SubnetCider)
	d.Set("create_time", blbDetail.CreateTime)
	d.Set("release_time", blbDetail.ReleaseTime)
	d.Set("performance_level", blbDetail.PerformanceLevel)
	d.Set("listener", appblbService.FlattenListenerModelToMap(blbDetail.Listener))
	if v, ok := d.GetOk("tags"); ok {
		if !slicesContainSameElements(blbDetail.Tags, tranceTagMapToModel(v.(map[string]interface{}))) {
			return WrapErrorf(Error("Tags bind failed."), DefaultErrorMsg, "baiducloud_appblb", action, BCESDKGoERROR)
		}
	}
	d.Set("tags", flattenTagsToMap(blbModel.Tags))
	d.Set("address", blbDetail.Address)

	securityIds, err := appblbService.getAppBlbSecurityGroupIds(d.Id(), meta)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appblb", action, BCESDKGoERROR)
	}
	d.Set("security_groups", securityIds)

	enterpriseSecurityIds, err := appblbService.getAppBlbEnterpriseSecurityGroupIds(d.Id(), meta)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appblb", action, BCESDKGoERROR)
	}
	d.Set("enterprise_security_groups", enterpriseSecurityIds)
	return nil
}

func resourceBaiduCloudAppBLBUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	appblbService := APPBLBService{client}

	blbId := d.Id()
	action := "Update APPBLB " + blbId

	update := false
	updateArgs := &appblb.UpdateLoadBalancerArgs{}

	if d.HasChange("name") {
		update = true
		updateArgs.Name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		update = true
		updateArgs.Description = d.Get("description").(string)
	}

	stateConf := buildStateConf(
		APPBLBProcessingStatus,
		APPBLBAvailableStatus,
		d.Timeout(schema.TimeoutUpdate),
		appblbService.APPBLBStateRefreshFunc(d.Id(), APPBLBFailedStatus))

	if update {
		d.Partial(true)

		updateArgs.ClientToken = buildClientToken()
		_, err := client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
			return blbId, client.UpdateLoadBalancer(blbId, updateArgs)
		})

		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appblb", action, BCESDKGoERROR)
		}

		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}

		d.SetPartial("name")
		d.SetPartial("description")
	}

	d.Partial(false)
	return resourceBaiduCloudAppBLBRead(d, meta)
}
func resourceBaiduCloudAppBLBDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	blbId := d.Id()
	action := "Delete APPBLB " + blbId

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
			return blbId, client.DeleteLoadBalancer(blbId)
		})
		addDebug(action, blbId)

		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		if IsExceptedErrors(err, ObjectNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appblb", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudCreateAppBlbArgs(d *schema.ResourceData) *appblb.CreateLoadBalancerArgs {
	result := &appblb.CreateLoadBalancerArgs{
		ClientToken: buildClientToken(),
	}

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		result.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		result.Description = v.(string)
	}

	if v, ok := d.GetOk("subnet_id"); ok && v.(string) != "" {
		result.SubnetId = v.(string)
	}

	if v, ok := d.GetOk("vpc_id"); ok && v.(string) != "" {
		result.VpcId = v.(string)
	}

	if v, ok := d.GetOk("tags"); ok {
		result.Tags = tranceTagMapToModel(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("eip"); ok {
		result.Eip = v.(string)
	}

	if v, ok := d.GetOk("address"); ok {
		result.Address = v.(string)
	}

	if v, ok := d.GetOk("performance_level"); ok {
		result.PerformanceLevel = v.(string)
	}

	if v, ok := d.GetOk("auto_renew_length"); ok {
		result.AutoRenewLength = v.(int)
		if result.AutoRenewLength > 0 {
			if v, ok := d.GetOk("auto_renew_time_unit"); ok {
				result.AutoRenewTimeUnit = v.(string)
			}
		}
	}

	return result
}
