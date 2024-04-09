package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/appblb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strconv"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

type APPBLBService struct {
	client *connectivity.BaiduClient
}

func (s *APPBLBService) APPBLBStateRefreshFunc(id string, failState []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		raw, err := s.client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
			return client.DescribeLoadBalancerDetail(id)
		})
		if err != nil {
			if NotFoundError(err) {
				return nil, string(appblb.BLBStatusCreating), nil
			}
			return nil, "", WrapError(err)
		}

		result := raw.(*appblb.DescribeLoadBalancerDetailResult)
		for _, statue := range failState {
			if string(result.Status) == statue {
				return result, string(result.Status), WrapError(Error(GetFailTargetStatus, result.Status))
			}
		}

		return result, string(result.Status), nil
	}
}

func (s *APPBLBService) GetAppBLBDetail(blbId string) (*appblb.AppBLBModel, *appblb.DescribeLoadBalancerDetailResult, error) {
	action := "Describe AppBLB " + blbId + " Detail"

	raw, err := s.client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
		return client.DescribeLoadBalancerDetail(blbId)
	})
	addDebug(action, raw)

	if err != nil {
		return nil, nil, WrapError(err)
	}
	blbDetail := raw.(*appblb.DescribeLoadBalancerDetailResult)

	action = "List AppBLB " + blbId
	listArgs := &appblb.DescribeLoadBalancersArgs{
		BlbId: blbDetail.BlbId,
	}
	raw, err = s.client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
		return client.DescribeLoadBalancers(listArgs)
	})
	addDebug(action, err)

	if err != nil {
		return nil, nil, WrapError(err)
	}
	blbModel := &raw.(*appblb.DescribeLoadBalancersResult).BlbList[0]

	return blbModel, blbDetail, nil
}

func (s *APPBLBService) ListAllAppBLB(args *appblb.DescribeLoadBalancersArgs) ([]appblb.AppBLBModel, map[string]appblb.DescribeLoadBalancerDetailResult, error) {
	action := "List all APPBLB"

	appblbModels := make([]appblb.AppBLBModel, 0)
	for {
		raw, err := s.client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
			return client.DescribeLoadBalancers(args)
		})

		if err != nil {
			return nil, nil, WrapError(err)
		}
		addDebug(action, raw)

		response := raw.(*appblb.DescribeLoadBalancersResult)
		appblbModels = append(appblbModels, response.BlbList...)

		if response.IsTruncated {
			args.Marker = response.NextMarker
			args.MaxKeys = response.MaxKeys
		} else {
			break
		}
	}

	appblbDetails := make(map[string]appblb.DescribeLoadBalancerDetailResult)
	for _, model := range appblbModels {
		raw, err := s.client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
			return client.DescribeLoadBalancerDetail(model.BlbId)
		})
		if err != nil {
			return nil, nil, WrapError(err)
		}

		appblbDetails[model.BlbId] = *raw.(*appblb.DescribeLoadBalancerDetailResult)
	}

	return appblbModels, appblbDetails, nil
}

func (s *APPBLBService) FlattenListenerModelToMap(listeners []appblb.ListenerModel) []map[string]interface{} {
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

func (s *APPBLBService) FlattenAppBLBDetailsToMap(models []appblb.AppBLBModel, details map[string]appblb.DescribeLoadBalancerDetailResult) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(models))
	for _, model := range models {
		detail := details[model.BlbId]

		result = append(result, map[string]interface{}{
			"blb_id":       model.BlbId,
			"name":         model.Name,
			"description":  model.Description,
			"address":      model.Address,
			"status":       model.Status,
			"vpc_id":       model.VpcId,
			"vpc_name":     detail.VpcName,
			"subnet_id":    model.SubnetId,
			"subnet_name":  detail.SubnetName,
			"subnet_cidr":  detail.SubnetCider,
			"public_ip":    model.PublicIp,
			"cidr":         detail.Cidr,
			"create_time":  detail.CreateTime,
			"release_time": detail.ReleaseTime,
			"listener":     s.FlattenListenerModelToMap(detail.Listener),
			"tags":         flattenTagsToMap(model.Tags),
		})
	}

	return result
}

func (s *APPBLBService) updateAppBlbSecurityGroups(d *schema.ResourceData, meta interface{}) error {
	return s.updateAppBlbSecurityGroupsGeneric(d, meta, false)
}

func (s *APPBLBService) updateAppBlbEnterpriseSecurityGroups(d *schema.ResourceData, meta interface{}) error {
	return s.updateAppBlbSecurityGroupsGeneric(d, meta, true)
}

func (s *APPBLBService) updateAppBlbSecurityGroupsGeneric(d *schema.ResourceData, meta interface{}, isEnterprise bool) error {
	var action string
	var err error

	if isEnterprise {
		action = "Update APPBLB Enterprise Security Groups"
		err := updateEnterpriseSecurityGroups("enterprise_security_groups", d, d.Id(), s)
		if err != nil {
			return err
		}
	} else {
		action = "Update APPBLB Security Groups"
		err := updateSecurityGroups("security_groups", d, d.Id(), s)
		if err != nil {
			return err
		}
	}

	addDebug(action, d.Id())

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appblb", action, BCESDKGoERROR)
	}
	return nil
}

