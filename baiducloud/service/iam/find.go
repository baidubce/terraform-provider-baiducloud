package iam

import (
	"github.com/baidubce/bce-sdk-go/services/iam"
	"github.com/baidubce/bce-sdk-go/services/iam/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func FindAccessKeys(conn *connectivity.BaiduClient, username string) ([]api.AccessKeyModel, error) {
	accessKeys := make([]api.AccessKeyModel, 0)
	raw, err := conn.WithIamClient(func(client *iam.Client) (interface{}, error) {
		return client.ListAccessKey(username)
	})
	if err != nil {
		return nil, err
	}

	result := raw.(*api.ListAccessKeyResult)
	for _, item := range result.AccessKeys {
		accessKeys = append(accessKeys, item)
	}

	return accessKeys, nil
}

func FindAccessKey(conn *connectivity.BaiduClient, username, id string) (*api.AccessKeyModel, error) {
	accessKeys, err := FindAccessKeys(conn, username)
	if err != nil {
		return nil, err
	}

	for _, accessKey := range accessKeys {
		if accessKey.Id == id {
			return &accessKey, nil
		}
	}

	return nil, &resource.NotFoundError{}
}
