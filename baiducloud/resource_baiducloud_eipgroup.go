/*
Provide a resource to create an EIP GROUP.

Example Usage

```hcl
resource "baiducloud_eipgroup" "default" {
  name              = "testEIPgroup"
  eip_count         = 2
  bandwidth_in_mbps = 100
  payment_timing    = "Postpaid"
  billing_method    = "ByTraffic"
}
```

Import

EIP group can be imported, e.g.

```hcl
$ terraform import baiducloud_eipgroup.default group_id
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

func resourceBaiduCloudEipgroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudEipgroupCreate,
		Read:   resourceBaiduCloudEipgroupRead,
		Update: resourceBaiduCloudEipgroupUpdate,
		Delete: resourceBaiduCloudEipgroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Description: "id of EIP group",
				Computed:    true,
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "Eip group name, length must be between 1 and 65 bytes",
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 65),
			},
			"eip_count": {
				Type:        schema.TypeInt,
				Description: "count of eip group",
				Required:    true,
			},
			"bandwidth_in_mbps": {
				Type: schema.TypeInt,
				Description: "Eip group bandwidth(Mbps), if payment_timing is Prepaid or billing_method is ByBandWidth" +
					", support between 1 and 200, if billing_method is ByTraffic, support between 1 and 1000",
				Required: true,
			},
			"payment_timing": {
				Type:         schema.TypeString,
				Description:  "Eip group payment timing, support Prepaid and Postpaid",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{PaymentTimingPrepaid, PaymentTimingPostpaid}, false),
			},
			"billing_method": {
				Type:        schema.TypeString,
				Description: "Eip group billing method, support ByTraffic or ByBandwidth",
				Required:    true,
				ForceNew:    true,
			},
			"reservation_length": {
				Type:             schema.TypeInt,
				Description:      "Eip group Prepaid billing reservation length, only useful when payment_timing is Prepaid",
				Optional:         true,
				Sensitive:        true,
				DiffSuppressFunc: postPaidDiffSuppressFunc,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				//ConflictsWith:    []string{"auto_renew_time", "auto_renew_time_unit"},
			},
			"reservation_time_unit": {
				Type:             schema.TypeString,
				Description:      "Eip group Prepaid billing reservation time unit, only useful when payment_timing is Prepaid",
				Optional:         true,
				Sensitive:        true,
				DiffSuppressFunc: postPaidDiffSuppressFunc,
				ValidateFunc:     validation.StringInSlice([]string{"month"}, false),
				//ConflictsWith:    []string{"auto_renew_time", "auto_renew_time_unit"},
			},
			"tags": tagsSchema(),
			"route_type": {
				Type:        schema.TypeString,
				Description: "Eip Group routeType",
				Computed:    true,
			},
			"idc": {
				Type:        schema.TypeString,
				Description: "idc of Eip group",
				Computed:    true,
			},
			"continuous": {
				Type:        schema.TypeBool,
				Description: "Eip group continuous",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Eip group status",
				Computed:    true,
			},
			"default_domestic_bandwidth": {
				Type:        schema.TypeInt,
				Description: "Eip group status",
				Computed:    true,
			},
			"bw_short_id": {
				Type:        schema.TypeString,
				Description: "Eip group status",
				Computed:    true,
			},
			"bw_bandwidth_in_mbps": {
				Type:        schema.TypeInt,
				Description: "Eip group status",
				Computed:    true,
			},
			"domestic_bw_short_id": {
				Type:        schema.TypeString,
				Description: "Eip group status",
				Computed:    true,
			},
			"domestic_bw_bandwidth_in_mbps": {
				Type:        schema.TypeInt,
				Description: "Eip group status",
				Computed:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "Eip group create time",
				Computed:    true,
			},
			"expire_time": {
				Type:        schema.TypeString,
				Description: "Eip group expire time",
				Computed:    true,
			},
			"region": {
				Type:        schema.TypeString,
				Description: "region of eip group",
				Computed:    true,
			},
			"eips": {
				Type:        schema.TypeList,
				Description: "Eip list",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"eip": {
							Type:        schema.TypeString,
							Description: "Eip address",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Eip name",
							Computed:    true,
						},
						"bandwidth_in_mbps": {
							Type:        schema.TypeInt,
							Description: "Eip bandwidth(Mbps)",
							Computed:    true,
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
							Type:        schema.TypeString,
							Description: "Eip payment timing",
							Computed:    true,
						},
						"billing_method": {
							Type:        schema.TypeString,
							Description: "Eip billing method",
							Computed:    true,
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
						"tags": tagsComputedSchema(),
					},
				},
			},
		},
	}
}

func resourceBaiduCloudEipgroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.BaiduClient)

	createEipArgs := buildBaiduCloudCreateEipgroupArgs(d)

	action := "Create EIP GROUP " + createEipArgs.Name

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {

		raw, err := client.WithEipClient(func(eipClient *eip.Client) (interface{}, error) {
			return eipClient.CreateEipGroup(createEipArgs)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		response, _ := raw.(*eip.CreateEipGroupResult)

		addDebug(action, raw)

		d.SetId(response.Id)
		d.Set("group_id", response.Id)

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eipgroup", action, BCESDKGoERROR)
	}

	stateConf := buildStateConf(EIPProcessingStatus,
		[]string{EIPStatusAvailable},
		d.Timeout(schema.TimeoutCreate),
		eipGroupStateRefreshFunc(client, d.Get("group_id").(string), EIPFailedStatus))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	return resourceBaiduCloudEipgroupRead(d, meta)
}

func resourceBaiduCloudEipgroupRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.BaiduClient)

	eipGroupId := d.Id()

	action := "Query EIP group groupid is " + eipGroupId

	raw, err := client.WithEipClient(func(eipClient *eip.Client) (interface{}, error) {
		return eipClient.EipGroupDetail(eipGroupId)
	})

	addDebug(action, raw)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eipgroup", action, BCESDKGoERROR)
	}

	result, _ := raw.(*eip.EipGroupModel)

	d.Set("group_id", eipGroupId)

	d.Set("status", result.Status)

	d.Set("default_domestic_bandwidth", result.DefaultDomesticBandwidth)

	d.Set("bw_short_id", result.BwShortId)

	d.Set("bw_bandwidth_in_mbps", result.BwBandwidthInMbps)

	d.Set("domestic_bw_short_id", result.DomesticBwShortId)

	d.Set("domestic_bw_bandwidth_in_mbps", result.DomesticBwBandwidthInMbps)

	d.Set("create_time", result.CreateTime)

	d.Set("expire_time", result.ExpireTime)

	d.Set("region", result.Region)

	eips := make([]map[string]interface{}, 0, len(result.Eips))

	for _, e := range result.Eips {

		eips = append(eips, map[string]interface{}{
			"eip":               e.Eip,
			"name":              e.Name,
			"status":            e.Status,
			"eip_instance_type": e.EipInstanceType,
			"share_group_id":    e.ShareGroupId,
			"bandwidth_in_mbps": e.BandWidthInMbps,
			"payment_timing":    e.PaymentTiming,
			"billing_method":    e.BillingMethod,
			"create_time":       e.CreateTime,
			"expire_time":       e.ExpireTime,
			"tags":              flattenTagsToMap(e.Tags),
		})

	}

	d.Set("eips", eips)

	return nil
}

func resourceBaiduCloudEipgroupUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.BaiduClient)

	eipGroupid := d.Id()

	action := "Update EIP Group" + eipGroupid

	if d.HasChange("bandwidth_in_mbps") {

		bandWidthRequest := &eip.ResizeEipGroupArgs{}

		bandWidthRequest.BandWidthInMbps = d.Get("bandwidth_in_mbps").(int)

		bandWidthRequest.ClientToken = buildClientToken()

		_, resizeErr := client.WithEipClient(func(eipClient *eip.Client) (interface{}, error) {
			return nil, eipClient.ResizeEipGroupBandWidth(eipGroupid, bandWidthRequest)
		})

		if resizeErr != nil {
			if IsExceptedErrors(resizeErr, ObjectNotFound) {
				return nil
			}
			return WrapErrorf(resizeErr, DefaultErrorMsg, "baiducloud_eipgroup", action, BCESDKGoERROR)

		}
	}

	if d.HasChange("name") {

		nameRequest := &eip.RenameEipGroupArgs{}

		nameRequest.Name = d.Get("name").(string)

		nameRequest.ClientToken = buildClientToken()

		_, renameErr := client.WithEipClient(func(eipClient *eip.Client) (interface{}, error) {
			return nil, eipClient.RenameEipGroup(eipGroupid, nameRequest)
		})

		if renameErr != nil {
			if IsExceptedErrors(renameErr, ObjectNotFound) {
				return nil
			}
			return WrapErrorf(renameErr, DefaultErrorMsg, "baiducloud_eipgroup", action, BCESDKGoERROR)

		}
	}

	return resourceBaiduCloudEipgroupRead(d, meta)
}

func resourceBaiduCloudEipgroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.BaiduClient)

	eipGroupId := d.Id()

	action := "Delete EIP GROUP ID IS " + eipGroupId

	_, err := client.WithEipClient(func(eipClient *eip.Client) (interface{}, error) {
		return nil, eipClient.DeleteEipGroup(eipGroupId, buildClientToken())
	})

	if err != nil {
		if IsExceptedErrors(err, ObjectNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eipgroup", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudCreateEipgroupArgs(d *schema.ResourceData) *eip.CreateEipGroupArgs {

	request := &eip.CreateEipGroupArgs{}

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		request.Name = v.(string)
	}

	if v := d.Get("eip_count").(int); v != 0 {
		request.EipCount = v
	}

	if v := d.Get("bandwidth_in_mbps").(int); v != 0 {
		request.BandWidthInMbps = v
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
	}

	if v, ok := d.GetOk("tags"); ok {
		request.Tags = tranceTagMapToModel(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("route_type"); ok && v.(string) != "" {
		request.RouteType = v.(string)
	}

	if v, ok := d.GetOk("idc"); ok && v.(string) != "" {
		request.Idc = v.(string)
	}

	if v, ok := d.GetOk("continuous"); ok {
		request.Continuous = v.(bool)
	}

	request.ClientToken = buildClientToken()

	return request
}

func eipGroupStateRefreshFunc(client *connectivity.BaiduClient, id string, failState []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		raw, err := client.WithEipClient(func(eipClient *eip.Client) (interface{}, error) {
			return eipClient.EipGroupDetail(id)
		})

		if err != nil {
			return nil, "", WrapError(err)
		}
		result, _ := raw.(*eip.EipGroupModel)

		for _, status := range failState {
			if result.Status == status {
				return result, result.Status, WrapError(Error(GetFailTargetStatus, result.Status))
			}
		}

		return result, result.Status, nil
	}
}
