package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/blb"
)

func (s *BLBService) BackendServerList(blbId string) ([]map[string]interface{}, error) {
	describeArgs := &blb.DescribeBackendServersArgs{}

	serversResult := make([]blb.BackendServerModel, 0)

	for {
		raw, err := s.client.WithBLBClient(func(client *blb.Client) (i interface{}, e error) {
			return client.DescribeBackendServers(blbId, describeArgs)
		})
		if err != nil {
			return nil, WrapError(err)
		}
		response := raw.(*blb.DescribeBackendServersResult)

		serversResult = append(serversResult, response.BackendServerList...)

		if response.IsTruncated {
			describeArgs.Marker = response.Marker
			describeArgs.MaxKeys = response.MaxKeys
		} else {
			break
		}

	}

	result := make([]map[string]interface{}, 0, len(serversResult))
	for _, server := range serversResult {
		result = append(result, map[string]interface{}{
			"instance_id": server.InstanceId,
			"weight":      server.Weight,
			"private_ip":  server.PrivateIp,
		})
	}

	return result, nil
}
