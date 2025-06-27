package flex

import (
	"encoding/json"

	"github.com/baidubce/bce-sdk-go/model"
	"github.com/baidubce/bce-sdk-go/util"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func ExpandStringValueList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, v.(string))
		}
	}
	return vs
}

func FlattenStringValueList(list []string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, v)
	}
	return vs
}

func ExpandStringValueSet(configured *schema.Set) []string {
	return ExpandStringValueList(configured.List())
}

func FlattenStringValueSet(list []string) *schema.Set {
	return schema.NewSet(schema.HashString, FlattenStringValueList(list))
}

func TagsSchema() *schema.Schema {
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

func FlattenTagsToMap(tags []model.TagModel) map[string]string {
	tagMap := make(map[string]string)
	for _, tag := range tags {
		tagMap[tag.TagKey] = tag.TagValue
	}

	return tagMap
}

func FlattenTagModelToMap(tags interface{}) map[string]string {
	data, _ := json.Marshal(tags)
	var items []map[string]string
	_ = json.Unmarshal(data, &items)

	result := make(map[string]string)
	for _, item := range items {
		result[item["tagKey"]] = item["tagValue"]
	}
	return result
}

func TranceTagMapToModel(tagMaps map[string]interface{}) []model.TagModel {
	tags := make([]model.TagModel, 0, len(tagMaps))
	for k, v := range tagMaps {
		tags = append(tags, model.TagModel{
			TagKey:   k,
			TagValue: v.(string),
		})
	}

	return tags
}

func ExpandMapToTagModel[T any](tagMap map[string]interface{}) []T {
	var tempList []map[string]string
	for k, v := range tagMap {
		tempList = append(tempList, map[string]string{
			"tagKey":   k,
			"tagValue": v.(string),
		})
	}
	jsonBytes, _ := json.Marshal(tempList)

	var result []T
	_ = json.Unmarshal(jsonBytes, &result)
	return result
}

// 判断两个tag切片是否包含相同的元素
func SlicesContainSameElements(a, b []model.TagModel) bool {
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

func MergeSchema(origin map[string]*schema.Schema, adding map[string]*schema.Schema) {
	for k, v := range adding {
		origin[k] = v
	}
}

func DiffMaps(o, n map[string]interface{}) (added, removed map[string]interface{}) {
	added = map[string]interface{}{}
	removed = map[string]interface{}{}

	if o == nil {
		o = map[string]interface{}{}
	}
	if n == nil {
		n = map[string]interface{}{}
	}

	for k, vA := range n {
		if vB, ok := o[k]; !ok || vA != vB {
			added[k] = vA
		}
	}
	for k, vB := range o {
		if vA, ok := n[k]; !ok || vA != vB {
			removed[k] = vB
		}
	}

	return
}

func DoNothing(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func BuildClientToken() string {
	return util.NewUUID()
}
