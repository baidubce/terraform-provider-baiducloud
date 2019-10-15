package baiducloud

import (
	"strconv"

	"github.com/baidubce/bce-sdk-go/services/appblb"
	"github.com/hashicorp/terraform/helper/resource"

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
