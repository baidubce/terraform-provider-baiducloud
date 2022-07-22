/*
Provide a resource to create a security group rule.

Example Usage

```hcl
resource "baiducloud_security_group" "default" {
  name = "my-sg"
  description = "default"
}

resource "baiducloud_security_group_rule" "default" {
  security_group_id = "${baiducloud_security_group.default.id}"
  remark            = "remark"
  protocol          = "udp"
  port_range        = "1-65523"
  direction         = "ingress"
}
```
*/
package baiducloud

import (
	"fmt"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudSecurityGroupRuleCreate,
		Read:   resourceBaiduCloudSecurityGroupRuleRead,
		Delete: resourceBaiduCloudSecurityGroupRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:        schema.TypeString,
				Description: "SecurityGroup rule's security group id",
				Required:    true,
				ForceNew:    true,
			},
			"remark": {
				Type:        schema.TypeString,
				Description: "SecurityGroup rule's remark",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"direction": {
				Type:         schema.TypeString,
				Description:  "SecurityGroup rule's direction, support ingress/egress",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ingress", "egress"}, false),
			},
			"ether_type": {
				Type:         schema.TypeString,
				Description:  "SecurityGroup rule's ether type, support IPv4/IPv6",
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"IPv4", "IPv6"}, false),
			},
			"port_range": {
				Type:        schema.TypeString,
				Description: "SecurityGroup rule's port range, you can set single port like 80, or set a port range, like 1-65535, default 1-65535. If protocol is all, only support 1-65535",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"protocol": {
				Type:         schema.TypeString,
				Description:  "SecurityGroup rule's protocol, support tcp/udp/icmp/all, default all",
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"tcp", "udp", "icmp", "all"}, false),
			},
			"source_group_id": {
				Type:          schema.TypeString,
				Description:   "SecurityGroup rule's source group id, source_group_id and source_ip can not set in the same time",
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source_ip"},
			},
			"source_ip": {
				Type:          schema.TypeString,
				Description:   "SecurityGroup rule's source ip, source_group_id and source_ip can not set in the same time",
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source_group_id"},
			},
			"dest_group_id": {
				Type:          schema.TypeString,
				Description:   "SecurityGroup rule's destination group id, dest_group_id and dest_ip can not set in the same time",
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"dest_ip"},
			},
			"dest_ip": {
				Type:          schema.TypeString,
				Description:   "SecurityGroup rule's destination ip, dest_group_id and dest_ip can not set in the same time",
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"dest_group_id"},
			},
		},
	}
}

func resourceBaiduCloudSecurityGroupRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bccService := BccService{client}

	singleRule, err := buildSingleSecurityGroupRuleModel(d)
	if err != nil {
		return WrapError(err)
	}

	ruleId, err := bccService.buildSecurityGroupRuleId(singleRule)
	if err != nil {
		return WrapError(err)
	}

	args := &api.AuthorizeSecurityGroupArgs{
		Rule:        singleRule,
		ClientToken: buildClientToken(),
	}
	action := "Authorize SecurityGroup Rules " + singleRule.SecurityGroupId

	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err := client.WithBccClient(func(bccClient *bcc.Client) (interface{}, error) {
			return nil, bccClient.AuthorizeSecurityGroupRule(singleRule.SecurityGroupId, args)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, args)

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_security_group_rule", action, BCESDKGoERROR)
	}

	d.SetId(ruleId)
	return resourceBaiduCloudSecurityGroupRuleRead(d, meta)
}

func resourceBaiduCloudSecurityGroupRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bccService := BccService{client}

	id := d.Id()
	action := "Query Security Group Rule " + id

	sgRule, err := bccService.GetSecurityGroupRule(id)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}

		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_security_group_rule", action, BCESDKGoERROR)
	}

	d.Set("security_group_id", sgRule.SecurityGroupId)
	d.Set("direction", sgRule.Direction)
	d.Set("remark", sgRule.Remark)
	d.Set("ether_type", sgRule.Ethertype)
	d.Set("port_range", sgRule.PortRange)
	d.Set("protocol", sgRule.Protocol)
	d.Set("source_group_id", sgRule.SourceGroupId)
	d.Set("source_ip", sgRule.SourceIp)
	d.Set("dest_group_id", sgRule.DestGroupId)
	d.Set("dest_ip", sgRule.DestIp)

	return nil
}

func resourceBaiduCloudSecurityGroupRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	action := "Delete Securitt Group Rule " + d.Id()
	singleRule, err := buildSingleSecurityGroupRuleModel(d)
	if err != nil {
		return WrapError(err)
	}

	revokeArgs := &api.RevokeSecurityGroupArgs{
		Rule: singleRule,
	}
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := client.WithBccClient(func(client *bcc.Client) (i interface{}, e error) {
			return nil, client.RevokeSecurityGroupRule(singleRule.SecurityGroupId, revokeArgs)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, revokeArgs)

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_security_group_rule", action, BCESDKGoERROR)
	}

	return nil
}

func buildSingleSecurityGroupRuleModel(d *schema.ResourceData) (*api.SecurityGroupRuleModel, error) {
	singleRule := &api.SecurityGroupRuleModel{}
	if v, ok := d.GetOk("remark"); ok && v.(string) != "" {
		singleRule.Remark = v.(string)
	}
	if v, ok := d.GetOk("direction"); ok && v.(string) != "" {
		singleRule.Direction = v.(string)
	}
	if v, ok := d.GetOk("ether_type"); ok && v.(string) != "" {
		singleRule.Ethertype = v.(string)
	}
	if v, ok := d.GetOk("port_range"); ok && v.(string) != "" {
		singleRule.PortRange = v.(string)
	}
	if v, ok := d.GetOk("protocol"); ok && v.(string) != "" {
		singleRule.Protocol = v.(string)
	}
	if v, ok := d.GetOk("source_group_id"); ok && v.(string) != "" {
		singleRule.SourceGroupId = v.(string)
	}
	if v, ok := d.GetOk("source_ip"); ok && v.(string) != "" {
		singleRule.SourceIp = v.(string)
	}
	if v, ok := d.GetOk("dest_group_id"); ok && v.(string) != "" {
		singleRule.DestGroupId = v.(string)
	}
	if v, ok := d.GetOk("dest_ip"); ok && v.(string) != "" {
		singleRule.DestIp = v.(string)
	}
	if v, ok := d.GetOk("security_group_id"); ok && v.(string) != "" {
		singleRule.SecurityGroupId = v.(string)
	}

	if singleRule.Protocol == "all" && !stringInSlice([]string{"", "1-65535"}, singleRule.PortRange) {
		return nil, fmt.Errorf("if protocol is all, port_range only support [\"\", \"1-65535\"], but now is %s",
			singleRule.PortRange)
	}

	return singleRule, nil
}
