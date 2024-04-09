package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/blb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strconv"
	"time"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

type BLBService struct {
	client *connectivity.BaiduClient
}

func (s *BLBService) BLBStateRefreshFunc(id string, failState []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		time.Sleep(time.Second * time.Duration(10))
		raw, err := s.client.WithBLBClient(func(client *blb.Client) (i interface{}, e error) {
			return client.DescribeLoadBalancerDetail(id)
		})
		if err != nil {
			if NotFoundError(err) {
				return nil, string(blb.BLBStatusCreating), nil
			}
			return nil, "", WrapError(err)
		}

		result := raw.(*blb.DescribeLoadBalancerDetailResult)
		for _, statue := range failState {
			if string(result.Status) == statue {
				return result, string(result.Status), WrapError(Error(GetFailTargetStatus, result.Status))
			}
		}

		return result, string(result.Status), nil
	}
}

func (s *BLBService) GetBLBDetail(blbId string) (*blb.BLBModel, *blb.DescribeLoadBalancerDetailResult, error) {
	action := "Describe BLB " + blbId + " Detail"

	raw, err := s.client.WithBLBClient(func(client *blb.Client) (i interface{}, e error) {
		return client.DescribeLoadBalancerDetail(blbId)
	})
	addDebug(action, raw)

	if err != nil {
		return nil, nil, WrapError(err)
	}
	blbDetail := raw.(*blb.DescribeLoadBalancerDetailResult)

	action = "List BLB " + blbId
	listArgs := &blb.DescribeLoadBalancersArgs{
		BlbId: blbDetail.BlbId,
	}
	raw, err = s.client.WithBLBClient(func(client *blb.Client) (i interface{}, e error) {
		return client.DescribeLoadBalancers(listArgs)
	})
	addDebug(action, err)

	if err != nil {
		return nil, nil, WrapError(err)
	}
	blbModel := &raw.(*blb.DescribeLoadBalancersResult).BlbList[0]

	return blbModel, blbDetail, nil
}

func (s *BLBService) ListAllBLB(args *blb.DescribeLoadBalancersArgs) ([]blb.BLBModel, map[string]blb.DescribeLoadBalancerDetailResult, error) {
	action := "List all BLB"

	blbModels := make([]blb.BLBModel, 0)
	for {
		raw, err := s.client.WithBLBClient(func(client *blb.Client) (i interface{}, e error) {
			return client.DescribeLoadBalancers(args)
		})

		if err != nil {
			return nil, nil, WrapError(err)
		}
		addDebug(action, raw)

		response := raw.(*blb.DescribeLoadBalancersResult)
		blbModels = append(blbModels, response.BlbList...)

		if response.IsTruncated {
			args.Marker = response.NextMarker
			args.MaxKeys = response.MaxKeys
		} else {
			break
		}
	}

	blbDetails := make(map[string]blb.DescribeLoadBalancerDetailResult)
	for _, model := range blbModels {
		raw, err := s.client.WithBLBClient(func(client *blb.Client) (i interface{}, e error) {
			return client.DescribeLoadBalancerDetail(model.BlbId)
		})
		if err != nil {
			return nil, nil, WrapError(err)
		}

		blbDetails[model.BlbId] = *raw.(*blb.DescribeLoadBalancerDetailResult)
	}

	return blbModels, blbDetails, nil
}

func (s *BLBService) FlattenListenerModelToMap(listeners []blb.ListenerModel) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(listeners))
	for _, l := range listeners {
		port, _ := strconv.Atoi(l.Port)
		result = append(result, map[string]interface{}{
			"port": port,
			"type": l.Type,
		})
	}

	return result
}

func (s *BLBService) FlattenBLBDetailsToMap(models []blb.BLBModel, details map[string]blb.DescribeLoadBalancerDetailResult) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(models))
	for _, model := range models {
		detail := details[model.BlbId]

		result = append(result, map[string]interface{}{
			"blb_id":      model.BlbId,
			"name":        model.Name,
			"description": model.Description,
			"address":     model.Address,
			"status":      model.Status,
			"vpc_id":      model.VpcId,
			"vpc_name":    detail.VpcName,
			"subnet_id":   model.SubnetId,
			"public_ip":   model.PublicIp,
			"cidr":        detail.Cidr,
			"create_time": detail.CreateTime,
			"listener":    s.FlattenListenerModelToMap(detail.Listener),
			"tags":        flattenTagsToMap(model.Tags),
		})
	}

	return result
}

// updateBlbSecurityGroups 更新指定BLB实例的安全组
//
// 参数：
// s *BLBService - BLB服务的指针
// blbId string - BLB实例ID
// securityGroups []string - 要更新的安全组列表
// meta interface{} - 额外的元数据参数
//
// 返回值：
// error - 错误信息，如果成功则为nil
func (s *BLBService) updateBlbSecurityGroups(d *schema.ResourceData, meta interface{}) error {
	return s.updateBlbSecurityGroupsGeneric(d, meta, false)
}

// updateBlbEnterpriseSecurityGroups 更新指定BLB实例的企业安全组
//
// 参数：
// s *BLBService - BLB服务的指针
// blbId string - BLB实例ID
// securityGroups []string - 要更新的企业安全组列表
// meta interface{} -
//
// 返回值：
// error - 错误信息，如果成功则为nil
func (s *BLBService) updateBlbEnterpriseSecurityGroups(d *schema.ResourceData, meta interface{}) error {
	return s.updateBlbSecurityGroupsGeneric(d, meta, true)
}

