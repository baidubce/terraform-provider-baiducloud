package baiducloud

import (
	"encoding/json"
	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/dts"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"strconv"
)

type DtsService struct {
	client *connectivity.BaiduClient
}

func (e *DtsService) Config(taskId string, configArgs *dts.ConfigArgs) (*dts.ConfigDtsResult, error) {
	action := "config dts task " + taskId

	raw, err := e.client.WithDtsClient(func(client *dts.Client) (i interface{}, e error) {
		return client.ConfigDts(taskId, configArgs)
	})

	if err != nil {
		if !IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
			return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_dts", "", BCESDKGoERROR)
		}
	}

	response := raw.(*dts.ConfigDtsResult)

	addDebug(action, configArgs)
	return response, nil
}

func (e *DtsService) PreCheck(taskId string) (*dts.PreCheckResult, error) {
	action := "precheck dts task " + taskId

	raw, err := e.client.WithDtsClient(func(client *dts.Client) (i interface{}, e error) {
		return client.PreCheck(taskId)
	})

	if err != nil {
		if !IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
			return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_dts", "", BCESDKGoERROR)
		}
	}

	response := raw.(*dts.PreCheckResult)

	for {
		result, _ := e.GetTaskDetail(taskId)
		if result.Status != DTSStatusChecking {
			break
		}
	}

	addDebug(action, action)
	return response, nil
}

func (e *DtsService) GetPreCheck(taskId string) ([]dts.CheckResult, error) {
	action := "get precheck result, dts task " + taskId

	raw, err := e.client.WithDtsClient(func(client *dts.Client) (i interface{}, e error) {
		return client.GetPreCheck(taskId)
	})

	if err != nil {
		return nil, WrapError(err)
	}

	response := raw.(*dts.GetPreCheckResult)

	result := make([]dts.CheckResult, 0)
	result = append(result, response.Result...)

	addDebug(action, action)
	return result, nil
}

func (e *DtsService) StartDts(taskId string) error {
	action := "start dts task " + taskId

	_, err := e.client.WithDtsClient(func(client *dts.Client) (i interface{}, e error) {
		return nil, client.StartDts(taskId)
	})

	if err != nil {
		if !IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dts", "", BCESDKGoERROR)
		}
	}

	addDebug(action, action)
	return nil
}

func (e *DtsService) PauseDts(taskId string) error {
	action := "pause dts task " + taskId

	_, err := e.client.WithDtsClient(func(client *dts.Client) (i interface{}, e error) {
		return nil, client.PauseDts(taskId)
	})

	if err != nil {
		if !IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dts", "", BCESDKGoERROR)
		}
	}

	addDebug(action, action)
	return nil
}

func (e *DtsService) ShutdownDts(taskId string) error {
	action := "shut down dts task " + taskId

	_, err := e.client.WithDtsClient(func(client *dts.Client) (i interface{}, e error) {
		return nil, client.ShutdownDts(taskId)
	})

	if err != nil {
		if !IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dts", "", BCESDKGoERROR)
		}
	}

	addDebug(action, action)
	return nil
}

func (e *DtsService) ListAllDtss(listArgs *dts.ListDtsArgs) ([]dts.DtsTaskMeta, error) {
	result := make([]dts.DtsTaskMeta, 0)
	for {
		raw, err := e.client.WithDtsClient(func(client *dts.Client) (interface{}, error) {
			return client.ListDts(listArgs)
		})

		if err != nil {
			return nil, WrapError(err)
		}

		response := raw.(*dts.ListDtsResult)
		result = append(result, response.Task...)

		isTruncated := response.IsTruncated
		if isTruncated {
			listArgs.MaxKeys = response.MaxKeys
			listArgs.Marker = response.Marker
		} else {
			return result, nil
		}
	}
}

func (e *DtsService) TaskStateRefresh(taskId string, failState []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		result, err := e.GetTaskDetail(taskId)
		if err != nil {
			return nil, "", WrapError(err)
		}

		for _, statue := range failState {
			if result.Status == statue {
				return result, result.Status, WrapError(Error(GetFailTargetStatus, result.Status))
			}
		}

		return result, result.Status, nil
	}
}

