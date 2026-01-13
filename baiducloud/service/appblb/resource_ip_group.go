package appblb

import (
	"errors"
	"fmt"
	"strings"

	"github.com/baidubce/bce-sdk-go/services/appblb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/flex"
)

func ResourceIpGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to manage Application LoadBalancer instance's IP group.",

		Create: resourceIpGroupCreate,
		Read:   resourceIpGroupRead,
		Update: resourceIpGroupUpdate,
		Delete: resourceIpGroupDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"blb_id": {
				Type:        schema.TypeString,
				Description: "ID of the Application LoadBalancer instance.",
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Name of the IP group. Maximum `65` characters, starts with a letter, may contain letters, digits, `-/.` characters. Auto-generated if not set.",
				ValidateFunc: validation.StringLenBetween(1, 65),
			},
			"desc": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Description of the IP group. Maximum `200` characters.",
				ValidateFunc: validation.StringLenBetween(1, 200),
			},
			"backend_policy_list": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "Associated backend policy list of the IP group.",
				Set:         backendPolicyHash,
				Elem: &schema.Resource{
					Schema: backendPolicySchema(),
				},
			},
			"member_list": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "Member list of the IP group.",
				Set:         memberHash,
				Elem: &schema.Resource{
					Schema: memberSchema(),
				},
			},
		},

		CustomizeDiff: validateIpGroupDiff,
	}
}

func backendPolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Backend policy ID of the IP group.",
		},
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "Protocol type of the IP group policy. Valid values: `TCP`, `HTTP`, `HTTPS`, `UDP`.",
			ValidateFunc: validation.StringInSlice(IpGroupPolicyTypes, false),
		},
		"enable_health_check": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Whether to enable health check. Default: `false`.",
		},
		"health_check": {
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			Description:      "Health check type. Defaults to policy type. Valid values: `TCP`, `HTTP`, `HTTPS`, `UDP`, `ICMP`. Allowed combinations: policy type `TCP` allows `TCP`; `UDP` allows `UDP`/`ICMP`; `HTTP` allows `TCP`/`HTTP`; `HTTPS` allows `TCP`/`HTTP`/`HTTPS`.",
			ValidateFunc:     validation.StringInSlice(IpGroupHealthCheckTypes, false),
			DiffSuppressFunc: ipGroupPolicyHealthCheckDisabledSuppressFunc,
		},
		"health_check_port": {
			Type:             schema.TypeInt,
			Optional:         true,
			Computed:         true,
			Description:      "Health check port. Required when policy type is `HTTP` or `HTTPS`. Range `1-65535`.",
			ValidateFunc:     validation.IntBetween(1, 65535),
			DiffSuppressFunc: ipGroupPolicyHealthCheckPortSuppressFunc,
		},
		"health_check_url_path": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "/",
			Description:      "Health check URL path. Default: `/`. Effective when health check type is `HTTP`.",
			DiffSuppressFunc: ipGroupPolicyHealthCheckHTTPSuppressFunc,
		},
		"health_check_timeout_in_second": {
			Type:             schema.TypeInt,
			Optional:         true,
			Default:          3,
			Description:      "Health check timeout in seconds. Range `1-60`. Default: `3`.",
			ValidateFunc:     validation.IntBetween(1, 60),
			DiffSuppressFunc: ipGroupPolicyHealthCheckDisabledSuppressFunc,
		},
		"health_check_interval_in_second": {
			Type:             schema.TypeInt,
			Optional:         true,
			Default:          3,
			Description:      "Health check interval in seconds. Range `1-10`. Default: `3`.",
			ValidateFunc:     validation.IntBetween(1, 10),
			DiffSuppressFunc: ipGroupPolicyHealthCheckDisabledSuppressFunc,
		},
		"health_check_down_retry": {
			Type:             schema.TypeInt,
			Optional:         true,
			Default:          3,
			Description:      "Unhealthy threshold. Number of consecutive failed health checks before marking the backend as unavailable. Range `2-5`. Default: `3`.",
			ValidateFunc:     validation.IntBetween(2, 5),
			DiffSuppressFunc: ipGroupPolicyHealthCheckDisabledSuppressFunc,
		},
		"health_check_up_retry": {
			Type:             schema.TypeInt,
			Optional:         true,
			Default:          3,
			Description:      "Healthy threshold. Number of consecutive successful health checks before marking the backend as available. Range `2-5`. Default: `3`.",
			ValidateFunc:     validation.IntBetween(2, 5),
			DiffSuppressFunc: ipGroupPolicyHealthCheckDisabledSuppressFunc,
		},
		"health_check_normal_status": {
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			Description:      "Health check normal status for HTTP checks. Supported values: `http_2xx`, `http_3xx`, `http_4xx`, `http_5xx`, or a `|`-separated combination (for example, `http_2xx|http_3xx`). Default: `http_2xx|http_3xx`. Effective when health check type is `HTTP`.",
			ValidateFunc:     validateHealthCheckNormalStatus(),
			DiffSuppressFunc: ipGroupPolicyHealthCheckHTTPSuppressFunc,
		},
		"health_check_host": {
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			Description:      "Host header for L7 health checks (for example, `localhost`). Effective when health check type is `HTTP`.",
			DiffSuppressFunc: ipGroupPolicyHealthCheckHTTPSuppressFunc,
		},
		"udp_health_check_string": {
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			Description:      "UDP health check string. Required when health check type is `UDP`.",
			DiffSuppressFunc: ipGroupPolicyHealthCheckUDPSuppressFunc,
		},
	}
}

func memberSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"member_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IP group member ID.",
		},
		"ip": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "IPv4 address.",
		},
		"port": {
			Type:         schema.TypeInt,
			Required:     true,
			Description:  "Port number. Range `1-65535`.",
			ValidateFunc: validation.IntBetween(1, 65535),
		},
		"weight": {
			Type:         schema.TypeInt,
			Required:     true,
			Description:  "Weight. Range `0-100`.",
			ValidateFunc: validation.IntBetween(0, 100),
		},
	}
}

func resourceIpGroupCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	blbId := d.Get("blb_id").(string)

	args := &appblb.CreateAppIpGroupArgs{
		ClientToken: flex.BuildClientToken(),
	}
	if v, ok := d.GetOk("name"); ok {
		args.Name = v.(string)
	}
	if v, ok := d.GetOk("desc"); ok {
		args.Desc = v.(string)
	}
	if v, ok := d.GetOk("member_list"); ok {
		args.MemberList = expandMemberList(v.(*schema.Set).List())
	}

	raw, err := conn.WithAppBLBClient(func(client *appblb.Client) (interface{}, error) {
		return client.CreateAppIpGroup(blbId, args)
	})
	if err != nil {
		return fmt.Errorf("error creating appblb ip group: %w", err)
	}
	result := raw.(*appblb.CreateAppIpGroupResult)
	d.SetId(result.Id)

	if v, ok := d.GetOk("backend_policy_list"); ok {
		if err := createIpGroupBackendPolicies(conn, blbId, d.Id(), v.(*schema.Set).List()); err != nil {
			return err
		}
	}

	return resourceIpGroupRead(d, meta)
}

func resourceIpGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	blbId := d.Get("blb_id").(string)
	ipGroupId := d.Id()

	group, err := findIpGroup(conn, blbId, ipGroupId)
	if err != nil {
		if errors.Is(err, errIpGroupNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error reading appblb ip group %s: %w", ipGroupId, err)
	}

	if err := d.Set("backend_policy_list", flattenBackendPolicies(group.BackendPolicyList)); err != nil {
		return fmt.Errorf("error setting backend_policy_list: %w", err)
	}

	members, err := listIpGroupMembers(conn, blbId, ipGroupId)
	if err != nil {
		return err
	}
	if err := d.Set("member_list", flattenMembers(members)); err != nil {
		return fmt.Errorf("error setting member_list: %w", err)
	}

	d.Set("name", group.Name)
	d.Set("desc", group.Desc)
	d.Set("blb_id", blbId)

	return nil
}

func resourceIpGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	blbId := d.Get("blb_id").(string)
	ipGroupId := d.Id()

	if d.HasChange("name") || d.HasChange("desc") {
		args := &appblb.UpdateAppIpGroupArgs{
			IpGroupId:   ipGroupId,
			ClientToken: flex.BuildClientToken(),
		}
		if v, ok := d.GetOk("name"); ok {
			args.Name = v.(string)
		}
		if v, ok := d.GetOk("desc"); ok {
			args.Desc = v.(string)
		}

		_, err := conn.WithAppBLBClient(func(client *appblb.Client) (interface{}, error) {
			return nil, client.UpdateAppIpGroup(blbId, args)
		})
		if err != nil {
			return fmt.Errorf("error updating appblb ip group %s: %w", ipGroupId, err)
		}
	}

	if d.HasChange("backend_policy_list") {
		var list []interface{}
		if v, ok := d.GetOk("backend_policy_list"); ok {
			list = v.(*schema.Set).List()
		}
		if err := updateIpGroupBackendPolicies(conn, blbId, ipGroupId, list); err != nil {
			return err
		}
	}

	if d.HasChange("member_list") {
		var list []interface{}
		if v, ok := d.GetOk("member_list"); ok {
			list = v.(*schema.Set).List()
		}
		if err := updateIpGroupMembers(conn, blbId, ipGroupId, list); err != nil {
			return err
		}
	}

	return resourceIpGroupRead(d, meta)
}

func resourceIpGroupDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	blbId := d.Get("blb_id").(string)
	ipGroupId := d.Id()

	if _, err := findIpGroup(conn, blbId, ipGroupId); err != nil {
		if errors.Is(err, errIpGroupNotFound) {
			return nil
		}
		return fmt.Errorf("error checking appblb ip group %s: %w", ipGroupId, err)
	}

	args := &appblb.DeleteAppIpGroupArgs{
		IpGroupId:   ipGroupId,
		ClientToken: flex.BuildClientToken(),
	}

	_, err := conn.WithAppBLBClient(func(client *appblb.Client) (interface{}, error) {
		return nil, client.DeleteAppIpGroup(blbId, args)
	})
	if err != nil {
		return fmt.Errorf("error deleting appblb ip group %s: %w", ipGroupId, err)
	}

	return nil
}

func createIpGroupBackendPolicies(conn *connectivity.BaiduClient, blbId, ipGroupId string, list []interface{}) error {
	if len(list) == 0 {
		return nil
	}

	for _, item := range list {
		policy := expandBackendPolicy(item.(map[string]interface{}))
		args := buildCreateBackendPolicyArgs(ipGroupId, policy)
		_, err := conn.WithAppBLBClient(func(client *appblb.Client) (interface{}, error) {
			return nil, client.CreateAppIpGroupBackendPolicy(blbId, args)
		})
		if err != nil {
			return fmt.Errorf("error creating appblb ip group backend policy (%s): %w", policy.Type, err)
		}
	}

	return nil
}

func updateIpGroupBackendPolicies(conn *connectivity.BaiduClient, blbId, ipGroupId string, list []interface{}) error {
	desired := map[string]appblb.AppIpGroupBackendPolicy{}
	for _, item := range list {
		policy := expandBackendPolicy(item.(map[string]interface{}))
		desired[policy.Type] = policy
	}

	group, err := findIpGroup(conn, blbId, ipGroupId)
	if err != nil {
		return err
	}

	existing := map[string]appblb.AppIpGroupBackendPolicy{}
	for _, policy := range group.BackendPolicyList {
		existing[policy.Type] = policy
	}

	deleteIDs := make([]string, 0)
	for policyType, policy := range existing {
		if _, ok := desired[policyType]; !ok {
			deleteIDs = append(deleteIDs, policy.Id)
		}
	}
	if len(deleteIDs) > 0 {
		args := &appblb.DeleteAppIpGroupBackendPolicyArgs{
			IpGroupId:           ipGroupId,
			BackendPolicyIdList: deleteIDs,
			ClientToken:         flex.BuildClientToken(),
		}
		_, err := conn.WithAppBLBClient(func(client *appblb.Client) (interface{}, error) {
			return nil, client.DeleteAppIpGroupBackendPolicy(blbId, args)
		})
		if err != nil {
			return fmt.Errorf("error deleting appblb ip group backend policy: %w", err)
		}
	}

	for policyType, cfg := range desired {
		if policy, ok := existing[policyType]; ok {
			args := buildUpdateBackendPolicyArgs(ipGroupId, policy.Id, cfg)
			_, err := conn.WithAppBLBClient(func(client *appblb.Client) (interface{}, error) {
				return nil, client.UpdateAppIpGroupBackendPolicy(blbId, args)
			})
			if err != nil {
				return fmt.Errorf("error updating appblb ip group backend policy (%s): %w", policyType, err)
			}
		} else {
			args := buildCreateBackendPolicyArgs(ipGroupId, cfg)
			_, err := conn.WithAppBLBClient(func(client *appblb.Client) (interface{}, error) {
				return nil, client.CreateAppIpGroupBackendPolicy(blbId, args)
			})
			if err != nil {
				return fmt.Errorf("error creating appblb ip group backend policy (%s): %w", policyType, err)
			}
		}
	}

	return nil
}

