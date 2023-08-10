package bcc

import (
	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func FindKeyPairs(conn *connectivity.BaiduClient, name string) ([]api.KeypairModel, error) {
	keyPairs := make([]api.KeypairModel, 0)
	args := &api.ListKeypairArgs{Name: name}
	for {
		raw, err := conn.WithBccClient(func(client *bcc.Client) (interface{}, error) {
			return client.ListKeypairs(args)
		})
		if err != nil {
			return nil, err
		}

		result := raw.(*api.ListKeypairResult)
		for _, item := range result.Keypairs {
			keyPairs = append(keyPairs, item)
		}

		if !result.IsTruncated {
			break
		}
		args.Marker = result.NextMarker
	}
	return keyPairs, nil
}
