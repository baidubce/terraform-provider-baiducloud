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

// 判断两个tag切片是否包含相同的元素
func slicesContainSameElements(a, b []model.TagModel) bool {
	if len(a) != len(b) {
		return false
	}
	// 创建映射来存储每个 TagModel 出现的次数
	counts := make(map[model.TagModel]int)
	// 计算第一个切片中每个元素出现的次数
	for _, item := range a {
		counts[item]++
	}
	// 减去第二个切片中每个元素出现的次数
	for _, item := range b {
		counts[item]--
		if counts[item] < 0 {
			// 如果某个元素在第二个切片中出现的次数多于第一个切片，返回 false
			return false
		}
	}
	// 检查所有元素的计数是否为 0
	for _, count := range counts {
		if count != 0 {
			return false
		}
	}
	return true
}
