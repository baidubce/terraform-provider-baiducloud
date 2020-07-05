package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/cce"
	"github.com/hashicorp/terraform/helper/resource"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

type CceService struct {
	client *connectivity.BaiduClient
}

func (s *CceService) ClusterStateRefresh(clusterUuid string, failState []cce.ClusterStatus) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		action := "Query CCE Cluster " + clusterUuid
		raw, err := s.client.WithCCEClient(func(cceClient *cce.Client) (i interface{}, e error) {
			return cceClient.GetCluster(clusterUuid)
		})
		addDebug(action, raw)
		if err != nil {
			if NotFoundError(err) {
				return 0, string(cce.ClusterStatusDeleted), nil
			}
			return nil, "", WrapErrorf(err, DefaultErrorMsg, "baiducloud_cce_cluster", action, BCESDKGoERROR)
		}
		if raw == nil {
			return nil, "", nil
		}

		result := raw.(*cce.GetClusterResult)
		for _, statue := range failState {
			if result.Status == statue {
				return result, string(result.Status), WrapError(Error(GetFailTargetStatus, result.Status))
			}
		}

		addDebug(action, raw)
		return result, string(result.Status), nil
	}
}
