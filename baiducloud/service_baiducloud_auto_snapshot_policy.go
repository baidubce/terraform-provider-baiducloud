package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
)

func (s *BccService) FlattenAutoSnapshotPolicyModelToMap(aspList []api.AutoSnapshotPolicyModel) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(aspList))

	for _, asp := range aspList {
		result = append(result, map[string]interface{}{
			"id":                asp.Id,
			"name":              asp.Name,
			"time_points":       asp.TimePoints,
			"repeat_weekdays":   asp.RepeatWeekdays,
			"status":            asp.Status,
			"retention_days":    asp.RetentionDays,
			"created_time":      asp.CreatedTime,
			"updated_time":      asp.UpdatedTime,
			"deleted_time":      asp.DeletedTime,
			"last_execute_time": asp.LastExecuteTime,
			"volume_count":      asp.VolumeCount,
		})
	}

	return result
}

func (s *BccService) ListAllAutoSnapshotPolicies(args *api.ListASPArgs) ([]api.AutoSnapshotPolicyModel, error) {
	result := make([]api.AutoSnapshotPolicyModel, 0)

	action := "Query All Auto Snapshot Policies"
	for {
		raw, err := s.client.WithBccClient(func(client *bcc.Client) (i interface{}, e error) {
			return client.ListAutoSnapshotPolicy(args)
		})
		addDebug(action, raw)

		if err != nil {
			return nil, WrapError(err)
		}
		response, _ := raw.(*api.ListASPResult)
		result = append(result, response.AutoSnapshotPolicys...)

		if response.IsTruncated {
			args.Marker = response.Marker
			args.MaxKeys = response.MaxKeys
		} else {
			return result, nil
		}
	}
}
