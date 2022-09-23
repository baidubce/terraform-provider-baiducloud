package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
)

func (s *BccService) ListAllDeploySets() ([]api.DeploySetModel, error) {
	action := "List all deploy sets"

	raw, err := s.client.WithBccClient(func(client *bcc.Client) (i interface{}, e error) {
		return client.ListDeploySets()
	})
	if err != nil {
		return nil, WrapError(err)
	}
	addDebug(action, raw)
	response := raw.(*api.ListDeploySetsResult)

	return response.DeploySetList, nil
}
