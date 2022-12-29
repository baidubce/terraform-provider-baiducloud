package bec

import (
	"fmt"
	"log"

	"github.com/baidubce/bce-sdk-go/services/bec"
	"github.com/baidubce/bce-sdk-go/services/bec/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func DataSourceNodes() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to query BEC node list. \n\n",

		Read: dataSourceNodesRead,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Description:  "Node type. Valid values: `gpu`, `container`, `vm`, `lb`, `bm`",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"gpu", "container", "vm", "lb", "bm"}, false),
			},
			"region_list": {
				Type:        schema.TypeList,
				Description: "Region list.",
				Computed:    true,
				Elem:        RegionSchema(),
			},
		},
	}
}

func RegionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Description: "English name of the region.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Chinese name of the region.",
				Computed:    true,
			},
			"country": {
				Type:        schema.TypeString,
				Description: "English name of the country to which the region belongs.",
				Computed:    true,
			},
			"country_name": {
				Type:        schema.TypeString,
				Description: "Chinese name of the country to which the region belongs.",
				Computed:    true,
			},
			"city_list": {
				Type:        schema.TypeList,
				Description: "City list of the region.",
				Computed:    true,
				Elem:        CitySchema(),
			},
		},
	}
}

func CitySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"city": {
				Type:        schema.TypeString,
				Description: "English name of the city.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Chinese name of the city.",
				Computed:    true,
			},
			"service_provider_list": {
				Type:        schema.TypeList,
				Description: "Service provider list of the city.",
				Computed:    true,
				Elem:        ServiceProviderSchema(),
			},
		},
	}
}

func ServiceProviderSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"service_provider": {
				Type:        schema.TypeString,
				Description: "English name of the service provider.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Chinese name of the service provider.",
				Computed:    true,
			},
			"region_id": {
				Type:        schema.TypeString,
				Description: "Full ID of the node.",
				Computed:    true,
			},
		},
	}
}

func dataSourceNodesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	raw, err := conn.WithBECClient(func(client *bec.Client) (interface{}, error) {
		return client.GetBecAvailableNodeInfoVo(d.Get("type").(string))
	})
	log.Printf("[DEBUG] Read BEC node list result: %+v", raw)
	if err != nil {
		return fmt.Errorf("error reading BEC node list: %w", err)
	}

	response := raw.(*api.GetBecAvailableNodeInfoVoResult)

	if err := d.Set("region_list", flattenRegionList(response.RegionList)); err != nil {
		return fmt.Errorf("error setting region_list: %w", err)
	}

	d.SetId(resource.UniqueId())
	return nil
}
