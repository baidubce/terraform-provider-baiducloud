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
				ImageType: imageType,
				Marker:    marker,
				MaxKeys:   100,
			}
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

func FindInstances(conn *connectivity.BaiduClient, queryArgs api.ListHpasByMakerReq) ([]api.HpasResponse, error) {
	result := []api.HpasResponse{}
	marker := ""
	for {
		raw, err := conn.WithHPASClient(func(client *hpas.Client) (interface{}, error) {
			args := &api.ListHpasByMakerReq{
				HpasIds:    queryArgs.HpasIds,
				Name:       queryArgs.Name,
				ZoneName:   queryArgs.ZoneName,
				HpasStatus: queryArgs.HpasStatus,
				AppType:    queryArgs.AppType,
				Marker:     marker,
				MaxKeys:    100,
			}
			return client.DescribeHPASInstancesByMaker(args)
		})

		log.Printf("[DEBUG] Read HPAS instance list result: %+v", raw)
		if err != nil {
			return nil, fmt.Errorf("error reading HPAS instance list: %w", err)
		}
		response := raw.(*api.ListHpasByMakerResp)
		result = append(result, response.Hpass...)

		if len(response.NextMarker) > 0 {
			marker = response.NextMarker
		} else {
			break
		}
	}
	return result, nil
}

func FindInstance(conn *connectivity.BaiduClient, instanceID string) (*api.HpasResponse, error) {
	args := api.ListHpasByMakerReq{
		HpasIds: []string{instanceID},
	}
	instances, err := FindInstances(conn, args)
	if err != nil {
		return nil, err
	}

	if len(instances) == 0 {
		return nil, fmt.Errorf("instance %s not found", instanceID)
	}

	return &instances[0], nil
}