func updateIpGroupMembers(conn *connectivity.BaiduClient, blbId, ipGroupId string, list []interface{}) error {
	existingMembers, err := listIpGroupMembers(conn, blbId, ipGroupId)
	if err != nil {
		return err
	}

	existing := map[string]appblb.AppIpGroupMember{}
	for _, member := range existingMembers {
		existing[memberKey(member.Ip, member.Port)] = member
	}

	desired := map[string]appblb.AppIpGroupMember{}
	for _, item := range list {
		member := expandMember(item.(map[string]interface{}))
		desired[memberKey(member.Ip, member.Port)] = member
	}

	deleteIDs := make([]string, 0)
	for key, member := range existing {
		if _, ok := desired[key]; !ok {
			deleteIDs = append(deleteIDs, member.MemberId)
		}
	}
	if len(deleteIDs) > 0 {
		args := &appblb.DeleteAppIpGroupMemberArgs{
			IpGroupId:    ipGroupId,
			MemberIdList: deleteIDs,
			ClientToken:  flex.BuildClientToken(),
		}
		_, err := conn.WithAppBLBClient(func(client *appblb.Client) (interface{}, error) {
			return nil, client.DeleteAppIpGroupMember(blbId, args)
		})
		if err != nil {
			return fmt.Errorf("error deleting appblb ip group members: %w", err)
		}
	}

	updateList := make([]appblb.AppIpGroupMember, 0)
	for key, desiredMember := range desired {
		if member, ok := existing[key]; ok {
			if weightValue(member.Weight) != weightValue(desiredMember.Weight) {
				updateList = append(updateList, appblb.AppIpGroupMember{
					MemberId: member.MemberId,
					Port:     desiredMember.Port,
					Weight:   desiredMember.Weight,
				})
			}
		}
	}
	if len(updateList) > 0 {
		args := &appblb.UpdateAppIpGroupMemberArgs{
			AppIpGroupMemberWriteOpArgs: appblb.AppIpGroupMemberWriteOpArgs{
				IpGroupId:   ipGroupId,
				MemberList:  updateList,
				ClientToken: flex.BuildClientToken(),
			},
		}
		_, err := conn.WithAppBLBClient(func(client *appblb.Client) (interface{}, error) {
			return nil, client.UpdateAppIpGroupMember(blbId, args)
		})
		if err != nil {
			return fmt.Errorf("error updating appblb ip group members: %w", err)
		}
	}

	addList := make([]appblb.AppIpGroupMember, 0)
	for key, desiredMember := range desired {
		if _, ok := existing[key]; !ok {
			addList = append(addList, desiredMember)
		}
	}
	if len(addList) > 0 {
		args := &appblb.CreateAppIpGroupMemberArgs{
			AppIpGroupMemberWriteOpArgs: appblb.AppIpGroupMemberWriteOpArgs{
				IpGroupId:   ipGroupId,
				MemberList:  addList,
				ClientToken: flex.BuildClientToken(),
			},
		}
		_, err := conn.WithAppBLBClient(func(client *appblb.Client) (interface{}, error) {
			return nil, client.CreateAppIpGroupMember(blbId, args)
		})
		if err != nil {
			return fmt.Errorf("error creating appblb ip group members: %w", err)
		}
	}

	return nil
}

func buildCreateBackendPolicyArgs(ipGroupId string, policy appblb.AppIpGroupBackendPolicy) *appblb.CreateAppIpGroupBackendPolicyArgs {
	enableHealthCheck := &policy.EnableHealthCheck
	return &appblb.CreateAppIpGroupBackendPolicyArgs{
		ClientToken:                 flex.BuildClientToken(),
		IpGroupId:                   ipGroupId,
		Type:                        policy.Type,
		EnableHealthCheck:           enableHealthCheck,
		HealthCheck:                 policy.HealthCheck,
		HealthCheckPort:             policy.HealthCheckPort,
		HealthCheckUrlPath:          policy.HealthCheckUrlPath,
		HealthCheckTimeoutInSecond:  policy.HealthCheckTimeoutInSecond,
		HealthCheckIntervalInSecond: policy.HealthCheckIntervalInSecond,
		HealthCheckDownRetry:        policy.HealthCheckDownRetry,
		HealthCheckUpRetry:          policy.HealthCheckUpRetry,
		HealthCheckNormalStatus:     policy.HealthCheckNormalStatus,
		HealthCheckHost:             policy.HealthCheckHost,
		UdpHealthCheckString:        policy.UdpHealthCheckString,
	}
}

