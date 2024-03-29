package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/model"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func tagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeMap,
		Description: "Tags, do not support modify",
		Optional:    true,
		ForceNew:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func tagsComputedSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeMap,
		Description: "Tags",
		Computed:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func tagsCreationSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeMap,
		Description: "Tags, support setting when creating instance, do not support modify",
		Optional:    true,
		Computed:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func flattenTagsToMap(tags []model.TagModel) map[string]string {
	tagMap := make(map[string]string)
	for _, tag := range tags {
		tagMap[tag.TagKey] = tag.TagValue
	}

	return tagMap
}

func tranceTagMapToModel(tagMaps map[string]interface{}) []model.TagModel {
	tags := make([]model.TagModel, 0, len(tagMaps))
	for k, v := range tagMaps {
		tags = append(tags, model.TagModel{
			TagKey:   k,
			TagValue: v.(string),
		})
	}

	return tags
}
