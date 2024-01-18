/*
Provide a resource to create an EIP BP.

Example Usage

```hcl
resource "baiducloud_eipbp" "default" {
  name              = "testEIPbp"
  eip               = 10.23.42.12
  bandwidth_in_mbps = 100
  eip_group_id      = "xxx"
}
```

Import

EIP bp can be imported, e.g.

```hcl
$ terraform import baiducloud_eipbp.default bp_id
```
*/
package baiducloud

import (
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/eip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudEipbp() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudEipbpCreate,
		Read:   resourceBaiduCloudEipbpRead,
		Update: resourceBaiduCloudEipbpUpdate,
		Delete: resourceBaiduCloudEipbpDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"bp_id": {
				Type:        schema.TypeString,
				Description: "id of EIP bp",
				Computed:    true,
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "Eip bp name, length must be between 1 and 65 bytes",
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 65),
			},
			"eip": {
				Type:        schema.TypeString,
				Description: "eip of eip bp",
				Required:    true,
				ForceNew:    true,
			},
			"eip_group_id": {
				Type:        schema.TypeString,
				Description: "eip group id of eip bp",
				Required:    true,
				ForceNew:    true,
			},
			"bandwidth_in_mbps": {
				Type: schema.TypeInt,
				Description: "Eip bp bandwidth(Mbps), if payment_timing is Prepaid or billing_method is ByBandWidth" +
					", support between 1 and 200, if billing_method is ByTraffic, support between 1 and 1000",
				Required: true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Eip bp type",
				Optional:    true,
			},
			"auto_release_time": {
				Type:        schema.TypeString,
				Description: "Eip bp auto release time",
				Optional:    true,
			},
			"bind_type": {
				Type:        schema.TypeString,
				Description: "Eip bp bind type",
				Computed:    true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Description: "Eip bp instance id",
				Computed:    true,
			},
			"eips": {
				Type:        schema.TypeSet,
				Description: "Eip bp eips",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"instance_bandwidth_in_mbps": {
				Type:        schema.TypeInt,
				Description: "Eip bp instance bandwidth in mbps",
				Computed:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "Eip bp create_time",
				Computed:    true,
			},
			"region": {
				Type:        schema.TypeString,
				Description: "Eip bp region",
				Computed:    true,
			},
		},
	}
}

func resourceBaiduCloudEipbpCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.BaiduClient)

	createEipArgs := buildBaiduCloudCreateEipbpArgs(d)

	action := "Create EIP bp " + createEipArgs.Name

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {

		raw, err := client.WithEipClient(func(eipClient *eip.Client) (interface{}, error) {
			return eipClient.CreateEipBp(createEipArgs)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		response, _ := raw.(*eip.CreateEipBpResult)

		addDebug(action, raw)

		d.SetId(response.Id)
		d.Set("bp_id", response.Id)

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eipbp", action, BCESDKGoERROR)
	}

	//stateConf := buildStateConf(EIPProcessingStatus,
	//	[]string{EIPStatusAvailable},
	//	d.Timeout(schema.TimeoutCreate),
	//	eipbpStateRefreshFunc(client, d.Get("bp_id").(string), EIPFailedStatus))

	//if _, err := stateConf.WaitForState(); err != nil {
	//	return WrapError(err)
	//}

	return resourceBaiduCloudEipbpRead(d, meta)
}

func resourceBaiduCloudEipbpRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.BaiduClient)

	eipbpId := d.Id()

	action := "Query EIP bp bpid is " + eipbpId

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {

		raw, err := client.WithEipClient(func(eipClient *eip.Client) (interface{}, error) {
			return eipClient.GetEipBp(eipbpId, buildClientToken())
		})

		if err != nil {
			return resource.RetryableError(err)
		}

		addDebug(action, raw)

		result, _ := raw.(*eip.EipBpDetail)

		d.Set("bp_id", eipbpId)

		d.Set("bind_type", result.BindType)

		d.Set("instance_id", result.InstanceId)

		d.Set("eips", result.Eips)

		d.Set("instance_bandwidth_in_mbps", result.InstanceBandwidthInMbps)

		d.Set("create_time", result.CreateTime)

		d.Set("region", result.Region)

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eipbp", action, BCESDKGoERROR)
	}

	return nil
}

func resourceBaiduCloudEipbpUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.BaiduClient)

	eipbpid := d.Id()

	action := "Update EIP bp" + eipbpid

	if d.HasChange("name") {

		nameArgs := buildBaiduCloudUpdateNameArgs(d)

		_, renameErr := client.WithEipClient(func(eipClient *eip.Client) (interface{}, error) {
			return nil, eipClient.UpdateEipBpName(eipbpid, nameArgs)
		})

		if renameErr != nil {
			if IsExceptedErrors(renameErr, ObjectNotFound) {
				return nil
			}
			return WrapErrorf(renameErr, DefaultErrorMsg, "baiducloud_eipbp", action, BCESDKGoERROR)

		}
	}

	if d.HasChange("bandwidth_in_mbps") {

		resizeArgs := buildBaiduCloudUpdateBandwidthArgs(d)

		_, resizeErr := client.WithEipClient(func(eipClient *eip.Client) (interface{}, error) {
			return nil, eipClient.ResizeEipBp(eipbpid, resizeArgs)
		})

		if resizeErr != nil {
			if IsExceptedErrors(resizeErr, ObjectNotFound) {
				return nil
			}
			return WrapErrorf(resizeErr, DefaultErrorMsg, "baiducloud_eipbp", action, BCESDKGoERROR)

		}
	}

	if d.HasChange("auto_release_time") {

		releaseArgs := buildBaiduCloudUpdateAutoReleaseTimeArgs(d)

		_, releaseErr := client.WithEipClient(func(eipClient *eip.Client) (interface{}, error) {
			return nil, eipClient.UpdateEipBpAutoReleaseTime(eipbpid, releaseArgs)
		})

		if releaseErr != nil {
			if IsExceptedErrors(releaseErr, ObjectNotFound) {
				return nil
			}
			return WrapErrorf(releaseErr, DefaultErrorMsg, "baiducloud_eipbp", action, BCESDKGoERROR)

		}
	}

	return resourceBaiduCloudEipbpRead(d, meta)
}

func resourceBaiduCloudEipbpDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.BaiduClient)

	eipbpId := d.Id()

	action := "Delete EIP bp ID IS " + eipbpId

	_, err := client.WithEipClient(func(eipClient *eip.Client) (interface{}, error) {
		return nil, eipClient.DeleteEipBp(eipbpId, buildClientToken())
	})

	if err != nil {
		if IsExceptedErrors(err, ObjectNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eipbp", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudCreateEipbpArgs(d *schema.ResourceData) *eip.CreateEipBpArgs {

	request := &eip.CreateEipBpArgs{}

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		request.Name = v.(string)
	}

	if v, ok := d.GetOk("eip"); ok && v.(string) != "" {
		request.Eip = v.(string)
	}

	if v, ok := d.GetOk("eip_group_id"); ok && v.(string) != "" {
		request.EipGroupId = v.(string)
	}

	if v := d.Get("bandwidth_in_mbps").(int); v != 0 {
		request.BandwidthInMbps = v
	}

	if v, ok := d.GetOk("type"); ok && v.(string) != "" {
		request.Type = v.(string)
	}

	if v, ok := d.GetOk("auto_release_time"); ok && v.(string) != "" {
		request.AutoReleaseTime = v.(string)
	}

	request.ClientToken = buildClientToken()

	return request
}

//func eipbpStateRefreshFunc(client *connectivity.BaiduClient, id string, failState []string) resource.StateRefreshFunc {
//	return func() (interface{}, string, error) {
//		raw, err := client.WithEipClient(func(eipClient *eip.Client) (interface{}, error) {
//			return eipClient.EipGroupDetail(id)
//		})
//
//		if err != nil {
//			return nil, "", WrapError(err)
//		}
//		result, _ := raw.(*eip.EipbpModel)
//
//		for _, status := range failState {
//			if result.Status == status {
//				return result, result.Status, WrapError(Error(GetFailTargetStatus, result.Status))
//			}
//		}
//
//		return result, result.Status, nil
//	}
//}

func buildBaiduCloudUpdateNameArgs(d *schema.ResourceData) *eip.UpdateEipBpNameArgs {

	request := &eip.UpdateEipBpNameArgs{}

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		request.Name = v.(string)
	}

	request.ClientToken = buildClientToken()

	return request
}

func buildBaiduCloudUpdateBandwidthArgs(d *schema.ResourceData) *eip.ResizeEipBpArgs {

	request := &eip.ResizeEipBpArgs{}

	if v := d.Get("bandwidth_in_mbps").(int); v != 0 {
		request.BandwidthInMbps = v
	}

	request.ClientToken = buildClientToken()

	return request
}

func buildBaiduCloudUpdateAutoReleaseTimeArgs(d *schema.ResourceData) *eip.UpdateEipBpAutoReleaseTimeArgs {

	request := &eip.UpdateEipBpAutoReleaseTimeArgs{}

	if v, ok := d.GetOk("auto_release_time"); ok && v.(string) != "" {
		request.AutoReleaseTime = v.(string)
	}

	request.ClientToken = buildClientToken()

	return request
}
