package baiducloud

import (
	"encoding/base64"
	"encoding/json"

	"github.com/baidubce/bce-sdk-go/services/bcc/api"
)

func (s *BccService) FlattenSecurityGroupRuleModelsToMap(list []api.SecurityGroupRuleModel) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))

	for _, r := range list {
		rule := map[string]interface{}{
			"remark":          r.Remark,
			"direction":       r.Direction,
			"ether_type":      r.Ethertype,
			"port_range":      r.PortRange,
			"protocol":        r.Protocol,
			"source_group_id": r.SourceGroupId,
			"source_ip":       r.SourceIp,
			"dest_group_id":   r.DestGroupId,
			"dest_ip":         r.DestIp,
		}

		result = append(result, rule)
	}

	return result
}

func (s *BccService) buildSecurityGroupRuleId(rule *api.SecurityGroupRuleModel) (string, error) {
	ruleByte, err := json.Marshal(rule)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ruleByte), nil
}

func (s *BccService) parseSecurityGroupRuleId(ruleId string) (*api.SecurityGroupRuleModel, error) {
	ruleByte, err := base64.StdEncoding.DecodeString(ruleId)
	if err != nil {
		return nil, err
	}

	rule := &api.SecurityGroupRuleModel{}
	if err := json.Unmarshal(ruleByte, rule); err != nil {
		return nil, err
	}

	return rule, nil
}

func (s *BccService) GetSecurityGroupRule(ruleId string) (*api.SecurityGroupRuleModel, error) {
	ruleInfo, err := s.parseSecurityGroupRuleId(ruleId)
	if err != nil {
		return nil, err
	}

	listArgs := &api.ListSecurityGroupArgs{}
	sgList, err := s.ListAllSecurityGroups(listArgs)
	if err != nil {
		return nil, WrapError(err)
	}

	var securityGroup api.SecurityGroupModel
	for _, sg := range sgList {
		if sg.Id == ruleInfo.SecurityGroupId {
			securityGroup = sg
			break
		}
	}

	if securityGroup.Id == "" {
		return nil, WrapError(Error(ResourceNotFound))
	}

	for _, r := range securityGroup.Rules {
		if compareSecurityGroupRule(&r, ruleInfo) {
			return &r, nil
		}
	}

	return nil, WrapError(Error(ResourceNotFound))
}

func compareSecurityGroupRule(r1, r2 *api.SecurityGroupRuleModel) bool {
	if r1 == r2 {
		return true
	}

	if r1 == nil || r2 == nil {
		return false
	}

	if r1.SecurityGroupId != r2.SecurityGroupId {
		return false
	}

	if r1.Direction != r2.Direction {
		return false
	}

	if !stringEqualWithDefault(r1.PortRange, r2.PortRange, []string{"", "1-65535"}) {
		return false
	}

	if !stringEqualWithDefault(r1.Protocol, r2.Protocol, []string{"", "all"}) {
		return false
	}

	if r1.Direction == "ingress" {
		if !stringEqualWithDefault(r1.SourceIp, r2.SourceIp, []string{"", "all"}) {
			return false
		}

		if r1.SourceGroupId != r2.SourceGroupId {
			return false
		}
	}

	if r1.Direction == "engress" {
		if !stringEqualWithDefault(r1.DestIp, r2.DestIp, []string{"", "all"}) {
			return false
		}

		if r1.DestGroupId != r2.DestGroupId {
			return false
		}
	}

	if !stringEqualWithDefault(r1.Ethertype, r2.Ethertype, []string{"", "IPv4"}) {
		return false
	}

	return true
}
