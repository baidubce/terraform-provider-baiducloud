package eip

import (
	"fmt"

	"github.com/baidubce/bce-sdk-go/services/eip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/flex"
)

func ResourceEipGroupDetachment() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to detach one or more EIPs from an EIP Group. \n\n" +
			"More information can be found in the [Developer Guide](https://cloud.baidu.com/doc/EIP/s/Qkoslycn3). \n\n",

		Create: resourceEipGroupDetachmentCreate,
		Read:   flex.DoNothing,
		Delete: flex.DoNothing,

		Schema: map[string]*schema.Schema{
			"eip_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the EIP Group.",
			},
			"move_out_eips": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MinItems:    1,
				Elem:        MoveOutEipSchema(),
				Description: "The list of EIPs to be detached, including both IPv4 and IPv6 addresses.",
			},
		},
	}
}

func MoveOutEipSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"eip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The EIP address to be detached",
			},
			"bandwidth_in_mbps": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The bandwidth value (in Mbps) for the EIP after detachment. Required when detaching an EIP that was originally created through the EIP Group.",
			},
			"payment_timing": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{flex.PaymentTimingPostpaid}, false),
				Description:  "Payment timing of billing. Valid value: `Postpaid`. Required when detaching an EIP that was originally created through the EIP Group.",
			},
			"billing_method": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ByTraffic", "ByBandwidth"}, false),
				Description:  "The billing method for the EIP. Valid values: `ByTraffic`, `ByBandwidth`. Required when detaching an EIP that was originally created through the EIP Group.",
			},
		},
	}

}

func resourceEipGroupDetachmentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	eipGroupID := d.Get("eip_group_id").(string)

	_, err := conn.WithEipClient(func(client *eip.Client) (interface{}, error) {
		args := eip.EipGroupMoveOutArgs{
			MoveOutEips: expandMoveOutEips(d.Get("move_out_eips").([]interface{})),
		}
		return nil, client.EipGroupMoveOut(eipGroupID, &args)
	})
	if err != nil {
		return fmt.Errorf("error detaching eips from eip group (%s): %w", eipGroupID, err)
	}
	d.SetId(resource.UniqueId())

	return nil
}
