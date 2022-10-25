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

func flattenDomainCertificate(certificate *api.CertificateDetail) interface{} {
	return []map[string]interface{}{
		{
			"cert_id":          certificate.CertId,
			"cert_name":        certificate.CertName,
			"cert_common_name": certificate.CommonName,
			"cert_dns_names":   certificate.DNSNames,
			"cert_start_time":  certificate.StartTime,
			"cert_stop_time":   certificate.StopTime,
			"cert_create_time": certificate.CreateTime,
			"cert_update_time": certificate.UpdateTime,
		},
	}
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
			"weight": v.Weight,
			"isp":    v.ISP,
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
			Weight: tfMap["weight"].(int),
			ISP:    tfMap["isp"].(string),
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

//<editor-fold desc="ACLConfig">

//<editor-fold desc="RefererACL">
func flattenRefererACL(refererACL *api.RefererACL) []interface{} {
	return []interface{}{map[string]interface{}{
		"black_list":  flex.FlattenStringValueSet(refererACL.BlackList),
		"white_list":  flex.FlattenStringValueSet(refererACL.WhiteList),
		"allow_empty": refererACL.AllowEmpty,
	}}
}

func expandRefererACL(tfList []interface{}) *api.RefererACL {
	refererACL := &api.RefererACL{
		BlackList:  []string{},
		WhiteList:  []string{},
		AllowEmpty: true,
	}
	if len(tfList) == 0 || tfList[0] == nil {
		return refererACL
	}
	tfMap := tfList[0].(map[string]interface{})
	refererACL.AllowEmpty = tfMap["allow_empty"].(bool)
	refererACL.BlackList = flex.ExpandStringValueSet(tfMap["black_list"].(*schema.Set))
	refererACL.WhiteList = flex.ExpandStringValueSet(tfMap["white_list"].(*schema.Set))
	return refererACL
}

//</editor-fold>

//<editor-fold desc="IpACL">
func flattenIpACL(ipACL *api.IpACL) []interface{} {
	if ipACL == nil {
		return nil
	}
	return []interface{}{map[string]interface{}{
		"black_list": flex.FlattenStringValueSet(ipACL.BlackList),
		"white_list": flex.FlattenStringValueSet(ipACL.WhiteList),
	}}
}

func expandIpACL(tfList []interface{}) *api.IpACL {
	ipACL := &api.IpACL{
		BlackList: []string{},
		WhiteList: []string{},
	}
	if len(tfList) == 0 || tfList[0] == nil {
		return ipACL
	}
	tfMap := tfList[0].(map[string]interface{})
	ipACL.BlackList = flex.ExpandStringValueSet(tfMap["black_list"].(*schema.Set))
	ipACL.WhiteList = flex.ExpandStringValueSet(tfMap["white_list"].(*schema.Set))
	return ipACL
}

//</editor-fold>

//<editor-fold desc="UaACL">
func flattenUaACL(uaACL *api.UaACL) []interface{} {
	return []interface{}{map[string]interface{}{
		"black_list": flex.FlattenStringValueSet(uaACL.BlackList),
		"white_list": flex.FlattenStringValueSet(uaACL.WhiteList),
	}}
}

func expandUaACL(tfList []interface{}) *api.UaACL {
	uaACL := &api.UaACL{
		BlackList: []string{},
		WhiteList: []string{},
	}
	if len(tfList) == 0 || tfList[0] == nil {
		return uaACL
	}
	tfMap := tfList[0].(map[string]interface{})
	uaACL.BlackList = flex.ExpandStringValueSet(tfMap["black_list"].(*schema.Set))
	uaACL.WhiteList = flex.ExpandStringValueSet(tfMap["white_list"].(*schema.Set))
	return uaACL
}

//</editor-fold>

//<editor-fold desc="Cors">
func flattenCors(cors *api.CorsCfg) []interface{} {
	tfMap := map[string]interface{}{
		"allow":       "off",
		"origin_list": flex.FlattenStringValueSet(cors.Origins),
	}
	if cors.IsAllow {
		tfMap["allow"] = "on"
		tfMap["origin_list"] = flex.FlattenStringValueSet(cors.Origins)
	}
	return []interface{}{tfMap}
}

func expandCors(tfList []interface{}) api.CorsCfg {
	cors := api.CorsCfg{}

	if len(tfList) == 0 || tfList[0] == nil {
		return cors
	}
	tfMap := tfList[0].(map[string]interface{})

	if tfMap["allow"].(string) == "on" {
		cors.IsAllow = true
	}
	if cors.IsAllow {
		cors.Origins = flex.ExpandStringValueSet(tfMap["origin_list"].(*schema.Set))
	}

	return cors
}

//</editor-fold>

//<editor-fold desc="AccessLimit">
func flattenAccessLimit(accessLimit *api.AccessLimit) []interface{} {
	tfMap := map[string]interface{}{
		"enabled": accessLimit.Enabled,
	}
	if accessLimit.Enabled {
		tfMap["limit"] = accessLimit.Limit
	}
	return []interface{}{tfMap}
}

func expandAccessLimit(tfList []interface{}) *api.AccessLimit {
	accessLimit := &api.AccessLimit{}

	if len(tfList) == 0 || tfList[0] == nil {
		return accessLimit
	}
	tfMap := tfList[0].(map[string]interface{})

	accessLimit.Enabled = tfMap["enabled"].(bool)
	if accessLimit.Enabled {
		accessLimit.Limit = tfMap["limit"].(int)
	}

	return accessLimit
}

//</editor-fold>

//<editor-fold desc="TrafficLimit">
func flattenTrafficLimit(trafficLimit *api.TrafficLimit) []interface{} {
	tfMap := map[string]interface{}{
		"enable": trafficLimit.Enabled,
	}
	if trafficLimit.Enabled {
		tfMap["limit_rate"] = trafficLimit.LimitRate
		tfMap["limit_start_hour"] = trafficLimit.LimitStartHour
		tfMap["limit_end_hour"] = trafficLimit.LimitEndHour
	}

	return []interface{}{tfMap}
}

func expandTrafficLimit(tfList []interface{}) *api.TrafficLimit {
	trafficLimit := &api.TrafficLimit{}

	if len(tfList) == 0 || tfList[0] == nil {
		return trafficLimit
	}
	tfMap := tfList[0].(map[string]interface{})

	trafficLimit.Enabled = tfMap["enable"].(bool)
	if trafficLimit.Enabled {
		trafficLimit.LimitRate = tfMap["limit_rate"].(int)
		if v, ok := tfMap["limit_start_hour"]; ok {
			trafficLimit.LimitStartHour = v.(int)
		}
		if v, ok := tfMap["limit_end_hour"]; ok {
			trafficLimit.LimitEndHour = v.(int)
		}
	}

	return trafficLimit
}

//</editor-fold>

//<editor-fold desc="RequestAuth">
func flattenRequestAuth(requestAuth *api.RequestAuth) []interface{} {
	if requestAuth == nil {
		return nil
	}
	return []interface{}{map[string]interface{}{
		"type":             requestAuth.Type,
		"key1":             requestAuth.Key1,
		"key2":             requestAuth.Key2,
		"timeout":          requestAuth.Timeout,
		"timestamp_metric": requestAuth.TimestampMetric,
	}}
}

func expandRequestAuth(tfList []interface{}) *api.RequestAuth {
	if len(tfList) == 0 || tfList[0] == nil {
		return nil
	}
	requestAuth := &api.RequestAuth{}
	tfMap := tfList[0].(map[string]interface{})
	requestAuth.Type = tfMap["type"].(string)
	requestAuth.Key1 = tfMap["key1"].(string)
	requestAuth.Key2 = tfMap["key2"].(string)
	requestAuth.Timeout = tfMap["timeout"].(int)
	requestAuth.TimestampMetric = tfMap["timestamp_metric"].(int)
	return requestAuth
}

//</editor-fold>

//</editor-fold>

//<editor-fold desc="Advanced">

//<editor-fold desc="IPv6Dispatch">
func flattenIPv6Dispatch(ipv6Enabled bool) []interface{} {
	return []interface{}{map[string]interface{}{
		"enable": ipv6Enabled,
	}}
}

func expandIPv6Dispatch(tfList []interface{}) bool {
	if len(tfList) == 0 || tfList[0] == nil {
		return false
	}
	tfMap := tfList[0].(map[string]interface{})
	return tfMap["enable"].(bool)
}

//</editor-fold>

//<editor-fold desc="HttpHeader">
func flattenHttpHeaders(httpHeaders []api.HttpHeader) *schema.Set {
	tfSet := schema.NewSet(httpHeaderHash, nil)
	for _, v := range httpHeaders {
		tfSet.Add(map[string]interface{}{
			"type":     v.Type,
			"header":   v.Header,
			"value":    v.Value,
			"action":   v.Action,
			"describe": v.Describe,
		})
	}
	return tfSet
}

func expandHttpHeaders(tfSet *schema.Set) []api.HttpHeader {
	list := []api.HttpHeader{}
	for _, v := range tfSet.List() {
		tfMap := v.(map[string]interface{})
		list = append(list, api.HttpHeader{
			Type:     tfMap["type"].(string),
			Header:   tfMap["header"].(string),
			Value:    tfMap["value"].(string),
			Action:   tfMap["action"].(string),
			Describe: tfMap["describe"].(string),
		})
	}
	return list
}

func httpHeaderHash(v interface{}) int {
	tfMap := v.(map[string]interface{})
	var s []string

	if v, ok := tfMap["type"]; ok {
		s = append(s, v.(string))
	}
	if v, ok := tfMap["header"]; ok {
		s = append(s, v.(string))
	}
	if v, ok := tfMap["value"]; ok {
		s = append(s, v.(string))
	}
	if v, ok := tfMap["action"]; ok {
		s = append(s, v.(string))
	}
	if v, ok := tfMap["describe"]; ok {
		s = append(s, v.(string))
	}
	return hashcode.Strings(s)
}

//</editor-fold>

//<editor-fold desc="MediaDrag">
func flattenMediaDragConf(mediaDragConf *api.MediaDragConf) []interface{} {
	tfMap := map[string]interface{}{}
	if mediaDragConf == nil {
		return []interface{}{tfMap}
	}
	if mediaDragConf.Mp4 != nil {
		tfMap["mp4"] = flattenMediaCfg(mediaDragConf.Mp4)
	}
	if mediaDragConf.Flv != nil {
		tfMap["flv"] = flattenMediaCfg(mediaDragConf.Flv)
	}
	return []interface{}{tfMap}
}

func expandMediaDragConf(tfList []interface{}) *api.MediaDragConf {
	mediaDragConf := &api.MediaDragConf{}
	if len(tfList) == 0 || tfList[0] == nil {
		return mediaDragConf
	}
	tfMap := tfList[0].(map[string]interface{})
	if mediaCfg, ok := tfMap["mp4"].([]interface{}); ok {
		mediaDragConf.Mp4 = expandMediaCfg(mediaCfg)
	}
	if mediaCfg, ok := tfMap["flv"].([]interface{}); ok {
		mediaDragConf.Flv = expandMediaCfg(mediaCfg)
	}
	return mediaDragConf
}

func flattenMediaCfg(mediaCfg *api.MediaCfg) []interface{} {
	return []interface{}{map[string]interface{}{
		"file_suffix":    flex.FlattenStringValueSet(mediaCfg.FileSuffix),
		"start_arg_name": mediaCfg.StartArgName,
		"end_arg_name":   mediaCfg.EndArgName,
		"drag_mode":      mediaCfg.DragMode,
	}}
}

func expandMediaCfg(tfList []interface{}) *api.MediaCfg {
	if len(tfList) == 0 || tfList[0] == nil {
		return nil
	}
	tfMap := tfList[0].(map[string]interface{})
	return &api.MediaCfg{
		FileSuffix:   flex.ExpandStringValueSet(tfMap["file_suffix"].(*schema.Set)),
		StartArgName: tfMap["start_arg_name"].(string),
		EndArgName:   tfMap["end_arg_name"].(string),
		DragMode:     tfMap["drag_mode"].(string),
	}
}

//</editor-fold>

//<editor-fold desc="SeoSwitch">
func flattenSeoSwitch(seo *api.SeoSwitch) []interface{} {
	return []interface{}{map[string]interface{}{
		"directly_origin": seo.DirectlyOrigin,
	}}
}

func expandSeoSwitch(tfList []interface{}) *api.SeoSwitch {
	seo := &api.SeoSwitch{
		DirectlyOrigin: "OFF",
		PushRecord:     "OFF",
	}

	if len(tfList) == 0 || tfList[0] == nil {
		return seo
	}
	tfMap := tfList[0].(map[string]interface{})
	seo.DirectlyOrigin = tfMap["directly_origin"].(string)
	return seo
}

//</editor-fold>

//<editor-fold desc="Compress">
func flattenCompress(compressType string) []interface{} {
	tfMap := map[string]interface{}{
		"allow": false,
	}
	if len(compressType) > 0 {
		tfMap["allow"] = true
		tfMap["type"] = compressType
	}
	return []interface{}{tfMap}
}

func expandCompress(tfList []interface{}) (bool, string) {
	if len(tfList) == 0 || tfList[0] == nil {
		return false, ""
	}
	tfMap := tfList[0].(map[string]interface{})
	return tfMap["allow"].(bool), tfMap["type"].(string)
}

//</editor-fold>

//</editor-fold>

//<editor-fold desc="Https">
func flattenHttps(https *api.HTTPSConfig) []interface{} {
	tfMap := map[string]interface{}{
		"enabled": https.Enabled,
		"cert_id": https.CertId,
	}
	if https.Enabled {
		tfMap["http_redirect"] = https.HttpRedirect
		tfMap["http_redirect_code"] = https.HttpRedirectCode
		tfMap["https_redirect"] = https.HttpsRedirect
		tfMap["https_redirect_code"] = https.HttpsRedirectCode
		tfMap["http2_enabled"] = https.Http2Enabled
		tfMap["verify_client"] = https.VerifyClient
		tfMap["ssl_protocols"] = flex.FlattenStringValueSet(https.SslProtocols)
	}

	return []interface{}{tfMap}
}

func expandHttps(tfList []interface{}) *api.HTTPSConfig {
	https := &api.HTTPSConfig{}

	if len(tfList) == 0 || tfList[0] == nil {
		return https
	}
	tfMap := tfList[0].(map[string]interface{})

	https.Enabled = tfMap["enabled"].(bool)
	if https.Enabled {
		https.CertId = tfMap["cert_id"].(string)
		https.HttpRedirect = tfMap["http_redirect"].(bool)
		if https.HttpRedirect {
			https.HttpRedirectCode = tfMap["http_redirect_code"].(int)
		}
		https.HttpsRedirect = tfMap["https_redirect"].(bool)
		if https.HttpsRedirect {
			https.HttpsRedirectCode = tfMap["https_redirect_code"].(int)
		}
		https.Http2Enabled = tfMap["http2_enabled"].(bool)
		https.VerifyClient = tfMap["verify_client"].(bool)
		https.SslProtocols = flex.ExpandStringValueSet(tfMap["ssl_protocols"].(*schema.Set))
	}

	return https
}

//</editor-fold>

//<editor-fold desc="Origin">
func flattenClientIp(clientIp *api.ClientIp) []interface{} {
	tfMap := map[string]interface{}{
		"enabled": false,
	}
	if clientIp != nil && clientIp.Enabled {
		tfMap["enabled"] = true
		tfMap["name"] = clientIp.Name
	}
	return []interface{}{tfMap}
}

func expandClientIp(tfList []interface{}) *api.ClientIp {
	clientIp := &api.ClientIp{}
	if len(tfList) == 0 || tfList[0] == nil {
		return clientIp
	}
	tfMap := tfList[0].(map[string]interface{})
	clientIp.Enabled = tfMap["enabled"].(bool)
	if clientIp.Enabled {
		clientIp.Name = tfMap["name"].(string)
	}
	return clientIp
}

func flattenOriginProtocol(originProtocol string) []interface{} {
	return []interface{}{map[string]interface{}{
		"value": originProtocol,
	}}
}

func expandOriginProtocol(tfList []interface{}) string {
	if len(tfList) == 0 || tfList[0] == nil {
		return "http"
	}
	tfMap := tfList[0].(map[string]interface{})
	return tfMap["value"].(string)
}

//</editor-fold>
