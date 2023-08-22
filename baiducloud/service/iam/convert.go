package iam

import (
	"time"

	"github.com/baidubce/bce-sdk-go/services/iam/api"
)

func flattenAccessKeyList(accessKeys []api.AccessKeyModel) interface{} {
	tfList := []map[string]interface{}{}
	for _, v := range accessKeys {
		tfMap := map[string]interface{}{
			"access_key_id":  v.Id,
			"enabled":        v.Enabled,
			"create_time":    v.CreateTime.Format(time.RFC3339),
			"last_used_time": v.LastUsedTime.Format(time.RFC3339),
		}
		tfList = append(tfList, tfMap)
	}
	return tfList
}
