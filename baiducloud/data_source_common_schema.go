package baiducloud

import (
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type DataSourceFilter map[string][]FilterConfig

type FilterConfig struct {
	strFilterValue string
	regFilterValue *regexp.Regexp
}

func dataSourceFiltersSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "only support filter string/int/bool value",
		Optional:    true,
		ForceNew:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Description: "filter variable name",
					Required:    true,
				},

				"values": {
					Type:        schema.TypeList,
					Description: "filter variable value list",
					Required:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

func NewDataSourceFilter(d *schema.ResourceData) DataSourceFilter {
	fConfig, ok := d.GetOk("filter")
	if !ok {
		return nil
	}

	fList := fConfig.(*schema.Set).List()
	result := make(DataSourceFilter)
	for _, f := range fList {
		filter := f.(map[string]interface{})
		filterValue := filter["values"].([]interface{})

		filterConfig := make([]FilterConfig, 0)
		for _, v := range filterValue {
			fConfigValue := FilterConfig{
				strFilterValue: v.(string),
			}

			reg, err := regexp.Compile(fConfigValue.strFilterValue)
			if err == nil {
				fConfigValue.regFilterValue = reg
			}
			filterConfig = append(filterConfig, fConfigValue)
		}

		result[filter["name"].(string)] = filterConfig
	}

	return result
}

func (f *DataSourceFilter) checkFilter(data map[string]interface{}) bool {
	for key, fValue := range *f {
		if d, ok := data[key]; ok {
			checkResult := false
			for _, fConfig := range fValue {
				if fConfig.checkValue(d) {
					checkResult = true
					break
				}
			}

			if !checkResult {
				return false
			}
		}
	}

	return true
}

func (c *FilterConfig) checkValue(value interface{}) bool {
	checkValue := ""
	switch v := value.(type) {
	case string:
		checkValue = v
	case int:
		checkValue = strconv.Itoa(v)
	case bool:
		checkValue = strconv.FormatBool(v)
	case int32:
		checkValue = strconv.Itoa(int(v))
	default:
		return true
	}

	if checkValue == c.strFilterValue {
		return true
	}

	if c.regFilterValue != nil && c.regFilterValue.MatchString(checkValue) {
		return true
	}

	return false
}

func FilterDataSourceResult(d *schema.ResourceData, result *[]map[string]interface{}) {
	filter := NewDataSourceFilter(d)
	for index := 0; index < len(*result); {
		if !filter.checkFilter((*result)[index]) {
			*result = append((*result)[:index], (*result)[index+1:]...)
			continue
		}
		index++
	}
}
