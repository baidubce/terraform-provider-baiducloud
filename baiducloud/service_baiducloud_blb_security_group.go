package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/blb"
)

func (s *BLBService) GetBlbSecurityGroup(blbId string) (*blb.DescribeSecurityGroupsResult, error) {
	action := "Query blb security groups " + blbId

	raw, err := s.client.WithBLBClient(func(blbClient *blb.Client) (i interface{}, e error) {
		return blbClient.DescribeSecurityGroups(blbId)
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb_security_groups", action, BCESDKGoERROR)
	}

	return raw.(*blb.DescribeSecurityGroupsResult), nil
}
