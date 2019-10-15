/*
Use this data source to query Security Group list.

Example Usage

```hcl
data "baiducloud_security_group_rules" "default" {}

output "security_group_rules" {
 value = "${data.baiducloud_security_group_rules.default.rules}"
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudSecurityGroupRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudSecurityGroupRulesRead,

		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:        schema.TypeString,
				Description: "Security Group ID",
				Required:    true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Description: "Security Group attached instance ID",
				Optional:    true,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "Security Group attached vpc id",
				Optional:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "Security Group search result output file",
				Optional:    true,
				ForceNew:    true,
			},
			"rules": {
				Type:        schema.TypeList,
				Description: "Security Group rules",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"remark": {
							Type:        schema.TypeString,
							Description: "SecurityGroup rule's remark",
							Computed:    true,
						},
						"direction": {
							Type:        schema.TypeString,
							Description: "SecurityGroup rule's direction",
							Computed:    true,
						},
						"ether_type": {
							Type:        schema.TypeString,
							Description: "SecurityGroup rule's ether type",
							Computed:    true,
						},
						"port_range": {
							Type:        schema.TypeString,
							Description: "SecurityGroup rule's port range",
							Computed:    true,
						},
						"protocol": {
							Type:        schema.TypeString,
							Description: "SecurityGroup rule's protocol",
							Computed:    true,
						},
						"source_group_id": {
							Type:        schema.TypeString,
							Description: "SecurityGroup rule's source group id",
							Computed:    true,
						},
						"source_ip": {
							Type:        schema.TypeString,
							Description: "SecurityGroup rule's source ip",
							Computed:    true,
						},
						"dest_group_id": {
							Type:        schema.TypeString,
							Description: "SecurityGroup rule's destination group id",
							Computed:    true,
						},
						"dest_ip": {
							Type:        schema.TypeString,
							Description: "SecurityGroup rule's destination ip",
							Computed:    true,
						},
						"security_group_id": {
							Type:        schema.TypeString,
							Description: "SecurityGroup rule's security group id",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudSecurityGroupRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bccService := BccService{client}

	listSgArgs := &api.ListSecurityGroupArgs{}
	if v, ok := d.GetOk("instance_id"); ok && v.(string) != "" {
		listSgArgs.InstanceId = v.(string)
	}
	if v, ok := d.GetOk("vpc_id"); ok && v.(string) != "" {
		listSgArgs.VpcId = v.(string)
	}

	sgId := d.Get("security_group_id").(string)
	action := "Data Source Query All Security Group " + sgId + " rules"
	sgList, err := bccService.ListAllSecurityGroups(listSgArgs)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_security_group_rules", action, BCESDKGoERROR)
	}

	var ruleMap []map[string]interface{}
	for _, sg := range sgList {
		if sg.Id == sgId {
			ruleMap = bccService.FlattenSecurityGroupRuleModelsToMap(sg.Rules)
			if err := d.Set("rules", ruleMap); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "baiducloud_security_group_rules", action, BCESDKGoERROR)
			}
			d.SetId(resource.UniqueId())

			if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
				if err := writeToFile(v.(string), ruleMap); err != nil {
					return WrapErrorf(err, DefaultErrorMsg, "baiducloud_security_group_rules", action, BCESDKGoERROR)
				}
			}

			return nil
		}
	}

	// no such security group
	return WrapErrorf(Error("No such Security Group"), DefaultErrorMsg, "baiducloud_security_group_rules", action, BCESDKGoERROR)
}
