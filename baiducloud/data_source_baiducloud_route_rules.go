/*
Use this data source to query route rule list.

Example Usage

```hcl
data "baiducloud_route_rules" "default" {
 vpc_id = "vpc-y4p102r3mz6m"
}

output "route_rules" {
 value = "${data.baiducloud_route_rules.default.route_rules}"
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

func dataSourceBaiduCloudRouteRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudRouteRulesRead,

		Schema: map[string]*schema.Schema{
			"route_table_id": {
				Type:        schema.TypeString,
				Description: "Routing table ID for the routing rules to retrieve.",
				Optional:    true,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "VPC ID for the routing rules.",
				Optional:    true,
			},
			"route_rule_id": {
				Type:        schema.TypeString,
				Description: "ID of the routing rule to be retrieved.",
				Optional:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "Output file for saving result.",
				Optional:    true,
				ForceNew:    true,
			},

			// Attributes used for result
			"route_rules": {
				Type:        schema.TypeList,
				Description: "Result of the routing rules.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"route_table_id": {
							Type:        schema.TypeString,
							Description: "Routing table ID of the routing rule.",
							Computed:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "Description of the routing rule.",
							Computed:    true,
						},
						"next_hop_id": {
							Type:        schema.TypeString,
							Description: "Next hop ID of the routing rule.",
							Computed:    true,
						},
						"destination_address": {
							Type:        schema.TypeString,
							Description: "Destination address of the routing rule.",
							Computed:    true,
						},
						"source_address": {
							Type:        schema.TypeString,
							Description: "Source address of the routing rule.",
							Computed:    true,
						},
						"route_rule_id": {
							Type:        schema.TypeString,
							Description: "ID of the routing rule.",
							Computed:    true,
						},
						"next_hop_type": {
							Type:        schema.TypeString,
							Description: "Next hop type of the routing rule.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudRouteRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	var (
		routeTableID string
		vpcID        string
		routeRuleID  string
		outputFile   string
	)
	if v, ok := d.GetOk("route_table_id"); ok {
		routeTableID = v.(string)
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		vpcID = v.(string)
	}
	if v, ok := d.GetOk("route_rule_id"); ok {
		routeRuleID = v.(string)
	}
	if v, ok := d.GetOk("output_file"); ok {
		outputFile = v.(string)
	}

	action := "Query route rules " + routeTableID + "_" + vpcID + "_" + routeRuleID

	if routeTableID == "" && vpcID == "" {
		err := fmt.Errorf("The route talbe id and VPC id cannot be empty at the same time.")
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_route_rules", action, BCESDKGoERROR)
	}

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
		return vpcClient.GetRouteTableDetail(routeTableID, vpcID)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_route_rules", action, BCESDKGoERROR)
	}

	result, _ := raw.(*vpc.GetRouteTableResult)
	d.Set("route_table_id", result.RouteTableId)
	d.Set("vpc_id", result.VpcId)

	routeRulesResult := make([]interface{}, 0)
	for _, rule := range result.RouteRules {
		if routeRuleID != "" && routeRuleID != rule.RouteRuleId {
			continue
		}

		ruleMap := flattenRouteRule(&rule)
		routeRulesResult = append(routeRulesResult, ruleMap)
	}

	d.Set("route_rules", routeRulesResult)
	d.SetId(resource.UniqueId())

	if outputFile != "" {
		if err := writeToFile(outputFile, routeRulesResult); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_route_rules", action, BCESDKGoERROR)
		}
	}

	return nil
}

func flattenRouteRule(rule *vpc.RouteRule) map[string]interface{} {
	routeMap := make(map[string]interface{})

	routeMap["route_table_id"] = rule.RouteTableId
	routeMap["description"] = rule.Description
	routeMap["next_hop_id"] = rule.NexthopId
	routeMap["destination_address"] = rule.DestinationAddress
	routeMap["source_address"] = rule.SourceAddress
	routeMap["route_rule_id"] = rule.RouteRuleId
	routeMap["next_hop_type"] = rule.NexthopType

	return routeMap
}
