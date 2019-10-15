package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/hashicorp/terraform/helper/resource"
)

func (s *BccService) SnapshotStateRefreshFunc(id string, failState []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		raw, err := s.client.WithBccClient(func(client *bcc.Client) (i interface{}, e error) {
			return client.GetSnapshotDetail(id)
		})

		if err != nil {
			return nil, "", WrapError(err)
		}

		status := raw.(*api.GetSnapshotDetailResult).Snapshot.Status
		for _, state := range failState {
			if string(status) == state {
				return string(status), string(status), WrapError(Error(GetFailTargetStatus, string(status)))
			}
		}

		return raw, string(status), nil
	}
}

func (s *BccService) ListAllSnapshots(volumeId string) ([]api.SnapshotModel, error) {
	result := make([]api.SnapshotModel, 0)

	args := &api.ListSnapshotArgs{
		VolumeId: volumeId,
	}

	action := "Query all Snapshots"
	for {
		raw, err := s.client.WithBccClient(func(bccClient *bcc.Client) (interface{}, error) {
			return bccClient.ListSnapshot(args)
		})
		addDebug(action, raw)

		if err != nil {
			return nil, WrapError(err)
		}
		response, _ := raw.(*api.ListSnapshotResult)
		result = append(result, response.Snapshots...)

		if response.IsTruncated {
			args.Marker = response.Marker
			args.MaxKeys = response.MaxKeys
		} else {
			return result, nil
		}
	}
}

func (s *BccService) FlattenSnapshotModelToMap(snapshots []api.SnapshotModel) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(snapshots))

	for _, sp := range snapshots {
		result = append(result, map[string]interface{}{
			"id":            sp.Id,
			"name":          sp.Name,
			"size_in_gb":    sp.SizeInGB,
			"create_time":   sp.CreateTime,
			"status":        sp.Status,
			"create_method": sp.CreateMethod,
			"volume_id":     sp.VolumeId,
			"description":   sp.Description,
		})
	}

	return result
}
