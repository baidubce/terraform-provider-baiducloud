/*
Use this resource to create a Blb SecurityGroup.

~> **NOTE:** The terminate operation of SecurityGroup does NOT take effect immediately，maybe takes for several minites.

Example Usage

```hcl
resource "baiducloud_blb_securitygroup" "my-server" {
 blb_id = "xxxx"
 security_group_ids = ["xxxxxx"]
}
```

Import

Blb SecurityGroup can be imported, e.g.

```hcl
$ terraform import baiducloud_blb_securitygroup.my-server id
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/blb"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudBlbSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudBlbSecurityGroupCreate,
		Read:   resourceBaiduCloudBlbSecurityGroupRead,
		Delete: resourceBaiduCloudBlbSecurityGroupDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"blb_id": {
				Type:        schema.TypeString,
				Description: "id of the blb",
				ForceNew:    true,
				Required:    true,
			},
			"security_group_ids": {
				Type:        schema.TypeSet,
				Description: "ids of the security.",
				ForceNew:    true,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"bind_security_groups": {
				Type:        schema.TypeList,
				Description: "blb bind security_groups",
				Computed:    true,
				MinItems:    1,
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
							MinItems:    1,
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

func resourceBaiduCloudBlbSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	blbId := d.Get("blb_id").(string)
	args := buildBaiduCloudBlbSecurityGroupArgs(d)

	action := "bind blb security group ids " + blbId

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithBLBClient(func(blbClient *blb.Client) (i interface{}, e error) {
			return nil, blbClient.BindSecurityGroups(blbId, args)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)
		d.SetId(blbId)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb_securitygroup", action, BCESDKGoERROR)
	}
	return resourceBaiduCloudBlbSecurityGroupRead(d, meta)

}

func resourceBaiduCloudBlbSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	blbService := BLBService{client}

	blbId := d.Id()

	action := "Query blb security groups " + blbId

	blbSecurityGroup, err := blbService.GetBlbSecurityGroup(blbId)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb_securitygroup", action, BCESDKGoERROR)
	}
	addDebug(action, blbSecurityGroup)

	groups := blbSecurityGroup.BlbSecurityGroups
	bindSecurityGroups := make([]interface{}, 0, len(groups))
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
	d.Set("bind_security_groups", bindSecurityGroups)

	return nil
}

func resourceBaiduCloudBlbSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	blbId := d.Id()

	action := "Unbind blb security groups " + blbId
	args := buildBaiduCloudBlbUnbindSecurityGroupArgs(d)

	_, err := client.WithBLBClient(func(blbClient *blb.Client) (i interface{}, e error) {
		return nil, blbClient.UnbindSecurityGroups(blbId, args)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb_securitygroup", action, BCESDKGoERROR)
	}
	addDebug(action, nil)

	return nil
}

func buildBaiduCloudBlbSecurityGroupArgs(d *schema.ResourceData) *blb.UpdateSecurityGroupsArgs {
	args := &blb.UpdateSecurityGroupsArgs{
		ClientToken: buildClientToken(),
	}

	securityGroupIds := make([]string, 0)
	ids, ok := d.GetOk("security_group_ids")
	if ok {
		for _, id := range ids.(*schema.Set).List() {
			securityGroupIds = append(securityGroupIds, id.(string))
		}
		args.SecurityGroupIds = securityGroupIds
	}

	return args
}

func buildBaiduCloudBlbUnbindSecurityGroupArgs(d *schema.ResourceData) *blb.UpdateSecurityGroupsArgs {
	args := &blb.UpdateSecurityGroupsArgs{
		ClientToken: buildClientToken(),
	}

	securityGroupIds := make([]string, 0)
	ids, ok := d.GetOk("security_group_ids")
	if ok {
		for _, id := range ids.(*schema.Set).List() {
			securityGroupIds = append(securityGroupIds, id.(string))
		}
		args.SecurityGroupIds = securityGroupIds
	}

	return args
}
