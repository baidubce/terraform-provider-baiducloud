/*
Provide a resource to manage an et gateway association.

Example Usage

```hcl
resource "baiducloud_et_gateway_association" "default" {
  et_gateway_id = "xxx"
  et_id = "xxx"
  channel_id = "xxx"
  local_cidrs = ["192.168.0.0/20"]
}
```

Import

ET Gateway Association can be imported, e.g.

```hcl
$ terraform import baiducloud_et_gateway_association.default et_gateway_id
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/etGateway"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudEtGatewayAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudEtGatewayAssociationCreate,
		Read:   resourceBaiduCloudEtGatewayAssociationRead,
		Delete: resourceBaiduCloudEtGatewayAssociationDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"et_gateway_id": {
				Type:        schema.TypeString,
				Description: "ID of et gateway.",
				Required:    true,
				ForceNew:    true,
			},
			"et_id": {
				Type:        schema.TypeString,
				Description: "et id of the et gateway",
				Optional:    true,
				ForceNew:    true,
			},
			"channel_id": {
				Type:        schema.TypeString,
				Description: "channel id of the et gateway",
				Optional:    true,
				ForceNew:    true,
			},
			"local_cidrs": {
				Type:        schema.TypeSet,
				Description: "local cidrs of the et gateway",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name": {
				Type:        schema.TypeString,
				Description: "name of et gateway.",
				Computed:    true,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "vpc id of et gateway.",
				Computed:    true,
			},
			"speed": {
				Type:        schema.TypeInt,
				Description: "speed of et gateway.",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "status of et gateway.",
				Computed:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "create_time of et gateway.",
				Computed:    true,
			},
			"health_check_source_ip": {
				Type:        schema.TypeString,
				Description: "health_check_source_ip of et gateway.",
				Computed:    true,
			},
			"health_check_dest_ip": {
				Type:        schema.TypeString,
				Description: "health_check_dest_ip of et gateway.",
				Computed:    true,
			},
			"health_check_type": {
				Type:        schema.TypeString,
				Description: "health_check_type of et gateway.",
				Computed:    true,
			},
			"health_check_interval": {
				Type:        schema.TypeInt,
				Description: "health_check_interval of et gateway.",
				Computed:    true,
			},
			"health_threshold": {
				Type:        schema.TypeInt,
				Description: "health_threshold of et gateway.",
				Computed:    true,
			},
			"unhealth_threshold": {
				Type:        schema.TypeInt,
				Description: "unhealth_threshold of et gateway.",
				Computed:    true,
			},
		},
	}
}

func resourceBaiduCloudEtGatewayAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	createArgs, err := buildBaiduCloudEtGatewayAssociationArgs(d, meta)

	if err != nil {
		return WrapError(err)
	}

	action := "Create ET gateway Association" + createArgs.EtGatewayId

	addDebug(action, createArgs)

	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err := client.WithEtGatewayClient(func(etGatewayClient *etGateway.Client) (interface{}, error) {

			return nil, etGatewayClient.BindEt(createArgs)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, err)

		d.SetId(createArgs.EtGatewayId)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_et_gateway_association", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudEtGatewayRead(d, meta)
}

func resourceBaiduCloudEtGatewayAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	etGatewayId := d.Id()
	action := "Query et gatewaty association info etGatewayId is " + etGatewayId

	raw, err := client.WithEtGatewayClient(func(etGatewayClient *etGateway.Client) (interface{}, error) {
		return etGatewayClient.GetEtGatewayDetail(etGatewayId)
	})

	addDebug(action, raw)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_et_gateway_association", action, BCESDKGoERROR)
	}

	result, _ := raw.(*etGateway.EtGatewayDetail)

	d.Set("status", result.Status)

	d.Set("speed", result.Speed)

	d.Set("description", result.Description)

	d.Set("vpc_id", result.VpcId)

	d.Set("create_time", result.CreateTime)

	d.Set("health_check_source_ip", result.HealthCheckSourceIp)

	d.Set("health_check_dest_ip", result.HealthCheckDestIp)

	d.Set("health_check_type", result.HealthCheckType)

	d.Set("health_check_interval", result.HealthCheckInterval)

	d.Set("health_threshold", result.HealthThreshold)

	d.Set("unhealth_threshold", result.UnhealthThreshold)

	return nil
}

func resourceBaiduCloudEtGatewayAssociationDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.BaiduClient)

	etGatewayId := d.Id()

	action := "Delete et gateway etGatewayId id is" + etGatewayId

	addDebug(action, "")

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := client.WithEtGatewayClient(func(etGatewayClient *etGateway.Client) (interface{}, error) {

			return nil, etGatewayClient.DeleteEtGateway(etGatewayId, buildClientToken())
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, err)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_et_gateway_association", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudEtGatewayAssociationArgs(d *schema.ResourceData, meta interface{}) (*etGateway.BindEtArgs, error) {
	request := &etGateway.BindEtArgs{
		ClientToken: buildClientToken(),
	}

	if v := d.Get("et_gateway_id").(string); v != "" {
		request.EtGatewayId = v
	}

	if v := d.Get("et_id").(string); v != "" {
		request.EtId = v
	}

	if v := d.Get("channel_id").(string); v != "" {
		request.ChannelId = v
	}

	if localCidrs, ok := d.GetOk("local_cidrs"); ok {

		cidrs := make([]string, 0)

		for _, ip := range localCidrs.(*schema.Set).List() {
			cidrs = append(cidrs, ip.(string))
		}
		request.LocalCidrs = cidrs
	}

	return request, nil

}
