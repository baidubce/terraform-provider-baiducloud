package hpas

import (
	"fmt"
	"log"

	"github.com/baidubce/bce-sdk-go/services/hpas"
	"github.com/baidubce/bce-sdk-go/services/hpas/api"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func FindImages(conn *connectivity.BaiduClient, imageType string) ([]api.ImageResponse, error) {
	result := []api.ImageResponse{}
	marker := ""
	for {
		raw, err := conn.WithHPASClient(func(client *hpas.Client) (interface{}, error) {
			args := &api.DescribeHpasImageReq{
				Marker:  marker,
				MaxKeys: 100,
			}
			args.ImageType = imageType
			return client.ImageList(args)
		})

		log.Printf("[DEBUG] Read HPAS image list result: %+v", raw)
		if err != nil {
			return nil, fmt.Errorf("error reading HPAS image list: %w", err)
		}
		response := raw.(*api.DescribeHpasImageResp)
		result = append(result, response.Images...)

		if len(response.NextMarker) > 0 {
			marker = response.NextMarker
		} else {
			break
		}
	}
	return result, nil

}
