/*
Use this data source to query localdns privatezone.

Example Usage

```hcl
data "baiducloud_localdns_privatezones" "default" {}

output "zones" {
   value = "${data.baiducloud_localdns_privatezones.default.zones}"
}
```
*/
package baiducloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudLocalDnsPrivateZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudDnsLocalPrivateZoneRead,

		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:        schema.TypeString,
				Description: "local dns privatezone search result output file",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),

			"zones": {
				Type:        schema.TypeList,
				Description: "zone list",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:        schema.TypeString,
							Description: "id of the DNS local PrivateZone",
							Computed:    true,
						},
						"zone_name": {
							Type:        schema.TypeString,
							Description: "name of the DNS local PrivateZone",
							Computed:    true,
						},
						"record_count": {
							Type:        schema.TypeInt,
							Description: "record_count of the DNS local PrivateZone.",
							Computed:    true,
						},
						"create_time": {
							Type:        schema.TypeString,
							Description: "Creation time of the DNS local PrivateZone.",
							Computed:    true,
						},
						"update_time": {
							Type:        schema.TypeString,
							Description: "update time of the DNS local PrivateZone.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudDnsLocalPrivateZoneRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	localDnsService := LocalDnsService{client}

	action := "Query localdns privatezone list"

	zoneMap := make([]map[string]interface{}, 0)
	zonelist, err := localDnsService.GetPrivateZoneList()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_localdns_privatezones", action, BCESDKGoERROR)
	}
	addDebug(action, zonelist)

	for _, zone := range zonelist.Zones {
		zoneMap = append(zoneMap, map[string]interface{}{
			"zone_id":      zone.ZoneId,
			"zone_name":    zone.ZoneName,
			"record_count": zone.RecordCount,
			"create_time":  zone.CreateTime,
			"update_time":  zone.UpdateTime,
		})
	}

	FilterDataSourceResult(d, &zoneMap)

	if err := d.Set("zones", zoneMap); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_localdns_privatezones", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), zoneMap); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_localdns_privatezones", action, BCESDKGoERROR)
		}
	}

	return nil
}
