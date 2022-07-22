/*
Provide a resource to create an ACL Rule.

Example Usage

```hcl
resource "baiducloud_acl" "default" {
  subnet_id = "sbn-86c3v6pnt8b4"
  protocol = "tcp"
  source_ip_address = "192.168.0.0/24"
  destination_ip_address = "192.168.1.0/24"
  source_port = "8888"
  destination_port = "9999"
  position = 20
  direction = "ingress"
  action = "allow"
}
```
*/
package baiducloud

import (
	"fmt"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudAclCreate,
		Read:   resourceBaiduCloudAclRead,
		Update: resourceBaiduCloudAclUpdate,
		Delete: resourceBaiduCloudAclDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"subnet_id": {
				Type:        schema.TypeString,
				Description: "Subnet ID of the acl.",
				Required:    true,
				ForceNew:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the acl.",
				Optional:    true,
			},
			"protocol": {
				Type:         schema.TypeString,
				Description:  "Protocol of the acl, available values are all, tcp, udp and icmp.",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"all", "tcp", "udp", "icmp"}, false),
			},
			"source_ip_address": {
				Type:        schema.TypeString,
				Description: "Source ip address of the acl.",
				Required:    true,
			},
			"destination_ip_address": {
				Type:        schema.TypeString,
				Description: "Destination ip address of the acl.",
				Required:    true,
			},
			"source_port": {
				Type:        schema.TypeString,
				Description: "Source port of the acl.",
				Required:    true,
			},
			"destination_port": {
				Type:        schema.TypeString,
				Description: "Destination port of the acl.",
				Required:    true,
			},
			"position": {
				Type:         schema.TypeInt,
				Description:  "Position of the acl, representing the priority of the acl rule. The value should be 1-5000 and cannot be duplicated with existing entries. The smaller the value, the higher the priority, and the rule matching order is to match the priority from high to low.",
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 5000),
			},
			"direction": {
				Type:         schema.TypeString,
				Description:  "Direction of the acl. Valid values are ingress and egress, respectively indicating the inbound of the rule and the outbound rule.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ingress", "egress"}, false),
			},
			"action": {
				Type:         schema.TypeString,
				Description:  "Action of the acl. Valid values are allow and deny.",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"allow", "deny"}, false),
			},
		},
	}
}

func resourceBaiduCloudAclCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	args := buildBaiduCloudAclCreateArgs(d)
	action := "Create ACL Rule with Subnet ID: " + args.AclRules[0].SubnetId

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return nil, vpcClient.CreateAclRule(args)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, nil)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_subnet", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudAclRead(d, meta)
}

func resourceBaiduCloudAclRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	aclService := &VpcService{client}

	var (
		subnetId string
		position int
	)
	if v := d.Get("subnet_id").(string); v != "" {
		subnetId = v
	}
	if v := d.Get("position").(int); v != 0 {
		position = v
	}
	action := "Query ACL Rule for Subnet ID: " + subnetId

	aclRules, err := aclService.DescribeAclRules(subnetId)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_acl", action, BCESDKGoERROR)
	}
	if len(aclRules) == 0 {
		return WrapErrorf(fmt.Errorf("there is no ACL Rule for Subnet %s", subnetId), DefaultErrorMsg, "baiducloud_acl", action, BCESDKGoERROR)
	}
	for _, aclRule := range aclRules {
		if aclRule.SubnetId == subnetId && aclRule.Position == position {
			d.SetId(aclRule.Id)
			d.Set("subnet_id", aclRule.SubnetId)
			d.Set("description", aclRule.Description)
			d.Set("protocol", aclRule.Protocol)
			d.Set("source_ip_address", aclRule.SourceIpAddress)
			d.Set("destination_ip_address", aclRule.DestinationIpAddress)
			d.Set("source_port", aclRule.SourcePort)
			d.Set("destination_port", aclRule.DestinationPort)
			d.Set("position", aclRule.Position)
			d.Set("direction", aclRule.Direction)
			d.Set("action", aclRule.Action)

			return nil
		}
	}
	d.SetId("")
	return nil
}

func resourceBaiduCloudAclUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	aclRuleId := d.Id()
	action := "Update ACL Rule " + aclRuleId
	args := buildBaiduCloudAclUpdateArgs(d)

	_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
		return nil, vpcClient.UpdateAclRule(aclRuleId, args)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_acl", action, BCESDKGoERROR)
	}
	addDebug(action, nil)

	return resourceBaiduCloudAclRead(d, meta)
}

func resourceBaiduCloudAclDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	aclRuleId := d.Id()
	action := "Delete ACL Rule " + aclRuleId

	_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
		return nil, vpcClient.DeleteAclRule(aclRuleId, buildClientToken())
	})
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_acl", action, BCESDKGoERROR)
	}
	addDebug(action, nil)
	return nil
}

func buildBaiduCloudAclCreateArgs(d *schema.ResourceData) *vpc.CreateAclRuleArgs {
	request := vpc.AclRuleRequest{}

	if v := d.Get("subnet_id").(string); v != "" {
		request.SubnetId = v
	}
	if v := d.Get("description").(string); v != "" {
		request.Description = v
	}
	if v := d.Get("protocol").(string); v != "" {
		request.Protocol = vpc.AclRuleProtocolType(v)
	}
	if v := d.Get("source_ip_address").(string); v != "" {
		request.SourceIpAddress = v
	}
	if v := d.Get("destination_ip_address").(string); v != "" {
		request.DestinationIpAddress = v
	}
	if v := d.Get("source_port").(string); v != "" {
		request.SourcePort = v
	}
	if v := d.Get("destination_port").(string); v != "" {
		request.DestinationPort = v
	}
	if v := d.Get("position").(int); v != 0 {
		request.Position = v
	}
	if v := d.Get("direction").(string); v != "" {
		request.Direction = vpc.AclRuleDirectionType(v)
	}
	if v := d.Get("action").(string); v != "" {
		request.Action = vpc.AclRuleActionType(v)
	}

	return &vpc.CreateAclRuleArgs{
		ClientToken: buildClientToken(),
		AclRules:    []vpc.AclRuleRequest{request},
	}
}

func buildBaiduCloudAclUpdateArgs(d *schema.ResourceData) *vpc.UpdateAclRuleArgs {
	request := &vpc.UpdateAclRuleArgs{
		ClientToken: buildClientToken(),
	}

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
	}
	if d.HasChange("protocol") {
		request.Protocol = vpc.AclRuleProtocolType(d.Get("protocol").(string))
	}
	if d.HasChange("source_ip_address") {
		request.SourceIpAddress = d.Get("source_ip_address").(string)
	}
	if d.HasChange("destination_ip_address") {
		request.DestinationIpAddress = d.Get("destination_ip_address").(string)
	}
	if d.HasChange("source_port") {
		request.SourcePort = d.Get("source_port").(string)
	}
	if d.HasChange("destination_port") {
		request.DestinationPort = d.Get("destination_port").(string)
	}
	if d.HasChange("position") {
		request.Position = d.Get("position").(int)
	}
	if d.HasChange("action") {
		request.Action = vpc.AclRuleActionType(d.Get("action").(string))
	}

	return request
}
