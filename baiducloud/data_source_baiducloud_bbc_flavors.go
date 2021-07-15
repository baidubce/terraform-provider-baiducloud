/*
Use this data source to query bbc flavors list.

Example Usage

```hcl
data "baiducloud_bbc_flavors" "default" {}

output "flavors" {
  value = "${data.baiducloud_bbc_flavors.default.flavors}"
}
```
*/
package baiducloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudBbcFlavors() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudBbcFlavorRead,

		Schema: map[string]*schema.Schema{
			"flavor_id": {
				Type:        schema.TypeString,
				Description: " flavor id",
				Optional:    true,
				ForceNew:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: " flavors result output file",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),

			"flavors": {
				Type:        schema.TypeList,
				Description: "flavor list",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"flavor_id": {
							Type:        schema.TypeString,
							Description: "flavor id",
							Computed:    true,
						},
						"cpu_count": {
							Type:        schema.TypeInt,
							Description: "cpu count",
							Computed:    true,
						},
						"cpu_type": {
							Type:        schema.TypeString,
							Description: "cpu type",
							Computed:    true,
						},
						"memory_capacity_in_gb": {
							Type:        schema.TypeInt,
							Description: "memory size in gb",
							Computed:    true,
						},
						"disk": {
							Type:        schema.TypeString,
							Description: "disk description",
							Computed:    true,
						},
						"network_card": {
							Type:        schema.TypeString,
							Description: "network card",
							Computed:    true,
						},
						"os_name": {
							Type:        schema.TypeString,
							Description: "Image os name",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudBbcFlavorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bbcService := BbcService{client}
	action := "Query All Flavors"
	flavorMap := make([]map[string]interface{}, 0)
	if v, ok := d.GetOk("flavor_id"); ok {
		result, err := bbcService.GetFlavorDetail(v.(string))
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_flavors", action, BCESDKGoERROR)
		}
		flavorMap = bbcService.FlattenFlavorDetailToMap(result)
	} else {
		result, err := bbcService.GetFlavors()
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_flavors", action, BCESDKGoERROR)
		}
		flavorMap = bbcService.FlattenFlavorsToMap(result.Flavors)
	}

	addDebug(action, flavorMap)

	FilterDataSourceResult(d, &flavorMap)

	if err := d.Set("flavors", flavorMap); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_flavors", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())
	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), flavorMap); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_flavors", action, BCESDKGoERROR)
		}
	}

	return nil
}
