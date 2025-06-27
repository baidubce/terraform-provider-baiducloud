package eip

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/baidubce/bce-sdk-go/model"
	"github.com/baidubce/bce-sdk-go/services/eip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/flex"
)

func ResourceEipGroup() *schema.Resource {
	return &schema.Resource{

		Description: "Use this resource to manage an EIP Group. \n\n" +
			"More information can be found in the [Developer Guide](https://cloud.baidu.com/doc/EIP/s/ijwvz2zq8). \n\n",

		Create: resourceEipGroupCreate,
		Read:   resourceEipGroupRead,
		Update: resourceEipGroupUpdate,
		Delete: resourceEipGroupDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"payment_timing": flex.SchemaPaymentTiming(),
			"billing_method": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ByTraffic", "ByBandwidth", "ByPeak95"}, false),
				Description:  "The billing method for the EIP. Valid values: `ByTraffic`, `ByBandwidth`, `ByPeak95`.",
			},
			"reservation_length":    flex.SchemaReservationLength(),
			"reservation_time_unit": flex.SchemaReservationTimeUnit(),
			"tags":                  flex.TagsSchema(),
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 65),
				Description: "The name of the EIP Group. 1–65 characters. Must start with a letter and may contain letters, digits, " +
					"hyphens, underscores, dots, or slashes. Auto-generated if not set.",
			},
			"route_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "BGP",
				Description: "BGP line type. Valid values: `BGP` (standard BGP), `BGP_S` (enhanced BGP). Defaults to `BGP`.",
			},
			"eip_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of IPv4 EIPs in the EIP Group.",
			},
			"eipv6_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of IPv6 EIPs in the EIP Group.",
			},
			"bandwidth_in_mbps": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Total bandwidth in Mbps. Standard BGP supports 20–500, enhanced BGP supports 100–5000.",
			},
			"idc": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Single Line ID to be used by the EIP Group.",
			},

			// computed fields
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the EIP Group. Possible values: `creating`, `available`, `binded`, `binding`, `unbinding`, `updating`, `paused`, `unavailable`.",
			},
			"default_domestic_bandwidth": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The default bandwidth for cross-border acceleration, in Mbps.",
			},
			"bw_short_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the associated bandwidth package.",
			},
			"bw_bandwidth_in_mbps": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The bandwidth of the associated bandwidth package, in Mbps.",
			},
			"domestic_bw_short_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the associated cross-border acceleration package.",
			},
			"domestic_bw_bandwidth_in_mbps": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The bandwidth of the cross-border acceleration package, in Mbps.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the EIP Group.",
			},
			"expire_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The expiration time of the EIP Group. Only available for prepaid resources.",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The region of the EIP Group",
			},
			"eips": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        eipSchema(),
				Description: "A list of IPv4 EIP instances attached to the EIP Group.",
			},
			"eipv6s": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        eipSchema(),
				Description: "A list of IPv6 EIP instances attached to the EIP Group.",
			},
			// deprecated fields
			"group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This field is deprecated and will be removed in a future release. Please use `id` instead.",
				Deprecated:  "This field is deprecated and will be removed in a future release. Please use `id` instead.",
			},
			"continuous": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "This field is deprecated and will be removed in a future release.",
				Deprecated:  "This field is deprecated and will be removed in a future release.",
			},
		},
	}
}

func eipSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"eip_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of EIP",
			},
			"eip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP address of the EIP",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the EIP.",
			},
			"bandwidth_in_mbps": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The public bandwidth of the EIP, in Mbps.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the EIP. Possible values: `creating`, `available`, `binded`, `binding`, `unbinding`, `updating`, `paused`, `unavailable`.",
			},
			"eip_instance_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of EIP. Possible values: `normal`(regular), `shared`(in an EIP Group).",
			},
			"share_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the EIP Group the EIP is part of. Empty if the EIP is not associated with a EIP Group.",
			},
			"instance_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the resource the EIP is attached to. Empty if the EIP is not bound.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the resource the EIP is attached to. Empty if the EIP is not bound.",
			},

			"payment_timing": flex.ComputedSchemaPaymentTiming(),
			"billing_method": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The billing method. Possible values: `ByTraffic`, `ByBandwidth`.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the EIP.",
			},
			"expire_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The expiration time of the EIP. Available only for prepaid resources.",
			},
			"tags": flex.ComputedSchemaTags(),
		},
	}
}

func resourceEipGroupCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	raw, err := conn.WithEipClient(func(eipClient *eip.Client) (interface{}, error) {
		args := buildCreationArgs(d)
		data, _ := json.Marshal(args)
		log.Printf("[DEBUG] Create EIP Group: %s", string(data))
		return eipClient.CreateEipGroup(args)
	})
	log.Printf("[DEBUG] Create EIP Group result: %+v", raw)

	if err != nil {
		return fmt.Errorf("error creating EIP Group: %w", err)
	}

	response := raw.(*eip.CreateEipGroupResult)
	if response.Id == "" {
		return fmt.Errorf("error creating EIP Group: %+v", raw)
	}

	d.SetId(response.Id)

	if _, err = waitEipGroupAvailable(conn, d.Timeout(schema.TimeoutCreate), d.Id()); err != nil {
		return fmt.Errorf("error waiting for EIP Group (%s) to become available: %w", d.Id(), err)
	}

	return resourceEipGroupRead(d, meta)
}

func resourceEipGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	detail, err := FindEipGroup(conn, d.Id())
	log.Printf("[DEBUG] Read EIP Group (%s) result: %+v", d.Id(), detail)
	if err != nil {
		return fmt.Errorf("error reading EIP Group (%s): %w", d.Id(), err)
	}
	if detail == nil {
		log.Printf("[WARN] EIP Group (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.SetId(detail.Id)

	if err := d.Set("group_id", detail.Id); err != nil {
		return fmt.Errorf("error setting group_id: %w", err)
	}
	if err := d.Set("name", detail.Name); err != nil {
		return fmt.Errorf("error setting name: %w", err)
	}
	if err := d.Set("payment_timing", detail.PaymentTiming); err != nil {
		return fmt.Errorf("error setting payment_timing: %w", err)
	}
	if err := d.Set("billing_method", detail.BillingMethod); err != nil {
		return fmt.Errorf("error setting billing_method: %w", err)
	}
	if err := d.Set("route_type", detail.RouteType); err != nil {
		return fmt.Errorf("error setting route_type: %w", err)
	}
	if err := d.Set("bandwidth_in_mbps", detail.BandWidthInMbps); err != nil {
		return fmt.Errorf("error setting bandwidth_in_mbps: %w", err)
	}
	if err := d.Set("tags", flex.FlattenTagModelToMap(detail.Tags)); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	// computed fields
	if err := d.Set("status", detail.Status); err != nil {
		return fmt.Errorf("error setting status: %w", err)
	}
	if err := d.Set("default_domestic_bandwidth", detail.DefaultDomesticBandwidth); err != nil {
		return fmt.Errorf("error setting default_domestic_bandwidth: %w", err)
	}
	if err := d.Set("bw_short_id", detail.BwShortId); err != nil {
		return fmt.Errorf("error setting bw_short_id: %w", err)
	}
	if err := d.Set("bw_bandwidth_in_mbps", detail.BwBandwidthInMbps); err != nil {
		return fmt.Errorf("error setting bw_bandwidth_in_mbps: %w", err)
	}
	if err := d.Set("domestic_bw_short_id", detail.DomesticBwShortId); err != nil {
		return fmt.Errorf("error setting domestic_bw_short_id: %w", err)
	}
	if err := d.Set("domestic_bw_bandwidth_in_mbps", detail.DomesticBwBandwidthInMbps); err != nil {
		return fmt.Errorf("error setting domestic_bw_bandwidth_in_mbps: %w", err)
	}
	if err := d.Set("create_time", detail.CreateTime); err != nil {
		return fmt.Errorf("error setting create_time: %w", err)
	}
	if err := d.Set("expire_time", detail.ExpireTime); err != nil {
		return fmt.Errorf("error setting expire_time: %w", err)
	}
	if err := d.Set("region", detail.Region); err != nil {
		return fmt.Errorf("error setting region: %w", err)
	}
	if err := d.Set("eips", flattenEipList(detail.Eips)); err != nil {
		return fmt.Errorf("error setting eips: %w", err)
	}
	if err := d.Set("eip_count", len(detail.Eips)); err != nil {
		return fmt.Errorf("error setting eip_count: %w", err)
	}
	if err := d.Set("eipv6s", flattenEipList(detail.Eipv6s)); err != nil {
		return fmt.Errorf("error setting epipv6s: %w", err)
	}
	if err := d.Set("eipv6_count", len(detail.Eipv6s)); err != nil {
		return fmt.Errorf("error setting eipv6_count: %w", err)
	}

	return nil
}

func resourceEipGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	if err := updateName(d, conn); err != nil {
		return fmt.Errorf("error updating EIP Group (%s) name: %w", d.Id(), err)
	}
	if err := updateBandwidth(d, conn); err != nil {
		return fmt.Errorf("error updating EIP Group (%s) bandwidth: %w", d.Id(), err)
	}
	if err := updateIpCount(d, conn); err != nil {
		return fmt.Errorf("error updating EIP Group (%s) IP count: %w", d.Id(), err)
	}

	return resourceEipGroupRead(d, meta)
}

func resourceEipGroupDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	_, err := conn.WithEipClient(func(client *eip.Client) (interface{}, error) {
		if d.Get("payment_timing") == "Prepaid" {
			log.Printf("release prepaid respurce")
			return nil, client.RefundEipGroup(d.Id(), flex.BuildClientToken())
		}

		return nil, client.DeleteEipGroup(d.Id(), flex.BuildClientToken())
	})
	log.Printf("[DEBUG] Delete EIP Group (%s)", d.Id())

	if err != nil {
		return fmt.Errorf("error deleting EIP Group (%s): %w", d.Id(), err)
	}

	err = waitEipGroupDeleted(conn, d.Timeout(schema.TimeoutDelete), d.Id())
	if err != nil {
		return fmt.Errorf("error waiting for EIP Group (%s) to be deleted: %w", d.Id(), err)
	}
	return nil
}

func buildCreationArgs(d *schema.ResourceData) *eip.CreateEipGroupArgs {
	billing := &eip.Billing{
		PaymentTiming: d.Get("payment_timing").(string),
		BillingMethod: d.Get("billing_method").(string),
	}

	if billing.PaymentTiming == "Prepaid" {
		billing.Reservation = &eip.Reservation{
			ReservationLength:   d.Get("reservation_length").(int),
			ReservationTimeUnit: d.Get("reservation_time_unit").(string),
		}
	}

	args := &eip.CreateEipGroupArgs{
		Name:            d.Get("name").(string),
		EipCount:        d.Get("eip_count").(int),
		Eipv6Count:      d.Get("eipv6_count").(int),
		BandWidthInMbps: d.Get("bandwidth_in_mbps").(int),
		Billing:         billing,
		Tags:            flex.ExpandMapToTagModel[model.TagModel](d.Get("tags").(map[string]interface{})),
		RouteType:       d.Get("route_type").(string),
		Idc:             d.Get("idc").(string),
		Continuous:      d.Get("continuous").(bool),
		ClientToken:     flex.BuildClientToken(),
	}

	return args
}

func updateName(d *schema.ResourceData, conn *connectivity.BaiduClient) error {
	if d.HasChange("name") {
		_, err := conn.WithEipClient(func(client *eip.Client) (interface{}, error) {
			args := &eip.RenameEipGroupArgs{
				Name:        d.Get("name").(string),
				ClientToken: flex.BuildClientToken(),
			}
			return nil, client.RenameEipGroup(d.Id(), args)
		})

		if err != nil {
			return fmt.Errorf("error waiting for EIP Group (%s) to be deleted: %w", d.Id(), err)
		}
	}
	return nil
}

func updateBandwidth(d *schema.ResourceData, conn *connectivity.BaiduClient) error {
	if d.HasChange("bandwidth_in_mbps") {
		_, err := conn.WithEipClient(func(client *eip.Client) (interface{}, error) {
			args := &eip.ResizeEipGroupArgs{
				BandWidthInMbps: d.Get("bandwidth_in_mbps").(int),
				ClientToken:     flex.BuildClientToken(),
			}
			return nil, client.ResizeEipGroupBandWidth(d.Id(), args)
		})

		if err != nil {
			return err
		}

		_, err = waitEipGroupAvailable(conn, d.Timeout(schema.TimeoutUpdate), d.Id())
		if err != nil {
			return fmt.Errorf("error waiting for EIP Group (%s) to become available: %w", d.Id(), err)
		}

	}
	return nil
}

func updateIpCount(d *schema.ResourceData, conn *connectivity.BaiduClient) error {
	if d.HasChanges("eip_count", "eipv6_count") {
		args := &eip.GroupAddEipCountArgs{
			ClientToken: flex.BuildClientToken(),
		}
		if d.HasChange("eip_count") {
			o, n := d.GetChange("eip_count")
			args.EipAddCount = n.(int) - o.(int)
		}
		if d.HasChange("eipv6_count") {
			o, n := d.GetChange("eipv6_count")
			args.Eipv6AddCount = n.(int) - o.(int)
		}

		if args.EipAddCount < 0 || args.Eipv6AddCount < 0 {
			return fmt.Errorf("please use the `baiducloud_eipgroup_detachment` resource to detach EIPs from the EIP Group")
		}

		_, err := conn.WithEipClient(func(client *eip.Client) (interface{}, error) {
			return nil, client.EipGroupAddEipCount(d.Id(), args)
		})
		if err != nil {
			return err
		}
	}
	return nil
}
