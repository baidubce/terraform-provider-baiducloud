package cdn

import (
	"github.com/baidubce/bce-sdk-go/services/cdn/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/flex"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/hashcode"
	"strconv"
)

//<editor-fold desc="DomainStatus">

func flattenDomainStatuses(domainStatuses []api.DomainStatus) interface{} {
	tfList := []map[string]interface{}{}
	for _, v := range domainStatuses {
		tfList = append(tfList, map[string]interface{}{
			"domain": v.Domain,
			"status": v.Status,
		})
	}
	return tfList
}

//</editor-fold>

//<editor-fold desc="OriginPeer">
func flattenOriginPeers(originPeers []api.OriginPeer) interface{} {
	tfList := []map[string]interface{}{}
	for _, v := range originPeers {
		tfList = append(tfList, map[string]interface{}{
			"peer":   v.Peer,
			"host":   v.Host,
			"backup": v.Backup,
		})
	}
	return tfList
}

func expandOriginPeers(tfList []interface{}) []api.OriginPeer {
	originPeers := []api.OriginPeer{}
	for _, v := range tfList {
		tfMap := v.(map[string]interface{})
		originPeers = append(originPeers, api.OriginPeer{
			Peer:   tfMap["peer"].(string),
			Host:   tfMap["host"].(string),
			Backup: tfMap["backup"].(bool),
		})
	}
	return originPeers
}

//</editor-fold>

//<editor-fold desc="CacheConfig">

//<editor-fold desc="CacheTTL">
func flattenCacheTTLs(cacheTTLs []api.CacheTTL) *schema.Set {
	tfSet := schema.NewSet(cacheTTLHash, nil)
	for _, v := range cacheTTLs {
		tfSet.Add(map[string]interface{}{
			"type":   v.Type,
			"value":  v.Value,
			"weight": v.Weight,
			"ttl":    v.TTL,
		})
	}
	return tfSet
}

func expandCacheTTLs(tfSet *schema.Set) []api.CacheTTL {
	list := []api.CacheTTL{}
	for _, v := range tfSet.List() {
		tfMap := v.(map[string]interface{})
		list = append(list, api.CacheTTL{
			Type:   tfMap["type"].(string),
			Value:  tfMap["value"].(string),
			Weight: tfMap["weight"].(int),
			TTL:    tfMap["ttl"].(int),
		})
	}
	return list
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
	return hashcode.Strings(s)
}

//</editor-fold>

//<editor-fold desc="CacheUrlArgs">
func flattenCacheUrlArgs(cacheUrlArgs *api.CacheUrlArgs) []interface{} {
	if cacheUrlArgs == nil {
		return []interface{}{}
	}
	tfMap := map[string]interface{}{
		"cache_full_url": cacheUrlArgs.CacheFullUrl,
	}
	if !cacheUrlArgs.CacheFullUrl {
		tfMap["cache_url_args"] = flex.FlattenStringValueSet(cacheUrlArgs.CacheUrlArgs)
	}
	return []interface{}{tfMap}
}

func expandCacheUrlArgs(tfList []interface{}) *api.CacheUrlArgs {
	cacheUrlArgs := &api.CacheUrlArgs{
		CacheFullUrl: true,
	}
	if len(tfList) == 0 || tfList[0] == nil {
		return cacheUrlArgs
	}
	tfMap := tfList[0].(map[string]interface{})

	if v, ok := tfMap["cache_full_url"].(bool); ok {
		cacheUrlArgs.CacheFullUrl = v
	}
	if !cacheUrlArgs.CacheFullUrl {
		cacheUrlArgs.CacheUrlArgs = flex.ExpandStringValueSet(tfMap["cache_url_args"].(*schema.Set))
	}
	return cacheUrlArgs
}

//</editor-fold>

//<editor-fold desc="ErrorPage">
func flattenErrorPages(errorPages []api.ErrorPage) *schema.Set {
	tfSet := schema.NewSet(errorPageHash, nil)
	for _, v := range errorPages {
		tfSet.Add(map[string]interface{}{
			"code": v.Code,
			"url":  v.Url,
		})
	}
	return tfSet
}

func expandErrorPages(tfSet *schema.Set) []api.ErrorPage {
	list := []api.ErrorPage{}
	for _, v := range tfSet.List() {
		tfMap := v.(map[string]interface{})
		list = append(list, api.ErrorPage{
			Code: tfMap["code"].(int),
			Url:  tfMap["url"].(string),
		})
	}
	return list
}

func errorPageHash(v interface{}) int {
	item := v.(map[string]interface{})
	var s []string

	if v, ok := item["code"]; ok {
		s = append(s, strconv.Itoa(v.(int)))
	}
	if v, ok := item["redirect_code"]; ok {
		s = append(s, strconv.Itoa(v.(int)))
	}
	if v, ok := item["url"]; ok {
		s = append(s, v.(string))
	}

	return hashcode.Strings(s)
}

//</editor-fold>

//<editor-fold desc="CacheShare">
func flattenCacheShare(cacheShard *api.CacheShared) []interface{} {
	if cacheShard == nil {
		return []interface{}{}
	}

	return []interface{}{map[string]interface{}{
		"enabled": cacheShard.Enabled,
		"domain":  cacheShard.SharedWith,
	}}
}

func expandCacheShare(tfList []interface{}) *api.CacheShared {
	cacheShared := &api.CacheShared{}

	if len(tfList) == 0 || tfList[0] == nil {
		return cacheShared
	}
	tfMap := tfList[0].(map[string]interface{})

	if v, ok := tfMap["enabled"].(bool); ok {
		cacheShared.Enabled = v
	}
	if cacheShared.Enabled {
		cacheShared.SharedWith = tfMap["domain"].(string)
	}

	return cacheShared
}

//</editor-fold>

//<editor-fold desc="MobileAccess">
func flattenMobileAccess(distinguishClient bool) []interface{} {
	return []interface{}{map[string]interface{}{
		"distinguish_client": distinguishClient,
	}}
}

func expandMobileAccess(tfList []interface{}) bool {
	if len(tfList) == 0 || tfList[0] == nil {
		return false
	}
	tfMap := tfList[0].(map[string]interface{})

	if v, ok := tfMap["distinguish_client"].(bool); ok {
		return v
	}
	return false
}

//</editor-fold>

//</editor-fold>
