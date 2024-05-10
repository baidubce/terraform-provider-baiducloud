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
	"github.com/baidubce/bce-sdk-go/services/resmanager"
	"strconv"
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
				Optional:    true,
			},
			"ipv6_address": {
				Type:        schema.TypeString,
				Description: "LoadBalance instance's ipv6 ip address",
				Computed:    true,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "The VPC short ID to which the LoadBalance instance belongs",
				Required:    true,
				ForceNew:    true,
			},
			"eip": {
				Type:        schema.TypeString,
				Description: "eip of the LoadBalance",
				Optional:    true,
			},
			"payment_timing": {
				Type:         schema.TypeString,
				Description:  "Payment timing of billing, which can be Prepaid or Postpaid. The default is Postpaid.",
				Optional:     true,
				ForceNew:     true,
				Default:      PaymentTimingPostpaid,
				ValidateFunc: validatePaymentTiming(),
			},
			"reservation": {
				Type:             schema.TypeMap,
				Description:      "Reservation of the BLB.",
				Optional:         true,
				DiffSuppressFunc: postPaidDiffSuppressFunc,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reservation_length": {
							Type: schema.TypeInt,
							Description: "The reservation length that you will pay for your resource. " +
								"It is valid when payment_timing is Prepaid. " +
								"Valid values: [1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36].",
							Required:         true,
							Default:          1,
							ValidateFunc:     validateReservationLength(),
							DiffSuppressFunc: postPaidDiffSuppressFunc,
						},
						"reservation_time_unit": {
							Type: schema.TypeString,
							Description: "The reservation time unit that you will pay for your resource. " +
								"It is valid when payment_timing is Prepaid. " +
								"The value can only be month currently, which is also the default value.",
							Required:         true,
							Default:          "month",
							ValidateFunc:     validateReservationUnit(),
							DiffSuppressFunc: postPaidDiffSuppressFunc,
						},
					},
				},
			},
			"auto_renew_length": {
				Type:         schema.TypeInt,
				Description:  "The automatic renewal time is 1-9 per month and 1-3 per year.",
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 9),
			},
			"auto_renew_time_unit": {
				Type:         schema.TypeString,
				Description:  "Monthly payment or annual payment, month is \"month\" and year is \"year\".",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"month", "year"}, false),
			},
			"security_groups": {
				Type:        schema.TypeSet,
				Description: "security groups of the LoadBalance.",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enterprise_security_groups": {
				Type:        schema.TypeSet,
				Description: "enterprise security group ids of the LoadBalance",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"performance_level": {
				Type:         schema.TypeString,
				Description:  "performance level, available values are small1, small2, medium1, medium2, large1, large2, large3",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"small1", "small2", "medium1", "medium2", "large1", "large2", "large3",}, false),
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
			"allow_delete": {
				Type:        schema.TypeBool,
				Default:     true,
				Description: "Whether to allow deletion, default value is true. ",
				Optional:    true,
			},
			"allocate_ipv6": {
				Type:        schema.TypeBool,
				Default:     false,
				Description: "Whether to allocated ipv6, default value is false, do not support modify",
				Optional:    true,
				ForceNew:    true,
			},
			"tags": tagsSchema(),
			"resource_group_id": {
				Type:        schema.TypeString,
				Description: "Resource group id, support setting when creating instance, do not support modify!",
				Optional:    true,
				ForceNew:    true,
			},
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
		addDebug(action, createArgs)
		response, _ := raw.(*blb.CreateLoadBalancerResult)
		d.SetId(response.BlbId)
		d.Set("address", response.Address)
		d.Set("ipv6_address", response.Ipv6)

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

	if _, ok := d.GetOk("security_groups"); ok {
		err := blbService.updateBlbSecurityGroups(d, meta)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb", action, BCESDKGoERROR)
		}
	}

	if _, ok := d.GetOk("enterprise_security_groups"); ok {
		err := blbService.updateBlbEnterpriseSecurityGroups(d, meta)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb", action, BCESDKGoERROR)
		}
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
	d.Set("performance_level", blbDetail.PerformanceLevel)
	d.Set("listener", blbService.FlattenListenerModelToMap(blbDetail.Listener))
	d.Set("payment_timing", blbDetail.PaymentTiming)
	d.Set("allow_delete", blbModel.AllowDelete)
	if d.HasChange("tags") {
		if v, ok := d.GetOk("tags"); ok {
			if !slicesContainSameElements(blbDetail.Tags, tranceTagMapToModel(v.(map[string]interface{}))) {
				return WrapErrorf(Error("Tags bind failed."), DefaultErrorMsg, "baiducloud_blb", action, BCESDKGoERROR)
			}
		}
	}
	resourceGroupId, err := getResourceGroup(d, client, action, "baiducloud_blb")
	if err != nil {
		return err
	}
	if d.HasChange("resource_group_id") {
		if v, ok := d.GetOk("resource_group_id"); ok {
			if resourceGroupId != v.(string) {
				return WrapErrorf(Error("Resource group bind failed."), DefaultErrorMsg, "baiducloud_blb", action, BCESDKGoERROR)
			}
		}
	}
	d.Set("resource_group_id", resourceGroupId)
	d.Set("tags", flattenTagsToMap(blbModel.Tags))
	securityIds, err := blbService.getBlbSecurityGroupIds(d.Id(), meta)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb", action, BCESDKGoERROR)
	}
	d.Set("security_groups", securityIds)

	enterpriseSecurityIds, err := blbService.getBlbEnterpriseSecurityGroupIds(d.Id(), meta)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb", action, BCESDKGoERROR)
	}
	d.Set("enterprise_security_groups", enterpriseSecurityIds)

	return nil
}

