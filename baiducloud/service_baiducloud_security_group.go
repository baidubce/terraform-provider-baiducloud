package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
)

func (s *BccService) FlattenSecurityGroupModelToMap(sgList []api.SecurityGroupModel) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(sgList))

	for _, sg := range sgList {
		result = append(result, map[string]interface{}{
			"id":          sg.Id,
			"name":        sg.Name,
			"vpc_id":      sg.VpcId,
			"description": sg.Desc,
			"tags":        flattenTagsToMap(sg.Tags),
		})
	}

	return result
}

func (s *BccService) ListAllSecurityGroups(args *api.ListSecurityGroupArgs) ([]api.SecurityGroupModel, error) {
	result := make([]api.SecurityGroupModel, 0)

	action := "Query all Security Groups"
	for {
		raw, err := s.client.WithBccClient(func(bccClient *bcc.Client) (interface{}, error) {
			return bccClient.ListSecurityGroup(args)
		})
		addDebug(action, raw)

		if err != nil {
			return nil, WrapError(err)
		}
		response, _ := raw.(*api.ListSecurityGroupResult)
		result = append(result, response.SecurityGroups...)

		if response.IsTruncated {
			args.Marker = response.Marker
			args.MaxKeys = response.MaxKeys
		} else {
			return result, nil
		}
	}
}
