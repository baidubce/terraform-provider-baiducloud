/*
Use this data source to query et gateways.

Example Usage

```hcl
data "baiducloud_et_gateways" "default" {
	vpc_id = "xxxxx"
}

output "gateways" {
 value = "${data.baiducloud_et_gateways.default.gateways}"
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

func dataSourceBaiduCloudEtGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudEtGatewayRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "ID of the instance",
				Required:    true,
				ForceNew:    true,
			},
			"et_gateway_id": {
				Type:        schema.TypeString,
				Description: "ID of et gateway.",
				Optional:    true,
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "name of et gateway.",
				Optional:    true,
				ForceNew:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "status of et gateway.",
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

			"gateways": {
				Type:        schema.TypeList,
				Description: "et gateway",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"et_gateway_id": {
							Type:        schema.TypeString,
							Description: "ID of et gateway.",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "name of the et gateway",
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
						"description": {
							Type:        schema.TypeString,
							Description: "description of the et gateway",
							Computed:    true,
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Description: "vpc id of the et gateway",
							Computed:    true,
						},
						"et_id": {
							Type:        schema.TypeString,
							Description: "et id of the et gateway",
							Computed:    true,
						},
						"channel_id": {
							Type:        schema.TypeString,
							Description: "channel id of the et gateway",
							Required:    true,
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
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudEtGatewayRead(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*connectivity.BaiduClient)

	vpcId := d.Get("vpc_id").(string)
	action := "List all etGateways " + vpcId

	gatewayArgs, err := buildBaiduCloudEtGatewayListArgs(d, meta)

	if err != nil {
		return WrapError(err)
	}

	etGateways, err := listAllEtGateways(gatewayArgs, meta)

	addDebug(action, err)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_et_gateways", action, BCESDKGoERROR)
	}

	gatewaysResult := make([]map[string]interface{}, 0)
	for _, gateway := range etGateways {

		gatewayMap := make(map[string]interface{})
		gatewayMap["et_gateway_id"] = gateway.EtGatewayId
		gatewayMap["name"] = gateway.Name
		gatewayMap["status"] = gateway.Status
		gatewayMap["speed"] = gateway.Speed
		gatewayMap["create_time"] = gateway.CreateTime
		gatewayMap["description"] = gateway.Description
		gatewayMap["vpc_id"] = gateway.VpcId
		gatewayMap["et_id"] = gateway.EtId
		gatewayMap["channel_id"] = gateway.ChannelId
		gatewayMap["local_cidrs"] = gateway.LocalCidrs

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_et_gateways", action, BCESDKGoERROR)
		}

		gatewaysResult = append(gatewaysResult, gatewayMap)
	}
	addDebug(action, gatewaysResult)

	FilterDataSourceResult(d, &gatewaysResult)

	if err := d.Set("gateways", gatewaysResult); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_et_gateways", action, BCESDKGoERROR)
	}

	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), gatewaysResult); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_et_gateways", action, BCESDKGoERROR)
		}
	}
	return nil
}

func listAllEtGateways(args *etGateway.ListEtGatewayArgs, meta interface{}) ([]etGateway.EtGateway, error) {
	client := meta.(*connectivity.BaiduClient)
	action := "List all etGateways"

	etGateways := make([]etGateway.EtGateway, 0)
	for {
		raw, err := client.WithEtGatewayClient(func(etGatewayClient *etGateway.Client) (interface{}, error) {
			return etGatewayClient.ListEtGateway(args)
		})

		if err != nil {
			return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_et_gateways", action, BCESDKGoERROR)
		}

		result, _ := raw.(*etGateway.ListEtGatewayResult)
		etGateways = append(etGateways, result.EtGateways...)
		if !result.IsTruncated {
			break
		}
		args.Marker = result.NextMarker
		args.MaxKeys = result.MaxKeys
	}

	return etGateways, nil
}

func buildBaiduCloudEtGatewayListArgs(d *schema.ResourceData, meta interface{}) (*etGateway.ListEtGatewayArgs, error) {
	request := &etGateway.ListEtGatewayArgs{}

	if v := d.Get("vpc_id").(string); v != "" {
		request.VpcId = v
	}

	if v := d.Get("et_gateway_id").(string); v != "" {
		request.EtGatewayId = v
	}

	if v := d.Get("name").(string); v != "" {
		request.Name = v
	}

	if v := d.Get("status").(string); v != "" {
		request.Status = v
	}

	return request, nil

}
