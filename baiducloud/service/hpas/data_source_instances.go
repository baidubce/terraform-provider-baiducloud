package hpas

import (
	"fmt"

	"github.com/baidubce/bce-sdk-go/services/hpas/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/flex"
)

func DataSourceInstances() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to query HPAS instance list. \n\n",

		Read: dataSourceInstancesRead,

		Schema: map[string]*schema.Schema{
			"hpas_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of instance IDs.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the instance.",
			},
			"zone_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Zone information, e.g., `cn-bj-a`.",
			},
			"hpas_status": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Status of the instance. Valid values: `Creating`, `Active`, `Expired`, `Error`, `Stopping`, `Starting`, " +
					"`Stopped`, `Reboot`, `Rebuild`, `Password`, `ChangeVpc`, `ChangeSubnet`, `Template`.",
			},
			"app_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Application type. e.g., `llama2_7B_train`.",
			},
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        InstanceSchema(),
				Description: "Instance list.",
			},
		},
	}
}

func InstanceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"payment_timing": flex.ComputedSchemaPaymentTiming(),
			"tags":           flex.ComputedSchemaTags(),
			"hpas_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the instance.",
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
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the instance.",
			},
			"zone_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Zone information.",
			},
			"image_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Image ID used for the application.",
			},
			"image_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the image.",
			},
			"internal_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Internal IP addresses.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Subnet ID.",
			},
			"subnet_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the subnet.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VPC ID.",
			},
			"vpc_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the VPC.",
			},
			"vpc_cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CIDR block of the VPC.",
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
			"security_group_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Security group type. Possible values: `normal`, `enterprise`.",
			},
			"security_group_ids": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of security group IDs",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "Status of the instance. Possible values: `Creating`, `Active`, `Expired`, `Error`, `Stopping`, `Starting`, " +
					"`Stopped`, `Reboot`, `Rebuild`, `Password`, `ChangeVpc`, `ChangeSubnet`, `Template`.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the instance.",
			},
		},
	}
}

func dataSourceInstancesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	args := buildInstanceListArgs(d)
	result, err := FindInstances(conn, args)
	if err != nil {
		return err
	}

	if err := d.Set("instance_list", flattenInstanceList(result)); err != nil {
		return fmt.Errorf("error setting instance_list: %w", err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func buildInstanceListArgs(d *schema.ResourceData) api.ListHpasByMakerReq {
	args := api.ListHpasByMakerReq{}

	if v, ok := d.GetOk("hpas_ids"); ok && v.(*schema.Set).Len() > 0 {
		args.HpasIds = flex.ExpandStringValueSet(d.Get("hpas_ids").(*schema.Set))
	}

	if v, ok := d.GetOk("name"); ok {
		args.Name = v.(string)
	}

	if v, ok := d.GetOk("zone_name"); ok {
		args.ZoneName = v.(string)
	}

	if v, ok := d.GetOk("hpas_status"); ok {
		args.HpasStatus = v.(string)
	}

	if v, ok := d.GetOk("app_type"); ok {
		args.AppType = v.(string)
	}

	return args
}
