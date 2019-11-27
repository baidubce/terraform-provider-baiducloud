package baiducloud

import "github.com/baidubce/bce-sdk-go/services/cfc/api"

var BosTriggerEventTypeList = []string{
	string(api.BosEventTypeAppendObject),
	string(api.BosEventTypeCompleteMultipartObject),
	string(api.BosEventTypeCopyObject),
	string(api.BosEventTypePostObject),
	string(api.BosEventTypePutObject),
}

var CDNTriggerEventTypeList = []string{
	string(api.CDNEventTypeCachedObjectsPushed),
	string(api.CDNEventTypeCachedObjectsBlocked),
	string(api.CDNEventTypeCachedObjectsRefreshed),
	string(api.CDNEventTypeCdnDomainCreated),
	string(api.CDNEventTypeCdnDomainDeleted),
	string(api.CDNEventTypeLogFileCreated),
	string(api.CDNEventTypeCdnDomainStarted),
	string(api.CDNEventTypeCdnDomainStopped),
}

var SourceTypeMap = map[string]api.SourceType{
	"bos":     api.SourceType("bos"),
	"http":    api.SourceTypeHTTP,
	"crontab": api.SourceTypeCrontab,
	"dueros":  api.SourceTypeDuerOS,
	"duedge":  api.SourceTypeDuEdge,
	"cdn":     api.SourceTypeCDN,
}
