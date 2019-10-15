package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/model"
	"github.com/hashicorp/terraform/helper/schema"
)

func tagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Tags",
		Optional:    true,
		ForceNew:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"tag_key": {
					Type:        schema.TypeString,
					Description: "Tag's key",
					Required:    true,
				},
				"tag_value": {
					Type:        schema.TypeString,
					Description: "Tag's value",
					Required:    true,
				},
			},
		},
	}
}

func tagsComputedSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "Tags",
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"tag_key": {
					Type:        schema.TypeString,
					Description: "Tag's key",
					Computed:    true,
				},
				"tag_value": {
					Type:        schema.TypeString,
					Description: "Tag's value",
					Computed:    true,
				},
			},
		},
	}
}

func flattenTagsToMap(tags []model.TagModel) []map[string]string {
	tagMap := make([]map[string]string, 0, len(tags))
	for _, tag := range tags {
		tagMap = append(tagMap, map[string]string{
			"tag_key":   tag.TagKey,
			"tag_value": tag.TagValue,
		})
	}

	return tagMap
}

func tranceTagMapToModel(tagMaps []interface{}) []model.TagModel {
	tags := make([]model.TagModel, 0, len(tagMaps))
	for _, t := range tagMaps {
		tag := t.(map[string]interface{})
		tags = append(tags, model.TagModel{
			TagKey:   tag["tag_key"].(string),
			TagValue: tag["tag_value"].(string),
		})
	}

	return tags
}
