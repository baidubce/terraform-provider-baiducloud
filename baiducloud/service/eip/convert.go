package eip

import (
	"github.com/baidubce/bce-sdk-go/services/eip"
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
