/*
Use this data source to query Dns zone list.

Example Usage

```hcl
data "baiducloud_dns_zones" "default" {
	name = "xxxx"
}

output "zones" {
 value = "${data.baiducloud_dns_zones.default.zones}"
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/dns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudDnsZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudDnszonesRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name of DNS ZONE",
				Optional:    true,
				ForceNew:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "DNS Zones search result output file",
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
						"name": {
							Type:        schema.TypeString,
							Description: "Dns zone name",
							Required:    true,
							ForceNew:    true,
						},
						"zone_id": {
							Type:        schema.TypeString,
							Description: "Dns zone id",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Dns zone status",
							Computed:    true,
						},
						"product_version": {
							Type:        schema.TypeString,
							Description: "Dns zone product_version",
							Computed:    true,
						},
						"create_time": {
							Type:        schema.TypeString,
							Description: "Dns zone create_time",
							Computed:    true,
						},
						"expire_time": {
							Type:        schema.TypeString,
							Description: "Dns zone expire_time",
							Computed:    true,
						},
						"tags": tagsSchema(),
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudDnszonesRead(d *schema.ResourceData, meta interface{}) error {

	name := d.Get("name").(string)

	action := "List all dns zone name is " + name

	dnsZoneArgs, err := buildBaiduCloudDnszoneListArgs(d)

	if err != nil {
		return WrapError(err)
	}

	zones, err := listAllDnsZones(dnsZoneArgs, meta)

	addDebug(action, zones)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dns_zones", action, BCESDKGoERROR)
	}

	zonesResult := make([]map[string]interface{}, 0)

	for _, zone := range zones {

		innerMap := make(map[string]interface{})
		innerMap["zone_id"] = zone.Id
		innerMap["name"] = zone.Name
		innerMap["status"] = zone.Status
		innerMap["product_version"] = zone.ProductVersion
		innerMap["create_time"] = zone.CreateTime
		innerMap["expire_time"] = zone.ExpireTime
		innerMap["tags"] = zoneTagsToMap(zone.Tags)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dns_zones", action, BCESDKGoERROR)
		}

		zonesResult = append(zonesResult, innerMap)
	}

	addDebug(action, zonesResult)

	FilterDataSourceResult(d, &zonesResult)

	if err := d.Set("zones", zonesResult); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dns_zones", action, BCESDKGoERROR)
	}

	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), zonesResult); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dns_zones", action, BCESDKGoERROR)
		}
	}
	return nil
}

func buildBaiduCloudDnszoneListArgs(d *schema.ResourceData) (*dns.ListZoneRequest, error) {

	request := &dns.ListZoneRequest{}

	if v, ok := d.GetOk("name"); ok && len(v.(string)) > 0 {
		request.Name = v.(string)
	}

	return request, nil
}

func listAllDnsZones(args *dns.ListZoneRequest, meta interface{}) ([]dns.Zone, error) {
	client := meta.(*connectivity.BaiduClient)
	action := "List all dns zones "

	zones := make([]dns.Zone, 0)

	for {
		raw, err := client.WithDNSClient(func(dnsClient *dns.Client) (interface{}, error) {
			return dnsClient.ListZone(args)
		})

		if err != nil {
			return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_dns_zone", action, BCESDKGoERROR)
		}

		result, _ := raw.(*dns.ListZoneResponse)
		zones = append(zones, result.Zones...)

		if !result.IsTruncated {
			break
		}

		args.Marker = result.NextMarker
		args.MaxKeys = int(result.MaxKeys)
	}

	return zones, nil
}

func zoneTagsToMap(tags []dns.TagModel) map[string]string {

	tagMap := make(map[string]string)

	for _, tag := range tags {
		tagMap[*tag.TagKey] = *tag.TagValue
	}

	return tagMap

}
