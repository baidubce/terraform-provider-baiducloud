package cdn

import (
	"github.com/baidubce/bce-sdk-go/services/cdn"
	"github.com/baidubce/bce-sdk-go/services/cdn/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func FindDomainStatusByName(conn *connectivity.BaiduClient, name string) (*api.DomainStatus, error) {
	domains, err := FindDomainsStatus(conn, "ALL", name)
	if err != nil {
		return nil, err
	}

	for _, v := range domains {
		if v.Domain == name {
			return &v, nil
		}
	}

	return nil, &resource.NotFoundError{}
}

func FindDomainsStatus(conn *connectivity.BaiduClient, status string, rule string) ([]api.DomainStatus, error) {
	raw, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
		return client.GetDomainStatus(status, rule)
	})

	if err != nil {
		return nil, err
	}

	return raw.([]api.DomainStatus), nil
}

func FindDomainConfigByName(conn *connectivity.BaiduClient, domainName string) (*api.DomainConfig, error) {
	raw, err := conn.WithCdnClient(func(client *cdn.Client) (interface{}, error) {
		return client.GetDomainConfig(domainName)
	})

	if err != nil {
		return nil, err
	}

	return raw.(*api.DomainConfig), nil
}
