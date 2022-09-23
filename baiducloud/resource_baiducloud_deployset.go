/*
Use this resource to creat a deployset.

Example Usage

```hcl
resource "baiducloud_deployset" "default" {
  name     = "terraform-test"
  desc     = "test desc"
  strategy = "HOST_HA"
}
```

Import

deployset can be imported, e.g.

```hcl
$ terraform import baiducloud_deployset.default id
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudDeploySet() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudDeploySetCreate,
		Read:   resourceBaiduCloudDeploySetRead,
		Update: resourceBaiduCloudDeploySetUpdate,
		Delete: resourceBaiduCloudDeploySetDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the deployset. Support for uppercase and lowercase letters, numbers, Chinese and special characters, such as \"-\",\"_\",\"/\",\".\", the value must start with a letter, length 1-65.",
				Optional:    true,
			},
			"desc": {
				Type:        schema.TypeString,
				Description: "Description of the deployset.",
				Optional:    true,
			},
			"strategy": {
				Type:        schema.TypeString,
				Description: "Strategy of deployset.Available values are HOST_HA, RACK_HA and TOR_HA",
				Optional:    true,
			},
			"short_id": {
				Type:        schema.TypeString,
				Description: "deployset short id.",
				Computed:    true,
			},
			"uuid": {
				Type:        schema.TypeString,
				Description: "deployset uuid.",
				Computed:    true,
			},
			"concurrency": {
				Type:        schema.TypeInt,
				Description: "concurrency of deployset.",
				Computed:    true,
			},
			"az_intstance_statis_list": {
				Type:        schema.TypeList,
				Description: "Availability Zone Instance Statistics List.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_count": {
							Type:        schema.TypeInt,
							Description: "Count of instance which is in the deployset.",
							Computed:    true,
						},
						"bcc_instance_cnt": {
							Type:        schema.TypeInt,
							Description: "Count of BCC instance which is in the deployset.",
							Computed:    true,
						},
						"bbc_instance_cnt": {
							Type:        schema.TypeInt,
							Description: "Count of BBC instance which is in the deployset.",
							Computed:    true,
						},
						"instance_total": {
							Type:        schema.TypeInt,
							Description: "Total of instance which is in the deployset.",
							Computed:    true,
						},
						"zone_name": {
							Type:        schema.TypeString,
							Description: "Zone name of deployset.",
							Computed:    true,
						},
						"instance_ids": {
							Type:        schema.TypeSet,
							Description: "IDs of instance which is in the deployset.",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"bcc_instance_ids": {
							Type:        schema.TypeSet,
							Description: "IDs of BCC instance which is in the deployset.",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"bbc_instance_ids": {
							Type:        schema.TypeSet,
							Description: "IDs of BBC instance which is in the deployset..",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}
func resourceBaiduCloudDeploySetCreate(d *schema.ResourceData, meta interface{}) error {
	action := "Create deploy set"
	client := meta.(*connectivity.BaiduClient)
	createDeploySetArgs := &api.CreateDeploySetArgs{
		ClientToken: buildClientToken(),
	}
	if v, ok := d.GetOk("name"); ok {
		createDeploySetArgs.Name = v.(string)
	}
	if v, ok := d.GetOk("desc"); ok {
		createDeploySetArgs.Desc = v.(string)
	}
	if v, ok := d.GetOk("strategy"); ok {
		createDeploySetArgs.Strategy = v.(string)
	}
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		res, err := client.WithBccClient(func(bccClient *bcc.Client) (i interface{}, e error) {
			return bccClient.CreateDeploySet(createDeploySetArgs)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, res)
		d.SetId(res.(*api.CreateDeploySetResult).DeploySetId)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_deployset", action, BCESDKGoERROR)
	}
	return resourceBaiduCloudDeploySetRead(d, meta)
}
func resourceBaiduCloudDeploySetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	action := "Query deploy set detail."
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		res, err := client.WithBccClient(func(bccClient *bcc.Client) (i interface{}, e error) {
			return bccClient.GetDeploySet(d.Id())
		})
		addDebug(action, res)
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		deploySet := res.(*api.DeploySetResult)
		d.Set("name", deploySet.Name)
		d.Set("desc", deploySet.Desc)
		d.Set("strategy", deploySet.Strategy)
		d.Set("concurrency", deploySet.Concurrency)
		intstanceStatisMap := make([]map[string]interface{}, 0, len(deploySet.InstanceList))

		for _, ins := range deploySet.InstanceList {
			intstanceStatisMap = append(intstanceStatisMap, map[string]interface{}{
				"bcc_instance_cnt": ins.BccCount,
				"bbc_instance_cnt": ins.BbcCount,
				"instance_count":   ins.Count,
				"instance_total":   ins.Total,
				"zone_name":        ins.ZoneName,
				"instance_ids":     ins.InstanceIds,
				"bcc_instance_ids": ins.BccInstanceIds,
				"bbc_instance_ids": ins.BbcInstanceIds,
			})
		}
		d.Set("az_intstance_statis_list", intstanceStatisMap)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_deployset", action, BCESDKGoERROR)
	}
	return nil
}
func resourceBaiduCloudDeploySetUpdate(d *schema.ResourceData, meta interface{}) error {
	action := "Update deploy set attribute "
	client := meta.(*connectivity.BaiduClient)
	args := &api.ModifyDeploySetArgs{
		ClientToken: buildClientToken(),
	}
	if d.HasChange("name") {
		args.Name = d.Get("name").(string)
	}
	if d.HasChange("desc") {
		args.Desc = d.Get("desc").(string)
	}
	err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		_, err := client.WithBccClient(func(bccClient *bcc.Client) (interface{}, error) {
			err, _ := bccClient.ModifyDeploySet(d.Id(), args)
			return nil, err
		})
		if err != nil {
			if IsExceptedErrors(err, []string{OperationDenied, bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, args)
		return nil
	})
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_deployset", action, BCESDKGoERROR)
	}
	return resourceBaiduCloudDeploySetRead(d, meta)
}
func resourceBaiduCloudDeploySetDelete(d *schema.ResourceData, meta interface{}) error {
	action := "delete deploy set"
	client := meta.(*connectivity.BaiduClient)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err := client.WithBccClient(func(bccClient *bcc.Client) (i interface{}, e error) {
			return nil, bccClient.DeleteDeploySet(d.Id())
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
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_deployset", action, BCESDKGoERROR)
	}
	return nil
}