func buildUpdateBackendPolicyArgs(ipGroupId, policyId string, policy appblb.AppIpGroupBackendPolicy) *appblb.UpdateAppIpGroupBackendPolicyArgs {
	enableHealthCheck := &policy.EnableHealthCheck
	return &appblb.UpdateAppIpGroupBackendPolicyArgs{
		ClientToken:                 flex.BuildClientToken(),
		IpGroupId:                   ipGroupId,
		Id:                          policyId,
		EnableHealthCheck:           enableHealthCheck,
		HealthCheck:                 policy.HealthCheck,
		HealthCheckPort:             policy.HealthCheckPort,
		HealthCheckHost:             policy.HealthCheckHost,
		HealthCheckUrlPath:          policy.HealthCheckUrlPath,
		HealthCheckTimeoutInSecond:  policy.HealthCheckTimeoutInSecond,
		HealthCheckIntervalInSecond: policy.HealthCheckIntervalInSecond,
		HealthCheckDownRetry:        policy.HealthCheckDownRetry,
		HealthCheckUpRetry:          policy.HealthCheckUpRetry,
		HealthCheckNormalStatus:     policy.HealthCheckNormalStatus,
		UdpHealthCheckString:        policy.UdpHealthCheckString,
	}
}

func validateIpGroupDiff(diff *schema.ResourceDiff, _ interface{}) error {
	if v, ok := diff.GetOk("backend_policy_list"); ok {
		seen := map[string]struct{}{}
		for _, item := range v.(*schema.Set).List() {
			policy := expandBackendPolicy(item.(map[string]interface{}))
			if _, exists := seen[policy.Type]; exists {
				return fmt.Errorf("backend_policy_list has duplicate type %q", policy.Type)
			}
			seen[policy.Type] = struct{}{}

			if !policy.EnableHealthCheck {
				continue
			}
			checkType := policy.HealthCheck
			if checkType == "" {
				checkType = policy.Type
			}
			if !ipGroupPolicyHealthCheckAllowed(policy.Type, checkType) {
				return fmt.Errorf("backend_policy_list health_check %q is not allowed for type %q", checkType, policy.Type)
			}
			if policy.Type == ProtocolHTTP || policy.Type == ProtocolHTTPS {
				if policy.HealthCheckPort == 0 {
					return fmt.Errorf("backend_policy_list health_check_port is required for type %q", policy.Type)
				}
			}
			if checkType == ProtocolUDP && policy.UdpHealthCheckString == "" {
				return fmt.Errorf("backend_policy_list udp_health_check_string is required when health_check is %q", checkType)
			}
		}
	}

	if v, ok := diff.GetOk("member_list"); ok {
		seen := map[string]struct{}{}
		for _, item := range v.(*schema.Set).List() {
			member := expandMember(item.(map[string]interface{}))
			key := memberKey(member.Ip, member.Port)
			if _, exists := seen[key]; exists {
				return fmt.Errorf("member_list has duplicate ip+port %q", key)
			}
			seen[key] = struct{}{}
		}
	}

	return nil
}

func ipGroupPolicyHealthCheckAllowed(policyType, checkType string) bool {
	switch policyType {
	case ProtocolTCP:
		return checkType == ProtocolTCP
	case ProtocolHTTP:
		return checkType == ProtocolHTTP || checkType == ProtocolTCP
	case ProtocolHTTPS:
		return checkType == ProtocolHTTP || checkType == ProtocolTCP || checkType == ProtocolHTTPS
	case ProtocolUDP:
		return checkType == ProtocolUDP || checkType == ProtocolICMP
	default:
		return false
	}
}

func validateHealthCheckNormalStatus() schema.SchemaValidateFunc {
	allowed := map[string]struct{}{
		"http_2xx": {},
		"http_3xx": {},
		"http_4xx": {},
		"http_5xx": {},
	}

	return func(value interface{}, key string) ([]string, []error) {
		raw := value.(string)
		if raw == "" {
			return nil, nil
		}

		parts := strings.Split(raw, "|")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				return nil, []error{fmt.Errorf("%s has empty status value", key)}
			}
			if _, ok := allowed[part]; !ok {
				return nil, []error{fmt.Errorf("%s has invalid status value %q", key, part)}
			}
		}

		return nil, nil
	}
}
