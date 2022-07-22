/*
Provide a resource to create a NAT Gateway SNAT Rule.

Example Usage

```hcl
resource "baiducloud_nat_snat_rule" "default" {
  nat_id = "nat-brkztytqzbh0"
  rule_name = "test"
  public_ips_address = ["100.88.14.90"]
  source_cidr = "192.168.1.0/24"
}
```
*/
package baiducloud

import (
	"log"
	"strings"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudNatSnatRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudNatSnatRuleCreate,
		Read:   resourceBaiduCloudNatSnatRuleRead,
		Update: resourceBaiduCloudNatSnatRuleUpdate,
		Delete: resourceBaiduCloudNatSnatRuleDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"nat_id": {
				Type:        schema.TypeString,
				Description: "ID of NAT Gateway.",
				Required:    true,
			},
			"rule_name": {
				Type:        schema.TypeString,
				Description: "Rule name, consisting of uppercase and lowercase letters、 numbers and special characters, such as \"-\"_\"/\".\". The value must start with a letter, and the length should between 1-65.",
				Required:    true,
			},
			"public_ips_address": {
				Type:        schema.TypeSet,
				Description: "Public network IPs, EIPs associated on the NAT gateway SNAT or IPs in the shared bandwidth.",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"source_cidr": {
				Type:        schema.TypeString,
				Description: "Intranet IP/segment.",
				Required:    true,
			},
		},
	}
}

func resourceBaiduCloudNatSnatRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	natId := d.Get("nat_id").(string)

	args := buildBaiduCloudNatSnatRuleArgs(d)
	action := "Create NAT " + natId + " SNAT Rule " + args.RuleName

	if err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return vpcClient.CreateNatGatewaySnatRule(natId, args)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)
		result, _ := raw.(*vpc.CreateNatGatewaySnatRuleResult)
		d.SetId(getNatSnatRuleResourceId(result.RuleId, natId))
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_nat_snat_rule", action, BCESDKGoERROR)
	}
	return resourceBaiduCloudNatSnatRuleRead(d, meta)

}

func resourceBaiduCloudNatSnatRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	vpcService := VpcService{client}

	natId, snatRuleId, err := parseNatSnatRuleResourceId(d.Id())
	if err != nil {
		return err
	}

	action := "Query NAT " + natId + " SNAT Rule " + snatRuleId

	snatRules, err := vpcService.ListAllNatSnatRulesWithNatID(natId)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_nat_snat_rule", action, BCESDKGoERROR)
	}
	addDebug(action, snatRules)

	var found bool
	for _, p := range snatRules {
		if p.RuleId == snatRuleId {
			found = true
			d.Set("rule_id", p.RuleId)
			d.Set("rule_name", p.RuleName)
			d.Set("public_ips_address", p.PublicIpAddresses)
			d.Set("source_cidr", p.SourceCIDR)
			d.Set("status", p.Status)
			return nil
		}
	}
	if !found {
		log.Printf("[WARN] Unable to find SNAT rule for NAT %s with SNAT Rule %s", natId, snatRuleId)
		d.SetId("")
	}

	return nil
}

func resourceBaiduCloudNatSnatRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	natId, snatRuleId, err := parseNatSnatRuleResourceId(d.Id())
	if err != nil {
		return err
	}

	action := "Update NAT " + natId + " SNAT Rule " + snatRuleId

	d.Partial(true)
	if d.HasChange("rule_name") || d.HasChange("public_ips_address") || d.HasChange("source_cidr") {
		args := &vpc.UpdateNatGatewaySnatRuleArgs{}
		if v := d.Get("rule_name").(string); v != "" {
			args.RuleName = v
		}
		if v := d.Get("public_ips_address"); v != nil {
			for _, id := range v.(*schema.Set).List() {
				args.PublicIpAddresses = append(args.PublicIpAddresses, id.(string))
			}
		}
		if v := d.Get("source_cidr").(string); v != "" {
			args.SourceCIDR = v
		}

		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return nil, vpcClient.UpdateNatGatewaySnatRule(natId, snatRuleId, args)
		})
		addDebug(action, args)
		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_nat_snat_rule", action, BCESDKGoERROR)
		}

		if len(args.RuleName) > 0 {
			d.SetPartial("rule_name")
		}
		if len(args.PublicIpAddresses) > 0 {
			d.SetPartial("public_ips_address")
		}
		if len(args.SourceCIDR) > 0 {
			d.SetPartial("source_cidr")
		}
	}
	d.Partial(false)

	return resourceBaiduCloudNatSnatRuleRead(d, meta)
}

func resourceBaiduCloudNatSnatRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	natId, snatRuleId, err := parseNatSnatRuleResourceId(d.Id())
	if err != nil {
		return err
	}

	action := "Delete NAT " + natId + " SNAT Rule " + snatRuleId

	clientToken := buildClientToken()

	_, err = client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
		return nil, vpcClient.DeleteNatGatewaySnatRule(natId, snatRuleId, clientToken)
	})
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_nat_SnatRule", action, BCESDKGoERROR)
	}
	addDebug(action, nil)

	return nil
}

func buildBaiduCloudNatSnatRuleArgs(d *schema.ResourceData) *vpc.CreateNatGatewaySnatRuleArgs {
	args := &vpc.CreateNatGatewaySnatRuleArgs{
		ClientToken: buildClientToken(),
	}

	if v := d.Get("rule_name").(string); v != "" {
		args.RuleName = v
	}
	if v, ok := d.GetOk("public_ips_address"); ok {
		for _, id := range v.(*schema.Set).List() {
			args.PublicIpAddresses = append(args.PublicIpAddresses, id.(string))
		}
	}
	if v := d.Get("source_cidr").(string); v != "" {
		args.SourceCIDR = v
	}

	return args
}

func parseNatSnatRuleResourceId(id string) (string, string, error) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return "", "", WrapErrorf(nil, DefaultErrorMsg, "baiducloud_nat_snat_rule",
			"parse Nat Snat rule resource id", BCESDKGoERROR)
	}
	return parts[0], parts[1], nil
}

func getNatSnatRuleResourceId(snatRuleId string, natId string) string {
	return strings.Join([]string{natId, snatRuleId}, ":")
}
