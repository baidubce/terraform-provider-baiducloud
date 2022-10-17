/*
Use this data source to query blb SecurityGroups.

Example Usage

```hcl
data "baiducloud_blb_securitygroups" "default" {
   blb_id = "lb-0d29axxx6"
}

output "security_groups" {
   value = "${data.baiducloud_blb_securitygroups.default.bind_security_groups}"
}
```
*/
package baiducloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudBLBSecurityGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudBlbSecurityGroupRead,

		Schema: map[string]*schema.Schema{
			"blb_id": {
				Type:        schema.TypeString,
				Description: "id of the blb",
				ForceNew:    true,
				Required:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "blb securitygroup search result output file",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),

			"bind_security_groups": {
				Type:        schema.TypeList,
				Description: "blb bind security_groups",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_group_id": {
							Type:        schema.TypeString,
							Description: "bind security id",
							Computed:    true,
						},
						"security_group_name": {
							Type:        schema.TypeString,
							Description: "name of security group",
							Computed:    true,
						},
						"security_group_desc": {
							Type:        schema.TypeString,
							Description: "desc of security group",
							Computed:    true,
						},
						"vpc_name": {
							Type:        schema.TypeString,
							Description: "name of vpc",
							Computed:    true,
						},
						"security_group_rules": {
							Type:        schema.TypeList,
							Description: "rules of security groups",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dest_group_id": {
										Type:        schema.TypeString,
										Description: "dest group id",
										Computed:    true,
									},
									"dest_ip": {
										Type:        schema.TypeString,
										Description: "dest ip",
										Computed:    true,
									},
									"direction": {
										Type:        schema.TypeString,
										Description: "direction",
										Computed:    true,
									},
									"ethertype": {
										Type:        schema.TypeString,
										Description: "ethertype",
										Computed:    true,
									},
									"port_range": {
										Type:        schema.TypeString,
										Description: "portRange",
										Computed:    true,
									},
									"protocol": {
										Type:        schema.TypeString,
										Description: "protocol",
										Computed:    true,
									},
									"security_group_rule_id": {
										Type:        schema.TypeString,
										Description: "id of security group rule",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudBlbSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	blbService := BLBService{client}

	blbId := d.Get("blb_id").(string)

	action := "Query blb securitygroup list" + blbId

	blbSecurityGroup, err := blbService.GetBlbSecurityGroup(blbId)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb_securitygroups", action, BCESDKGoERROR)
	}
	addDebug(action, blbSecurityGroup)

	groups := blbSecurityGroup.BlbSecurityGroups
	bindSecurityGroups := make([]map[string]interface{}, 0, len(groups))
	for _, group := range groups {
		securityMap := make(map[string]interface{})
		securityMap["security_group_id"] = group.SecurityGroupId
		securityMap["security_group_name"] = group.SecurityGroupName
		securityMap["security_group_desc"] = group.SecurityGroupDesc
		securityMap["vpc_name"] = group.VpcName

		models := group.SecurityGroupRules
		securityGroupRules := make([]interface{}, 0, len(models))
		for _, model := range models {
			ruleMap := make(map[string]interface{})
			ruleMap["dest_group_id"] = model.DestGroupId
			ruleMap["dest_ip"] = model.DestIp
			ruleMap["direction"] = model.Direction
			ruleMap["ethertype"] = model.Ethertype
			ruleMap["port_range"] = model.PortRange
			ruleMap["protocol"] = model.Protocol
			ruleMap["security_group_rule_id"] = model.SecurityGroupRuleId
			securityGroupRules = append(securityGroupRules, ruleMap)
		}
		securityMap["security_group_rules"] = securityGroupRules

		bindSecurityGroups = append(bindSecurityGroups, securityMap)
	}

	FilterDataSourceResult(d, &bindSecurityGroups)

	if err := d.Set("bind_security_groups", bindSecurityGroups); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb_securitygroups", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), bindSecurityGroups); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb_securitygroups", action, BCESDKGoERROR)
		}
	}

	return nil
}
