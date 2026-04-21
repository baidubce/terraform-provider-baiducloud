package hpas

import (
	"fmt"

	"github.com/baidubce/bce-sdk-go/services/hpas/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/flex"
)

func DataSourceReservedInstances() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to query HPAS reserved instance list. \n\n",

		Read: dataSourceReservedInstancesRead,

		Schema: map[string]*schema.Schema{
			"reserved_hpas_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of reserved instance IDs to filter by.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the reserved instance to filter by.",
			},
			"zone_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Zone information to filter by, e.g., `cn-bj-a`.",
			},
			"reserved_hpas_status": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Status to filter by. Valid values: `Creating`, `Active`, `Pending`, " +
					"`Expired`, `Recharge`, `Deleted`.",
			},
			"app_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Application type to filter by. e.g., `llama2_7B_train`.",
			},
			"hpas_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by the ID of the HPAS instance deducted in the previous billing period.",
			},
			"reserved_instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        ReservedInstanceSchema(),
				Description: "List of reserved instances.",
			},
		},
	}
}

func ReservedInstanceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"reserved_hpas_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the reserved instance.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the reserved instance.",
			},
			"zone_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Zone information.",
			},
			"app_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Application type.",
			},
			"app_performance_level": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Performance level of the application.",
			},
			"payment_timing": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Payment timing of billing.",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "Status of the reserved instance. Possible values: `Creating`, `Active`, `Pending`, " +
					"`Expired`, `Recharge`, `Deleted`.",
			},
			"period": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Reservation period in months.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the reserved instance.",
			},
			"expire_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expiration time of the reserved instance.",
			},
			"tags": flex.ComputedSchemaTags(),
			"hpas_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the HPAS instance deducted in the previous billing period.",
			},
			"hpas_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the HPAS instance deducted in the previous billing period.",
			},
			"deduct_instance": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether there is a deducting instance.",
			},
			"ehc_cluster_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "EHC cluster ID.",
			},
			"ehc_cluster_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the EHC cluster.",
			},
		},
	}
}

func dataSourceReservedInstancesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	args := buildReservedInstanceListArgs(d)
	result, err := FindReservedInstances(conn, args)
	if err != nil {
		return err
	}

	if err := d.Set("reserved_instance_list", flattenReservedInstanceList(result)); err != nil {
		return fmt.Errorf("error setting reserved_instance_list: %w", err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func buildReservedInstanceListArgs(d *schema.ResourceData) api.ListReservedHpasByMakerReq {
	args := api.ListReservedHpasByMakerReq{}

	if v, ok := d.GetOk("reserved_hpas_ids"); ok && v.(*schema.Set).Len() > 0 {
		args.ReservedHpasIds = flex.ExpandStringValueSet(v.(*schema.Set))
	}

	if v, ok := d.GetOk("name"); ok {
		args.Name = v.(string)
	}

	if v, ok := d.GetOk("zone_name"); ok {
		args.ZoneName = v.(string)
	}

	if v, ok := d.GetOk("reserved_hpas_status"); ok {
		args.ReservedHpasStatus = v.(string)
	}

	if v, ok := d.GetOk("app_type"); ok {
		args.AppType = v.(string)
	}

	if v, ok := d.GetOk("hpas_id"); ok {
		args.HpasId = v.(string)
	}

	return args
}
