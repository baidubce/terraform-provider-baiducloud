package baiducloud

import (
	"reflect"

	"github.com/baidubce/bce-sdk-go/services/bos"
	"github.com/baidubce/bce-sdk-go/services/bos/api"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	BOS_BUCKET_ACL_PRIVATE           = "private"
	BOS_BUCKET_ACL_PUBLIC_READ       = "public-read"
	BOS_BUCKET_ACL_PUBLIC_READ_WRITE = "public-read-write"

	BOS_BUCKET_OBJECT_CACHE_CONTROL_PRIVATE         = "private"
	BOS_BUCKET_OBJECT_CACHE_CONTROL_NO_CACHE        = "no-cache"
	BOS_BUCKET_OBJECT_CACHE_CONTROL_MAX_AGE         = "max-age"
	BOS_BUCKET_OBJECT_CACHE_CONTROL_MUST_REVALIDATE = "must-revalidate"

	BOS_BUCKET_OBJECT_CONTENT_DISPOSITION_INLINE     = "inline"
	BOS_BUCKET_OBJECT_CONTENT_DISPOSITION_ATTACHMENT = "attachment"
)

type BosService struct {
	client *connectivity.BaiduClient
}

func (s *BosService) ListAllObjects(bucket, prefix string) ([]api.ObjectSummaryType, error) {
	args := &api.ListObjectsArgs{
		Prefix: prefix,
	}
	action := "List All Objects for bucket " + bucket

	objects := make([]api.ObjectSummaryType, 0)
	for {
		raw, err := s.client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
			return bosClient.ListObjects(bucket, args)
		})
		if err != nil {
			return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket", action, BCESDKGoERROR)
		}

		result, _ := raw.(*api.ListObjectsResult)
		objects = append(objects, result.Contents...)
		if !result.IsTruncated {
			break
		}
		args.Marker = result.Marker
	}

	return objects, nil
}

func (s *BosService) resourceBaiduCloudBosBucketReadAcl(bucket string) (string, error) {
	action := "read bos bucket acl " + bucket

	raw, err := s.client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
		return bosClient.GetBucketAcl(bucket)
	})
	addDebug(action, raw)
	if err != nil {
		return "", err
	}

	result, _ := raw.(*api.GetBucketAclResult)
	aclResult := getAclByAccessControlList(result.AccessControlList)

	return aclResult, nil
}

func (s *BosService) resourceBaiduCloudBosBucketReadReplicationConfigure(bucket string) ([]map[string]interface{}, error) {
	action := "read bos bucket replication configuration " + bucket

	raw, err := s.client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
		return bosClient.GetBucketReplication(bucket, "")
	})
	if err != nil {
		if !IsExceptedErrors(err, ReplicationConfigurationNotFound) {
			return nil, err
		}
	}
	addDebug(action, raw)

	replicationConfiguration := make([]map[string]interface{}, 0)
	if getBucketReplicationResult, ok := raw.(*api.GetBucketReplicationResult); ok {
		replicationConfiguration = flattenBaiduCloudBucketReplicationConfiguration(getBucketReplicationResult)
	}

	return replicationConfiguration, nil
}

func (s *BosService) resourceBaiduCloudBosBucketReadLogging(bucket string) ([]map[string]interface{}, error) {
	action := "read bos bucket logging " + bucket

	raw, err := s.client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
		return bosClient.GetBucketLogging(bucket)
	})
	if err != nil {
		return nil, err
	}
	addDebug(action, raw)

	logging := make([]map[string]interface{}, 0, 1)
	if result, ok := raw.(*api.GetBucketLoggingResult); ok && result.TargetBucket != "" {
		l := make(map[string]interface{})
		l["target_bucket"] = result.TargetBucket
		l["target_prefix"] = result.TargetPrefix
		logging = append(logging, l)
	}

	return logging, nil
}

func (s *BosService) resourceBaiduCloudBosBucketReadLifecycle(bucket string) ([]interface{}, error) {
	action := "read bos bucket lifecycle " + bucket

	raw, err := s.client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
		return bosClient.GetBucketLifecycle(bucket)
	})
	if err != nil {
		if !IsExceptedErrors(err, []string{"NoLifecycleConfiguration"}) {
			return nil, err
		}
		raw = nil
	}
	addDebug(action, raw)

	rules := make([]interface{}, 0)
	if getBucketLifecycleResult, ok := raw.(*api.GetBucketLifecycleResult); ok {
		for _, lifecycleRule := range getBucketLifecycleResult.Rule {
			rule := make(map[string]interface{})

			rule["id"] = lifecycleRule.Id
			rule["status"] = lifecycleRule.Status
			rule["resource"] = lifecycleRule.Resource

			// condition
			conditions := make([]interface{}, 0, 1)
			condition := make(map[string]interface{})
			timeConditions := make([]interface{}, 0, 1)
			timeCondition := make(map[string]interface{})

			timeCondition["date_greater_than"] = lifecycleRule.Condition.Time.DateGreaterThan
			timeConditions = append(timeConditions, timeCondition)
			condition["time"] = timeConditions
			conditions = append(conditions, condition)
			rule["condition"] = conditions

			// action
			actions := make([]interface{}, 0, 1)
			action := make(map[string]interface{})
			action["name"] = lifecycleRule.Action.Name
			action["storage_class"] = lifecycleRule.Action.StorageClass
			actions = append(actions, action)
			rule["action"] = actions

			rules = append(rules, rule)
		}
	}

	return rules, nil
}

