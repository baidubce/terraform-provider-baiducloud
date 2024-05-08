package abroad

import (
	"github.com/baidubce/bce-sdk-go/model"
	"github.com/baidubce/bce-sdk-go/services/cdn/abroad"
	"github.com/baidubce/bce-sdk-go/services/cdn/abroad/api"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func FindAbroadDomainConfigByName(conn *connectivity.BaiduClient, domainName string) (*api.DomainConfig, error) {
	raw, err := conn.WithAbroadCdnClient(func(client *abroad.Client) (interface{}, error) {
		return client.GetDomainConfig(domainName)
	})

	if err != nil {
		return nil, err
	}

	return raw.(*api.DomainConfig), nil
}

func FindAbroadDomainTagsByName(conn *connectivity.BaiduClient, domainName string) ([]model.TagModel, error) {
	raw, err := conn.WithAbroadCdnClient(func(client *abroad.Client) (interface{}, error) {
		return client.GetTags(domainName)
	})

	if err != nil {
		return nil, err
	}
	return raw.([]model.TagModel), nil
}
