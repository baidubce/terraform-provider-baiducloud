package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/rds"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

type RdsService struct {
	client *connectivity.BaiduClient
}

func (s *RdsService) ListAllInstances(args *rds.ListRdsArgs) ([]rds.Instance, error) {
	result := make([]rds.Instance, 0)

	action := "List all RDS instance "
	for {
		raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.ListRds(args)
		})
		if err != nil {
			return nil, err
		}
		addDebug(action, raw)

		response := raw.(*rds.ListRdsResult)
		result = append(result, response.Instances...)

		if response.IsTruncated {
			args.Marker = response.NextMarker
			args.MaxKeys = response.MaxKeys
		} else {
			return result, nil
		}
	}
}

func (s *RdsService) InstanceStateRefresh(instanceId string, failState []string) resource.StateRefreshFunc {
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

func (s *RdsService) GetInstanceDetail(instanceID string) (*rds.Instance, error) {
	action := "Get RDS instance detail " + instanceID
	raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.GetDetail(instanceID)
	})
	addDebug(action, raw)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_rds", action, BCESDKGoERROR)
	}

	result, _ := raw.(*rds.Instance)
	return result, nil
}

func (e *RdsService) FlattenRdsModelsToMap(rdss []rds.Instance) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(rdss))

	for _, e := range rdss {
		result = append(result, map[string]interface{}{
			"instance_id":        e.InstanceId,
			"instance_name":      e.InstanceName,
			"instance_status":    e.InstanceStatus,
			"source_instance_id": e.SourceInstanceId,
			"source_region":      e.SourceRegion,
			"engine":             e.Engine,
			"engine_version":     e.EngineVersion,
			"category":           e.Category,
			"cpu_count":          e.CpuCount,
			"memory_capacity":    e.MemoryCapacity,
			"volume_capacity":    e.VolumeCapacity,
			"node_amount":        e.NodeAmount,
			"used_storage":       e.UsedStorage,
			"create_time":        e.InstanceCreateTime,
			"expire_time":        e.InstanceExpireTime,
			"address":            e.Endpoint.Address,
			"port":               e.Endpoint.Port,
			"v_net_ip":           e.Endpoint.VnetIp,
			"region":             e.Region,
			"instance_type":      e.InstanceType,
			"payment_timing":     e.PaymentTiming,
			"zone_names":         e.ZoneNames,
		})
	}
	return result
}
