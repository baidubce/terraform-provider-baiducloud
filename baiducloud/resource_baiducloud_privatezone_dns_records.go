package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/localDns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"time"
)

const (
	MXType = "MX"
)

func resourceBaiduCloudPrivateZoneRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudPrivateZoneDNSRecordCreate,
		Read:   resourceBaiduCloudPrivateZoneRecordRead,
		Update: resourceBaiduCloudPrivateZoneRecordUpdate,
		Delete: resourceBaiduCloudPrivateZoneRecordDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Description: "Id of private zone.",
				Required:    true,
			},
			"record_id": {
				Type:        schema.TypeString,
				Description: "Id of private zone dns record",
				Computed:    true,
			},
			"rr": {
				Type:        schema.TypeString,
				Description: "rr",
				Required:    true,
			},
			"value": {
				Type:        schema.TypeString,
				Description: "record value",
				Required:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "record type",
				Required:    true,
			},
			"priority": {
				Type:        schema.TypeInt,
				Description: "priority,when the record type is NOT MX,priority should be 0",
				Optional:    true,
			},
			"ttl": {
				Type:        schema.TypeInt,
				Description: "ttl",
				Optional:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "description of record",
				Optional:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "status of record,enable or disable",
				Optional:    true,
			},
		},
	}
}
func resourceBaiduCloudPrivateZoneDNSRecordCreate(d *schema.ResourceData, meta interface{}) error {
	action := "Create Private Zone DNS Record"
	client := meta.(*connectivity.BaiduClient)
	err := inputFieldCheck(d)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_private_zone_dns_record", action, BCESDKGoERROR)
	}
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithLocalDnsClient(func(recordClient *localDns.Client) (interface{}, error) {
			request := &localDns.AddRecordRequest{
				ClientToken: buildClientToken(),
			}
			var zoneId string
			//var createArgs interface{}
			if v, ok := d.GetOk("zone_id"); ok {
				zoneId = v.(string)
			}
			if v, ok := d.GetOk("rr"); ok {
				request.Rr = v.(string)
			}
			if v, ok := d.GetOk("value"); ok {
				request.Value = v.(string)
			}
			if v, ok := d.GetOk("type"); ok {
				request.Type = v.(string)
			}
			if v, ok := d.GetOk("description"); ok {
				request.Description = v.(string)
			}
			if v, ok := d.GetOk("ttl"); ok {
				request.Ttl = int32(v.(int))
			}
			if v, ok := d.GetOk("priority"); ok {
				if request.Type == MXType {
					request.Priority = int32(v.(int))
				} else {
					request.Priority = 0
				}
			}
			resp, err := recordClient.AddRecord(zoneId, request)
			if err != nil {
				return resp, err
			}
			recordId := resp.RecordId
			d.SetId(recordId)
			var status string
			if v, ok := d.GetOk("status"); ok {
				status = v.(string)
			}
			return resp, UpdateStatus(recordClient, recordId, status)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)

		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_private_zone_dns_record", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudPrivateZoneRecordRead(d, meta)
}
func resourceBaiduCloudPrivateZoneRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	recordId := d.Id()
	action := "Query local DNS record " + recordId

	raw, err := client.WithLocalDnsClient(func(localDnsClient *localDns.Client) (i interface{}, e error) {
		var zoneId string
		//var createArgs interface{}
		if v, ok := d.GetOk("zone_id"); ok {
			zoneId = v.(string)
		}
		return localDnsClient.ListRecord(zoneId, nil)
	})
	addDebug(action, raw)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_private_zone_dns_record", action, BCESDKGoERROR)
	}

	for _, record := range raw.(*localDns.ListRecordResponse).Records {
		if record.RecordId == recordId {
			d.Set("rr", record.Rr)
			d.Set("value", record.Value)
			d.Set("type", record.Type)
			d.Set("description", record.Description)
			d.Set("priority", record.Priority)
			d.Set("ttl", record.Ttl)
			d.Set("record_id", record.RecordId)
			d.Set("status", record.Status)
			break
		}
	}
	return nil
}
func resourceBaiduCloudPrivateZoneRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	recordId := d.Id()

	action := "Delete records" + recordId
	_, err := client.WithLocalDnsClient(func(localDnsClient *localDns.Client) (interface{}, error) {
		return recordId, localDnsClient.DeleteRecord(recordId, buildClientToken())
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dns_record", action, BCESDKGoERROR)
	}
	return nil
}
func resourceBaiduCloudPrivateZoneRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	recordId := d.Id()

	action := "Update record " + recordId
	err := inputFieldCheck(d)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_private_zone_dns_record", action, BCESDKGoERROR)
	}
	_, err = client.WithLocalDnsClient(func(localDnsClient *localDns.Client) (interface{}, error) {
		request := &localDns.UpdateRecordRequest{
			ClientToken: buildClientToken(),
		}
		if v, ok := d.GetOk("rr"); ok {
			request.Rr = v.(string)
		}
		if v, ok := d.GetOk("value"); ok {
			request.Value = v.(string)
		}
		if v, ok := d.GetOk("type"); ok {
			request.Type = v.(string)
		}
		if v, ok := d.GetOk("description"); ok {
			request.Description = v.(string)
		}
		if v, ok := d.GetOk("ttl"); ok {
			request.Ttl = int32(v.(int))
		}
		if v, ok := d.GetOk("priority"); ok {
			request.Priority = int32(v.(int))
		}
		var status string
		if v, ok := d.GetOk("status"); ok {
			status = v.(string)
		}
		err := localDnsClient.UpdateRecord(recordId, request)
		if err != nil {
			return recordId, err
		}
		return recordId, UpdateStatus(localDnsClient, recordId, status)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dns_record", action, BCESDKGoERROR)
	}
	return nil
}
