/*
Use this resource to get information about a Local Dns VPC.

~> **NOTE:** The terminate operation of vpc does NOT take effect immediately，maybe takes for several minites.

Example Usage

```hcl
resource "baiducloud_localdns_vpc" "my-server" {
 zone_name = "terrraform.com"
}
```

Import

Local Dns vpc can be imported, e.g.

```hcl
$ terraform import baiducloud_localdns_vpc.my-server id
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

func resourceBaiduCloudLocalDnsVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudLocalDnsVpcCreate,
		Read:   resourceBaiduCloudLocalDnsVpcRead,
		Delete: resourceBaiduCloudLocalDnsVpcDelete,

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
				Description: "zone_id of the DNS privatezone ",
				ForceNew:    true,
				Required:    true,
			},
			"region": {
				Type:        schema.TypeString,
				Description: "region of the DNS  vpc",
				ForceNew:    true,
				Required:    true,
			},
			"vpc_ids": {
				Type:        schema.TypeSet,
				Description: "vpc_ids  of the DNS  vpc.",
				ForceNew:    true,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
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

func resourceBaiduCloudLocalDnsVpcCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	zoneId := d.Get("zone_id").(string)
	args := buildBaiduCloudLocalDnsVpcArgs(d)

	action := "bind local dns Private Zone vpcs " + zoneId

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithLocalDnsClient(func(localDnsClient *localDns.Client) (i interface{}, e error) {
			return nil, localDnsClient.BindVpc(zoneId, args)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)
		d.SetId(zoneId)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_localdns_pravitezone", action, BCESDKGoERROR)
	}
	return resourceBaiduCloudLocalDnsVpcRead(d, meta)

}

func resourceBaiduCloudLocalDnsVpcRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	localDnsService := LocalDnsService{client}

	zoneId := d.Id()

	action := "Query DNS Local VPCS " + zoneId

	zone, err := localDnsService.GetPrivateZoneDetail(zoneId)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_localdns_vpc", action, BCESDKGoERROR)
	}
	addDebug(action, zone)

	vpcs := zone.BindVpcs
	bindVpcs := make([]interface{}, 0, len(vpcs))
	for _, vpc := range vpcs {
		vpcMap := make(map[string]interface{})
		vpcMap["vpc_id"] = vpc.VpcId
		vpcMap["vpc_name"] = vpc.VpcName
		vpcMap["vpc_region"] = vpc.VpcRegion

		bindVpcs = append(bindVpcs, vpcMap)
	}
	d.Set("bind_vpcs", bindVpcs)

	return nil
}

func resourceBaiduCloudLocalDnsVpcDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	zoneId := d.Id()

	action := "Unbind local dns vpcs " + zoneId
	args := buildBaiduCloudLocalDnsUnbindVpcArgs(d)

	_, err := client.WithLocalDnsClient(func(localDnsClient *localDns.Client) (i interface{}, e error) {
		return nil, localDnsClient.UnbindVpc(zoneId, args)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_localdns_vpc", action, BCESDKGoERROR)
	}
	addDebug(action, nil)

	return nil
}

func buildBaiduCloudLocalDnsVpcArgs(d *schema.ResourceData) *localDns.BindVpcRequest {
	args := &localDns.BindVpcRequest{
		ClientToken: buildClientToken(),
	}

	if v := d.Get("region").(string); v != "" {
		args.Region = v
	}

	vpcIds := make([]string, 0)
	ids, ok := d.GetOk("vpc_ids")
	if ok {
		for _, id := range ids.(*schema.Set).List() {
			vpcIds = append(vpcIds, id.(string))
		}
		args.VpcIds = vpcIds
	}

	return args
}

func buildBaiduCloudLocalDnsUnbindVpcArgs(d *schema.ResourceData) *localDns.UnbindVpcRequest {
	args := &localDns.UnbindVpcRequest{
		ClientToken: buildClientToken(),
	}

	if v := d.Get("region").(string); v != "" {
		args.Region = v
	}

	vpcIds := make([]string, 0)
	ids, ok := d.GetOk("vpc_ids")
	if ok {
		for _, id := range ids.(*schema.Set).List() {
			vpcIds = append(vpcIds, id.(string))
		}
		args.VpcIds = vpcIds
	}

	return args
}
