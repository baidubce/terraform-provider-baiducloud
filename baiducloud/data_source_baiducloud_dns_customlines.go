/*
Use this data source to query Dns customline list.

Example Usage

```hcl
data "baiducloud_dns_customlines" "default" {
	name = "xxxx"
}

output "customlines" {
 value = "${data.baiducloud_dns_customlines.default.customlines}"
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

func dataSourceBaiduCloudDnscustomlines() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudDnscustomlinesRead,

		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:        schema.TypeString,
				Description: "DNS customlines search result output file",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),

			"customlines": {
				Type:        schema.TypeList,
				Description: "customline list",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Dns customline name",
							Computed:    true,
						},
						"lines": {
							Type:        schema.TypeSet,
							Description: "lines of dns ",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"line_id": {
							Type:        schema.TypeString,
							Description: "Dns customline id",
							Computed:    true,
						},
						"related_zone_count": {
							Type:        schema.TypeInt,
							Description: "Dns customline related zone count",
							Computed:    true,
						},
						"related_record_count": {
							Type:        schema.TypeInt,
							Description: "Dns customline related record count",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudDnscustomlinesRead(d *schema.ResourceData, meta interface{}) error {

	action := "List all dns customline name "

	dnscustomlineArgs := buildBaiduCloudCreatednscustomlineListArgs(d)

	customlines, err := listAllcustomlineList(dnscustomlineArgs, meta)

	addDebug(action, customlines)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dns_customlines", action, BCESDKGoERROR)
	}

	customlinesResult := make([]map[string]interface{}, 0)

	for _, customline := range customlines {

		innerMap := make(map[string]interface{})
		innerMap["line_id"] = customline.Id
		innerMap["name"] = customline.Name
		innerMap["lines"] = customline.Lines
		innerMap["related_zone_count"] = customline.RelatedZoneCount
		innerMap["related_record_count"] = customline.RelatedRecordCount

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dns_customlines", action, BCESDKGoERROR)
		}

		customlinesResult = append(customlinesResult, innerMap)
	}

	addDebug(action, customlinesResult)

	FilterDataSourceResult(d, &customlinesResult)

	if err := d.Set("customlines", customlinesResult); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dns_customlines", action, BCESDKGoERROR)
	}

	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), customlinesResult); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dns_customlines", action, BCESDKGoERROR)
		}
	}
	return nil
}

func listAllcustomlineList(args *dns.ListLineGroupRequest, meta interface{}) ([]dns.Line, error) {
	client := meta.(*connectivity.BaiduClient)

	action := "List all dns customlines "

	customlines := make([]dns.Line, 0)

	for {
		raw, err := client.WithDNSClient(func(dnsClient *dns.Client) (interface{}, error) {
			return dnsClient.ListLineGroup(args)
		})

		if err != nil {
			return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_dns_customline", action, BCESDKGoERROR)
		}

		result, _ := raw.(*dns.ListLineGroupResponse)
		customlines = append(customlines, result.LineList...)

		if !*result.IsTruncated {
			break
		}

		args.Marker = *result.NextMarker
		args.MaxKeys = int(*result.MaxKeys)
	}

	return customlines, nil
}
func buildBaiduCloudCreatednscustomlineListArgs(d *schema.ResourceData) *dns.ListLineGroupRequest {

	request := &dns.ListLineGroupRequest{}

	return request
}
