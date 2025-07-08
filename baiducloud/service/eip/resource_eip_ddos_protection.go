package eip

import (
	"fmt"

	"github.com/baidubce/bce-sdk-go/services/eip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func ResourceEipDDosProtection() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to manage the basic DDoS protection threshold of a specified public IP. \n\n" +
			"More information can be found in the [Developer Guide](https://cloud.baidu.com/doc/EIP/s/alhagbhi0). \n\n",

		Create: resourceEipDDosProtectionCreate,
		Read:   resourceEipDDosProtectionRead,
		Update: resourceEipDDosProtectionUpdate,
		Delete: schema.RemoveFromState,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The public IP address to apply the DDoS protection threshold to.",
			},
			"threshold_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"bandwidth", "auto", "manual"}, false),
				Description:  "Threshold mode for basic DDoS protection. Valid values: `bandwidth` (based on bandwidth cap), `auto` (intelligent threshold), and `manual` (manual setting).",
			},
			"ip_clean_mbps": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				Description:      "Cleaning threshold for traffic in Mbps. Required when `threshold_type` is `manual`. Valid range: 120–5000 Mbps.",
				ValidateFunc:     validation.IntBetween(120, 5000),
				DiffSuppressFunc: suppressWhenNotManual,
			},
			"ip_clean_pps": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				Description:      "Cleaning threshold for packets per second (PPS). Required when `threshold_type` is `manual`. Valid range: 58594–4882813 PPS.",
				ValidateFunc:     validation.IntBetween(58594, 4882813),
				DiffSuppressFunc: suppressWhenNotManual,
			},
		},

		CustomizeDiff: func(diff *schema.ResourceDiff, i interface{}) error {
			if diff.Get("threshold_type") == "manual" {
				if diff.Get("ip_clean_mbps") == 0 || diff.Get("ip_clean_pps") == 0 {
					return fmt.Errorf("ip_clean_mbps and ip_clean_pps must be set when threshold_type is manual")
				}
			}
			return nil
		},
	}
}

func suppressWhenNotManual(k, old, new string, d *schema.ResourceData) bool {
	return d.Get("threshold_type").(string) != "manual"
}

func resourceEipDDosProtectionCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	err := modifyEipDDosProtection(d, conn)
	if err != nil {
		return err
	}

	d.SetId(d.Get("ip").(string))

	return resourceEipDDosProtectionRead(d, meta)
}

func resourceEipDDosProtectionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	detail, err := FindEipDDosProtection(conn, d.Id())
	if err != nil {
		return err
	}

	if detail == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("ip", detail.Ip); err != nil {
		return fmt.Errorf("error setting ip: %w", err)
	}
	if err := d.Set("threshold_type", detail.ThresholdType); err != nil {
		return fmt.Errorf("error setting threshold_type: %w", err)
	}
	if err := d.Set("ip_clean_mbps", detail.IpCleanMbps); err != nil {
		return fmt.Errorf("error setting ip_clean_mbps: %w", err)
	}
	if err := d.Set("ip_clean_pps", detail.IpCleanPps); err != nil {
		return fmt.Errorf("error setting ip_clean_pps: %w", err)
	}

	return nil
}

func resourceEipDDosProtectionUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	if d.HasChanges("threshold_type", "ip_clean_mbps", "ip_clean_pps") {
		err := modifyEipDDosProtection(d, conn)
		if err != nil {
			return err
		}
	}

	return resourceEipDDosProtectionRead(d, meta)
}

func modifyEipDDosProtection(d *schema.ResourceData, conn *connectivity.BaiduClient) error {
	if err := checkEip(conn, d.Get("ip").(string)); err != nil {
		return err
	}

	args := eip.ModifyDdosThresholdRequest{
		Ip:            d.Get("ip").(string),
		ThresholdType: d.Get("threshold_type").(string),
	}
	if args.ThresholdType == "manual" {
		args.IpCleanMbps = int64(d.Get("ip_clean_mbps").(int))
		args.IpCleanPps = int64(d.Get("ip_clean_pps").(int))
	}

	_, err := conn.WithEipClient(func(client *eip.Client) (interface{}, error) {
		return nil, client.ModifyDdosThreshold(&args)
	})
	if err != nil {
		return fmt.Errorf("error modifying DDoS protection threshold for eip (%s): %w", args.Ip, err)
	}
	return nil
}

func checkEip(conn *connectivity.BaiduClient, ip string) error {
	detail, err := FindEip(conn, ip)
	if err != nil {
		return err
	}
	if detail == nil {
		return fmt.Errorf("eip (%s) not found", ip)
	}
	return nil
}
