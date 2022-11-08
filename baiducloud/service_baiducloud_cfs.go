package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/cfs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	CfsMountTargetDeleting = "CfsMountTargetDeleting"
	CfsMountTargetDeleted  = "CfsMountTargetDeleted"
)

type CfsService struct {
	Client *connectivity.BaiduClient
}

func (s *CfsService) cfsStateRefresh(fsId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		action := "Refresh CFS State " + fsId
		raw, err := s.Client.WithCfsClient(func(cfsClient *cfs.Client) (i interface{}, e error) {
			return cfsClient.DescribeFS(&cfs.DescribeFSArgs{
				FSID: fsId,
			})
		})
		addDebug(action, raw)
		if err != nil {
			if NotFoundError(err) {
				return 0, string(cfs.FSStatusUnavailable), nil
			}
			return nil, "", WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfs", action, BCESDKGoERROR)
		}
		if raw == nil {
			return nil, "", nil
		}
		addDebug(action, raw)
		result, _ := raw.(*cfs.DescribeFSResult)
		return result, string(result.FSList[0].Status), nil
	}
}

func (s *CfsService) cfsMountTargetCountRefresh(fsId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		action := "Refresh CFS State " + fsId
		raw, err := s.Client.WithCfsClient(func(cfsClient *cfs.Client) (i interface{}, e error) {
			return cfsClient.DescribeFS(&cfs.DescribeFSArgs{
				FSID: fsId,
			})
		})
		addDebug(action, raw)
		if err != nil {
			if NotFoundError(err) {
				return 0, string(cfs.FSStatusUnavailable), nil
			}
			return nil, "", WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfs", action, BCESDKGoERROR)
		}
		if raw == nil {
			return nil, "", nil
		}
		addDebug(action, raw)
		result, _ := raw.(*cfs.DescribeFSResult)
		if len(result.FSList[0].MoutTargets) == 0 {
			return result, CfsMountTargetDeleted, nil
		}
		return result, CfsMountTargetDeleting, nil
	}
}

func (s *CfsService) GetCfsDetail(fsId string) (model *cfs.FSModel, e error) {
	action := "Query CFS detail " + fsId
	raw, err := s.Client.WithCfsClient(func(client *cfs.Client) (i interface{}, e error) {
		return client.DescribeFS(&cfs.DescribeFSArgs{
			FSID: fsId,
		})
	})

	if err != nil {
		return nil, err
	}

	addDebug(action, raw)
	response, _ := raw.(*cfs.DescribeFSResult)
	return &response.FSList[0], nil
}

func (s *CfsService) ListCfs() (model []cfs.FSModel, e error) {
	action := "Query CFS detail "
	args := &cfs.DescribeFSArgs{}
	result := make([]cfs.FSModel, 0)
	for {
		raw, err := s.Client.WithCfsClient(func(client *cfs.Client) (i interface{}, e error) {
			return client.DescribeFS(args)
		})

		if err != nil {
			return nil, err
		}

		addDebug(action, raw)
		response, _ := raw.(*cfs.DescribeFSResult)
		result = append(result, response.FSList...)
		if response.IsTruncated {
			args.Marker = response.Marker
			args.MaxKeys = response.MaxKeys
		} else {
			break
		}
	}
	return result, nil
}

func (s *CfsService) ListCfsMountTarget(fsId string) (model []cfs.MoutTargetModel, e error) {
	action := "Query CFS Mount Target List "
	args := &cfs.DescribeMountTargetArgs{}
	args.FSID = fsId
	result := make([]cfs.MoutTargetModel, 0)
	for {
		raw, err := s.Client.WithCfsClient(func(client *cfs.Client) (i interface{}, e error) {
			return client.DescribeMountTarget(args)
		})

		if err != nil {
			return nil, err
		}

		addDebug(action, raw)
		response, _ := raw.(*cfs.DescribeMountTargetResult)
		result = append(result, response.MountTargetList...)
		if response.IsTruncated {
			args.Marker = response.Marker
			args.MaxKeys = response.MaxKeys
		} else {
			break
		}
	}
	return result, nil
}
