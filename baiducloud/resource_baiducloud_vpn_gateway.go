/*
Provide a resource to create a VPN gateway.

Example Usage

```hcl
resource "baiducloud_vpn_gateway" "default" {
  vpn_name       = "test_vpn_gateway"
  vpc_id         = "vpc-65cz3hu92kz2"
  description    = "test desc"
  payment_timing = "Postpaid"
}
```

Import

VPN gateway can be imported, e.g.

```hcl
$ terraform import baiducloud_vpn_gateway.default vpn_gateway_id
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/vpn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"strconv"
	"time"
)

func resourceBaiduCloudVpnGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudVpnGatewayCreate,
		Read:   resourceBaiduCloudVpnGatewayRead,
		Update: resourceBaiduCloudVpnGatewayUpdate,
		Delete: resourceBaiduCloudVpnGatewayDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "ID of the VPC which vpn gateway belong to.",
				Required:    true,
			},
			"vpn_name": {
				Type:        schema.TypeString,
				Description: "Name of the VPN gateway, which cannot take the value \"default\", the length is no more than 65 characters, and the value can be composed of numbers, characters and underscores.",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the VPN. The value is no more than 200 characters.",
				Optional:    true,
			},
			"eip": {
				Type:        schema.TypeString,
				Description: "Eip address.",
				Optional:    true,
			},
			"payment_timing": {
				Type:         schema.TypeString,
				Description:  "Payment timing of billing, which can be Prepaid or Postpaid. The default is Postpaid.",
				Required:     true,
				ValidateFunc: validatePaymentTiming(),
			},
			"reservation": {
				Type:             schema.TypeMap,
				Description:      "Reservation of the VPN gateway.",
				Optional:         true,
				DiffSuppressFunc: postPaidDiffSuppressFunc,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reservation_length": {
							Type:             schema.TypeInt,
							Description:      "The reservation length that you will pay for your resource. It is valid when payment_timing is Prepaid. Valid values: [1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36].",
							Required:         true,
							Default:          1,
							ValidateFunc:     validateReservationLength(),
							DiffSuppressFunc: postPaidDiffSuppressFunc,
						},
						"reservation_time_unit": {
							Type:             schema.TypeString,
							Description:      "The reservation time unit that you will pay for your resource. It is valid when payment_timing is Prepaid. The value can only be month currently, which is also the default value.",
							Required:         true,
							Default:          "Month",
							ValidateFunc:     validateReservationUnit(),
							DiffSuppressFunc: postPaidDiffSuppressFunc,
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeString,
				Description: "VPN gateway status.",
				Computed:    true,
			},
			"expired_time": {
				Type:        schema.TypeString,
				Description: "Expired time of VPN gateway.",
				Computed:    true,
			},
			"bandwidth_in_mbps": {
				Type:        schema.TypeString,
				Description: "Eip bandwidth(Mbps), if payment_timing is Prepaid or billing_method is ByBandWidth, support between 1 and 200, if billing_method is ByTraffic, support between 1 and 1000",
				Computed:    true,
			},
			"vpn_conn_num": {
				Type:        schema.TypeInt,
				Description: "Number of VPN tunnels.",
				Computed:    true,
			},
			"vpn_conns": {
				Type:        schema.TypeList,
				Description: "ID list of VPN tunnels.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}
func resourceBaiduCloudVpnGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	createVpnGatewayArgs, err := buildBaiduCloudVpnGatewayArgs(d, meta)
	vpnService := VpnService{client}

	action := "Create VPN Gateway " + createVpnGatewayArgs.VpnName

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_vpn", action, BCESDKGoERROR)
	}
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithVPNClient(func(vpnClient *vpn.Client) (i interface{}, e error) {
			return vpnClient.CreateVpnGateway(createVpnGatewayArgs)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)
		result, _ := raw.(*vpn.CreateVpnGatewayResult)
		d.SetId(result.VpnId)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_vpn", action, BCESDKGoERROR)
	}
	stateConf := buildStateConf(
		[]string{string(vpn.VPN_STATUS_BUILDING)},
		[]string{string(vpn.VPN_STATUS_UNCONFIGURED), string(vpn.VPN_STATUS_ACTIVE)},
		d.Timeout(schema.TimeoutCreate),
		vpnService.VpnGatewayStateRefresh(d.Id()),
	)
	if _, err = stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instance", action, BCESDKGoERROR)
	}
	return resourceBaiduCloudVpnGatewayRead(d, meta)
}
func resourceBaiduCloudVpnGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	vpnService := VpnService{client}

	vpnGatewayID := d.Id()
	action := "Query vpn gateway " + vpnGatewayID
	vpnRes, err := vpnService.VpnGatewayDetail(vpnGatewayID)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_vpn_gateway", action, BCESDKGoERROR)
	}
	d.Set("vpc_id", vpnRes.VpcId)
	d.Set("status", vpnRes.Status)
	d.Set("vpn_name", vpnRes.Name)
	d.Set("description", vpnRes.Description)
	d.Set("vpn_conn_num", vpnRes.VpnConnNum)
	d.Set("bandwidth_in_mbps", vpnRes.BandwidthInMbps)
	d.Set("eip", vpnRes.Eip)
	d.Set("expired_time", vpnRes.ExpiredTime)

	conns := make([]string, len(vpnRes.VpnConns))
	for _, conn := range vpnRes.VpnConns {
		conns = append(conns, conn.VpnConnId)
	}
	d.Set("vpn_conns", conns)

	return nil
}
func resourceBaiduCloudVpnGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	action := "Update VPN gateway attribute "
	client := meta.(*connectivity.BaiduClient)
	args := &vpn.UpdateVpnGatewayArgs{
		ClientToken: buildClientToken(),
	}
	if d.HasChange("vpn_name") {
		args.Name = d.Get("vpn_name").(string)
	}
	err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		_, err := client.WithVPNClient(func(vpnClient *vpn.Client) (interface{}, error) {
			err := vpnClient.UpdateVpnGateway(d.Id(), args)
			if err != nil {
				return nil, err
			}
			if d.HasChange("eip") {
				eip := d.Get("eip")
				if len(eip.(string)) == 0 {
					// 解绑
					err = vpnClient.UnBindEip(d.Id(), buildClientToken())
				} else {
					// 换绑
					eipArgs := &vpn.BindEipArgs{
						ClientToken: buildClientToken(),
						Eip:         eip.(string),
					}
					err = vpnClient.BindEip(d.Id(), eipArgs)
				}
			}
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
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_vpn_gateway", action, BCESDKGoERROR)
	}
	return resourceBaiduCloudVpnGatewayRead(d, meta)
}
func resourceBaiduCloudVpnGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	vpnGatewayId := d.Id()
	action := "Delete VPN gateway " + vpnGatewayId

	// Delete VPN Gateway
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, _ := client.WithVPNClient(func(vpnClient *vpn.Client) (interface{}, error) {
			return vpnGatewayId, vpnClient.DeleteVpnGateway(vpnGatewayId, buildClientToken())
		})
		addDebug(action, raw)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_vpn_gateway", action, BCESDKGoERROR)
	}
	return nil
}

func buildBaiduCloudVpnGatewayArgs(d *schema.ResourceData, meta interface{}) (*vpn.CreateVpnGatewayArgs, error) {
	res := &vpn.CreateVpnGatewayArgs{
		ClientToken: buildClientToken(),
	}
	if vpcID, ok := d.GetOk("vpc_id"); ok {
		res.VpcId = vpcID.(string)
	}
	if vpnName, ok := d.GetOk("vpn_name"); ok {
		res.VpnName = vpnName.(string)
	}
	if desc, ok := d.GetOk("description"); ok {
		res.Description = desc.(string)
	}
	// build billing
	billingRequest := &vpn.Billing{
		PaymentTiming: vpn.PaymentTimingType(""),
		Reservation:   nil,
	}
	if p, ok := d.GetOk("payment_timing"); ok {
		paymentTiming := vpn.PaymentTimingType(p.(string))
		billingRequest.PaymentTiming = paymentTiming
	}
	if p, ok := d.GetOk("eip"); ok {
		res.Eip = p.(string)
	}
	if billingRequest.PaymentTiming == vpn.PAYMENT_TIMING_PREPAID {
		if r, ok := d.GetOk("reservation"); ok {
			reservation := r.(map[string]interface{})
			if reservationLength, ok := reservation["reservation_length"]; ok {
				reservationLengthInt, err := strconv.Atoi(reservationLength.(string))
				billingRequest.Reservation.ReservationLength = reservationLengthInt
				if err != nil {
					return nil, err
				}
			}
			if reservationTimeUnit, ok := reservation["reservation_time_unit"]; ok {
				billingRequest.Reservation.ReservationTimeUnit = reservationTimeUnit.(string)
			}
		}
	}
	res.Billing = billingRequest
	return res, nil
}
