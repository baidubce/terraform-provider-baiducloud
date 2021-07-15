package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/bbc"
)

func (s *BbcService) FlattenBbcImageModelToMap(images []bbc.ImageModel) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(images))

	for _, image := range images {
		result = append(result, map[string]interface{}{
			"id":              image.Id,
			"name":            image.Name,
			"type":            image.Type,
			"os_type":         image.OsType,
			"os_version":      image.OsVersion,
			"os_arch":         image.OsArch,
			"os_name":         image.OsName,
			"os_build":        image.OsBuild,
			"create_time":     image.CreateTime,
			"status":          image.Status,
			"description":     image.Desc,
			"special_version": image.SpecialVersion,
		})
	}

	return result
}

func (s *BbcService) ListAllBbcImages(args *bbc.ListImageArgs) ([]bbc.ImageModel, error) {
	action := "List all " + args.ImageType + " images"

	result := make([]bbc.ImageModel, 0)
	for {
		raw, err := s.client.WithBbcClient(func(client *bbc.Client) (i interface{}, e error) {
			return client.ListImage(args)
		})
		if err != nil {
			return nil, WrapError(err)
		}
		addDebug(action, raw)

		response := raw.(*bbc.ListImageResult)
		result = append(result, response.Images...)

		if response.IsTruncated {
			args.Marker = response.NextMarker
			args.MaxKeys = response.MaxKeys
		} else {
			return result, nil
		}
	}
}