func (s *BosService) resourceBaiduCloudBosBucketReadWebsite(bucket string) ([]interface{}, error) {
	action := "read bos bucket website " + bucket

	raw, err := s.client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
		return bosClient.GetBucketStaticWebsite(bucket)
	})
	if err != nil {
		if !IsExceptedErrors(err, []string{"NoSuchBucketStaticWebSiteConfig"}) {
			return nil, err
		}
		raw = nil
	}
	addDebug(action, raw)

	website := make([]interface{}, 0, 1)
	if getBucketStaticWebsiteResult, ok := raw.(*api.GetBucketStaticWebsiteResult); ok {
		web := make(map[string]interface{})
		web["index_document"] = getBucketStaticWebsiteResult.Index
		web["error_document"] = getBucketStaticWebsiteResult.NotFound
		website = append(website, web)
	}

	return website, nil
}

func (s *BosService) resourceBaiduCloudBosBucketReadCors(bucket string) ([]interface{}, error) {
	action := "read bos bucket cors rules " + bucket

	raw, err := s.client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
		return bosClient.GetBucketCors(bucket)
	})
	if err != nil {
		if !IsExceptedErrors(err, []string{"NoSuchCORSConfiguration"}) {
			return nil, err
		}
		raw = nil
	}
	addDebug(action, raw)

	cors := make([]interface{}, 0)
	if getBucketCorsResult, ok := raw.(*api.GetBucketCorsResult); ok {
		for _, cor := range getBucketCorsResult.CorsConfiguration {
			corsMap := make(map[string]interface{})

			corsMap["allowed_headers"] = cor.AllowedHeaders
			corsMap["allowed_methods"] = cor.AllowedMethods
			corsMap["allowed_origins"] = cor.AllowedOrigins
			corsMap["allowed_expose_headers"] = cor.AllowedExposeHeaders
			corsMap["max_age_seconds"] = cor.MaxAgeSeconds

			cors = append(cors, corsMap)
		}
	}

	return cors, nil
}

func (s *BosService) resourceBaiduCloudBosBucketReadCopyright(bucket string) ([]interface{}, error) {
	action := "read bos bucket copyright protection " + bucket

	raw, err := s.client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
		return bosClient.GetBucketCopyrightProtection(bucket)
	})
	if err != nil {
		if !IsExceptedErrors(err, []string{"NoCopyrightProtectionConfiguration"}) {
			return nil, err
		}
		raw = nil
	}
	addDebug(action, raw)

	resource := make([]string, 0)
	if result, ok := raw.([]string); ok {
		resource = result
	}

	copyright := make([]interface{}, 0)
	if len(resource) != 0 {
		copyMap := map[string]interface{}{"resource": resource}
		copyright = append(copyright, copyMap)
	}

	return copyright, nil
}

func (s *BosService) resourceBaiduCloudBucketObjectReadAcl(bucket, key string) (string, error) {
	action := "read bos bucket object acl, bucket: " + bucket + ", key: " + key

	raw, err := s.client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
		return bosClient.GetObjectAcl(bucket, key)
	})
	addDebug(action, raw)
	if err != nil {
		return "", err
	}

	result, _ := raw.(*api.GetObjectAclResult)
	aclResult := getAclByAccessControlList(result.AccessControlList)

	return aclResult, nil
}

func getAclByAccessControlList(acList []api.GrantType) string {
	aclResult := BOS_BUCKET_ACL_PRIVATE

LOOPACL:
	for _, acl := range acList {
		for _, grantee := range acl.Grantee {
			if grantee.Id != "*" {
				continue
			}

			if reflect.DeepEqual(acl.Permission, []string{"READ", "WRITE"}) {
				aclResult = BOS_BUCKET_ACL_PUBLIC_READ_WRITE
			} else if reflect.DeepEqual(acl.Permission, []string{"READ"}) {
				aclResult = BOS_BUCKET_ACL_PUBLIC_READ
			}

			break LOOPACL
		}
	}

	return aclResult
}
