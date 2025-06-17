package hpas

import (
	"github.com/baidubce/bce-sdk-go/services/hpas/api"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/flex"
)

func flattenImageList(images []api.ImageResponse) interface{} {
	tfList := []map[string]interface{}{}
	for _, v := range images {
		tfMap := map[string]interface{}{
			"image_id":           v.ImageId,
			"name":               v.Name,
			"image_type":         v.ImageType,
			"image_status":       v.ImageStatus,
			"create_time":        v.CreateTime,
			"supported_app_type": flex.FlattenStringValueList(v.SupportedAppType),
		}
		tfList = append(tfList, tfMap)
	}
	return tfList
}
