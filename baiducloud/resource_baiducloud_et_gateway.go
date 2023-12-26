/*
Provide a resource to manage an et gateway.

Example Usage

```hcl
resource "baiducloud_et_gateway" "default" {
  name = "my_name"
  vpc_id = "vpc-xxx"
  speed = 200
  description = "description"
  et_id = "xxx"
  channel_id = "xxx"
  local_cidrs = ["192.168.0.0/20"]
}
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

func resourceBaiduCloudEtGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudEtGatewayCreate,
		Read:   resourceBaiduCloudEtGatewayRead,
		Update: resourceBaiduCloudEtGatewayUpdate,
		Delete: resourceBaiduCloudEtGatewayDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name of the et gateway",
				Required:    true,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "vpc id of the et gateway",
				Required:    true,
				ForceNew:    true,
			},
			"speed": {
				Type:        schema.TypeInt,
				Description: "speed of the et gateway (Mbps)",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "description of the et gateway",
				Optional:    true,
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
			"et_gateway_id": {
				Type:        schema.TypeString,
				Description: "ID of et gateway.",
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

func resourceBaiduCloudEtGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	createArgs, err := buildBaiduCloudEtGatewayArgs(d, meta)

	if err != nil {
		return WrapError(err)
	}

	action := "Create ET gateway " + createArgs.Name

	addDebug(action, createArgs)

	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithEtGatewayClient(func(etGatewayClient *etGateway.Client) (interface{}, error) {

			return etGatewayClient.CreateEtGateway(createArgs)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)

		result, _ := raw.(*etGateway.CreateEtGatewayResult)
		d.SetId(result.EtGatewayId)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_et_gateway", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudEtGatewayRead(d, meta)
}

func resourceBaiduCloudEtGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	etGatewayId := d.Id()
	action := "Query et gatewaty etGatewayId is " + etGatewayId

	raw, err := client.WithEtGatewayClient(func(etGatewayClient *etGateway.Client) (interface{}, error) {
		return etGatewayClient.GetEtGatewayDetail(etGatewayId)
	})

	addDebug(action, raw)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_et_gateway", action, BCESDKGoERROR)
	}

	result, _ := raw.(*etGateway.EtGatewayDetail)

	d.Set("name", result.Name)

	d.Set("status", result.Status)

	d.Set("speed", result.Speed)

	d.Set("create_time", result.CreateTime)

	d.Set("description", result.Description)

	d.Set("vpc_id", result.VpcId)

	d.Set("et_id", result.EtId)

	d.Set("channel_id", result.ChannelId)

	d.Set("local_cidrs", result.LocalCidrs)

	d.Set("health_check_source_ip", result.HealthCheckSourceIp)

	d.Set("health_check_dest_ip", result.HealthCheckDestIp)

	d.Set("health_check_type", result.HealthCheckType)

	d.Set("health_check_interval", result.HealthCheckInterval)

	d.Set("health_threshold", result.HealthThreshold)

	d.Set("unhealth_threshold", result.UnhealthThreshold)

	return nil
}

func resourceBaiduCloudEtGatewayUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.BaiduClient)
	etGatewayId := d.Id()

	updateArgs, err := buildBaiduCloudEtGatewayUpdateArgs(d, meta)

	if err != nil {
		return WrapError(err)
	}

	updateArgs.EtGatewayId = etGatewayId
	action := "Update et gateway etGatewayId is" + etGatewayId
	addDebug(action, updateArgs)

	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err := client.WithEtGatewayClient(func(etGatewayClient *etGateway.Client) (interface{}, error) {
			return nil, etGatewayClient.UpdateEtGateway(updateArgs)
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
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_et_gateway", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudEtGatewayRead(d, meta)
}

func resourceBaiduCloudEtGatewayDelete(d *schema.ResourceData, meta interface{}) error {
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
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_et_gateway", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudEtGatewayArgs(d *schema.ResourceData, meta interface{}) (*etGateway.CreateEtGatewayArgs, error) {
	request := &etGateway.CreateEtGatewayArgs{
		ClientToken: buildClientToken(),
	}

	if v := d.Get("name").(string); v != "" {
		request.Name = v
	}

	if v := d.Get("vpc_id").(string); v != "" {
		request.VpcId = v
	}

	if v := d.Get("speed").(int); v != 0 {
		request.Speed = v
	}

	if v := d.Get("description").(string); v != "" {
		request.Description = v
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

func buildBaiduCloudEtGatewayUpdateArgs(d *schema.ResourceData, meta interface{}) (*etGateway.UpdateEtGatewayArgs, error) {
	request := &etGateway.UpdateEtGatewayArgs{
		ClientToken: buildClientToken(),
	}

	if v := d.Get("name").(string); v != "" {
		request.Name = v
	}

	if v := d.Get("speed").(int); v != 0 {
		request.Speed = v
	}

	if v := d.Get("description").(string); v != "" {
		request.Description = v
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
