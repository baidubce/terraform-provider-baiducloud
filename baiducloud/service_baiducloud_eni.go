package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/baidubce/bce-sdk-go/services/eni"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

type EniService struct {
	client *connectivity.BaiduClient
}

func (s *EniService) eniStateRefresh(eniId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		action := "Refresh Eni State " + eniId
		raw, err := s.client.WithEniClient(func(eniClient *eni.Client) (i interface{}, e error) {
			return eniClient.GetEniDetail(eniId)
		})
		addDebug(action, raw)
		if err != nil {
			if NotFoundError(err) {
				return 0, string(api.InstanceStatusDeleted), nil
			}
			return nil, "", WrapErrorf(err, DefaultErrorMsg, "baiducloud_eni", action, BCESDKGoERROR)
		}
		if raw == nil {
			return nil, "", nil
		}
		addDebug(action, raw)
		result, _ := raw.(*eni.Eni)
		return result, result.Status, nil
	}
}

func (s *EniService) ListEnis(args *eni.ListEniArgs) ([]eni.Eni, error) {
	action := "List " + args.VpcId + " enis"

	result := make([]eni.Eni, 0)
	for {
		raw, err := s.client.WithEniClient(func(eniClient *eni.Client) (i interface{}, e error) {
			return eniClient.ListEni(args)
		})
		if err != nil {
			return nil, WrapError(err)
		}
		addDebug(action, raw)

		response := raw.(*eni.ListEniResult)
		result = append(result, response.Eni...)

		if response.IsTruncated {
			args.Marker = response.NextMarker
			args.MaxKeys = response.MaxKeys
		} else {
			return result, nil
		}
	}
}

func (s *EniService) eniToMap(item eni.Eni) map[string]interface{} {
	privateIpSetMap := make([]map[string]interface{}, 0)
	for _, privateIp := range item.PrivateIpSet {
		privateIpMap := map[string]interface{}{
			"public_ip_address":  privateIp.PublicIpAddress,
			"primary":            privateIp.Primary,
			"private_ip_address": privateIp.PrivateIpAddress,
		}
		privateIpSetMap = append(privateIpSetMap, privateIpMap)
	}
	res := map[string]interface{}{
		"eni_id":                        item.EniId,
		"name":                          item.Name,
		"zone_name":                     item.ZoneName,
		"description":                   item.Description,
		"instance_id":                   item.InstanceId,
		"mac_address":                   item.MacAddress,
		"subnet_id":                     item.SubnetId,
		"status":                        item.Status,
		"private_ip_set":                privateIpSetMap,
		"security_group_ids":            item.SecurityGroupIds,
		"enterprise_security_group_ids": item.EnterpriseSecurityGroupIds,
		"created_time":                  item.CreatedTime,
	}
	return res
}
