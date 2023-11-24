/*
Provides a resource to create a VPC routing rule.

Example Usage

```hcl
resource "baiducloud_route_rule" "default" {
  route_table_id = "rt-as4npcsp2hve"
  source_address = "192.168.0.0/24"
  destination_address = "192.168.1.0/24"
  next_hop_id = "i-BtXnDM6y"
  next_hop_type = "custom"
  description = "created by terraform"
}
```
*/
package baiducloud

import (
	"strings"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

// There are only create and delete api in bce-go-sdk.
// When the config file of routing rule is updated, we have to destroy the old rule and create a new one.
// In order to read the route rule data, we can use the api of routing table instead.
func resourceBaiduCloudRouteRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudRouteRuleCreate,
		Read:   resourceBaiduCloudRouteRuleRead,
		Delete: resourceBaiduCloudRouteRuleDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"route_table_id": {
				Type:        schema.TypeString,
				Description: "ID of the routing table.",
				Required:    true,
				ForceNew:    true,
			},
			"source_address": {
				Type:        schema.TypeString,
				Description: "Source CIDR block of the routing rule. The value can be all network segments 0.0.0.0/0, existing subnet segments in the VPC, or the network segment within the subnet.",
				Required:    true,
				ForceNew:    true,
			},
			"destination_address": {
				Type:        schema.TypeString,
				Description: "Destination CIDR block of the routing rule. The network segment can be 0.0.0.0/0, otherwise, the destination address cannot overlap with this VPC CIDR block(except when the destination network segment or the VPC CIDR is 0.0.0.0/0).",
				Required:    true,
				ForceNew:    true,
			},
			"next_hop_id": {
				Type:        schema.TypeString,
				Description: "Next-hop ID, this field must be filled when creating a single path route.",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			/**
			* todo 目前OpenAPI的多线路由只支持创建，查询和更新都不支持，所以暂时不开放，待API完善后开放
			 */
			//"next_hop_list": {
			//	Type:        schema.TypeList,
			//	Description: "Create a multi-path route based on the next hop information. This field is required when creating a multi-path route.",
			//	Optional:    true,
			//	Computed:    true,
			//	Elem: &schema.Resource{
			//		Schema: map[string]*schema.Schema{
			//			"next_hop_id": {
			//				Type:        schema.TypeString,
			//				Description: "Next-hop ID.",
			//				Required:    true,
			//				ForceNew:    true,
			//			},
			//			"next_hop_type": {
			//				Type:         schema.TypeString,
			//				Description:  "Routing type. Currently only the dedicated gateway type dcGateway is supported.",
			//				Required:     true,
			//				ForceNew:     true,
			//				ValidateFunc: validation.StringInSlice([]string{"dcGateway"}, false),
			//			},
			//			"path_type": {
			//				Type:         schema.TypeString,
			//				Description:  "Multi-line mode. The load balancing value is ecmp; the main backup mode value is ha:active, ha:standby, which represent the main and backup routes respectively.",
			//				Required:     true,
			//				ForceNew:     true,
			//				ValidateFunc: validation.StringInSlice([]string{"ecmp", "ha:active", "ha:standby"}, false),
			//			},
			//		},
			//	},
			//},
			"next_hop_type": {
				Type:         schema.TypeString,
				Description:  "Type of the next hop, available values are custom, vpn, nat and dcGateway.",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"custom", "vpn", "nat", "dcGateway", "peerConn", "ipv6gateway"}, false),
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the routing rule.",
				Optional:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceBaiduCloudRouteRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	createRouteRuleArgs := buildBaiduCloudRouteRuleArgs(d, meta)
	action := "Create Route Rule for Route Table " + createRouteRuleArgs.RouteTableId

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return vpcClient.CreateRouteRule(createRouteRuleArgs)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)
		result, _ := raw.(*vpc.CreateRouteRuleResult)
		d.SetId(result.RouteRuleId)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_route_rule", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudRouteRuleRead(d, meta)
}

func resourceBaiduCloudRouteRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	idParts := strings.Split(d.Id(), ":")
	routeRuleId := idParts[0]
	action := "Query Route Rule " + routeRuleId

	routeTableID := ""
	if len(idParts) > 1 {
		routeTableID = idParts[1]
		d.SetId(idParts[0])
	} else if v, ok := d.GetOk("route_table_id"); ok {
		routeTableID = v.(string)
	}
	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
		return vpcClient.GetRouteTableDetail(routeTableID, "")
	})
	addDebug(action, raw)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_route_rule", action, BCESDKGoERROR)
	}

	result, _ := raw.(*vpc.GetRouteTableResult)
	d.SetId(routeRuleId)
	for _, rule := range result.RouteRules {
		if rule.RouteRuleId == routeRuleId {
			d.Set("route_table_id", rule.RouteTableId)
			d.Set("source_address", rule.SourceAddress)
			d.Set("destination_address", rule.DestinationAddress)
			d.Set("next_hop_id", rule.NexthopId)
			d.Set("next_hop_type", rule.NexthopType)
			d.Set("description", rule.Description)
			return nil
		}
	}
	return nil
}

func resourceBaiduCloudRouteRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	routeRuleId := d.Id()
	action := "Delete Route Rule " + routeRuleId

	clientToken := buildClientToken()
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return nil, vpcClient.DeleteRouteRule(routeRuleId, clientToken)
		})
		addDebug(action, nil)
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_route_rule", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudRouteRuleArgs(d *schema.ResourceData, meta interface{}) *vpc.CreateRouteRuleArgs {
	request := &vpc.CreateRouteRuleArgs{
		ClientToken: buildClientToken(),
	}

	if v := d.Get("route_table_id").(string); v != "" {
		request.RouteTableId = v
	}
	if v := d.Get("source_address").(string); v != "" {
		request.SourceAddress = v
	}
	if v := d.Get("destination_address").(string); v != "" {
		request.DestinationAddress = v
	}
	if v := d.Get("next_hop_id").(string); v != "" {
		request.NexthopId = v
	}
	if v := d.Get("next_hop_type").(string); v != "" {
		request.NexthopType = vpc.NexthopType(v)
		if len(d.Get("next_hop_list").([]interface{})) > 0 {
			nextHopList := d.Get("next_hop_list").([]interface{})
			result := make([]vpc.NextHop, len(nextHopList))
			for _, item := range nextHopList {
				itemMap := item.(map[string]interface{})
				temp := vpc.NextHop{
					NexthopId:   itemMap["next_hop_id"].(string),
					NexthopType: vpc.NexthopType(itemMap["next_hop_type"].(string)),
					PathType:    itemMap["path_type"].(string),
				}
				result = append(result, temp)
			}
			request.NextHopList = result
		}
	}
	if v := d.Get("description").(string); v != "" {
		request.Description = v
	}

	return request
}
