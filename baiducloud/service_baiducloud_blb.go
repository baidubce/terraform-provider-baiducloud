package baiducloud

import (
	"strconv"

	"github.com/baidubce/bce-sdk-go/services/blb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

type BLBService struct {
	client *connectivity.BaiduClient
}

func (s *BLBService) BLBStateRefreshFunc(id string, failState []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		raw, err := s.client.WithBLBClient(func(client *blb.Client) (i interface{}, e error) {
			return client.DescribeLoadBalancerDetail(id)
		})
		if err != nil {
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
