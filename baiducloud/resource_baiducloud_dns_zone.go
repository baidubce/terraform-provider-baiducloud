/*
Provide a resource to create an Dns zone.

Example Usage

```hcl
resource "baiducloud_dns_zone" "default" {
  name              = "testDnsZone"
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

func resourceBaiduCloudDnsZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudDnszoneCreate,
		Read:   resourceBaiduCloudDnszoneRead,
		Delete: resourceBaiduCloudDnszoneDelete,

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
	}
}

func resourceBaiduCloudDnszoneCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.BaiduClient)

	createDnsArgs := buildBaiduCloudCreatednszoneArgs(d)

	action := "Create Dns zone " + createDnsArgs.Name

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {

		raw, err := client.WithDNSClient(func(dnsClient *dns.Client) (interface{}, error) {
			return nil, dnsClient.CreateZone(createDnsArgs, buildClientToken())
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
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dns_zone", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudDnszoneRead(d, meta)
}

func resourceBaiduCloudDnszoneRead(d *schema.ResourceData, meta interface{}) error {

	action := "Query DNS zone "

	queryArgs := buildBaiduCloudCreatednszoneQueryArgs(d)

	zones, err := listAllZones(queryArgs, meta)

	addDebug(action, zones)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dns_zone", action, BCESDKGoERROR)
	}

	var zone dns.Zone
	for _, z := range zones {
		if z.Name == d.Get("name") {
			zone = z
			break
		}
	}

	d.Set("zone_id", zone.Id)

	d.Set("name", zone.Name)

	d.Set("status", zone.Status)

	d.Set("product_version", zone.ProductVersion)

	d.Set("create_time", zone.CreateTime)

	d.Set("expire_time", zone.ExpireTime)

	d.Set("tags", tagsToMap(zone.Tags))

	return nil
}

func resourceBaiduCloudDnszoneDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.BaiduClient)

	zoneName := d.Get("name").(string)

	action := "Delete dns zone name IS " + zoneName

	_, err := client.WithDNSClient(func(dnsClient *dns.Client) (interface{}, error) {
		return nil, dnsClient.DeleteZone(zoneName, buildClientToken())
	})

	if err != nil {
		if IsExceptedErrors(err, ObjectNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dns_zone", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudCreatednszoneQueryArgs(d *schema.ResourceData) *dns.ListZoneRequest {

	request := &dns.ListZoneRequest{}

	if v, ok := d.GetOk("name"); ok && len(v.(string)) > 0 {
		request.Name = v.(string)
	}

	return request
}

func buildBaiduCloudCreatednszoneArgs(d *schema.ResourceData) *dns.CreateZoneRequest {

	request := &dns.CreateZoneRequest{}

	if v, ok := d.GetOk("name"); ok && len(v.(string)) > 0 {
		request.Name = v.(string)
	}

	return request
}

func listAllZones(args *dns.ListZoneRequest, meta interface{}) ([]dns.Zone, error) {
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

func tagsToMap(tags []dns.TagModel) map[string]string {
	tagMap := make(map[string]string)
	for _, tag := range tags {
		tagMap[*tag.TagKey] = *tag.TagValue
	}

	return tagMap
}
