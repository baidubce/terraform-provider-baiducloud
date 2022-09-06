/*
Provide a resource to create an BLB.

Example Usage

```hcl
resource "baiducloud_blb" "default" {
  name        = "testLoadBalance"
  description = "this is a test LoadBalance instance"
  vpc_id      = "vpc-xxxx"
  subnet_id   = "sbn-xxxx"

  tags = {
    "tagAKey" = "tagAValue"
    "tagBKey" = "tagBValue"
  }
}
```

Import

BLB can be imported, e.g.

```hcl
$ terraform import baiducloud_blb.default id
```
*/
package baiducloud

import (
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/blb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudBLB() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudBLBCreate,
		Read:   resourceBaiduCloudBLBRead,
		Update: resourceBaiduCloudBLBUpdate,
		Delete: resourceBaiduCloudBLBDelete,

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
				Computed:    true,
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
			"create_time": {
				Type:        schema.TypeString,
				Description: "LoadBalance instance's create time",
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

func resourceBaiduCloudBLBCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	blbService := BLBService{client}

	createArgs := buildBaiduCloudCreateBlbArgs(d)
	action := "Create BLB " + createArgs.Name

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithBLBClient(func(client *blb.Client) (i interface{}, e error) {
			return client.CreateLoadBalancer(createArgs)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		addDebug(action, raw)
		response, _ := raw.(*blb.CreateLoadBalancerResult)
		d.SetId(response.BlbId)
		d.Set("address", response.Address)

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb", action, BCESDKGoERROR)
	}

	stateConf := buildStateConf(
		BLBProcessingStatus,
		BLBAvailableStatus,
		d.Timeout(schema.TimeoutCreate),
		blbService.BLBStateRefreshFunc(d.Id(), BLBFailedStatus))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	return resourceBaiduCloudBLBRead(d, meta)
}
func resourceBaiduCloudBLBRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	blbService := BLBService{client}

	blbId := d.Id()
	action := "Query BLB " + blbId

	blbModel, blbDetail, err := blbService.GetBLBDetail(blbId)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb", action, BCESDKGoERROR)
	}

	d.Set("name", blbModel.Name)
	d.Set("status", blbDetail.Status)
	d.Set("address", blbDetail.Address)
	d.Set("description", blbDetail.Description)
	d.Set("vpc_id", blbModel.VpcId)
	d.Set("vpc_name", blbDetail.VpcName)
	d.Set("subnet_id", blbModel.SubnetId)
	d.Set("cidr", blbDetail.Cidr)
	d.Set("public_ip", blbDetail.PublicIp)
	d.Set("create_time", blbDetail.CreateTime)
	d.Set("listener", blbService.FlattenListenerModelToMap(blbDetail.Listener))
	d.Set("tags", flattenTagsToMap(blbModel.Tags))

	return nil
}

func resourceBaiduCloudBLBUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	blbService := BLBService{client}

	blbId := d.Id()
	action := "Update BLB " + blbId

	update := false
	updateArgs := &blb.UpdateLoadBalancerArgs{}

	if d.HasChange("name") {
		update = true
		updateArgs.Name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		update = true
		updateArgs.Description = d.Get("description").(string)
	}

	stateConf := buildStateConf(
		BLBProcessingStatus,
		BLBAvailableStatus,
		d.Timeout(schema.TimeoutUpdate),
		blbService.BLBStateRefreshFunc(d.Id(), BLBFailedStatus))

	if update {
		d.Partial(true)

		updateArgs.ClientToken = buildClientToken()
		_, err := client.WithBLBClient(func(client *blb.Client) (i interface{}, e error) {
			return blbId, client.UpdateLoadBalancer(blbId, updateArgs)
		})

		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb", action, BCESDKGoERROR)
		}

		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}

		d.SetPartial("name")
		d.SetPartial("description")
	}

	d.Partial(false)
	return resourceBaiduCloudBLBRead(d, meta)
}
func resourceBaiduCloudBLBDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	blbId := d.Id()
	action := "Delete BLB " + blbId

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := client.WithBLBClient(func(client *blb.Client) (i interface{}, e error) {
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
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudCreateBlbArgs(d *schema.ResourceData) *blb.CreateLoadBalancerArgs {
	result := &blb.CreateLoadBalancerArgs{
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

	return result
}
