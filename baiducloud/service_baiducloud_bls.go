package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/bls"
	"github.com/baidubce/bce-sdk-go/services/bls/api"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

type BLSService struct {
	client *connectivity.BaiduClient
}

func (s *BLSService) GetBLSLogStoreList(args *api.QueryConditions) (*api.ListLogStoreResult, error) {
	action := "List BLS LogStroe "

	raw, err := s.client.WithBLSClient(func(blsClient *bls.Client) (i interface{}, e error) {
		return blsClient.ListLogStore(args)
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_bls_log_store", action, BCESDKGoERROR)
	}

	return raw.(*api.ListLogStoreResult), nil
}

func (s *BLSService) GetBLSLogStoreDetail(logStoreName string) (*api.LogStore, error) {
	action := "Query BLS LogStore " + logStoreName

	raw, err := s.client.WithBLSClient(func(blsClient *bls.Client) (i interface{}, e error) {
		return blsClient.DescribeLogStore(logStoreName)
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_bls_log_store", action, BCESDKGoERROR)
	}

	return raw.(*api.LogStore), nil
}

func (e *BLSService) FlattenLogStoreModelsToMap(stores []api.LogStore) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(stores))
	for _, e := range stores {
		result = append(result, map[string]interface{}{
			"creation_date_time": e.CreationDateTime,
			"last_modified_time": e.LastModifiedTime,
			"log_store_name":     e.LogStoreName,
			"retention":          e.Retention,
		})
	}

	return result
}