func (s *BLBService) updateBlbSecurityGroupsGeneric(d *schema.ResourceData, meta interface{}, isEnterprise bool) error {
	var action string
	var err error

	if isEnterprise {
		action = "Update BLB Enterprise Security Groups"
		err := updateEnterpriseSecurityGroups("enterprise_security_groups", d, d.Id(), s)
		if err != nil {
			return err
		}
	} else {
		action = "Update BLB Security Groups"
		err := updateSecurityGroups("security_groups", d, d.Id(), s)
		if err != nil {
			return err
		}
	}

	addDebug(action, d.Id())

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb", action, BCESDKGoERROR)
	}
	return nil
}

// getBlbSecurityGroupIds 从BLB服务中获取指定BLB实例的安全组ID列表
//
// 参数：
// s *BLBService - BLB服务的指针
// blbId string - BLB实例ID
// meta interface{} - 额外的元数据参数
//
// 返回值：
// []string - 安全组ID列表
// error - 错误信息，如果成功则为nil
func (s *BLBService) getBlbSecurityGroupIds(blbId string, meta interface{}) ([]string, error) {
	return s.getBlbSecurityGroupIdsGeneric(blbId, meta, false)
}

// getBlbEnterpriseSecurityGroupIds 从BLB服务中获取指定BLB实例的企业安全组ID列表
//
// 参数：
// s *BLBService - BLB服务的指针
// blbId string - BLB实例ID
// meta interface{}
//
// 返回值：
// []string - 企业安全组ID列表
// error - 错误信息，如果成功则为nil
func (s *BLBService) getBlbEnterpriseSecurityGroupIds(blbId string, meta interface{}) ([]string, error) {
	return s.getBlbSecurityGroupIdsGeneric(blbId, meta, true)
}

func (s *BLBService) getBlbSecurityGroupIdsGeneric(blbId string, meta interface{}, isEnterprise bool) ([]string, error) {
	client := meta.(*connectivity.BaiduClient)
	var action string
	var raw interface{}
	var err error

	if isEnterprise {
		action = "Get BLB Enterprise Security Group ids"
		raw, err = client.WithBLBClient(func(client *blb.Client) (i interface{}, e error) {
			return client.DescribeEnterpriseSecurityGroups(blbId)
		})
	} else {
		action = "Get BLB Security Group ids"
		raw, err = client.WithBLBClient(func(client *blb.Client) (i interface{}, e error) {
			return client.DescribeSecurityGroups(blbId)
		})
	}

	addDebug(action, blbId)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb", action, BCESDKGoERROR)
	}

	var ids []string
	if isEnterprise {
		result := raw.(*blb.DescribeEnterpriseSecurityGroupsResult)
		for _, item := range result.BlbEnterpriseSecurityGroups {
			ids = append(ids, item.EnterpriseSecurityGroupId)
		}
	} else {
		result := raw.(*blb.DescribeSecurityGroupsResult)
		for _, item := range result.BlbSecurityGroups {
			ids = append(ids, item.SecurityGroupId)
		}
	}

	return ids, nil
}

// AddSecurityGroups implements the method to add security groups to an instance.
func (s *BLBService) AddSecurityGroups(instanceID string, securityGroupIDs []string) error {
	args := &blb.UpdateSecurityGroupsArgs{
		SecurityGroupIds: securityGroupIDs,
	}
	_, err := s.client.WithBLBClient(func(blbClient *blb.Client) (i interface{}, e error) {
		return nil, blbClient.BindSecurityGroups(instanceID, args)
	})
	return err
}

// RemoveSecurityGroups implements the method to remove security groups from an instance.
func (s *BLBService) RemoveSecurityGroups(instanceID string, securityGroupIDs []string) error {
	args := &blb.UpdateSecurityGroupsArgs{
		SecurityGroupIds: securityGroupIDs,
	}
	_, err := s.client.WithBLBClient(func(blbClient *blb.Client) (i interface{}, e error) {
		return nil, blbClient.UnbindSecurityGroups(instanceID, args)
	})
	return err
}

// AddEnterpriseSecurityGroups implements the method to add enterprise security groups to an instance.
func (s *BLBService) AddEnterpriseSecurityGroups(instanceID string, securityGroupIDs []string) error {
	args := &blb.UpdateEnterpriseSecurityGroupsArgs{
		EnterpriseSecurityGroupIds: securityGroupIDs,
	}
	_, err := s.client.WithBLBClient(func(blbClient *blb.Client) (i interface{}, e error) {
		return nil, blbClient.BindEnterpriseSecurityGroups(instanceID, args)
	})
	return err
}

// RemoveEnterpriseSecurityGroups implements the method to remove enterprise security groups from an instance.
func (s *BLBService) RemoveEnterpriseSecurityGroups(instanceID string, securityGroupIDs []string) error {
	args := &blb.UpdateEnterpriseSecurityGroupsArgs{
		EnterpriseSecurityGroupIds: securityGroupIDs,
	}
	_, err := s.client.WithBLBClient(func(blbClient *blb.Client) (i interface{}, e error) {
		return nil, blbClient.UnbindEnterpriseSecurityGroups(instanceID, args)
	})
	return err
}
