package bcc

import "github.com/baidubce/bce-sdk-go/services/bcc/api"

func flattenKeyPairList(keyPairs []api.KeypairModel) interface{} {
	tfList := []map[string]interface{}{}
	for _, v := range keyPairs {
		tfMap := map[string]interface{}{
			"keypair_id":     v.KeypairId,
			"name":           v.Name,
			"description":    v.Description,
			"created_time":   v.CreatedTime,
			"public_key":     v.PublicKey,
			"instance_count": v.InstanceCount,
			"region_id":      v.RegionId,
			"fingerprint":    v.FingerPrint,
		}
		tfList = append(tfList, tfMap)
	}
	return tfList
}
