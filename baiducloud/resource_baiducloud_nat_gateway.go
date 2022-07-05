/*
Provide a resource to create a NAT Gateway.

Example Usage

```hcl
resource "baiducloud_nat_gateway" "default" {
  name = "terraform-nat-gateway"
  vpc_id = "vpc-ggm7drdgyvha"
  spec = "medium"
  billing = {
    payment_timing = "Postpaid"
  }
}
```

Import

NAT Gateway instance can be imported, e.g.

```hcl
$ terraform import baiducloud_nat_gateway.default nat_gateway_id
```
*/
package baiducloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudNatGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudNatGatewayCreate,
		Read:   resourceBaiduCloudNatGatewayRead,
		Update: resourceBaiduCloudNatGatewayUpdate,
		Delete: resourceBaiduCloudNatGatewayDelete,

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
				Description: "Name of the NAT gateway, consisting of uppercase and lowercase letters、numbers and special characters, such as \"-\",\"_\",\"/\",\".\". The value must start with a letter, and the length should between 1-65.",
				Required:    true,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "VPC ID of the NAT gateway.",
				Required:    true,
				ForceNew:    true,
			},
			"spec": {
				Type:         schema.TypeString,
				Description:  "Specification of the NAT gateway, available values are small(supports up to 5 public IPs), medium(up to 10 public IPs) and large(up to 15 public IPs). Default to small.",
				Optional:     true,
				ForceNew:     true,
				Default:      "small",
				ValidateFunc: validation.StringInSlice([]string{"small", "medium", "large"}, false),
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Status of the NAT gateway.",
				Computed:    true,
			},
			"expired_time": {
				Type:        schema.TypeString,
				Description: "Expired time of the NAT gateway, which will be empty when the payment_timing is Postpaid.",
				Computed:    true,
			},
			"eips": {
				Type:        schema.TypeSet,
				Description: "One public network EIP associated with the NAT gateway or one or more EIPs in the shared bandwidth.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"billing": {
				Type:        schema.TypeMap,
				Description: "Billing information of the NAT gateway.",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"payment_timing": {
							Type:         schema.TypeString,
							Description:  "Payment timing of the billing, which can be Prepaid or Postpaid. The default is Postpaid.",
							Required:     true,
							ForceNew:     true,
							Default:      PaymentTimingPostpaid,
							ValidateFunc: validatePaymentTiming(),
						},
						"reservation": {
							Type:             schema.TypeMap,
							Description:      "Reservation of the NAT gateway.",
							Optional:         true,
							DiffSuppressFunc: postPaidDiffSuppressFunc,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"reservation_length": {
										Type:             schema.TypeInt,
										Description:      "Reservation length that you will pay for your resource. It is valid when payment_timing is Prepaid. Valid values: [1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36].",
										Optional:         true,
										Default:          1,
										ForceNew:         true,
										DiffSuppressFunc: postPaidDiffSuppressFunc,
										ValidateFunc:     validateReservationLength(),
									},
									"reservation_time_unit": {
										Type:             schema.TypeString,
										Description:      "Reservation time unit that you will pay for your resource. It is valid when payment_timing is Prepaid. The value can only be month currently, which is also the default value.",
										Optional:         true,
										Default:          "month",
										ValidateFunc:     validateReservationUnit(),
										DiffSuppressFunc: postPaidDiffSuppressFunc,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceBaiduCloudNatGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	vpcService := VpcService{client}

	args := buildBaiduCloudNatGatewayArgs(d)
	action := "Create NAT Gateway " + args.Name

	if err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return vpcClient.CreateNatGateway(args)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)
		result, _ := raw.(*vpc.CreateNatGatewayResult)
		d.SetId(result.NatId)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_nat_gateway", action, BCESDKGoERROR)
	}

	stateConf := buildStateConf(
		[]string{string(vpc.NAT_STATUS_BUILDING), string(vpc.NAT_STATUS_CONFIGURING)},
		[]string{string(vpc.NAT_STATUS_ACTIVE)},
		d.Timeout(schema.TimeoutCreate),
		vpcService.NatGatewayStateRefresh(d.Id()),
	)
	if args.Eips == nil || len(args.Eips) == 0 {
		stateConf.Target = []string{string(vpc.NAT_STATUS_UNCONFIGURED)}
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_nat_gateway", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudNatGatewayRead(d, meta)
}

func resourceBaiduCloudNatGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	vpcService := VpcService{client}

	natId := d.Id()
	action := "Query NAT Gateway " + natId

	result, state, err := vpcService.NatGatewayStateRefresh(natId)()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_nat_gateway", action, BCESDKGoERROR)
	}

	status := map[string]bool{
		string(vpc.NAT_STATUS_DELETING): true,
		string(vpc.NAT_STATUS_DELETED):  true,
		string(vpc.NAT_STATUS_DOWN):     true,
		string(vpc.NAT_STATUS_STOPPING): true,
	}
	if _, ok := status[strings.ToLower(state)]; ok || result == nil {
		d.SetId("")
		return nil
	}

	nat := result.(*vpc.NAT)
	d.Set("name", nat.Name)
	d.Set("vpc_id", nat.VpcId)
	d.Set("spec", nat.Spec)
	d.Set("eips", nat.Eips)

	billingMap := map[string]interface{}{"payment_timing": nat.PaymentTiming}
	d.Set("billing", billingMap)

	d.Set("expired_time", nat.ExpiredTime)
	d.Set("status", nat.Status)

	return nil
}

func resourceBaiduCloudNatGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	natId := d.Id()
	action := "Update NAT Gateway " + natId

	d.Partial(true)
	if d.HasChange("name") {
		args := &vpc.UpdateNatGatewayArgs{}
		if v := d.Get("name").(string); v != "" {
			args.Name = v
		} else {
			return WrapErrorf(fmt.Errorf("The name cannot be changed to empty."), DefaultErrorMsg, "baiducloud_nat_gateway", action, BCESDKGoERROR)
		}

		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return nil, vpcClient.UpdateNatGateway(natId, args)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_nat_gateway", action, BCESDKGoERROR)
		}
		d.SetPartial("name")
	}

	d.Partial(false)

	return resourceBaiduCloudNatGatewayRead(d, meta)
}

func resourceBaiduCloudNatGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	vpcService := VpcService{client}

	natId := d.Id()
	action := "Delete NAT Gateway " + natId

	clientToken := buildClientToken()
	_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
		return nil, vpcClient.DeleteNatGateway(natId, clientToken)
	})
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_nat_gateway", action, BCESDKGoERROR)
	}
	addDebug(action, nil)

	stateConf := buildStateConf(
		[]string{string(vpc.NAT_STATUS_DELETING), string(vpc.NAT_STATUS_ACTIVE), string(vpc.NAT_STATUS_UNCONFIGURED)},
		[]string{string(vpc.NAT_STATUS_DELETED)},
		d.Timeout(schema.TimeoutDelete),
		vpcService.NatGatewayStateRefresh(natId),
	)
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_nat_gateway", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudNatGatewayArgs(d *schema.ResourceData) *vpc.CreateNatGatewayArgs {
	args := &vpc.CreateNatGatewayArgs{
		ClientToken: buildClientToken(),
		Billing:     &vpc.Billing{},
	}

	if v := d.Get("name").(string); v != "" {
		args.Name = v
	}
	if v := d.Get("vpc_id").(string); v != "" {
		args.VpcId = v
	}
	if v := d.Get("spec").(string); v != "" {
		args.Spec = vpc.NatGatewaySpecType(v)
	}

	if v, ok := d.GetOk("billing"); ok {
		billing := v.(map[string]interface{})
		if p, ok := billing["payment_timing"]; ok {
			paymentTiming := vpc.PaymentTimingType(p.(string))
			args.Billing.PaymentTiming = paymentTiming
		}
		if args.Billing.PaymentTiming == PaymentTimingPrepaid {
			if r, ok := billing["reservation"]; ok {
				args.Billing.Reservation = &vpc.Reservation{}
				reservation := r.(map[string]interface{})
				if reservationLength, ok := reservation["reservation_length"]; ok {
					args.Billing.Reservation.ReservationLength = reservationLength.(int)
				}
				if reservationTimeUnit, ok := reservation["reservation_time_unit"]; ok {
					args.Billing.Reservation.ReservationTimeUnit = reservationTimeUnit.(string)
				}
			}
		}
	}

	return args
}