// getAppBlbSecurityGroupIds 从APPBLB服务中获取指定APPBLB实例的安全组ID列表
//
// 参数：
// s *APPBLBService - APPBLB服务的指针
// blbId string - APPBLB实例ID
// meta interface{} - 额外的元数据参数
//
// 返回值：
// []string - 安全组ID列表
// error - 错误信息，如果成功则为nil
func (s *APPBLBService) getAppBlbSecurityGroupIds(blbId string, meta interface{}) ([]string, error) {
	return s.getAppBlbSecurityGroupIdsGeneric(blbId, meta, false)
}

// getAppBlbEnterpriseSecurityGroupIds 从APPBLB服务中获取指定APPBLB实例的企业安全组ID列表
//
// 参数：
// s *APPBLBService - APPBLB服务的指针
// blbId string - APPBLB实例ID
// meta interface{}
//
// 返回值：
// []string - 企业安全组ID列表
// error - 错误信息，如果成功则为nil
func (s *APPBLBService) getAppBlbEnterpriseSecurityGroupIds(blbId string, meta interface{}) ([]string, error) {
	return s.getAppBlbSecurityGroupIdsGeneric(blbId, meta, true)
}

func (s *APPBLBService) getAppBlbSecurityGroupIdsGeneric(blbId string, meta interface{}, isEnterprise bool) ([]string, error) {
	client := meta.(*connectivity.BaiduClient)
	var action string
	var raw interface{}
	var err error

	if isEnterprise {
		action = "Get App BLB Enterprise Security Group ids"
		raw, err = client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
			return client.DescribeEnterpriseSecurityGroups(blbId)
		})
	} else {
		action = "Get App BLB Security Group ids"
		raw, err = client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
			return client.DescribeSecurityGroups(blbId)
		})
	}

	addDebug(action, blbId)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_appblb", action, BCESDKGoERROR)
	}

	var ids []string
	if isEnterprise {
		result := raw.(*appblb.DescribeEnterpriseSecurityGroupsResult)
		for _, item := range result.BlbEnterpriseSecurityGroups {
			ids = append(ids, item.EnterpriseSecurityGroupId)
		}
	} else {
		result := raw.(*appblb.DescribeSecurityGroupsResult)
		for _, item := range result.BlbSecurityGroups {
			ids = append(ids, item.SecurityGroupId)
		}
	}

	return ids, nil
}

// AddSecurityGroups implements the method to add security groups to an instance.
func (s *APPBLBService) AddSecurityGroups(instanceID string, securityGroupIDs []string) error {
	args := &appblb.UpdateSecurityGroupsArgs{
		SecurityGroupIds: securityGroupIDs,
	}
	_, err := s.client.WithAppBLBClient(func(appblbClient *appblb.Client) (i interface{}, e error) {
		return nil, appblbClient.BindSecurityGroups(instanceID, args)
	})
	return err
}

// RemoveSecurityGroups implements the method to remove security groups from an instance.
func (s *APPBLBService) RemoveSecurityGroups(instanceID string, securityGroupIDs []string) error {
	args := &appblb.UpdateSecurityGroupsArgs{
		SecurityGroupIds: securityGroupIDs,
	}
	_, err := s.client.WithAppBLBClient(func(appblbClient *appblb.Client) (i interface{}, e error) {
		return nil, appblbClient.UnbindSecurityGroups(instanceID, args)
	})
	return err
}

// AddEnterpriseSecurityGroups implements the method to add enterprise security groups to an instance.
func (s *APPBLBService) AddEnterpriseSecurityGroups(instanceID string, securityGroupIDs []string) error {
	args := &appblb.UpdateEnterpriseSecurityGroupsArgs{
		EnterpriseSecurityGroupIds: securityGroupIDs,
	}
	_, err := s.client.WithAppBLBClient(func(appblbClient *appblb.Client) (i interface{}, e error) {
		return nil, appblbClient.BindEnterpriseSecurityGroups(instanceID, args)
	})
	return err
}

// RemoveEnterpriseSecurityGroups implements the method to remove enterprise security groups from an instance.
func (s *APPBLBService) RemoveEnterpriseSecurityGroups(instanceID string, securityGroupIDs []string) error {
	args := &appblb.UpdateEnterpriseSecurityGroupsArgs{
		EnterpriseSecurityGroupIds: securityGroupIDs,
	}
	_, err := s.client.WithAppBLBClient(func(appblbClient *appblb.Client) (i interface{}, e error) {
		return nil, appblbClient.UnbindEnterpriseSecurityGroups(instanceID, args)
	})
	return err
}
