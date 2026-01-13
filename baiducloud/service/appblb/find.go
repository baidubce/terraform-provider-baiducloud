package appblb

import (
	"errors"
	"fmt"

	"github.com/baidubce/bce-sdk-go/services/appblb"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

var errIpGroupNotFound = errors.New("appblb ip group not found")

func findIpGroup(conn *connectivity.BaiduClient, blbId, ipGroupId string) (*appblb.AppIpGroup, error) {
	args := &appblb.DescribeAppIpGroupArgs{
		MaxKeys: 1000,
	}

	for {
		raw, err := conn.WithAppBLBClient(func(client *appblb.Client) (interface{}, error) {
			return client.DescribeAppIpGroup(blbId, args)
		})
		if err != nil {
			return nil, err
		}

		result := raw.(*appblb.DescribeAppIpGroupResult)
		for _, group := range result.AppIpGroupList {
			if group.Id == ipGroupId {
				return &group, nil
			}
		}

		if result.IsTruncated {
			args.Marker = result.NextMarker
			args.MaxKeys = result.MaxKeys
		} else {
			break
		}
	}

	return nil, errIpGroupNotFound
}

func listIpGroupMembers(conn *connectivity.BaiduClient, blbId, ipGroupId string) ([]appblb.AppIpGroupMember, error) {
	args := &appblb.DescribeAppIpGroupMemberArgs{
		IpGroupId: ipGroupId,
		MaxKeys:   1000,
	}

	result := make([]appblb.AppIpGroupMember, 0)
	for {
		raw, err := conn.WithAppBLBClient(func(client *appblb.Client) (interface{}, error) {
			return client.DescribeAppIpGroupMember(blbId, args)
		})
		if err != nil {
			return nil, fmt.Errorf("error reading appblb ip group members: %w", err)
		}

		response := raw.(*appblb.DescribeAppIpGroupMemberResult)
		result = append(result, response.MemberList...)

		if response.IsTruncated {
			args.Marker = response.NextMarker
			args.MaxKeys = response.MaxKeys
		} else {
			break
		}
	}

	return result, nil
}
