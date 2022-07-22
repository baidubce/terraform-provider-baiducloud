package baiducloud

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, v.(string))
		}
	}
	return vs
}

func expandStringSet(configured *schema.Set) []string {
	return expandStringList(configured.List())
}

func flattenStringListToInterface(sl []string) []interface{} {
	result := make([]interface{}, 0, len(sl))

	for _, v := range sl {
		result = append(result, v)
	}

	return result
}
