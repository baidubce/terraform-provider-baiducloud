package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/localDns"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

type LocalDnsService struct {
	client *connectivity.BaiduClient
}

func (s *LocalDnsService) GetPrivateZoneDetail(zoneId string) (*localDns.GetPrivateZoneResponse, error) {
	action := "Query Local Dns Private Zone " + zoneId

	raw, err := s.client.WithLocalDnsClient(func(localDnsClient *localDns.Client) (i interface{}, e error) {
		return localDnsClient.GetPrivateZone(zoneId)
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_localdns", action, BCESDKGoERROR)
	}

	return raw.(*localDns.GetPrivateZoneResponse), nil
}

func (s *LocalDnsService) GetPrivateZoneList() (*localDns.ListPrivateZoneResponse, error) {
	action := "Query Local Dns Private Zone List"
	listArgs := &localDns.ListPrivateZoneRequest{}
	raw, err := s.client.WithLocalDnsClient(func(localDnsClient *localDns.Client) (i interface{}, e error) {
		return localDnsClient.ListPrivateZone(listArgs)
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_localdns", action, BCESDKGoERROR)
	}

	return raw.(*localDns.ListPrivateZoneResponse), nil
}
