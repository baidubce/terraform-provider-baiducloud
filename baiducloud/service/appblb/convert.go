package appblb

import (
	"fmt"
	"strconv"

	"github.com/baidubce/bce-sdk-go/services/appblb"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/hashcode"
)

func memberKey(ip string, port int) string {
	return fmt.Sprintf("%s:%d", ip, port)
}

func backendPolicyHash(v interface{}) int {
	m := v.(map[string]interface{})
	policyType := m["type"].(string)
	enableHealthCheck := m["enable_health_check"].(bool)
	healthCheck := m["health_check"].(string)
	if healthCheck == "" {
		healthCheck = policyType
	}

	parts := []string{policyType, boolToString(enableHealthCheck)}
	if !enableHealthCheck {
		return hashcode.Strings(parts)
	}

	parts = append(parts, healthCheck)
	appendNonEmpty := func(value string) {
		if value != "" {
			parts = append(parts, value)
		}
	}
	appendNonZero := func(value int) {
		if value != 0 {
			parts = append(parts, intToString(value))
		}
	}

	if healthCheck != ProtocolICMP {
		appendNonZero(m["health_check_port"].(int))
	}
	appendNonZero(m["health_check_timeout_in_second"].(int))
	appendNonZero(m["health_check_interval_in_second"].(int))
	appendNonZero(m["health_check_down_retry"].(int))
	appendNonZero(m["health_check_up_retry"].(int))

	if healthCheck == ProtocolHTTP || healthCheck == ProtocolHTTPS {
		appendNonEmpty(m["health_check_url_path"].(string))
		appendNonEmpty(m["health_check_normal_status"].(string))
		appendNonEmpty(m["health_check_host"].(string))
	}
	if healthCheck == ProtocolUDP {
		appendNonEmpty(m["udp_health_check_string"].(string))
	}

	return hashcode.Strings(parts)
}

func memberHash(v interface{}) int {
	m := v.(map[string]interface{})
	return hashcode.Strings([]string{
		m["ip"].(string),
		intToString(m["port"].(int)),
		intToString(m["weight"].(int)),
	})
}

func expandBackendPolicy(m map[string]interface{}) appblb.AppIpGroupBackendPolicy {
	return appblb.AppIpGroupBackendPolicy{
		Type:                        m["type"].(string),
		EnableHealthCheck:           m["enable_health_check"].(bool),
		HealthCheck:                 m["health_check"].(string),
		HealthCheckPort:             m["health_check_port"].(int),
		HealthCheckUrlPath:          m["health_check_url_path"].(string),
		HealthCheckTimeoutInSecond:  m["health_check_timeout_in_second"].(int),
		HealthCheckIntervalInSecond: m["health_check_interval_in_second"].(int),
		HealthCheckDownRetry:        m["health_check_down_retry"].(int),
		HealthCheckUpRetry:          m["health_check_up_retry"].(int),
		HealthCheckNormalStatus:     m["health_check_normal_status"].(string),
		HealthCheckHost:             m["health_check_host"].(string),
		UdpHealthCheckString:        m["udp_health_check_string"].(string),
	}
}

func expandMemberList(list []interface{}) []appblb.AppIpGroupMember {
	if len(list) == 0 {
		return nil
	}

	members := make([]appblb.AppIpGroupMember, 0, len(list))
	for _, item := range list {
		members = append(members, expandMember(item.(map[string]interface{})))
	}

	return members
}

func expandMember(m map[string]interface{}) appblb.AppIpGroupMember {
	weight := m["weight"].(int)
	return appblb.AppIpGroupMember{
		Ip:     m["ip"].(string),
		Port:   m["port"].(int),
		Weight: &weight,
	}
}

func flattenBackendPolicies(list []appblb.AppIpGroupBackendPolicy) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, policy := range list {
		result = append(result, map[string]interface{}{
			"id":                              policy.Id,
			"type":                            policy.Type,
			"enable_health_check":             policy.EnableHealthCheck,
			"health_check":                    policy.HealthCheck,
			"health_check_port":               policy.HealthCheckPort,
			"health_check_url_path":           policy.HealthCheckUrlPath,
			"health_check_timeout_in_second":  policy.HealthCheckTimeoutInSecond,
			"health_check_interval_in_second": policy.HealthCheckIntervalInSecond,
			"health_check_down_retry":         policy.HealthCheckDownRetry,
			"health_check_up_retry":           policy.HealthCheckUpRetry,
			"health_check_normal_status":      policy.HealthCheckNormalStatus,
			"health_check_host":               policy.HealthCheckHost,
			"udp_health_check_string":         policy.UdpHealthCheckString,
		})
	}

	return result
}

func flattenMembers(list []appblb.AppIpGroupMember) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, member := range list {
		result = append(result, map[string]interface{}{
			"member_id": member.MemberId,
			"ip":        member.Ip,
			"port":      member.Port,
			"weight":    weightValue(member.Weight),
		})
	}

	return result
}

func weightValue(weight *int) int {
	if weight == nil {
		return 0
	}
	return *weight
}

func intToString(value int) string {
	return strconv.Itoa(value)
}

func boolToString(value bool) string {
	if value {
		return "true"
	}
	return "false"
}
