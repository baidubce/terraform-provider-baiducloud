/*
Use this data source to query zone list.

Example Usage

```hcl
data "baiducloud_zones" "default" {}

output "zone" {
  value = "${data.baiducloud_zones.default.zones}"
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudZonesRead,

		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:        schema.TypeString,
				Description: "Output file for saving result.",
				Optional:    true,
				ForceNew:    true,
			},

			// Attributes used for result
			"zones": {
				Type:        schema.TypeList,
				Description: "Useful zone list",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_name": {
							Type:        schema.TypeString,
							Description: "Useful zone name",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	action := "Query all zones"
	raw, err := client.WithBccClient(func(client *bcc.Client) (i interface{}, e error) {
		return client.ListZone()
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_zones", action, BCESDKGoERROR)
	}
	addDebug(action, raw)

	response := raw.(*api.ListZoneResult)
	zoneMap := make([]map[string]interface{}, 0, len(response.Zones))
	for _, zone := range response.Zones {
		zoneMap = append(zoneMap, map[string]interface{}{
			"zone_name": zone.ZoneName,
		})
	}

	if err := d.Set("zones", zoneMap); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_zones", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), zoneMap); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_zones", action, BCESDKGoERROR)
		}
	}

	return nil
}
