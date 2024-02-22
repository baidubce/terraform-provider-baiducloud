/*
Provide a resource to create an Dns record.

Example Usage

```hcl
resource "baiducloud_dns_record" "default" {
  zone_name              = "testZoneName"
  rr                     = "test"
  type                   = "test"
  value                  = "test"
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/dns"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudDnsRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudDnsrecordCreate,
		Read:   resourceBaiduCloudDnsrecordRead,
		Update: resourceBaiduCloudDnsrecordUpdate,
		Delete: resourceBaiduCloudDnsrecordDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

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
				Required:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Dns record type",
				Required:    true,
			},
			"value": {
				Type:        schema.TypeString,
				Description: "Dns record value",
				Required:    true,
			},
			"ttl": {
				Type:        schema.TypeInt,
				Description: "Dns record ttl",
				Optional:    true,
				Computed:    true,
			},
			"line": {
				Type:        schema.TypeString,
				Description: "Dns record line",
				Optional:    true,
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Dns record description",
				Optional:    true,
				Computed:    true,
			},
			"priority": {
				Type:        schema.TypeInt,
				Description: "Dns record priority",
				Optional:    true,
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
			"record_action": {
				Type:        schema.TypeString,
				Description: "Dns record action",
				Optional:    true,
			},
		},
	}
}

func resourceBaiduCloudDnsrecordCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.BaiduClient)

	createDnsArgs := buildBaiduCloudCreatednsrecordArgs(d)

	zoneName := d.Get("zone_name").(string)

	action := "Create Dns record zone name " + zoneName + " - " + createDnsArgs.Rr + " - " + createDnsArgs.Type + " - " + createDnsArgs.Value

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {

		raw, err := client.WithDNSClient(func(dnsClient *dns.Client) (interface{}, error) {
			return nil, dnsClient.CreateRecord(zoneName, createDnsArgs, buildClientToken())
		})

		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		addDebug(action, raw)

		d.SetId(resource.UniqueId())
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dns_record", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudDnsrecordRead(d, meta)
}

func resourceBaiduCloudDnsrecordRead(d *schema.ResourceData, meta interface{}) error {

	action := "Query DNS record "

	queryArgs := buildBaiduCloudCreatednsrecordQueryArgs(d)

	zoneName := d.Get("zone_name").(string)

	records, err := listAllrecords(zoneName, queryArgs, meta)

	addDebug(action, records)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dns_record", action, BCESDKGoERROR)
	}

	var record dns.Record
	for _, z := range records {
		if z.Rr == d.Get("rr") && z.Type == d.Get("type") && z.Value == d.Get("value") {
			record = z
			break
		}
	}

	d.Set("record_id", record.Id)

	d.Set("rr", record.Rr)

	d.Set("status", record.Status)

	d.Set("type", record.Type)

	d.Set("value", record.Value)

	d.Set("ttl", record.Ttl)

	d.Set("line", record.Line)

	d.Set("description", record.Description)

	d.Set("priority", record.Priority)

	return nil
}

func resourceBaiduCloudDnsrecordDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.BaiduClient)

	zoneName := d.Get("zone_name").(string)

	recordId := d.Get("record_id").(string)

	action := "Delete dns zoneName IS " + zoneName + " recordId is " + recordId

	_, err := client.WithDNSClient(func(dnsClient *dns.Client) (interface{}, error) {
		return nil, dnsClient.DeleteRecord(zoneName, recordId, buildClientToken())
	})

	if err != nil {
		if IsExceptedErrors(err, ObjectNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dns_record", action, BCESDKGoERROR)
	}

	addDebug(action, err)

	return nil
}

func buildBaiduCloudCreatednsrecordQueryArgs(d *schema.ResourceData) *dns.ListRecordRequest {

	request := &dns.ListRecordRequest{}

	return request
}

func buildBaiduCloudCreatednsrecordArgs(d *schema.ResourceData) *dns.CreateRecordRequest {

	request := &dns.CreateRecordRequest{}

	if v, ok := d.GetOk("rr"); ok && len(v.(string)) > 0 {
		request.Rr = v.(string)
	}

	if v, ok := d.GetOk("type"); ok && len(v.(string)) > 0 {
		request.Type = v.(string)
	}

	if v, ok := d.GetOk("value"); ok && len(v.(string)) > 0 {
		request.Value = v.(string)
	}

	if v, ok := d.GetOk("ttl"); ok {
		ttl := int32(v.(int))
		request.Ttl = &ttl
	}

	if v, ok := d.GetOk("line"); ok && len(v.(string)) > 0 {
		line := v.(string)
		request.Line = &line
	}

	if v, ok := d.GetOk("description"); ok && len(v.(string)) > 0 {
		description := v.(string)
		request.Description = &description
	}

	if v, ok := d.GetOk("priority"); ok && len(v.(string)) > 0 {
		priority := int32(v.(int))
		request.Priority = &priority
	}

	return request
}

func listAllrecords(zoneName string, args *dns.ListRecordRequest, meta interface{}) ([]dns.Record, error) {
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

func resourceBaiduCloudDnsrecordUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.BaiduClient)

	zoneName := d.Get("zone_name").(string)

	recordId := d.Get("record_id").(string)

	action := "Update Dns record recordId is " + recordId

	updateArgs := buildBaiduCloudDnsrecordUpdateArgs(d)

	_, updateErr := client.WithDNSClient(func(dnsClient *dns.Client) (interface{}, error) {
		return nil, dnsClient.UpdateRecord(zoneName, recordId, updateArgs, buildClientToken())
	})

	if updateErr != nil {

		readErr := resourceBaiduCloudDnsrecordRead(d, meta)
		addDebug(action, readErr)

		if IsExceptedErrors(updateErr, ObjectNotFound) {
			return nil
		}
		return WrapErrorf(updateErr, DefaultErrorMsg, "baiducloud_dns_record", action, BCESDKGoERROR)

	}

	addDebug(action, updateErr)

	readErr := resourceBaiduCloudDnsrecordRead(d, meta)
	addDebug(action, readErr)

	if d.HasChange("record_action") {

		recordAction := d.Get("record_action").(string)

		recordId := d.Get("record_id").(string)

		action := "Update Dns record action is " + recordAction + " record id is " + recordId

		if "enable" == recordAction {
			_, updateErr := client.WithDNSClient(func(dnsClient *dns.Client) (interface{}, error) {
				return nil, dnsClient.UpdateRecordEnable(zoneName, recordId, buildClientToken())

			})

			addDebug(action, updateErr)
			if updateErr != nil {

				readErr := resourceBaiduCloudDnsrecordRead(d, meta)
				addDebug(action, readErr)
				if IsExceptedErrors(updateErr, ObjectNotFound) {
					return nil
				}
				return WrapErrorf(updateErr, DefaultErrorMsg, "baiducloud_dns_record", action, BCESDKGoERROR)

			}

		} else if "disable" == recordAction {
			_, updateErr := client.WithDNSClient(func(dnsClient *dns.Client) (interface{}, error) {
				return nil, dnsClient.UpdateRecordDisable(zoneName, recordId, buildClientToken())

			})

			addDebug(action, updateErr)
			if updateErr != nil {

				readErr := resourceBaiduCloudDnsrecordRead(d, meta)
				addDebug(action, readErr)
				if IsExceptedErrors(updateErr, ObjectNotFound) {
					return nil
				}
				return WrapErrorf(updateErr, DefaultErrorMsg, "baiducloud_dns_record", action, BCESDKGoERROR)

			}
		}

	}

	return resourceBaiduCloudDnsrecordRead(d, meta)
}

func buildBaiduCloudDnsrecordUpdateArgs(d *schema.ResourceData) *dns.UpdateRecordRequest {

	request := &dns.UpdateRecordRequest{}

	if v, ok := d.GetOk("rr"); ok && len(v.(string)) > 0 {
		request.Rr = v.(string)
	}

	if v, ok := d.GetOk("type"); ok && len(v.(string)) > 0 {
		request.Type = v.(string)
	}

	if v, ok := d.GetOk("value"); ok && len(v.(string)) > 0 {
		request.Value = v.(string)
	}

	if v, ok := d.GetOk("ttl"); ok {
		ttl := int32(v.(int))
		request.Ttl = &ttl
	}

	if v, ok := d.GetOk("description"); ok && len(v.(string)) > 0 {
		description := v.(string)
		request.Description = &description
	}

	if v, ok := d.GetOk("priority"); ok && len(v.(string)) > 0 {
		priority := int32(v.(int))
		request.Priority = &priority
	}

	return request

}
