package appblb

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func ipGroupPolicyHealthCheckDisabledSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	_, enabled, ok := ipGroupPolicyEnableHealthCheck(k, d)
	if !ok {
		return false
	}
	return !enabled
}

func ipGroupPolicyHealthCheckHTTPSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	_, enabled, ok := ipGroupPolicyEnableHealthCheck(k, d)
	if ok && !enabled {
		return true
	}

	checkType := ipGroupPolicyHealthCheckType(k, d)
	return checkType != ProtocolHTTP && checkType != ProtocolHTTPS
}

func ipGroupPolicyHealthCheckUDPSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	_, enabled, ok := ipGroupPolicyEnableHealthCheck(k, d)
	if ok && !enabled {
		return true
	}

	return ipGroupPolicyHealthCheckType(k, d) != ProtocolUDP
}

func ipGroupPolicyHealthCheckPortSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	_, enabled, ok := ipGroupPolicyEnableHealthCheck(k, d)
	if ok && !enabled {
		return true
	}

	return ipGroupPolicyHealthCheckType(k, d) == ProtocolICMP
}

func ipGroupPolicyEnableHealthCheck(k string, d *schema.ResourceData) (string, bool, bool) {
	strs := strings.Split(k, ".")
	if len(strs) < 3 {
		return "", false, false
	}

	key := "backend_policy_list." + strs[1] + ".enable_health_check"
	value, ok := d.GetOkExists(key)
	if !ok {
		return strs[1], false, false
	}

	return strs[1], value.(bool), true
}

func ipGroupPolicyHealthCheckType(k string, d *schema.ResourceData) string {
	strs := strings.Split(k, ".")
	if len(strs) < 3 {
		return ""
	}

	key := "backend_policy_list." + strs[1] + ".health_check"
	if v, ok := d.GetOk(key); ok {
		return v.(string)
	}

	key = "backend_policy_list." + strs[1] + ".type"
	if v, ok := d.GetOk(key); ok {
		return v.(string)
	}

	return ""
}
