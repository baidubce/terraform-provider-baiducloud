package snic

import (
	"github.com/baidubce/bce-sdk-go/services/endpoint"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func FindSNIC(conn *connectivity.BaiduClient, snicId string) (*endpoint.Endpoint, error) {
	raw, err := conn.WithSNICClient(func(client *endpoint.Client) (interface{}, error) {
		return client.GetEndpointDetail(snicId)
	})
	if err != nil {
		return nil, err
	}
	return raw.(*endpoint.Endpoint), nil
}
