/*
Use this data source to query spec list.

Example Usage

```hcl
data "baiducloud_bcc_specs" "default" {
  zone_name = "cn-bj-d"
  output_file = "specs.json"

  filter {
    name = "cpu_count"
    values = ["^([1])$"]
  }
}

output "spec" {
  value = "${data.baiducloud_specs.default.specs}"
}
```
*/
package baiducloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudBccFlavors() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudBccFlavorsRead,

		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:        schema.TypeString,
				Description: "BCC Flavor search result output file",
				Optional:    true,
				ForceNew:    true,
			},
			"zone_name": {
				Type:        schema.TypeString,
				Description: "Zone name",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),
			"specs": {
				Type:        schema.TypeList,
				Description: "Specs list",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_name": {
							Type:        schema.TypeString,
							Description: "Zone name",
							Computed:    true,
						},
						"group_id": {
							Type:        schema.TypeString,
							Description: "Group id",
							Computed:    true,
						},
						"cpu_count": {
							Type:        schema.TypeInt,
							Description: "CPU count",
							Computed:    true,
						},
						"memory_capacity_in_gb": {
							Type:        schema.TypeInt,
							Description: "Memory capacity in GB",
							Computed:    true,
						},
						"ephemeral_disk_in_gb": {
							Type:        schema.TypeInt,
							Description: "Ephemeral disk size in gb",
							Computed:    true,
						},
						"ephemeral_disk_count": {
							Type:        schema.TypeInt,
							Description: "Count of ephemeral disk",
							Computed:    true,
						},
						"ephemeral_disk_type": {
							Type:        schema.TypeString,
							Description: "Type of ephemeral disk",
							Computed:    true,
						},
						"gpu_card_type": {
							Type:        schema.TypeString,
							Description: "Type of gpu card",
							Computed:    true,
						},
						"gpu_card_count": {
							Type:        schema.TypeInt,
							Description: "Count of gpu card",
							Computed:    true,
						},
						"fpga_card_type": {
							Type:        schema.TypeString,
							Description: "Type of FPGA card",
							Computed:    true,
						},
						"fpga_card_count": {
							Type:        schema.TypeInt,
							Description: "Count of FPGA card",
							Computed:    true,
						},
						"product_type": {
							Type:        schema.TypeString,
							Description: "Product type",
							Computed:    true,
						},
						"spec": {
							Type:        schema.TypeString,
							Description: "Spec name",
							Computed:    true,
						},
						"spec_id": {
							Type:        schema.TypeString,
							Description: "Spec id",
							Computed:    true,
						},
						"cpu_model": {
							Type:        schema.TypeString,
							Description: "CPU model name",
							Computed:    true,
						},
						"cpu_ghz": {
							Type:        schema.TypeString,
							Description: "CPU frequency",
							Computed:    true,
						},
						"network_bandwidth": {
							Type:        schema.TypeString,
							Description: "Network bandwidth",
							Computed:    true,
						},
						"network_package": {
							Type:        schema.TypeString,
							Description: "Network package",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudBccFlavorsRead(d *schema.ResourceData, meta interface{}) error {
	action := "List all specs"
	client := meta.(*connectivity.BaiduClient)
	bccService := BccService{client}
	zone := d.Get("zone_name").(string)
	specs, err := bccService.ListAllFlavors(zone)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bcc_specs", action, BCESDKGoERROR)
	}
	specsMap := bccService.FlattenBccFlavorsModelToMap(specs)

	FilterDataSourceResult(d, &specsMap)
	if err = d.Set("specs", specsMap); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bcc_specs", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), specsMap); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bcc_specs", action, BCESDKGoERROR)
		}
	}
	return nil
}
