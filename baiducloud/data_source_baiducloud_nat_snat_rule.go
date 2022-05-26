/*
Use this data source to query NAT Gateway SNAT rule list.

Example Usage

```hcl
data "baiducloud_nat_snat_rules" "default" {
 nat_id = "nat-brkztytqzbh0"
}

output "nat_snat_rules" {
 value = "${data.baiducloud_nat_snat_rules.default.nat_snat_rules}"
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

func dataSourceBaiduCloudNatSnatRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudNatSnatRulesRead,

		Schema: map[string]*schema.Schema{
			"nat_id": {
				Type:        schema.TypeString,
				Description: "ID of the NAT gateway to retrieve.",
				Required:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "Output file for saving result.",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),

			// Attributes used for result
			"nat_snat_rules": {
				Type:        schema.TypeList,
				Description: "The list of NAT SNAT rules.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:        schema.TypeString,
							Description: "ID of the NAT SNAT rule.",
							Computed:    true,
						},
						"rule_name": {
							Type:        schema.TypeString,
							Description: "Name of the NAT SNAT rule.",
							Computed:    true,
						},
						"public_ips_address": {
							Type:        schema.TypeList,
							Description: "Public network IPs, EIPs associated on the NAT gateway SNAT or IPs in the shared bandwidth.",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"source_cidr": {
							Type:        schema.TypeString,
							Description: "Intranet IP/segment.",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Status of the NAT SNAT rule.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudNatSnatRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	vpcService := VpcService{client}

	var (
		natID      string
		name       string
		outputFile string
	)
	if v, ok := d.GetOk("nat_id"); ok {
		natID = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}
	if v, ok := d.GetOk("output_file"); ok {
		outputFile = v.(string)
	}

	action := "Query NAT SNAT Rules " + natID + "_" + name

	if natID == "" {
		err := fmt.Errorf("The NAT ID cannot be empty.")
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_nat_snat_rules", action, BCESDKGoERROR)
	}

	snatRulesResult := make([]map[string]interface{}, 0)
	snatRules, err := vpcService.ListAllNatSnatRulesWithNatID(natID)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_nat_snat_rules", action, BCESDKGoERROR)
	}

	for _, snatRule := range snatRules {
		natMap := flattenNATSnatRule(&snatRule)
		snatRulesResult = append(snatRulesResult, natMap)
	}

	FilterDataSourceResult(d, &snatRulesResult)
	d.Set("nat_snat_rules", snatRulesResult)

	d.SetId(resource.UniqueId())

	if outputFile != "" {
		if err := writeToFile(outputFile, snatRulesResult); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_nat_snat_rules", action, BCESDKGoERROR)
		}
	}

	return nil
}

func flattenNATSnatRule(snatRule *vpc.SnatRule) map[string]interface{} {
	natMap := make(map[string]interface{})

	natMap["rule_id"] = snatRule.RuleId
	natMap["rule_name"] = snatRule.RuleName
	natMap["public_ips_address"] = snatRule.PublicIpAddresses
	natMap["source_cidr"] = snatRule.SourceCIDR
	natMap["status"] = snatRule.Status

	return natMap
}
