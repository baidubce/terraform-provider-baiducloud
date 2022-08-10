package baiducloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudPrivateZoneDNSRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudPrivateZoneDNSRecordsRead,

		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Description: "ID of the private zone.",
				Required:    true,
				ForceNew:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "Output file for saving result.",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),
			"records": {
				Type:        schema.TypeList,
				Description: "Result of records.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"record_id": {
							Type:        schema.TypeString,
							Description: "Id of private zone dns record",
							Computed:    true,
						},
						"rr": {
							Type:        schema.TypeString,
							Description: "rr",
							Computed:    true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "value",
							Computed:    true,
						},
						"type": {
							Type:        schema.TypeString,
							Description: "type",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "status",
							Computed:    true,
						},
						"priority": {
							Type:        schema.TypeInt,
							Description: "priority",
							Computed:    true,
						},
						"ttl": {
							Type:        schema.TypeInt,
							Description: "ttl",
							Computed:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "description",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudPrivateZoneDNSRecordsRead(d *schema.ResourceData, meta interface{}) error {
	action := "List all records"
	client := meta.(*connectivity.BaiduClient)
	privateZoneDnsService := PrivateZoneDnsService{client}
	var (
		zoneId string
	)
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}
	records, err := privateZoneDnsService.ListAlDnsRecords(zoneId)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_private_zone_dns_record", action, BCESDKGoERROR)
	}
	addDebug(action, records)
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}
	recordsResult := make([]map[string]interface{}, 0)
	for _, record := range records {
		localDnsMap := make(map[string]interface{})
		localDnsMap["record_id"] = record.RecordId
		localDnsMap["rr"] = record.Rr
		localDnsMap["value"] = record.Value
		localDnsMap["type"] = record.Type
		localDnsMap["ttl"] = record.Ttl
		localDnsMap["priority"] = record.Priority
		localDnsMap["description"] = record.Description
		localDnsMap["status"] = record.Status
		recordsResult = append(recordsResult, localDnsMap)
	}
	FilterDataSourceResult(d, &recordsResult)
	if err := d.Set("records", recordsResult); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_private_zone_dns_record", action, BCESDKGoERROR)
	}

	d.SetId(resource.UniqueId())
	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), recordsResult); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_private_zone_dns_record", action, BCESDKGoERROR)
		}
	}
	return nil
}
