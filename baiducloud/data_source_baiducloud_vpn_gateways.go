/*
Use this data source to query VPN gateway list.

Example Usage

```hcl
data "baiducloud_vpn_gateways" "default" {
  vpc_id = "vpc-65cz3hu92kz2"
}

output "vpns" {
  value = "${data.baiducloud_vpn_gateways.default.vpns}"
}
```
*/
package baiducloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudVpnGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudVpnGatewayRead,

		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:        schema.TypeString,
				Description: "Output file for saving result.",
				Optional:    true,
				ForceNew:    true,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "ID of the VPC which vpn gateway belong to.",
				Required:    true,
			},
			"filter": dataSourceFiltersSchema(),

			"vpn_gateways": {
				Type:        schema.TypeList,
				Description: "Result of VPCs.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpn_id": {
							Type:        schema.TypeString,
							Description: "ID of the VPN gateway.",
							Computed:    true,
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Description: "ID of the VPC which vpn gateway belong to.",
							Computed:    true,
						},
						"vpn_name": {
							Type:        schema.TypeString,
							Description: "Name of the VPN gateway.",
							Computed:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "Description of the VPN.",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Status of the VPN.",
							Computed:    true,
						},
						"expired_time": {
							Type:        schema.TypeString,
							Description: "Expired time.",
							Computed:    true,
						},
						"payment_timing": {
							Type:        schema.TypeString,
							Description: "Payment timing of billing, which can be Prepaid or Postpaid. The default is Postpaid.",
							Computed:    true,
						},
						"eip": {
							Type:        schema.TypeString,
							Description: "Eip address.",
							Computed:    true,
						},
						"band_width_in_mbps": {
							Type:        schema.TypeInt,
							Description: "Eip bandwidth(Mbps), if payment_timing is Prepaid or billing_method is ByBandWidth, support between 1 and 200, if billing_method is ByTraffic, support between 1 and 1000",
							Computed:    true,
						},
						"vpn_conn_num": {
							Type:        schema.TypeInt,
							Description: "Number of VPN tunnels.",
							Computed:    true,
						},
						"vpn_conns": {
							Type:        schema.TypeList,
							Description: "ID List of the VPN gateway tunnels.",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"max_connection": {
							Type:        schema.TypeString,
							Description: "Max connection of VPN gateway.",
							Computed:    true,
						},
						"create_time": {
							Type:        schema.TypeString,
							Description: "Create time of VPN gateway.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudVpnGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	service := VpnService{client: client}
	action := "Query all vpn gateways"

	var vpcId string
	var eip string
	if value, ok := d.GetOk("vpc_id"); ok {
		vpcId = value.(string)
	}
	if value, ok := d.GetOk("eip"); ok {
		eip = value.(string)
	}
	vpns, err := service.ListVpnGateways(vpcId, eip)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_vpn_gateway", action, BCESDKGoERROR)
	}
	vpnsResult := make([]map[string]interface{}, 0)
	for _, vpn := range vpns {
		vpcMap := make(map[string]interface{})
		vpcMap["vpc_id"] = vpn.VpcId
		vpcMap["vpn_id"] = vpn.VpnId
		vpcMap["vpn_name"] = vpn.Name
		vpcMap["description"] = vpn.Description
		vpcMap["status"] = vpn.Status
		vpcMap["expired_time"] = vpn.ExpiredTime
		vpcMap["payment_timing"] = vpn.ProductType
		vpcMap["eip"] = vpn.Eip
		vpcMap["band_width_in_mbps"] = vpn.BandwidthInMbps
		vpcMap["vpn_conn_num"] = vpn.VpnConnNum
		conns := make([]string, 0)
		for _, item := range vpn.VpnConns {
			conns = append(conns, item.VpnConnId)
		}
		vpcMap["vpn_conns"] = conns
		vpnsResult = append(vpnsResult, vpcMap)
	}
	addDebug(action, vpnsResult)

	FilterDataSourceResult(d, &vpnsResult)

	d.SetId(resource.UniqueId())

	if err := d.Set("vpn_gateways", vpnsResult); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_vpn_gateways", action, BCESDKGoERROR)
	}
	return nil
}
