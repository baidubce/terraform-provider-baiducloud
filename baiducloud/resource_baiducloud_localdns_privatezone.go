/*
Use this resource to get information about a Local Dns PrivateZone.

~> **NOTE:** The terminate operation of PrivateZone does NOT take effect immediately，maybe takes for several minites.

Example Usage

```hcl
resource "baiducloud_localdns_privatezone" "my-server" {
  zone_name = "terrraform.com"
}
```

Import

Local Dns PrivateZone can be imported, e.g.

```hcl
$ terraform import baiducloud_localdns_privatezone.my-server id
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/localDns"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudLocalDnsPrivateZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudDnsLocalPrivateZoneCreate,
		Read:   resourceBaiduCloudDnsLocalPrivateZoneRead,
		Delete: resourceBaiduCloudDnsLocalPrivateZoneDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"zone_name": {
				Type:        schema.TypeString,
				Description: "name of the DNS local PrivateZone",
				ForceNew:    true,
				Required:    true,
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
	}
}

func resourceBaiduCloudDnsLocalPrivateZoneCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	args := buildBaiduCloudDnsLocalPrivateZoneArgs(d)
	action := "Create DNS Local Private Zone " + args.ZoneName

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithLocalDnsClient(func(localDnsClient *localDns.Client) (i interface{}, e error) {
			return localDnsClient.CreatePrivateZone(args)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)
		result, _ := raw.(*localDns.CreatePrivateZoneResponse)
		d.SetId(result.ZoneId)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_localdns_pravitezone", action, BCESDKGoERROR)
	}
	return resourceBaiduCloudDnsLocalPrivateZoneRead(d, meta)

}

func resourceBaiduCloudDnsLocalPrivateZoneRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	localDnsService := LocalDnsService{client}

	zoneId := d.Id()

	action := "Query DNS Local Private Zone " + zoneId

	privateZone, err := localDnsService.GetPrivateZoneDetail(zoneId)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_localdns_privatezone", action, BCESDKGoERROR)
	}
	addDebug(action, privateZone)

	d.Set("record_count", privateZone.RecordCount)
	d.Set("create_time", privateZone.CreateTime)
	d.Set("update_time", privateZone.UpdateTime)

	return nil
}

func resourceBaiduCloudDnsLocalPrivateZoneDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	zoneId := d.Id()

	action := "Delete DNS Local Private Zone " + zoneId

	clientToken := buildClientToken()

	_, err := client.WithLocalDnsClient(func(localDnsClient *localDns.Client) (i interface{}, e error) {
		return nil, localDnsClient.DeletePrivateZone(zoneId, clientToken)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_localdns_privatezone", action, BCESDKGoERROR)
	}
	addDebug(action, nil)

	return nil
}

func buildBaiduCloudDnsLocalPrivateZoneArgs(d *schema.ResourceData) *localDns.CreatePrivateZoneRequest {
	args := &localDns.CreatePrivateZoneRequest{
		ClientToken: buildClientToken(),
	}

	if v := d.Get("zone_name").(string); v != "" {
		args.ZoneName = v
	}

	return args
}
