/*
Use this data source to query ACL list.

Example Usage

```hcl
data "baiducloud_acls" "default" {
 vpc_id = "vpc-y4p102r3mz6m"
}

output "acls" {
 value = "${data.baiducloud_acls.default.acls}"
}
```
*/
package baiducloud

import (
	"errors"
	"github.com/baidubce/bce-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudAcls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudAclsRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "VPC ID of the ACLs to retrieve.",
				Optional:    true,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Description: "Subnet ID of the ACLs to retrieve.",
				Optional:    true,
			},
			"acl_id": {
				Type:        schema.TypeString,
				Description: "ID of the ACL to retrieve.",
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
			"acls": {
				Type:        schema.TypeList,
				Description: "List of the ACLs.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"acl_id": {
							Type:        schema.TypeString,
							Description: "ID of the ACL.",
							Computed:    true,
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Description: "Subnet ID of the ACL.",
							Computed:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "Description of the ACL.",
							Computed:    true,
						},
						"protocol": {
							Type:        schema.TypeString,
							Description: "Protocol of the ACL.",
							Computed:    true,
						},
						"source_ip_address": {
							Type:        schema.TypeString,
							Description: "Source IP address of the ACL.",
							Computed:    true,
						},
						"destination_ip_address": {
							Type:        schema.TypeString,
							Description: "Destination IP address of the ACL.",
							Computed:    true,
						},
						"source_port": {
							Type:        schema.TypeString,
							Description: "Source port of the ACL.",
							Computed:    true,
						},
						"destination_port": {
							Type:        schema.TypeString,
							Description: "Destination port of the ACL.",
							Computed:    true,
						},
						"position": {
							Type:        schema.TypeInt,
							Description: "Position of the ACL.",
							Computed:    true,
						},
						"direction": {
							Type:        schema.TypeString,
							Description: "Direction of the ACL.",
							Computed:    true,
						},
						"action": {
							Type:        schema.TypeString,
							Description: "Action of the ACL.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudAclsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	vpcService := VpcService{client}

	var (
		vpcID      string
		subnetID   string
		aclID      string
		outputFile string
	)

	if v, ok := d.GetOk("vpc_id"); ok {
		vpcID = v.(string)
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		subnetID = v.(string)
	}
	if v, ok := d.GetOk("acl_id"); ok {
		aclID = v.(string)
	}
	if v, ok := d.GetOk("output_file"); ok {
		outputFile = v.(string)
	}

	action := "Query ACL " + vpcID + "_" + subnetID

	if vpcID == "" && subnetID == "" {
		err := errors.New("the VPC ID and Subnet ID cannot be empty at the same time")
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_acls", action, BCESDKGoERROR)
	}

	aclRules := make([]vpc.AclRule, 0)
	if vpcID != "" {
		aclEntrys, err := vpcService.ListAllAclEntrysWithVPCID(vpcID)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_acls", action, BCESDKGoERROR)
		}

		for _, aclEntry := range aclEntrys {
			if subnetID != "" && subnetID != aclEntry.SubnetId {
				continue
			}

			aclRules = append(aclRules, aclEntry.AclRules...)
		}
	} else if subnetID != "" {
		acls, err := vpcService.ListAllAclRulesWithSubnetID(subnetID)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_acls", action, BCESDKGoERROR)
		}

		aclRules = append(aclRules, acls...)
	}

	filter := NewDataSourceFilter(d)
	aclsResult := make([]interface{}, 0)
	for _, aclRule := range aclRules {
		if aclID != "" && aclID != aclRule.Id {
			continue
		}

		aclMap := flattenACL(&aclRule)
		if !filter.checkFilter(aclMap) {
			continue
		}

		aclsResult = append(aclsResult, aclMap)
	}

	d.Set("acls", aclsResult)
	d.SetId(resource.UniqueId())

	if outputFile != "" {
		if err := writeToFile(outputFile, aclsResult); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_acls", action, BCESDKGoERROR)
		}
	}

	return nil
}

func flattenACL(aclRule *vpc.AclRule) map[string]interface{} {
	aclMap := make(map[string]interface{})

	aclMap["acl_id"] = aclRule.Id
	aclMap["subnet_id"] = aclRule.SubnetId
	aclMap["description"] = aclRule.Description
	aclMap["protocol"] = aclRule.Protocol
	aclMap["source_ip_address"] = aclRule.SourceIpAddress
	aclMap["destination_ip_address"] = aclRule.DestinationIpAddress
	aclMap["source_port"] = aclRule.SourcePort
	aclMap["destination_port"] = aclRule.DestinationPort
	aclMap["position"] = aclRule.Position
	aclMap["direction"] = aclRule.Direction
	aclMap["action"] = aclRule.Action

	return aclMap
}
