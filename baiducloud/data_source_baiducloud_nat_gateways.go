/*
Use this data source to query NAT gateway list.

Example Usage

```hcl
data "baiducloud_nat_gateways" "default" {
 vpc_id = "vpc-y4p102r3mz6m"
}

output "nat_gateways" {
 value = "${data.baiducloud_nat_gateways.default.nat_gateways}"
}
```
*/
package baiducloud

import (
	"fmt"

	"github.com/baidubce/bce-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudNatGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudNatGatewaysRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "VPC ID where the NAT gateways located.",
				Optional:    true,
			},
			"nat_id": {
				Type:        schema.TypeString,
				Description: "ID of the NAT gateway to retrieve.",
				Optional:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the NAT gateway.",
				Optional:    true,
			},
			"ip": {
				Type:        schema.TypeString,
				Description: "Specify the EIP binded by the NAT gateway to retrieve.",
				Optional:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "Output file for saving result.",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),

			// Attributes used for result
			"nat_gateways": {
				Type:        schema.TypeList,
				Description: "The list of NAT gateways.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "ID of the NAT gateway.",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the NAT gateway.",
							Computed:    true,
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Description: "VPC ID of the NAT gateway.",
							Computed:    true,
						},
						"spec": {
							Type:        schema.TypeString,
							Description: "Spec of the NAT gateway.",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Status of the NAT gateway.",
							Computed:    true,
						},
						"eips": {
							Type:        schema.TypeList,
							Description: "EIP list of the NAT gateway.",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"payment_timing": {
							Type:        schema.TypeString,
							Description: "Payment timing of the NAT gateway.",
							Computed:    true,
						},
						"expired_time": {
							Type:        schema.TypeString,
							Description: "Expired time of the NAT gateway.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudNatGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	vpcService := VpcService{client}

	var (
		vpcID      string
		natID      string
		name       string
		ip         string
		outputFile string
	)
	if v, ok := d.GetOk("vpc_id"); ok {
		vpcID = v.(string)
	}
	if v, ok := d.GetOk("nat_id"); ok {
		natID = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}
	if v, ok := d.GetOk("ip"); ok {
		ip = v.(string)
	}
	if v, ok := d.GetOk("output_file"); ok {
		outputFile = v.(string)
	}

	action := "Query NAT Gateways " + vpcID + "_" + natID + "_" + name

	if vpcID == "" && natID == "" {
		err := fmt.Errorf("The VPC ID and NAT ID cannot be empty at the same time.")
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_nat_gateways", action, BCESDKGoERROR)
	}

	natsResult := make([]map[string]interface{}, 0)
	if vpcID != "" {
		args := &vpc.ListNatGatewayArgs{
			VpcId: vpcID,
			NatId: natID,
			Name:  name,
			Ip:    ip,
		}
		nats, err := vpcService.ListAllNatGateways(args)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_nat_gateways", action, BCESDKGoERROR)
		}
		for _, nat := range nats {
			natMap := flattenNAT(&nat)
			natsResult = append(natsResult, natMap)
		}
	} else if natID != "" {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return vpcClient.GetNatGatewayDetail(natID)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_nat_gateways", action, BCESDKGoERROR)
		}

		result, _ := raw.(*vpc.NAT)
		natMap := flattenNAT(result)
		natsResult = append(natsResult, natMap)
	}

	FilterDataSourceResult(d, &natsResult)
	d.Set("nat_gateways", natsResult)

	d.SetId(resource.UniqueId())

	if outputFile != "" {
		if err := writeToFile(outputFile, natsResult); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_nat_gateways", action, BCESDKGoERROR)
		}
	}

	return nil
}

func flattenNAT(nat *vpc.NAT) map[string]interface{} {
	natMap := make(map[string]interface{})

	natMap["id"] = nat.Id
	natMap["name"] = nat.Name
	natMap["vpc_id"] = nat.VpcId
	natMap["spec"] = nat.Spec
	natMap["status"] = nat.Status
	natMap["eips"] = nat.Eips
	natMap["payment_timing"] = nat.PaymentTiming
	natMap["expired_time"] = nat.ExpiredTime

	return natMap
}
