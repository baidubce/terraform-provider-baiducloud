package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/scs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"

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

func (s *ScsService) ResourceGroupBindStatusRefresh(instanceId string, resourceGroupId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		result, err := s.GetInstanceDetail(instanceId)
		if err != nil {
			return nil, "", WrapError(err)
		}
		if result.ResourceGroupId == resourceGroupId {
			return result, "true", nil
		}

		return result, "false", nil
	}
}

func (s *ScsService) checkScsTagsAndResourceGroupBind(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	instanceID := d.Id()
	action := "check SCS tags and resource group bind " + instanceID
	raw, err := client.WithScsClient(func(scsClient *scs.Client) (interface{}, error) {
		return scsClient.GetInstanceDetail(instanceID)
	})
	addDebug(action, raw)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs", action, BCESDKGoERROR)
	}
	response, _ := raw.(*scs.GetInstanceDetailResult)
	if _, ok := d.GetOk("tags"); ok {
		if response.Tags == nil || len(response.Tags) == 0 {
			// bind tags failed
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs", "tags bind", BCESDKGoERROR)
		}
	}
	scsService := ScsService{client}
	if v, ok := d.GetOk("resource_group_id"); ok {
		// scs的resource group 是在订单完成之后，这里刷新一分钟，如果没绑定成功则视为失败
		groupBindStatus := &resource.StateChangeConf{
			Delay:      10 * time.Second,
			Pending:    []string{"false"},
			Refresh:    scsService.ResourceGroupBindStatusRefresh(d.Id(), v.(string)),
			Target:     []string{"true"},
			Timeout:    1 * time.Minute,
			MinTimeout: 3 * time.Second,
		}
		if _, err := groupBindStatus.WaitForState(); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_scs", "resource group bind", BCESDKGoERROR)
		}
	}
	return nil
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

func (s *ScsService) GetSecurityGroups(instanceId string) ([]string, error) {
	result := make([]string, 0)

	action := "List all SCS instance security groups"
	raw, err := s.client.WithScsClient(func(scsClient *scs.Client) (interface{}, error) {
		return scsClient.ListSecurityGroupByInstanceId(instanceId)
	})
	if err != nil {
		return nil, err
	}
	addDebug(action, raw)

	response := raw.(*scs.ListSecurityGroupResult)
	for _, group := range response.Groups {
		result = append(result, group.SecurityGroupID)
	}
	addDebug(action, result)
	return result, nil
}
