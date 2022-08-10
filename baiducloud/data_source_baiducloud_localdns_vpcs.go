/*
Use this data source to query localdns VPCs.

Example Usage

```hcl
data "baiducloud_localdns_vpcs" "default" {}

output "vpcs" {
   value = "${data.baiducloud_localdns_vpcs.default.bind_vpcs}"
}
```
*/
package baiducloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudLocalDnsVpcs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudLocalDnsVpcRead,

		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Description: "zone_id of the DNS privatezone ",
				ForceNew:    true,
				Required:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "local dns vpc search result output file",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),

			"bind_vpcs": {
				Type:        schema.TypeList,
				Description: "privatezone bind vpcs",
				Computed:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Description: "bind vpc id",
							Computed:    true,
						},
						"vpc_name": {
							Type:        schema.TypeString,
							Description: "name of vpc",
							Computed:    true,
						},
						"vpc_region": {
							Type:        schema.TypeString,
							Description: "region of vpc",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudLocalDnsVpcRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	localDnsService := LocalDnsService{client}

	zoneId := d.Get("zone_id").(string)

	action := "Query localdns VPC list" + zoneId

	zone, err := localDnsService.GetPrivateZoneDetail(zoneId)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_localdns_vpcs", action, BCESDKGoERROR)
	}
	addDebug(action, zone)

	bindVpcs := make([]map[string]interface{}, 0)
	for _, vpc := range zone.BindVpcs {
		vpcMap := make(map[string]interface{})
		vpcMap["vpc_id"] = vpc.VpcId
		vpcMap["vpc_name"] = vpc.VpcName
		vpcMap["vpc_region"] = vpc.VpcRegion

		bindVpcs = append(bindVpcs, vpcMap)
	}

	FilterDataSourceResult(d, &bindVpcs)

	if err := d.Set("bind_vpcs", bindVpcs); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_localdns_vpcs", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), bindVpcs); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_localdns_vpcs", action, BCESDKGoERROR)
		}
	}

	return nil
}