func getResourceGroup(d *schema.ResourceData, client *connectivity.BaiduClient, action string, product string) (string, error) {
	region := string(client.Region)
	raw, err := client.WithResourceManagerClient(func(client *resmanager.Client) (i interface{}, e error) {
		args := &resmanager.ResGroupDetailRequest{
			ResourceBrief: []resmanager.ResourceBrief{
				{
					ResourceId:     d.Id(),
					ResourceType:   "BLB",
					ResourceRegion: region,
				},
			},
		}
		return client.GetResGroupBatch(args)
	})
	if err != nil {
		return "",WrapErrorf(err, DefaultErrorMsg, product, action, BCESDKGoERROR)
	}
	resp := raw.(*resmanager.ResGroupDetailResponse)
	if len(resp.ResourceGroupsDetailFull)>0 && len(resp.ResourceGroupsDetailFull[0].BindGroupInfo)> 0 {
		return resp.ResourceGroupsDetailFull[0].BindGroupInfo[0].GroupId, nil
	}
	return "", nil
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

	if d.HasChange("allow_delete") {
		update = true
		allowDelete := d.Get("allow_delete").(bool)
		updateArgs.AllowDelete = &allowDelete
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
	if d.HasChange("security_groups") {
		err := blbService.updateBlbSecurityGroups(d, meta)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb", action, BCESDKGoERROR)
		}
	}

	if d.HasChange("enterprise_security_groups") {
		err := blbService.updateBlbEnterpriseSecurityGroups(d, meta)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb", action, BCESDKGoERROR)
		}
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

	if v, ok := d.GetOk("eip"); ok {
		result.Eip = v.(string)
	}

	if v, ok := d.GetOk("address"); ok {
		result.Address = v.(string)
	}

	if v, ok := d.GetOk("performance_level"); ok {
		result.PerformanceLevel = v.(string)
	}

	if v, ok := d.GetOk("type"); ok {
		result.Type = v.(string)
	}

	if v, ok := d.GetOk("payment_timing"); ok {
		billingRequest := &blb.Billing{
			PaymentTiming: "",
		}
		paymentTiming := v.(string)
		billingRequest.PaymentTiming = paymentTiming
		if billingRequest.PaymentTiming == PaymentTimingPrepaid {
			if r, ok := d.GetOk("reservation"); ok {
				reservation := r.(map[string]interface{})
				billingRequest.Reservation = &blb.Reservation{}
				if reservationLength, ok := reservation["reservation_length"]; ok {
					reservationLengthInt, _ := strconv.Atoi(reservationLength.(string))
					billingRequest.Reservation.ReservationLength = reservationLengthInt
				}
				if reservationTimeUnit, ok := reservation["reservation_time_unit"]; ok {
					billingRequest.Reservation.ReservationTimeUnit = reservationTimeUnit.(string)
				}
			}
			if v, ok := d.GetOk("auto_renew_length"); ok {
				result.AutoRenewLength = v.(int)
				if result.AutoRenewLength > 0 {
					if v, ok := d.GetOk("auto_renew_time_unit"); ok {
						result.AutoRenewTimeUnit = v.(string)
					}
				}
			}
		}
		result.Billing = billingRequest
	}

	if v := d.Get("allow_delete"); true{
		allowDelete := v.(bool)
		result.AllowDelete = &allowDelete
	}

	if v := d.Get("allocate_ipv6"); true {
		allocateIpv6 := v.(bool)
		result.AllocateIpv6 = &allocateIpv6
	}

	if v := d.Get("resource_group_id"); true {
		result.ResourceGroupId = v.(string)
	}

	return result
}
