package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/scs"
	"github.com/hashicorp/terraform/helper/resource"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

type ScsService struct {
	client *connectivity.BaiduClient
}

func (s *ScsService) ListAllInstances(args *scs.ListInstancesArgs) ([]scs.InstanceModel, error) {
	result := make([]scs.InstanceModel, 0)

	action := "List all SCS instance "
	for {
		raw, err := s.client.WithScsClient(func(scsClient *scs.Client) (interface{}, error) {
			return scsClient.ListInstances(args)
		})
		if err != nil {
			return nil, err
		}
		addDebug(action, raw)

		response := raw.(*scs.ListInstancesResult)
		result = append(result, response.Instances...)

		if response.IsTruncated {
			args.Marker = response.NextMarker
			args.MaxKeys = response.MaxKeys
		} else {
			return result, nil
		}
	}
}

func (s *ScsService) InstanceStateRefresh(instanceId string, failState []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		result, err := s.GetInstanceDetail(instanceId)
		if err != nil {
			return nil, "", WrapError(err)
		}

		for _, statue := range failState {
			if result.InstanceStatus == statue {
				return result, result.InstanceStatus, WrapError(Error(GetFailTargetStatus, result.InstanceStatus))
			}
		}

		return result, result.InstanceStatus, nil
	}
}

func (s *ScsService) GetInstanceDetail(instanceID string) (*scs.GetInstanceDetailResult, error) {
	action := "Get SCS instance detail " + instanceID
	raw, err := s.client.WithScsClient(func(scsClient *scs.Client) (interface{}, error) {
		return scsClient.GetInstanceDetail(instanceID)
	})
	addDebug(action, raw)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs", action, BCESDKGoERROR)
	}

	result, _ := raw.(*scs.GetInstanceDetailResult)
	return result, nil
}

func (s *ScsService) GetNodeTypeList() (*scs.GetNodeTypeListResult, error) {
	action := "Get SCS nodetype list "
	raw, err := s.client.WithScsClient(func(scsClient *scs.Client) (interface{}, error) {
		return scsClient.GetNodeTypeList()
	})
	addDebug(action, raw)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs", action, BCESDKGoERROR)
	}

	result, _ := raw.(*scs.GetNodeTypeListResult)
	return result, nil
}

func (e *ScsService) FlattenScsModelsToMap(scss []scs.InstanceModel) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(scss))

	for _, e := range scss {
		result = append(result, map[string]interface{}{
			"instance_id":     e.InstanceID,
			"instance_name":   e.InstanceName,
			"instance_status": e.InstanceStatus,
			"cluster_type":    e.ClusterType,
			"engine":          e.Engine,
			"engine_version":  e.EngineVersion,
			"v_net_ip":        e.VnetIP,
			"domain":          e.Domain,
			"port":            e.Port,
			"create_time":     e.InstanceCreateTime,
			"capacity":        e.Capacity,
			"used_capacity":   e.UsedCapacity,
			"payment_timing":  e.PaymentTiming,
			"zone_names":      e.ZoneNames,
			"tags":            flattenTagsToMap(e.Tags),
		})
	}
	return result
}
