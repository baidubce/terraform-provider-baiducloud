/*
Provide a resource to create a DTS.

Example Usage

```hcl
resource "baiducloud_dts" "default" {
    product_type         = "postpay"
	type                 = "migration"
	standard             = "Large"
	source_instance_type = "public"
	target_instance_type = "public"
	cross_region_tag     = 0

    task_name            = "taskname"
	data_type			 = ["schema","base"]
    src_connection = {
        region          = "public"
		db_type			= "mysql"
		db_user			= "baidu"
		db_pass			= "password"
		db_port			= 3306
		db_host			= "106.12.174.191"
		instance_id		= "rds-lNy3KsQQ"
		instance_type	= "public"
    }
	dst_connection = {
        region          = "public"
		db_type			= "mysql"
		db_user			= "baidu"
		db_pass			= "password"
		db_port			= 3306
		db_host			= "106.12.174.191"
		instance_id		= "rds-lNy3KsQQ"
		instance_type	= "public"
    }
    schema_mapping {
			type        = "db"
			src			= "db1"
			dst			= "db2"
			where		= ""
	}
}
```

Import

DTS can be imported, e.g.

```hcl
$ terraform import baiducloud_dts.default dts
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/dts"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"strconv"
	"time"
)

func resourceBaiduCloudDts() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudDtsCreate,
		Read:   resourceBaiduCloudDtsRead,
		Update: resourceBaiduCloudDtsUpdate,
		Delete: resourceBaiduCloudDtsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"operation": {
				Type:         schema.TypeString,
				Description:  "operation of the task. Available values are precheck, getprecheck, start, pause, shutdown.",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"precheck", "getprecheck", "start", "pause", "shutdown"}, false),
			},
			"product_type": {
				Type:         schema.TypeString,
				Description:  "product type of the task. Available value is postpay.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"postpay"}, false),
			},
			"type": {
				Type:         schema.TypeString,
				Description:  "type of the task. Available values are migration, sync, subscribe.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"migration", "sync", "subscribe"}, false),
			},
			"standard": {
				Type:         schema.TypeString,
				Description:  "standard of the task. Available value is Large.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Large"}, false),
			},
			"source_instance_type": {
				Type:         schema.TypeString,
				Description:  "source instance type of the task. Available values are public, bcerds.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"public", "bcerds"}, false),
			},
			"target_instance_type": {
				Type:         schema.TypeString,
				Description:  "target instance type of the task. Available values are public, bcerds.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"public", "bcerds"}, false),
			},
			"cross_region_tag": {
				Type:         schema.TypeInt,
				Description:  "cross region tag of the task. Available value are 0, 1.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
			"dts_id": {
				Type:        schema.TypeString,
				Description: "Dts task id",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"task_name": {
				Type:        schema.TypeString,
				Description: "Dts task name",
				Required:    true,
			},
			"running_time": {
				Type:        schema.TypeInt,
				Description: "Dts task running time",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Dts task status",
				Computed:    true,
			},
			"data_type": {
				Type:        schema.TypeList,
				Description: "Dts task data type",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "Dts create time",
				Computed:    true,
			},
			"region": {
				Type:        schema.TypeString,
				Description: "Dts region",
				Computed:    true,
			},
			"src_connection": connectionSchema(),
			"dst_connection": connectionSchema(),
			"schema_mapping": {
				Type:        schema.TypeList,
				Description: "schema mapping",
				Required:    true,
				MinItems:    1,
				MaxItems:    5,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Description: "type",
							Optional:    true,
							Computed:    true,
						},
						"src": {
							Type:        schema.TypeString,
							Description: "src",
							Optional:    true,
							Computed:    true,
						},
						"dst": {
							Type:        schema.TypeString,
							Description: "dst",
							Optional:    true,
							Computed:    true,
						},
						"where": {
							Type:        schema.TypeString,
							Description: "where",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"sub_status": {
				Type:        schema.TypeList,
				Description: "sub status",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"s": {
							Type:        schema.TypeString,
							Description: "s",
							Computed:    true,
						},
						"b": {
							Type:        schema.TypeString,
							Description: "b",
							Computed:    true,
						},
						"i": {
							Type:        schema.TypeString,
							Description: "i",
							Computed:    true,
						},
					},
				},
			},
			"schema": schemaInfo(),
			"base":   schemaInfo(),
			"increment": {
				Type:        schema.TypeMap,
				Description: "increment",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"errmsg": {
				Type:        schema.TypeString,
				Description: "Dts error message",
				Computed:    true,
			},
			"sdk_realtime_progress": {
				Type:        schema.TypeString,
				Description: "Dts sdk realtime progress",
				Computed:    true,
			},
			"granularity": {
				Type:        schema.TypeString,
				Description: "Dts granularity",
				Optional:    true,
				Computed:    true,
			},
			"sub_start_time": {
				Type:        schema.TypeString,
				Description: "Dts subDataScope start time",
				Computed:    true,
			},
			"sub_end_time": {
				Type:        schema.TypeString,
				Description: "Dts subDataScope end time",
				Computed:    true,
			},
			"pay_create_time": {
				Type:        schema.TypeInt,
				Description: "Dts pay create time",
				Computed:    true,
			},
			"pay_end_time": {
				Type:        schema.TypeString,
				Description: "Dts pay end time",
				Computed:    true,
			},
			"init_position_type": {
				Type:        schema.TypeString,
				Description: "Dts init position type",
				Optional:    true,
				Computed:    true,
			},
			"init_position": {
				Type:        schema.TypeString,
				Description: "Dts init position",
				Optional:    true,
				Computed:    true,
			},
			"queue_type": {
				Type:        schema.TypeString,
				Description: "Dts queue type",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceBaiduCloudDtsCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	action := "Create DTS "

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		createDtsArgs := buildBaiduCloudCreateDtsArgs(d)

		createRaw, createErr := client.WithDtsClient(func(dtsClient *dts.Client) (interface{}, error) {
			return dtsClient.CreateDts(createDtsArgs)
		})
		if createErr != nil {
			if IsExceptedErrors(createErr, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(createErr)
			}
			return resource.NonRetryableError(createErr)
		}
		response, _ := createRaw.(*dts.CreateDtsResult)

		addDebug(action, createRaw)
		d.SetId(response.DtsTasks[0].DtsId)

		configDtsArgs := buildBaiduCloudConfigDtsArgs(d)
		configRaw, configErr := client.WithDtsClient(func(dtsClient *dts.Client) (interface{}, error) {
			return dtsClient.ConfigDts(d.Id(), configDtsArgs)
		})
		if configErr != nil {
			if IsExceptedErrors(configErr, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(configErr)
			}
			return resource.NonRetryableError(configErr)
		}
		addDebug("Config DTS", configRaw)

		preCheckRaw, preCheckErr := client.WithDtsClient(func(dtsClient *dts.Client) (interface{}, error) {
			return dtsClient.PreCheck(d.Id())
		})
		if preCheckErr != nil {
			if IsExceptedErrors(preCheckErr, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(preCheckErr)
			}
			return resource.NonRetryableError(preCheckErr)
		}
		addDebug("PreCheck DTS", preCheckRaw)

		_, startErr := client.WithDtsClient(func(dtsClient *dts.Client) (interface{}, error) {
			return nil, dtsClient.StartDts(d.Id())
		})
		if startErr != nil {
			if IsExceptedErrors(startErr, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(startErr)
			}
			return resource.NonRetryableError(startErr)
		}

		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dts", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudDtsRead(d, meta)
}

func resourceBaiduCloudDtsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	taskID := d.Id()
	action := "Query DTS Instance " + taskID

	raw, err := client.WithDtsClient(func(dtsClient *dts.Client) (interface{}, error) {
		return dtsClient.GetDetail(taskID)
	})

	addDebug(action, raw)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dts_instance", action, BCESDKGoERROR)
	}

	result, _ := raw.(*dts.DtsTaskMeta)
	schema, base := flattenSchemaInfoToMap(result.DynamicInfo)

	d.Set("dts_id", result.DtsId)
	d.Set("task_name", result.TaskName)
	d.Set("running_time", result.RunningTime)
	d.Set("status", result.Status)
	d.Set("data_type", result.DataType)
	d.Set("create_time", result.CreateTime)
	d.Set("region", result.Region)
	d.Set("src_connection", flattenConnectionToMap(result.SrcConnection))
	d.Set("dst_connection", flattenConnectionToMap(result.DstConnection))
	d.Set("schema_mapping", result.SchemaMapping)
	d.Set("sub_status", result.SubStatus)
	d.Set("schema", schema)
	d.Set("base", base)
	d.Set("increment", flattenIncrementToMap(result.DynamicInfo))
	d.Set("errmsg", result.Errmsg)
	d.Set("sdk_realtime_progress", result.SdkRealtimeProgress)
	d.Set("granularity", result.Granularity)
	d.Set("sub_start_time", result.SubDataScope.StartTime)
	d.Set("sub_end_time", result.SubDataScope.EndTime)
	d.Set("pay_create_time", result.PayInfo.CreateTime)
	d.Set("pay_end_time", result.PayInfo.EndTime)

	return nil
}

func resourceBaiduCloudDtsUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	dtsClient := DtsService{client}

	taskId := d.Id()
	action := "Update Fucntion " + taskId
	d.Partial(true)

	if updateConfig, updateConfigArgs := buildUpdateConfigArgs(d); updateConfig {
		if _, err := dtsClient.Config(taskId, updateConfigArgs); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dts", action, BCESDKGoERROR)
		}
	}

	if d.HasChange("operation") {

		switch d.Get("operation").(string) {

		case "precheck":

			stateConf := buildStateConf([]string{DTSStatusReady, DTSStatusChecking},
				[]string{DTSStatusCheckPass},
				d.Timeout(schema.TimeoutUpdate),
				dtsClient.TaskStateRefresh(taskId, []string{DTSStatusCheckFailed}))

			if _, err := dtsClient.PreCheck(taskId); err != nil {
				return WrapError(err)
			}

			if _, err := stateConf.WaitForState(); err != nil {
				return WrapError(err)
			}

			d.SetPartial("operation")

		case "getprecheck":

			if _, err := dtsClient.GetPreCheck(taskId); err != nil {
				return WrapError(err)
			}

			d.SetPartial("operation")

		case "start":

			stateConf := buildStateConf([]string{DTSStatusCheckPass, DTSStatusRunning},
				[]string{DTSStatusFinished},
				d.Timeout(schema.TimeoutUpdate),
				dtsClient.TaskStateRefresh(taskId, []string{DTSStatusRunFailed}))

			if err := dtsClient.StartDts(taskId); err != nil {
				return WrapError(err)
			}

			if _, err := stateConf.WaitForState(); err != nil {
				return WrapError(err)
			}

			d.SetPartial("operation")

		case "pause":

			stateConf := buildStateConf([]string{DTSStatusStopping, DTSStatusRunning},
				[]string{DTSStatusStopped},
				d.Timeout(schema.TimeoutUpdate),
				dtsClient.TaskStateRefresh(taskId, []string{}))

			if err := dtsClient.PauseDts(taskId); err != nil {
				return WrapError(err)
			}

			if _, err := stateConf.WaitForState(); err != nil {
				return WrapError(err)
			}

			d.SetPartial("operation")

		case "shutdown":

			stateConf := buildStateConf([]string{DTSStatusCheckFailed, DTSStatusCheckPass, DTSStatusRunning, DTSStatusStopped, DTSStatusRunFailed},
				[]string{DTSStatusFinished},
				d.Timeout(schema.TimeoutUpdate),
				dtsClient.TaskStateRefresh(taskId, []string{}))

			if err := dtsClient.ShutdownDts(taskId); err != nil {
				return WrapError(err)
			}

			if _, err := stateConf.WaitForState(); err != nil {
				return WrapError(err)
			}

			d.SetPartial("operation")
		}
	}

	d.Partial(false)
	return resourceBaiduCloudDtsRead(d, meta)
}

func resourceBaiduCloudDtsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	dtsClient := DtsService{client}

	taskId := d.Id()
	action := "Delete DTS Instance " + taskId

	for {
		result, _ := dtsClient.GetTaskDetail(taskId)
		if result.Status == DTSStatusFinished || result.Status == DTSStatusUnConfig ||
			result.Status == DTSStatusReady {
			break
		}
	}

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithDtsClient(func(dtsClient *dts.Client) (interface{}, error) {
			return taskId, dtsClient.DeleteDts(taskId)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InvalidInstanceStatus, bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)
		return nil
	})
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidInstanceStatus, InstanceNotExist, bce.EINTERNAL_ERROR}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dts_instance", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudCreateDtsArgs(d *schema.ResourceData) *dts.CreateDtsArgs {
	result := &dts.CreateDtsArgs{}

	if v, ok := d.GetOk("product_type"); ok && v.(string) != "" {
		result.ProductType = v.(string)
	}

	if v, ok := d.GetOk("type"); ok && v.(string) != "" {
		result.Type = v.(string)
	}

	if v, ok := d.GetOk("standard"); ok && v.(string) != "" {
		result.Standard = v.(string)
	}

	if v, ok := d.GetOk("source_instance_type"); ok && v.(string) != "" {
		result.SourceInstanceType = v.(string)
	}

	if v, ok := d.GetOk("target_instance_type"); ok && v.(string) != "" {
		result.TargetInstanceType = v.(string)
	}

	if v, ok := d.GetOk("cross_region_tag"); ok {
		result.CrossRegionTag = v.(int)
	}

	return result
}

func buildBaiduCloudConfigDtsArgs(d *schema.ResourceData) *dts.ConfigArgs {
	arg := &dts.ConfigArgs{}

	if v, ok := d.GetOk("type"); ok && v.(string) != "" {
		arg.Type = v.(string)
	}

	if v, ok := d.GetOk("task_name"); ok && v.(string) != "" {
		arg.TaskName = v.(string)
	}

	if v, ok := d.GetOk("data_type"); ok && v.([]interface{}) != nil {
		for _, e := range v.([]interface{}) {
			arg.DataType = append(arg.DataType, e.(string))
		}
	}

	if v, ok := d.GetOk("src_connection"); ok && v.(map[string]interface{}) != nil {
		srcConnection := v.(map[string]interface{})

		if region, ok := srcConnection["region"]; ok && region.(string) != "" {
			arg.SrcConnection.Region = region.(string)
		}
		if dbType, ok := srcConnection["db_type"]; ok && dbType.(string) != "" {
			arg.SrcConnection.DbType = dbType.(string)
		}
		if dbUser, ok := srcConnection["db_user"]; ok && dbUser.(string) != "" {
			arg.SrcConnection.DbUser = dbUser.(string)
		}
		if dbPass, ok := srcConnection["db_pass"]; ok && dbPass.(string) != "" {
			arg.SrcConnection.DbPass = dbPass.(string)
		}
		if port, ok := srcConnection["db_port"]; ok && port.(string) != "" {
			if dbPort, err := strconv.Atoi(port.(string)); err == nil {
				arg.SrcConnection.DbPort = dbPort
			}
		}
		if dbHost, ok := srcConnection["db_host"]; ok && dbHost.(string) != "" {
			arg.SrcConnection.DbHost = dbHost.(string)
		}
		if instanceId, ok := srcConnection["instance_id"]; ok && instanceId.(string) != "" {
			arg.SrcConnection.InstanceId = instanceId.(string)
		}
		if instanceType, ok := srcConnection["instance_type"]; ok && instanceType.(string) != "" {
			arg.SrcConnection.InstanceType = instanceType.(string)
		}
		if whitelist, ok := srcConnection["field_whitelist"]; ok && whitelist.(string) != "" {
			arg.SrcConnection.FieldWhitelist = whitelist.(string)
		}
		if blacklist, ok := srcConnection["field_blacklist"]; ok && blacklist.(string) != "" {
			arg.SrcConnection.FieldBlacklist = blacklist.(string)
		}
		if startTime, ok := srcConnection["start_time"]; ok && startTime.(string) != "" {
			arg.SrcConnection.StartTime = startTime.(string)
		}
		if endTime, ok := srcConnection["end_time"]; ok && endTime.(string) != "" {
			arg.SrcConnection.EndTime = endTime.(string)
		}
	}

	if v, ok := d.GetOk("dst_connection"); ok && v.(map[string]interface{}) != nil {
		dstConnection := v.(map[string]interface{})

		if region, ok := dstConnection["region"]; ok && region.(string) != "" {
			arg.DstConnection.Region = region.(string)
		}
		if dbType, ok := dstConnection["db_type"]; ok && dbType.(string) != "" {
			arg.DstConnection.DbType = dbType.(string)
		}
		if dbUser, ok := dstConnection["db_user"]; ok && dbUser.(string) != "" {
			arg.DstConnection.DbUser = dbUser.(string)
		}
		if dbPass, ok := dstConnection["db_pass"]; ok && dbPass.(string) != "" {
			arg.DstConnection.DbPass = dbPass.(string)
		}
		if port, ok := dstConnection["db_port"]; ok && port.(string) != "" {
			if dbPort, err := strconv.Atoi(port.(string)); err == nil {
				arg.DstConnection.DbPort = dbPort
			}
		}
		if dbHost, ok := dstConnection["db_host"]; ok && dbHost.(string) != "" {
			arg.DstConnection.DbHost = dbHost.(string)
		}
		if instanceId, ok := dstConnection["instance_id"]; ok && instanceId.(string) != "" {
			arg.DstConnection.InstanceId = instanceId.(string)
		}
		if instanceType, ok := dstConnection["instance_type"]; ok && instanceType.(string) != "" {
			arg.DstConnection.InstanceType = instanceType.(string)
		}
		if whitelist, ok := dstConnection["field_whitelist"]; ok && whitelist.(string) != "" {
			arg.DstConnection.FieldWhitelist = whitelist.(string)
		}
		if blacklist, ok := dstConnection["field_blacklist"]; ok && blacklist.(string) != "" {
			arg.DstConnection.FieldBlacklist = blacklist.(string)
		}
		if startTime, ok := dstConnection["start_time"]; ok && startTime.(string) != "" {
			arg.DstConnection.StartTime = startTime.(string)
		}
		if endTime, ok := dstConnection["end_time"]; ok && endTime.(string) != "" {
			arg.DstConnection.EndTime = endTime.(string)
		}
	}

	if v, ok := d.GetOk("schema_mapping"); ok && v.([]interface{}) != nil {
		schemaMapping := v.([]interface{})
		argSchemaMapping := make([]dts.Schema, 0)
		for _, e := range schemaMapping {
			schema := dts.Schema{}
			m := e.(map[string]interface{})
			schema.Type = m["type"].(string)
			schema.Src = m["src"].(string)
			schema.Dst = m["dst"].(string)
			schema.Where = m["where"].(string)
			argSchemaMapping = append(argSchemaMapping, schema)
		}
		arg.SchemaMapping = argSchemaMapping
	}

	if v, ok := d.GetOk("init_position"); ok && v.(string) != "" {
		arg.InitPosition.Position = v.(string)
	}

	if v, ok := d.GetOk("init_position_type"); ok && v.(string) != "" {
		arg.InitPosition.Type = v.(string)
	}

	if v, ok := d.GetOk("granularity"); ok && v.(string) != "" {
		arg.Granularity = v.(string)
	}

	if v, ok := d.GetOk("queue_type"); ok && v.(string) != "" {
		arg.QueueType = v.(string)
	}

	return arg
}

func buildUpdateConfigArgs(d *schema.ResourceData) (bool, *dts.ConfigArgs) {
	update := false
	arg := &dts.ConfigArgs{
		Type: "migration",
	}

	if d.HasChange("task_name") || d.HasChange("data_type") || d.HasChange("src_connection") ||
		d.HasChange("dst_connection") || d.HasChange("schema_mapping") || d.HasChange("init_position") ||
		d.HasChange("init_position_type") || d.HasChange("granularity") || d.HasChange("queue_type") {

		update = true
		arg.TaskName = d.Get("task_name").(string)

		dataType := d.Get("data_type").([]interface{})
		for _, e := range dataType {
			arg.DataType = append(arg.DataType, e.(string))
		}

		srcConnection := d.Get("src_connection").(map[string]interface{})
		if region, ok := srcConnection["region"]; ok {
			arg.SrcConnection.Region = region.(string)
		}
		if dbType, ok := srcConnection["db_type"]; ok {
			arg.SrcConnection.DbType = dbType.(string)
		}
		if dbUser, ok := srcConnection["db_user"]; ok {
			arg.SrcConnection.DbUser = dbUser.(string)
		}
		if dbPass, ok := srcConnection["db_pass"]; ok {
			arg.SrcConnection.DbPass = dbPass.(string)
		}
		if port, ok := srcConnection["db_port"]; ok {
			dbPort, _ := strconv.Atoi(port.(string))
			arg.SrcConnection.DbPort = dbPort
		}
		if dbHost, ok := srcConnection["db_host"]; ok {
			arg.SrcConnection.DbHost = dbHost.(string)
		}
		if instanceId, ok := srcConnection["instance_id"]; ok {
			arg.SrcConnection.InstanceId = instanceId.(string)
		}
		if instanceType, ok := srcConnection["instance_type"]; ok {
			arg.SrcConnection.InstanceType = instanceType.(string)
		}
		if whitelist, ok := srcConnection["field_whitelist"]; ok {
			arg.SrcConnection.FieldWhitelist = whitelist.(string)
		}
		if blacklist, ok := srcConnection["field_blacklist"]; ok {
			arg.SrcConnection.FieldBlacklist = blacklist.(string)
		}
		if startTime, ok := srcConnection["start_time"]; ok {
			arg.SrcConnection.StartTime = startTime.(string)
		}
		if endTime, ok := srcConnection["end_time"]; ok {
			arg.SrcConnection.EndTime = endTime.(string)
		}

		dstConnection := d.Get("dst_connection").(map[string]interface{})
		if region, ok := dstConnection["region"]; ok {
			arg.DstConnection.Region = region.(string)
		}
		if dbType, ok := dstConnection["db_type"]; ok {
			arg.DstConnection.DbType = dbType.(string)
		}
		if dbUser, ok := dstConnection["db_user"]; ok {
			arg.DstConnection.DbUser = dbUser.(string)
		}
		if dbPass, ok := dstConnection["db_pass"]; ok {
			arg.DstConnection.DbPass = dbPass.(string)
		}
		if port, ok := dstConnection["db_port"]; ok {
			dbPort, _ := strconv.Atoi(port.(string))
			arg.DstConnection.DbPort = dbPort
		}
		if dbHost, ok := dstConnection["db_host"]; ok {
			arg.DstConnection.DbHost = dbHost.(string)
		}
		if instanceId, ok := dstConnection["instance_id"]; ok {
			arg.DstConnection.InstanceId = instanceId.(string)
		}
		if instanceType, ok := dstConnection["instance_type"]; ok {
			arg.DstConnection.InstanceType = instanceType.(string)
		}
		if whitelist, ok := dstConnection["field_whitelist"]; ok {
			arg.DstConnection.FieldWhitelist = whitelist.(string)
		}
		if blacklist, ok := dstConnection["field_blacklist"]; ok {
			arg.DstConnection.FieldBlacklist = blacklist.(string)
		}
		if startTime, ok := dstConnection["start_time"]; ok {
			arg.DstConnection.StartTime = startTime.(string)
		}
		if endTime, ok := dstConnection["end_time"]; ok {
			arg.DstConnection.EndTime = endTime.(string)
		}

		schemaMapping := d.Get("schema_mapping").([]interface{})
		argSchemaMapping := make([]dts.Schema, 0)
		for _, e := range schemaMapping {
			schema := dts.Schema{}
			m := e.(map[string]interface{})
			schema.Type = m["type"].(string)
			schema.Src = m["src"].(string)
			schema.Dst = m["dst"].(string)
			schema.Where = m["where"].(string)
			argSchemaMapping = append(argSchemaMapping, schema)
		}
		arg.SchemaMapping = argSchemaMapping

		if v, ok := d.GetOk("init_position"); ok && v.(string) != "" {
			update = true
			arg.InitPosition.Position = v.(string)
		}

		if v, ok := d.GetOk("init_position_type"); ok && v.(string) != "" {
			update = true
			arg.InitPosition.Type = v.(string)
		}

		if v, ok := d.GetOk("granularity"); ok && v.(string) != "" {
			update = true
			arg.Granularity = v.(string)
		}

		if v, ok := d.GetOk("queue_type"); ok && v.(string) != "" {
			update = true
			arg.QueueType = v.(string)
		}
	}

	d.SetPartial("task_name")
	d.SetPartial("data_type")
	d.SetPartial("src_connection")
	d.SetPartial("dst_connection")
	d.SetPartial("schema_mapping")

	return update, arg
}
