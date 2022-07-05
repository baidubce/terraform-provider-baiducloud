/*
Provide a resource to create an EIP.

Example Usage

```hcl
resource "baiducloud_eip" "default" {
  name              = "testEIP"
  bandwidth_in_mbps = 100
  payment_timing    = "Postpaid"
  billing_method    = "ByTraffic"
}
```

Import

EIP can be imported, e.g.

```hcl
$ terraform import baiducloud_eip.default eip
```
*/
package baiducloud

import (
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/eip"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudEip() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudEipCreate,
		Read:   resourceBaiduCloudEipRead,
		Update: resourceBaiduCloudEipUpdate,
		Delete: resourceBaiduCloudEipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"eip": {
				Type:        schema.TypeString,
				Description: "Eip address",
				Computed:    true,
				ForceNew:    true,
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "Eip name, length must be between 1 and 65 bytes",
				Computed:     true,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 65),
			},
			"route_type": {
				Type:        schema.TypeString,
				Description: "EIP route type",
				Optional:    true,
			},
			"bandwidth_in_mbps": {
				Type:        schema.TypeInt,
				Description: "Eip bandwidth(Mbps), if payment_timing is Prepaid or billing_method is ByBandWidth, support between 1 and 200, if billing_method is ByTraffic, support between 1 and 1000",
				Required:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Eip status",
				Computed:    true,
			},
			"eip_instance_type": {
				Type:        schema.TypeString,
				Description: "Eip instance type",
				Computed:    true,
			},
			"share_group_id": {
				Type:        schema.TypeString,
				Description: "Eip share group id",
				Computed:    true,
			},
			"payment_timing": {
				Type:         schema.TypeString,
				Description:  "Eip payment timing, support Prepaid and Postpaid",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{PaymentTimingPrepaid, PaymentTimingPostpaid}, false),
			},
			"billing_method": {
				Type:        schema.TypeString,
				Description: "Eip billing method, support ByTraffic or ByBandwidth",
				Required:    true,
				ForceNew:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "Eip create time",
				Computed:    true,
			},
			"expire_time": {
				Type:        schema.TypeString,
				Description: "Eip expire time",
				Computed:    true,
			},
			"reservation_length": {
				Type:             schema.TypeInt,
				Description:      "Eip Prepaid billing reservation length, only useful when payment_timing is Prepaid",
				Optional:         true,
				Sensitive:        true,
				DiffSuppressFunc: postPaidDiffSuppressFunc,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				ConflictsWith:    []string{"auto_renew_time", "auto_renew_time_unit"},
			},
			"reservation_time_unit": {
				Type:             schema.TypeString,
				Description:      "Eip Prepaid billing reservation time unit, only useful when payment_timing is Prepaid",
				Optional:         true,
				Sensitive:        true,
				DiffSuppressFunc: postPaidDiffSuppressFunc,
				ValidateFunc:     validation.StringInSlice([]string{"month"}, false),
				ConflictsWith:    []string{"auto_renew_time", "auto_renew_time_unit"},
			},
			"auto_renew_time": {
				Type:             schema.TypeInt,
				Description:      "Eip auto renew time length, only useful when payment_timing is Prepaid. If auto_renew_time_unit is month, support 1-9, if auto_renew_time_unit is year, support 1-3.",
				Optional:         true,
				DiffSuppressFunc: postPaidDiffSuppressFunc,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
				ConflictsWith:    []string{"reservation_length", "reservation_time_unit"},
			},
			"auto_renew_time_unit": {
				Type:             schema.TypeString,
				Description:      "Eip auto renew time unit, only useful when payment_timing is Prepaid, support month/year",
				Optional:         true,
				DiffSuppressFunc: postPaidDiffSuppressFunc,
				ValidateFunc:     validation.StringInSlice([]string{"month", "year"}, false),
				ConflictsWith:    []string{"reservation_length", "reservation_time_unit"},
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceBaiduCloudEipCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	eipClient := EipService{client}

	createEipArgs := buildBaiduCloudCreateEipArgs(d)
	action := "Create EIP " + createEipArgs.Name

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithEipClient(func(eipClient *eip.Client) (interface{}, error) {
			return eipClient.CreateEip(createEipArgs)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response, _ := raw.(*eip.CreateEipResult)

		addDebug(action, raw)
		d.Set("eip", response.Eip)
		d.SetId(response.Eip)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eip", action, BCESDKGoERROR)
	}

	stateConf := buildStateConf(EIPProcessingStatus,
		[]string{EIPStatusAvailable},
		d.Timeout(schema.TimeoutCreate),
		eipClient.EipStateRefreshFunc(d.Get("eip").(string), EIPFailedStatus))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	return resourceBaiduCloudEipRead(d, meta)
}

func resourceBaiduCloudEipRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	eipClient := EipService{client}

	eipAddr := d.Id()

	action := "Query EIP " + eipAddr
	result, err := eipClient.EipGetDetail(eipAddr)
	addDebug(action, result)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			d.Set("eip", "")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eip", action, BCESDKGoERROR)
	}

	d.Set("name", result.Name)
	d.Set("bandwidth_in_mbps", result.BandWidthInMbps)
	d.Set("status", result.Status)
	d.Set("eip_instance_type", result.EipInstanceType)
	d.Set("share_group_id", result.ShareGroupId)
	d.Set("payment_timing", result.PaymentTiming)
	d.Set("billing_method", result.BillingMethod)
	d.Set("create_time", result.CreateTime)
	d.Set("expire_time", result.ExpireTime)
	d.Set("tags", flattenTagsToMap(result.Tags))
	d.Set("eip", result.Eip)

	return nil
}

func resourceBaiduCloudEipUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	eipClient := EipService{client}

	eipAddr := d.Id()
	stateConf := buildStateConf(EIPProcessingStatus,
		[]string{EIPStatusAvailable, EIPStatusBinded},
		d.Timeout(schema.TimeoutUpdate),
		eipClient.EipStateRefreshFunc(eipAddr, EIPFailedStatus))

	if d.HasChange("bandwidth_in_mbps") {

		if err := eipClient.EipResizeBandwidth(eipAddr, d.Get("bandwidth_in_mbps").(int)); err != nil {
			return WrapError(err)
		}

		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}

		d.SetPartial("bandwidth_in_mbps")
	}

	if d.Get("billing_method").(string) == PaymentTimingPrepaid &&
		(d.HasChange("auto_renew_time") || d.HasChange("auto_renew_time_unit")) {
		isStart, args := buildUpdateAutoRenewArgs(d)
		if isStart {
			if err := eipClient.StartAutoRenew(eipAddr, args); err != nil {
				return WrapError(err)
			}
		} else {
			if err := eipClient.StopAutoRenew(eipAddr); err != nil {
				return WrapError(err)
			}
		}
	}

	d.Partial(false)
	return resourceBaiduCloudEipRead(d, meta)
}

func resourceBaiduCloudEipDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	eipService := EipService{client: client}

	eipAddr := d.Id()
	action := "Delete EIP " + eipAddr

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		eipDetail, errGet := eipService.EipGetDetail(eipAddr)
		if errGet != nil {
			return resource.NonRetryableError(errGet)
		}

		if eipDetail.Status == EIPStatusBinded {
			return resource.RetryableError(eipStillInUsed)
		}

		raw, errDelete := client.WithEipClient(func(client *eip.Client) (i interface{}, e error) {
			return eipAddr, client.DeleteEip(eipAddr, buildClientToken())
		})
		if errDelete != nil {
			if IsExceptedErrors(errDelete, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(errDelete)
			}
			return resource.NonRetryableError(errDelete)
		}

		addDebug(action, raw)
		return nil
	})

	if err != nil {
		if IsExceptedErrors(err, ObjectNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eip", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudCreateEipArgs(d *schema.ResourceData) *eip.CreateEipArgs {
	request := &eip.CreateEipArgs{}

	if v, ok := d.GetOk("route_type"); ok && v.(string) != "" {
		request.RouteType = v.(string)
	}

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		request.Name = v.(string)
	}

	if v := d.Get("bandwidth_in_mbps").(int); v != 0 {
		request.BandWidthInMbps = v
	}

	if v, ok := d.GetOk("tags"); ok {
		request.Tags = tranceTagMapToModel(v.(map[string]interface{}))
	}
	request.Billing = &eip.Billing{
		PaymentTiming: d.Get("payment_timing").(string),
		BillingMethod: d.Get("billing_method").(string),
	}

	if request.Billing.PaymentTiming == "Prepaid" {
		request.Billing.Reservation = &eip.Reservation{}

		if v, ok := d.GetOk("reservation_length"); ok && v.(int) > 0 {
			request.Billing.Reservation.ReservationLength = v.(int)
		}

		if v, ok := d.GetOk("reservation_time_unit"); ok && len(v.(string)) > 0 {
			request.Billing.Reservation.ReservationTimeUnit = v.(string)
		}

		if v, ok := d.GetOk("auto_renew_time"); ok && v.(int) > 0 {
			request.AutoRenewTime = v.(int)
		}

		if v, ok := d.GetOk("auto_renew_time_unit"); ok && len(v.(string)) > 0 {
			request.AutoRenewTimeUnit = v.(string)
		}
	}

	request.ClientToken = buildClientToken()

	return request
}

func buildUpdateAutoRenewArgs(d *schema.ResourceData) (bool, *eip.StartAutoRenewArgs) {
	_, nTime := d.GetChange("auto_renew_time")
	_, nTimeUnit := d.GetChange("auto_renew_time_unit")

	newIsStop := nTime == nil || nTime.(int) == 0 || nTimeUnit == nil || !stringInSlice([]string{"month", "year"}, nTimeUnit.(string))

	if !newIsStop {
		args := eip.StartAutoRenewArgs{
			AutoRenewTime:     nTime.(int),
			AutoRenewTimeUnit: nTimeUnit.(string),
			ClientToken:       buildClientToken(),
		}

		return true, &args
	}

	return false, nil
}
