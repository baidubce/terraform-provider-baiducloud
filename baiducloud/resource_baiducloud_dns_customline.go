/*
Provide a resource to create an Dns customline.

Example Usage

```hcl
resource "baiducloud_dns_customline" "default" {
 name              = "testDnscustomline"
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

func resourceBaiduCloudDnsCustomline() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudDnscustomlineCreate,
		Read:   resourceBaiduCloudDnscustomlineRead,
		Update: resourceBaiduCloudDnscustomlineUpdate,
		Delete: resourceBaiduCloudDnscustomlineDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Dns customline name",
				Required:    true,
				ForceNew:    true,
			},
			"lines": {
				Type:        schema.TypeSet,
				Description: "lines of dns ",
				Required:    true,
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
				Type:        schema.TypeString,
				Description: "Dns customline related zone count",
				Computed:    true,
			},
			"related_record_count": {
				Type:        schema.TypeString,
				Description: "Dns customline related record count",
				Computed:    true,
			},
		},
	}
}

func resourceBaiduCloudDnscustomlineCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.BaiduClient)

	createArgs := buildBaiduCloudCreatednscustomlineArgs(d)

	action := "Create Dns customline " + createArgs.Name

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {

		raw, err := client.WithDNSClient(func(dnsClient *dns.Client) (interface{}, error) {
			return nil, dnsClient.AddLineGroup(createArgs, buildClientToken())
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
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dns_customline", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudDnscustomlineRead(d, meta)
}

func resourceBaiduCloudDnscustomlineRead(d *schema.ResourceData, meta interface{}) error {

	action := "Query DNS customline "

	queryArgs := buildBaiduCloudCreatednscustomlineQueryArgs(d)

	customlines, err := listAllcustomlines(queryArgs, meta)

	addDebug(action, customlines)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dns_customline", action, BCESDKGoERROR)
	}

	var customline dns.Line
	for _, z := range customlines {
		if z.Name == d.Get("name") {
			customline = z
			break
		}
	}

	d.Set("line_id", customline.Id)

	d.Set("related_zone_count", customline.RelatedZoneCount)

	d.Set("related_record_count", customline.RelatedRecordCount)

	return nil
}

func resourceBaiduCloudDnscustomlineUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.BaiduClient)

	if d.HasChange("lines") {

		lineId := d.Get("line_id").(string)

		action := "Update Dns customline lineId is " + lineId

		updateArgs := buildBaiduCloudDnsLineUpdateLinesArgs(d)

		_, updateErr := client.WithDNSClient(func(dnsClient *dns.Client) (interface{}, error) {
			return nil, dnsClient.UpdateLineGroup(lineId, updateArgs, buildClientToken())
		})

		if updateErr != nil {
			if IsExceptedErrors(updateErr, ObjectNotFound) {
				return nil
			}
			return WrapErrorf(updateErr, DefaultErrorMsg, "baiducloud_dns_customline", action, BCESDKGoERROR)

		}

		addDebug(action, updateErr)
	}

	if d.HasChange("name") {

		lineId := d.Get("line_id").(string)

		action := "Update Dns customline lineId is " + lineId

		updateArgs := buildBaiduCloudDnsLineUpdateNameArgs(d)

		_, updateErr := client.WithDNSClient(func(dnsClient *dns.Client) (interface{}, error) {
			return nil, dnsClient.UpdateLineGroup(lineId, updateArgs, buildClientToken())
		})

		if updateErr != nil {
			if IsExceptedErrors(updateErr, ObjectNotFound) {
				return nil
			}
			return WrapErrorf(updateErr, DefaultErrorMsg, "baiducloud_dns_customline", action, BCESDKGoERROR)

		}

		addDebug(action, updateErr)
	}

	return resourceBaiduCloudDnscustomlineRead(d, meta)
}

func resourceBaiduCloudDnscustomlineDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.BaiduClient)

	lineId := d.Get("line_id").(string)

	action := "Delete dns customline lineId IS " + lineId

	_, err := client.WithDNSClient(func(dnsClient *dns.Client) (interface{}, error) {
		return nil, dnsClient.DeleteLineGroup(lineId, buildClientToken())
	})

	if err != nil {
		if IsExceptedErrors(err, ObjectNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dns_customline", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudCreatednscustomlineQueryArgs(d *schema.ResourceData) *dns.ListLineGroupRequest {

	request := &dns.ListLineGroupRequest{}

	return request
}

func buildBaiduCloudCreatednscustomlineArgs(d *schema.ResourceData) *dns.AddLineGroupRequest {

	request := &dns.AddLineGroupRequest{}

	if v, ok := d.GetOk("name"); ok && len(v.(string)) > 0 {
		request.Name = v.(string)
	}

	if lines, ok := d.GetOk("lines"); ok {

		group_lines := make([]string, 0)
		for _, ip := range lines.(*schema.Set).List() {
			group_lines = append(group_lines, ip.(string))
		}
		request.Lines = group_lines
	}

	return request
}

func buildBaiduCloudDnsLineUpdateLinesArgs(d *schema.ResourceData) *dns.UpdateLineGroupRequest {

	request := &dns.UpdateLineGroupRequest{}

	if v, ok := d.GetOk("name"); ok && len(v.(string)) > 0 {
		request.Name = v.(string)
	}

	if lines, ok := d.GetOk("lines"); ok {

		group_lines := make([]string, 0)
		for _, ip := range lines.(*schema.Set).List() {
			group_lines = append(group_lines, ip.(string))
		}
		request.Lines = group_lines
	}

	return request
}

func buildBaiduCloudDnsLineUpdateNameArgs(d *schema.ResourceData) *dns.UpdateLineGroupRequest {

	request := &dns.UpdateLineGroupRequest{}

	if v, ok := d.GetOk("name"); ok && len(v.(string)) > 0 {
		request.Name = v.(string)
	}

	if lines, ok := d.GetOk("lines"); ok {

		group_lines := make([]string, 0)
		for _, ip := range lines.(*schema.Set).List() {
			group_lines = append(group_lines, ip.(string))
		}
		request.Lines = group_lines
	}

	return request
}

func listAllcustomlines(args *dns.ListLineGroupRequest, meta interface{}) ([]dns.Line, error) {
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
