package flex

import (
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
