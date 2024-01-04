/*
Use this data source to query et gateway associations.

Example Usage

```hcl
data "baiducloud_et_gateway_associations" "default" {
	et_gateway_id = "xxxxx"
}

output "gateway" {
 value = "${data.baiducloud_et_gateway_associations.default.gateway_associations}"
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/etGateway"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudEtGatewayAssociations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudEtGatewayAssociationsRead,

		Schema: map[string]*schema.Schema{
			"et_gateway_id": {
				Type:        schema.TypeString,
				Description: "ID of et gateway.",
				Optional:    true,
				ForceNew:    true,
			},

			"output_file": {
				Type:        schema.TypeString,
				Description: "Query result output file path",
				Optional:    true,
				ForceNew:    true,
			},

			"filter": dataSourceFiltersSchema(),

			"gateway_associations": {
				Type:        schema.TypeList,
				Description: "et gateway associations",
				Computed:    true,
				Elem: &schema.Resource{
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
						"status": {
							Type:        schema.TypeString,
							Description: "status of et gateway.",
							Computed:    true,
						},
						"speed": {
							Type:        schema.TypeInt,
							Description: "speed of the et gateway (Mbps)",
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
				},
			},
		},
	}
}

func dataSourceBaiduCloudEtGatewayAssociationsRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.BaiduClient)

	etGatewayId := d.Get("et_gateway_id").(string)

	action := "Query et gatewaty associations etGatewayId is " + etGatewayId

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

	gatewaysResult := make([]map[string]interface{}, 0)

	gatewayMap := make(map[string]interface{})

	gatewayMap["name"] = result.Name

	gatewayMap["status"] = result.Status

	gatewayMap["speed"] = result.Speed

	gatewayMap["create_time"] = result.CreateTime

	gatewayMap["description"] = result.Description

	gatewayMap["vpc_id"] = result.VpcId

	gatewayMap["et_id"] = result.EtId

	gatewayMap["channel_id"] = result.ChannelId

	gatewayMap["local_cidrs"] = result.LocalCidrs

	gatewayMap["health_check_source_ip"] = result.HealthCheckSourceIp

	gatewayMap["health_check_dest_ip"] = result.HealthCheckDestIp

	gatewayMap["health_check_type"] = result.HealthCheckType

	gatewayMap["health_check_interval"] = result.HealthCheckInterval

	gatewayMap["health_threshold"] = result.HealthThreshold

	gatewayMap["unhealth_threshold"] = result.UnhealthThreshold

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_et_gateway_associations", action, BCESDKGoERROR)
	}

	gatewaysResult = append(gatewaysResult, gatewayMap)

	addDebug(action, gatewaysResult)

	FilterDataSourceResult(d, &gatewaysResult)

	if err := d.Set("gateway_associations", gatewaysResult); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_et_gateway_associations", action, BCESDKGoERROR)
	}

	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), gatewaysResult); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_et_gateway_associations", action, BCESDKGoERROR)
		}
	}
	return nil
}

