/*
Use this data source to query BBC flavors list.

Example Usage

```hcl
data "baiducloud_bbc_flavors" "bbc_flavors" {

}

output "flavors" {
 value = "${data.baiducloud_bbc_flavors.bbc_flavors.flavors}"
}

```
*/
package baiducloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudBbcFlavors() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudBbcFlavorsRead,

		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:        schema.TypeString,
				Description: "Flavor search result output file",
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
							Description: "memory capacity in GB",
							Computed:    true,
						},
						"disk": {
							Type:        schema.TypeString,
							Description: "disk",
							Computed:    true,
						},
						"network_card": {
							Type:        schema.TypeString,
							Description: "network card",
							Computed:    true,
						},
						"others": {
							Type:        schema.TypeString,
							Description: "others",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudBbcFlavorsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bbcService := BbcService{client}
	action := "Query All Bbc Flavors"
	flavorsResult, err := bbcService.ListAllBbcFlavors()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_flavors", action, BCESDKGoERROR)
	}
	addDebug(action, flavorsResult)

	flavorsMap := bbcService.FlattenFlavorModelToMap(flavorsResult)
	FilterDataSourceResult(d, &flavorsMap)

	if err := d.Set("flavors", flavorsMap); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_flavors", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), flavorsMap); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_flavors", action, BCESDKGoERROR)
		}
	}
	return nil
}
