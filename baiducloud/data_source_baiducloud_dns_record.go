/*
Use this data source to query Dns record list.

Example Usage

```hcl
data "baiducloud_dns_records" "default" {
	zone_name = "xxxx"
}

output "records" {
 value = "${data.baiducloud_dns_records.default.records}"
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

func dataSourceBaiduCloudDnsrecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudDnsrecordsRead,

		Schema: map[string]*schema.Schema{
			"zone_name": {
				Type:        schema.TypeString,
				Description: "Dns record zone name",
				Required:    true,
				ForceNew:    true,
			},
			"rr": {
				Type:        schema.TypeString,
				Description: "Dns record rr",
				Optional:    true,
			},
			"record_id": {
				Type:        schema.TypeString,
				Description: "Dns record id",
				Optional:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "DNS records search result output file",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),

			"records": {
				Type:        schema.TypeList,
				Description: "record list",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rr": {
							Type:        schema.TypeString,
							Description: "Dns record rr",
							Computed:    true,
						},
						"type": {
							Type:        schema.TypeString,
							Description: "Dns record type",
							Computed:    true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "Dns record value",
							Computed:    true,
						},
						"ttl": {
							Type:        schema.TypeInt,
							Description: "Dns record ttl",
							Computed:    true,
						},
						"line": {
							Type:        schema.TypeString,
							Description: "Dns record line",
							Computed:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "Dns record description",
							Computed:    true,
						},
						"priority": {
							Type:        schema.TypeInt,
							Description: "Dns record priority",
							Computed:    true,
						},
						"record_id": {
							Type:        schema.TypeString,
							Description: "Dns record id",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Dns record status",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudDnsrecordsRead(d *schema.ResourceData, meta interface{}) error {

	action := "List all dns record name "

	dnsrecordArgs := buildBaiduCloudCreatednsrecordListArgs(d)

	zoneName := d.Get("zone_name").(string)

	records, err := listAllrecordList(zoneName, dnsrecordArgs, meta)

	addDebug(action, records)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dns_records", action, BCESDKGoERROR)
	}

	recordsResult := make([]map[string]interface{}, 0)

	for _, record := range records {

		innerMap := make(map[string]interface{})
		innerMap["record_id"] = record.Id
		innerMap["rr"] = record.Rr
		innerMap["status"] = record.Status
		innerMap["type"] = record.Type
		innerMap["value"] = record.Value
		innerMap["ttl"] = record.Ttl
		innerMap["line"] = record.Line
		innerMap["description"] = record.Description
		innerMap["priority"] = record.Priority

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dns_records", action, BCESDKGoERROR)
		}

		recordsResult = append(recordsResult, innerMap)
	}

	addDebug(action, recordsResult)

	FilterDataSourceResult(d, &recordsResult)

	if err := d.Set("records", recordsResult); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dns_records", action, BCESDKGoERROR)
	}

	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), recordsResult); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dns_records", action, BCESDKGoERROR)
		}
	}
	return nil
}

func listAllrecordList(zoneName string, args *dns.ListRecordRequest, meta interface{}) ([]dns.Record, error) {
	client := meta.(*connectivity.BaiduClient)

	action := "List all dns records "

	records := make([]dns.Record, 0)

	for {
		raw, err := client.WithDNSClient(func(dnsClient *dns.Client) (interface{}, error) {
			return dnsClient.ListRecord(zoneName, args)
		})

		if err != nil {
			return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_dns_record", action, BCESDKGoERROR)
		}

		result, _ := raw.(*dns.ListRecordResponse)
		records = append(records, result.Records...)

		if !result.IsTruncated {
			break
		}

		args.Marker = result.NextMarker
		args.MaxKeys = int(result.MaxKeys)
	}

	return records, nil
}
func buildBaiduCloudCreatednsrecordListArgs(d *schema.ResourceData) *dns.ListRecordRequest {

	request := &dns.ListRecordRequest{}

	if v, ok := d.GetOk("rr"); ok && len(v.(string)) > 0 {
		request.Rr = v.(string)
	}

	if v, ok := d.GetOk("record_id"); ok && len(v.(string)) > 0 {
		request.Id = v.(string)
	}

	return request
}
