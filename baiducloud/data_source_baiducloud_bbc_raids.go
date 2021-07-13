/*
Use this data source to query bbc flavors list.

Example Usage

```hcl
data "baiducloud_bbc_raids" "default" {
	flavor_id = "abcd"
}

output "flavors" {
  value = "${data.bbcbaiducloud_bbc_raids.default.raids}"
}
```
*/
package baiducloud

import (
	"errors"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudBbcRaids() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudBbcRaidRead,

		Schema: map[string]*schema.Schema{
			"flavor_id": {
				Type:        schema.TypeString,
				Description: " flavor id",
				Required:    true,
				ForceNew:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: " flavors result output file",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),
			"raids": {
				Type:        schema.TypeList,
				Description: "raids list",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"raid_id": {
							Type:        schema.TypeString,
							Description: "raid id",
							Computed:    true,
						},
						"raid": {
							Type:        schema.TypeString,
							Description: "raid",
							Computed:    true,
						},
						"sys_swap_size": {
							Type:        schema.TypeInt,
							Description: "swap size",
							Computed:    true,
						},
						"sys_root_size": {
							Type:        schema.TypeInt,
							Description: "root size",
							Computed:    true,
						},
						"sys_home_size": {
							Type:        schema.TypeInt,
							Description: "home size",
							Computed:    true,
						},
						"sys_disk_size": {
							Type:        schema.TypeInt,
							Description: "system disk size",
							Computed:    true,
						},
						"data_disk_size": {
							Type:        schema.TypeInt,
							Description: "data disk size",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudBbcRaidRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bbcService := BbcService{client}
	action := "Query All raids for flavor"
	raidsMap := make([]map[string]interface{}, 0, 0)
	if v, ok := d.GetOk("flavor_id"); ok {
		result, err := bbcService.GetRaids(v.(string))
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "bbcbaiducloud_bbc_raids", action, BCESDKGoERROR)
		}
		raidsMap = bbcService.FlattenRaidsToMap(result.Raids)
	} else {
		return WrapErrorf(errors.New("flavor_id is required"), DefaultErrorMsg, "bbcbaiducloud_bbc_raids", action, BCESDKGoERROR)
	}

	addDebug(action, raidsMap)

	FilterDataSourceResult(d, &raidsMap)

	if err := d.Set("raids", raidsMap); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "bbcbaiducloud_bbc_raids", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())
	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), raidsMap); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "bbcbaiducloud_bbc_raids", action, BCESDKGoERROR)
		}
	}

	return nil
}
