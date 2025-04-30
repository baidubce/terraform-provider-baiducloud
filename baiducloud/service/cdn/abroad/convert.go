package abroad

import (
	"github.com/baidubce/bce-sdk-go/services/cdn/abroad/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/flex"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/hashcode"
	"strconv"
)

func expandAbroadCacheTTLs(tfSet *schema.Set) []api.CacheTTL {
	var list []api.CacheTTL
	for _, v := range tfSet.List() {
		tfMap := v.(map[string]interface{})
		list = append(list, api.CacheTTL{
			Type:           tfMap["type"].(string),
			Value:          tfMap["value"].(string),
			Weight:         tfMap["weight"].(int),
			TTL:            tfMap["ttl"].(int),
			OverrideOrigin: tfMap["override_origin"].(bool),
		})
	}
	return list
}

func flattenAbroadCacheTTLs(cacheTTLs []api.CacheTTL) *schema.Set {
	tfSet := schema.NewSet(cacheTTLHash, nil)
	for _, v := range cacheTTLs {
		tfSet.Add(map[string]interface{}{
			"type":            v.Type,
			"value":           v.Value,
			"weight":          v.Weight,
			"ttl":             v.TTL,
			"override_origin": v.OverrideOrigin,
		})
	}
	return tfSet
}

func cacheTTLHash(v interface{}) int {
	tfMap := v.(map[string]interface{})
	var s []string

	if v, ok := tfMap["type"]; ok {
		s = append(s, v.(string))
	}
	if v, ok := tfMap["value"]; ok {
		s = append(s, v.(string))
	}
	if v, ok := tfMap["ttl"]; ok {
		s = append(s, strconv.Itoa(v.(int)))
	}
	if v, ok := tfMap["weight"]; ok {
		s = append(s, strconv.Itoa(v.(int)))
	}
	if v, ok := tfMap["override_origin"]; ok {
		s = append(s, strconv.FormatBool(v.(bool)))
	}
	return hashcode.Strings(s)
}

func flattenAbroadOriginPeers(originPeers []api.OriginPeer) interface{} {
	var tfList []map[string]interface{}
	for _, v := range originPeers {
		tfList = append(tfList, map[string]interface{}{
			"type":   v.Type,
			"addr":   v.Addr,
			"backup": v.Backup,
		})
	}
	return tfList
}

func expandRefererACL(tfList []interface{}) *api.RefererACL {
	refererACL := &api.RefererACL{
		AllowEmpty: true,
	}
	if len(tfList) == 0 || tfList[0] == nil {
		refererACL.WhiteList = make([]string, 0)
		return refererACL
	}
	tfMap := tfList[0].(map[string]interface{})
	whiteList := flex.ExpandStringValueSet(tfMap["white_list"].(*schema.Set))
	blackList := flex.ExpandStringValueSet(tfMap["black_list"].(*schema.Set))
	if len(whiteList) > 0 {
		refererACL.WhiteList = whiteList
	}
	if len(blackList) > 0 {
		refererACL.BlackList = blackList
	}
	return refererACL
}

func expandIpACL(tfList []interface{}) *api.IpACL {
	ipACL := &api.IpACL{}
	if len(tfList) == 0 || tfList[0] == nil {
		ipACL.WhiteList = nil
		return ipACL
	}
	tfMap := tfList[0].(map[string]interface{})
	whiteList := flex.ExpandStringValueSet(tfMap["white_list"].(*schema.Set))
	blackList := flex.ExpandStringValueSet(tfMap["black_list"].(*schema.Set))
	if len(whiteList) > 0 {
		ipACL.WhiteList = whiteList
	}
	if len(blackList) > 0 {
		ipACL.BlackList = blackList
	}
	return ipACL
}

func flattenRefererACL(refererACL *api.RefererACL) []interface{} {
	blackList := flex.FlattenStringValueSet(refererACL.BlackList)
	whiteList := flex.FlattenStringValueSet(refererACL.WhiteList)
	if blackList.Len() == 0 && whiteList.Len() == 0 {
		return nil
	}
	return []interface{}{map[string]interface{}{
		"black_list": blackList,
		"white_list": whiteList,
	}}
}

func flattenIpACL(ipACL *api.IpACL) []interface{} {
	if ipACL == nil {
		return nil
	}
	return []interface{}{map[string]interface{}{
		"black_list": flex.FlattenStringValueSet(ipACL.BlackList),
		"white_list": flex.FlattenStringValueSet(ipACL.WhiteList),
	}}
}
