package eip

import (
	"github.com/baidubce/bce-sdk-go/services/eip"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/flex"
)

func expandMoveOutEips(tfList []interface{}) []eip.MoveOutEip {
	moveOutEips := []eip.MoveOutEip{}
	for _, v := range tfList {
		tfMap := v.(map[string]interface{})
		moveOutEip := eip.MoveOutEip{Eip: tfMap["eip"].(string)}

		if v, ok := tfMap["bandwidth_in_mbps"]; ok {
			moveOutEip.BandWidthInMbps = v.(int)
		}

		if pt, ok := tfMap["payment_timing"]; ok && len(pt.(string)) > 0 {
			if bm, ok := tfMap["billing_method"]; ok && len(bm.(string)) > 0 {
				moveOutEip.Billing = &eip.Billing{
					PaymentTiming: pt.(string),
					BillingMethod: bm.(string),
				}
			}
		}
		moveOutEips = append(moveOutEips, moveOutEip)
	}
	return moveOutEips
}

func flattenEipList(eips []eip.EipModel) interface{} {
	tfList := []map[string]interface{}{}
	for _, v := range eips {
		tfMap := map[string]interface{}{
			"eip_id":            v.EipId,
			"eip":               v.Eip,
			"name":              v.Name,
			"bandwidth_in_mbps": v.BandWidthInMbps,
			"status":            v.Status,
			"eip_instance_type": v.EipInstanceType,
			"instance_type":     v.InstanceType,
			"instance_id":       v.InstanceId,
			"share_group_id":    v.ShareGroupId,
			"payment_timing":    v.PaymentTiming,
			"billing_method":    v.BillingMethod,
			"create_time":       v.CreateTime,
			"expire_time":       v.ExpireTime,
			"tags":              flex.FlattenTagModelToMap(v.Tags),
		}
		tfList = append(tfList, tfMap)
	}
	return tfList
}
