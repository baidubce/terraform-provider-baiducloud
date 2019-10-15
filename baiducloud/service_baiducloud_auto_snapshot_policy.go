package baiducloud

import (
	"strconv"

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

// TODO:  api.AutoSnapshotPolicyModel2 may be wrong, it should has same structure with api.AutoSnapshotPolicyModel,
//  remove it when this wrong which in sdk be fixed in the future
func (s *BccService) TransAutoSnapshotPolicyModel2ToAutoSnapshotPolicyModel(asp2 *api.AutoSnapshotPolicyModel2) *api.AutoSnapshotPolicyModel {
	result := &api.AutoSnapshotPolicyModel{
		Id:              asp2.Id,
		Name:            asp2.Name,
		Status:          asp2.Status,
		CreatedTime:     asp2.CreatedTime,
		UpdatedTime:     asp2.UpdatedTime,
		DeletedTime:     asp2.DeletedTime,
		LastExecuteTime: asp2.LastExecuteTime,
	}

	result.RetentionDays, _ = strconv.Atoi(asp2.RetentionDays)
	result.VolumeCount, _ = strconv.Atoi(asp2.VolumeCount)

	for _, pointStr := range asp2.TimePoints {
		point, _ := strconv.Atoi(pointStr)
		result.TimePoints = append(result.TimePoints, point)
	}

	for _, dayStr := range asp2.RepeatWeekdays {
		day, _ := strconv.Atoi(dayStr)
		result.RepeatWeekdays = append(result.RepeatWeekdays, day)
	}

	return result
}