func (e *DtsService) GetTaskDetail(taskID string) (*dts.DtsTaskMeta, error) {
	action := "Get DTS instance detail " + taskID
	raw, err := e.client.WithDtsClient(func(dtsClient *dts.Client) (interface{}, error) {
		return dtsClient.GetDetail(taskID)
	})
	addDebug(action, raw)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_dts", action, BCESDKGoERROR)
	}

	result, _ := raw.(*dts.DtsTaskMeta)
	return result, nil
}

func (e *DtsService) FlattenDtsModelsToMap(dtss []dts.DtsTaskMeta) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(dtss))

	for _, c := range dtss {
		result = append(result, map[string]interface{}{
			"dts_id":                c.DtsId,
			"running_time":          c.RunningTime,
			"task_name":             c.TaskName,
			"status":                c.Status,
			"data_type":             c.DataType,
			"region":                c.Region,
			"create_time":           c.CreateTime,
			"src_connection":        flattenConnectionToMap(c.SrcConnection),
			"dst_connection":        flattenConnectionToMap(c.DstConnection),
			"schema_mapping":        c.SchemaMapping,
			"sub_status":            c.SubStatus,
			"schema":                c.DynamicInfo.Schema,
			"base":                  c.DynamicInfo.Base,
			"increment":             c.DynamicInfo.Increment,
			"errmsg":                c.Errmsg,
			"sdk_realtime_progress": c.SdkRealtimeProgress,
			"granularity":           c.Granularity,
			"sub_start_time":        c.SubDataScope.StartTime,
			"sub_end_time":          c.SubDataScope.EndTime,
			"product_type":          c.PayInfo.ProductType,
			"source_instance_type":  c.PayInfo.SourceInstanceType,
			"target_instance_type":  c.PayInfo.TargetInstanceType,
			"cross_region_tag":      c.PayInfo.CrossRegionTag,
			"pay_create_time":       c.PayInfo.CreateTime,
			"standard":              c.PayInfo.Standard,
			"pay_end_time":          c.PayInfo.EndTime,
		})
	}
	return result
}

func flattenIncrementToMap(dynamicInfo dts.DynamicInfo) map[string]string {
	increment := dynamicInfo.Increment
	result := make(map[string]string)

	data, _ := json.Marshal(increment)
	json.Unmarshal(data, &result)

	return result
}

func connectionSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeMap,
		Description: "Connection",
		Optional:    true,
		Computed:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func flattenConnectionToMap(connection dts.Connection) map[string]string {
	ConnectionMap := make(map[string]string)
	ConnectionMap["instance_type"] = connection.InstanceType
	ConnectionMap["region"] = connection.Region
	ConnectionMap["db_type"] = connection.DbType
	ConnectionMap["db_user"] = connection.DbUser
	ConnectionMap["db_pass"] = ""
	ConnectionMap["db_port"] = strconv.Itoa(connection.DbPort)
	ConnectionMap["db_host"] = connection.DbHost
	ConnectionMap["instance_id"] = connection.InstanceId
	return ConnectionMap
}

func schemaInfo() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "schemaInfo",
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"current": {
					Type:        schema.TypeString,
					Description: "current",
					Computed:    true,
				},
				"count": {
					Type:        schema.TypeString,
					Description: "count",
					Computed:    true,
				},
				"speed": {
					Type:        schema.TypeString,
					Description: "speed",
					Computed:    true,
				},
				"expect_finish_time": {
					Type:        schema.TypeString,
					Description: "expect finish time",
					Computed:    true,
				},
			},
		},
	}
}

func flattenSchemaInfoToMap(dynamicInfo dts.DynamicInfo) ([]map[string]string, []map[string]string) {
	schema := make([]map[string]string, 0, len(dynamicInfo.Schema))
	for _, e := range dynamicInfo.Schema {
		schema = append(schema, map[string]string{
			"current":            e.Current,
			"count":              e.Count,
			"speed":              e.Speed,
			"expect_finish_time": e.ExpectFinishTime,
		})
	}

	base := make([]map[string]string, 0, len(dynamicInfo.Base))
	for _, e := range dynamicInfo.Base {
		base = append(base, map[string]string{
			"current":            e.Current,
			"count":              e.Count,
			"speed":              e.Speed,
			"expect_finish_time": e.ExpectFinishTime,
		})
	}

	return schema, base
}
