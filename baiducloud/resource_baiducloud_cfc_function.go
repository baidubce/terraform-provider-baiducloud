/*
Provide a resource to create an CFC Function.

Example Usage

```hcl
resource "baiducloud_cfc_function" "default" {
  function_name = "terraform-cfc"
  description   = "terraform create"
  handler       = "index.handler"
  memory_size   = 256
  runtime       = "nodejs8.5"
  time_out      = 20
  code_zip_file = "UEsDBBQACAAIAAyjX00AAAAAAAAAAAAAAAAIABAAaW5kZXguanNVWAwAsJ/ZW/ie2Vv6Z7qeS60oyC8qKdbLSMxLyUktUrBV0EgtS80r0VFIzs8rSa0AMRJzcpISk7M1FWztFKq5FIAAJqSRV5qTo6Og5JGak5OvUJ5flJOiqKRpzVVrDQBQSwcILzRMjVAAAABYAAAAUEsDBAoAAAAAAHCjX00AAAAAAAAAAAAAAAAJABAAX19NQUNPU1gvVVgMALSf2Vu0n9lb+me6nlBLAwQUAAgACAAMo19NAAAAAAAAAAAAAAAAEwAQAF9fTUFDT1NYLy5faW5kZXguanNVWAwAsJ/ZW/ie2Vv6Z7qeY2AVY2dgYmDwTUxW8A9WiFCAApAYAycQGwFxHRCD+BsYiAKOISFBUCZIxwIgFkBTwogQl0rOz9VLLCjISdXLSSwuKS1OTUlJLElVDggGKXw772Y0iO5J8tAH0QBQSwcIDgnJLFwAAACwAAAAUEsBAhUDFAAIAAgADKNfTS80TI1QAAAAWAAAAAgADAAAAAAAAAAAQKSBAAAAAGluZGV4LmpzVVgIALCf2Vv4ntlbUEsBAhUDCgAAAAAAcKNfTQAAAAAAAAAAAAAAAAkADAAAAAAAAAAAQP1BlgAAAF9fTUFDT1NYL1VYCAC0n9lbtJ/ZW1BLAQIVAxQACAAIAAyjX00OCcksXAAAALAAAAATAAwAAAAAAAAAAECkgc0AAABfX01BQ09TWC8uX2luZGV4LmpzVVgIALCf2Vv4ntlbUEsFBgAAAAADAAMA0gAAAHoBAAAAAA=="
}
```

Import

CFC can be imported, e.g.

```hcl
$ terraform import baiducloud_cfc_function.default functionName
```
*/
package baiducloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/cfc"
	"github.com/baidubce/bce-sdk-go/services/cfc/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudCFCFunction() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudCFCFunctionCreate,
		Read:   resourceBaiduCloudCFCFunctionRead,
		Update: resourceBaiduCloudCFCFunctionUpdate,
		Delete: resourceBaiduCloudCFCFunctionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"function_name": {
				Type:         schema.TypeString,
				Description:  "CFC function name, length must be between 1 and 64 bytes",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"description": {
				Type:         schema.TypeString,
				Description:  "Function description",
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 256),
			},
			"version": {
				Type:        schema.TypeString,
				Description: "Function version, should only be $LATEST",
				Computed:    true,
			},
			"environment": {
				Type:        schema.TypeMap,
				Description: "CFC Function environment variables",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"handler": {
				Type:        schema.TypeString,
				Description: "CFC Function execution handler",
				Required:    true,
			},
			"memory_size": {
				Type:         schema.TypeInt,
				Description:  "CFC Function memory size, should be an integer multiple of 128",
				Optional:     true,
				Default:      128,
				ValidateFunc: validateCFCMemorySize,
			},
			"runtime": {
				Type:        schema.TypeString,
				Description: "CFC Function runtime",
				Required:    true,
			},
			"time_out": {
				Type:         schema.TypeInt,
				Description:  "Function time out, support [1, 300]s",
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 300),
			},
			"code_file_dir": {
				Type:          schema.TypeString,
				Description:   "CFC Function Code local file dir",
				Optional:      true,
				ConflictsWith: []string{"code_file_name", "code_bos_bucket", "code_bos_object"},
			},
			"code_file_name": {
				Type:          schema.TypeString,
				Description:   "CFC Function Code local zip file name",
				Optional:      true,
				ConflictsWith: []string{"code_file_dir", "code_bos_bucket", "code_bos_object"},
			},
			"code_bos_bucket": {
				Type:          schema.TypeString,
				Description:   "CFC Function Code storage bos bucket name",
				Optional:      true,
				ConflictsWith: []string{"code_file_name", "code_file_dir"},
			},
			"code_bos_object": {
				Type:          schema.TypeString,
				Description:   "CFC Function Code storage bos object key",
				Optional:      true,
				ConflictsWith: []string{"code_file_name", "code_file_dir"},
			},
			"code_storage": {
				Type:        schema.TypeMap,
				Description: "CFC Code storage information",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"update_time": {
				Type:        schema.TypeString,
				Description: "Last update time",
				Computed:    true,
			},
			"last_modified": {
				Type:        schema.TypeString,
				Description: "The same as update_time",
				Computed:    true,
			},
			"code_sha256": {
				Type:        schema.TypeString,
				Description: "Function code sha256",
				Computed:    true,
			},
			"code_size": {
				Type:        schema.TypeString,
				Description: "Function code size",
				Computed:    true,
			},
			"function_brn": {
				Type:        schema.TypeString,
				Description: "Function brn",
				Computed:    true,
			},
			"function_arn": {
				Type:        schema.TypeString,
				Description: "The same as function brn",
				Computed:    true,
			},
			"commit_id": {
				Type:        schema.TypeString,
				Description: "Function commit id",
				Computed:    true,
			},
			"role": {
				Type:        schema.TypeString,
				Description: "Function exec role",
				Computed:    true,
			},
			"uid": {
				Type:        schema.TypeString,
				Description: "Function user uid",
				Computed:    true,
			},
			"region": {
				Type:        schema.TypeString,
				Description: "Function region",
				Computed:    true,
			},
			"vpc_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_ids": {
							Type:        schema.TypeSet,
							Description: "CFC Function bined VPC Subnet id list",
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
						},
						"security_group_ids": {
							Type:        schema.TypeSet,
							Description: "CFC Function binded Security group list",
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},

				// Suppress diffs if the VPC configuration is provided, but empty
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.Id() == "" || old == "1" || new == "0" {
						return false
					}

					if d.HasChange("vpc_config.0.security_group_ids") || d.HasChange("vpc_config.0.subnet_ids") {
						return false
					}

					return true
				},
			},
			"reserved_concurrent_executions": {
				Type:         schema.TypeInt,
				Description:  "Function reserved concurrent executions, support [0-90]",
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 90),
			},
			"log_type": {
				Type:         schema.TypeString,
				Description:  "Log save type, support bos/none",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"", "none", "bos"}, false),
			},
			"log_bos_dir": {
				Type:        schema.TypeString,
				Description: "Log save dir if log type is bos",
				Optional:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if value, ok := d.GetOk("log_type"); !ok && value.(string) != "bos" {
						return true
					}

					return false
				},
			},
			"source_tag": {
				Type:        schema.TypeString,
				Description: "CFC Function source tag",
				Computed:    true,
			},
			"code_id": {
				Type:        schema.TypeString,
				Description: "CFC Function code id",
				Computed:    true,
			},
		},
	}
}

func resourceBaiduCloudCFCFunctionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	cfcService := CFCService{client}

	createArgs, err := buildBaiduCloudCreateCFCFunctionArgs(d)
	if err != nil {
		return WrapError(err)
	}
	action := "Create CFC Function " + createArgs.FunctionName

	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithCFCClient(func(client *cfc.Client) (i interface{}, e error) {
			return client.CreateFunction(createArgs)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		addDebug(action, raw)
		response, _ := raw.(*api.CreateFunctionResult)
		d.SetId(response.FunctionName)

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfc_function", action, BCESDKGoERROR)
	}

	if value, ok := d.GetOk("reserved_concurrent_executions"); ok {
		if err := cfcService.CFCSetReservedConcurrent(d.Id(), value.(int)); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfc_function", action, BCESDKGoERROR)
		}
	}

	return resourceBaiduCloudCFCFunctionRead(d, meta)
}

func resourceBaiduCloudCFCFunctionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	cfcService := CFCService{client}

	functionName := d.Id()
	action := "Query CFC Function " + functionName

	raw, err := client.WithCFCClient(func(client *cfc.Client) (i interface{}, e error) {
		return client.GetFunction(
			&api.GetFunctionArgs{
				FunctionName: functionName,
			})
	})

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfc_function", action, BCESDKGoERROR)
	}

	response := raw.(*api.GetFunctionResult)
	d.Set("function_name", response.Configuration.FunctionName)
	d.Set("uid", response.Configuration.Uid)
	d.Set("description", response.Configuration.Description)
	d.Set("function_brn", response.Configuration.FunctionBrn)
	d.Set("function_arn", response.Configuration.FunctionBrn)
	d.Set("region", response.Configuration.Region)
	d.Set("time_out", response.Configuration.Timeout)
	d.Set("version_desc", response.Configuration.VersionDesc)
	d.Set("update_time", response.Configuration.UpdatedAt.String())
	d.Set("last_modified", response.Configuration.LastModified.String())
	d.Set("code_sha256", response.Configuration.CodeSha256)
	d.Set("code_size", response.Configuration.CodeSize)
	d.Set("handler", response.Configuration.Handler)
	d.Set("runtime", response.Configuration.Runtime)
	d.Set("memory_size", response.Configuration.MemorySize)
	d.Set("commit_id", response.Configuration.CommitID)
	d.Set("role", response.Configuration.Role)
	d.Set("source_type", response.Configuration.SourceTag)
	d.Set("code_id", response.Configuration.CodeID)
	d.Set("version", response.Configuration.Version)
	d.Set("log_type", response.Configuration.LogType)
	if response.Configuration.LogType == "bos" {
		if strings.HasPrefix(response.Configuration.LogBosDir, "bos://") {
			d.Set("log_bos_dir", response.Configuration.LogBosDir[6:])
		} else {
			d.Set("log_bos_dir", response.Configuration.LogBosDir)
		}
	}

	if response.Configuration.Environment != nil {
		d.Set("environment", response.Configuration.Environment.Variables)
	}

	codeStorage := make(map[string]string)
	codeStorage["location"] = response.Code.Location
	codeStorage["repository_type"] = response.Code.RepositoryType
	d.Set("code_storage", codeStorage)

	if response.Configuration.VpcConfig != nil {
		d.Set("vpc_config", []interface{}{cfcService.faltternCFCFunctionVpcConfigToMap(response.Configuration.VpcConfig)})
	}

	return nil
}

func resourceBaiduCloudCFCFunctionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	cfcService := CFCService{client}

	functionName := d.Id()
	action := "Update Fucntion " + functionName

	if d.HasChange("code_file_name") || d.HasChange("code_file_dir") || d.HasChange("code_bos_bucket") || d.HasChange("code_bos_object") {
		if err := cfcService.CFCUpdateFunctionCode(functionName, d); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfc_function", action, BCESDKGoERROR)
		}
	}

	if d.HasChange("reserved_concurrent_executions") {
		if newConcurrent, ok := d.GetOk("reserved_concurrent_executions"); !ok || (newConcurrent.(int) == 0) {
			if err := cfcService.CFCDeleteReservedConcurrent(functionName); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfc_function", action, BCESDKGoERROR)
			}
		} else {
			if err := cfcService.CFCSetReservedConcurrent(functionName, newConcurrent.(int)); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfc_function", action, BCESDKGoERROR)
			}
		}
	}

	if updateConfig, updateConfigArgs := buildBaiduCloudUpdateCFCFunctionConfigArgs(d); updateConfig {
		if err := cfcService.CFCUpdateFunctionConfigure(updateConfigArgs); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfc_function", action, BCESDKGoERROR)
		}
	}

	return resourceBaiduCloudCFCFunctionRead(d, meta)
}

func resourceBaiduCloudCFCFunctionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	functionName := d.Id()
	action := "Delete CFC Function " + functionName

	err := resource.Retry(d.Timeout(schema.TimeoutDefault), func() *resource.RetryError {
		_, err := client.WithCFCClient(func(client *cfc.Client) (i interface{}, e error) {
			return nil, client.DeleteFunction(
				&api.DeleteFunctionArgs{
					FunctionName: functionName,
				})
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		addDebug(action, functionName)
		return nil
	})

	if err != nil {
		if IsExceptedErrors(err, ObjectNotFound) {
			return nil
		}

		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfc_function", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudCreateCFCFunctionArgs(d *schema.ResourceData) (*api.CreateFunctionArgs, error) {
	result := &api.CreateFunctionArgs{
		FunctionName: d.Get("function_name").(string),
		MemorySize:   d.Get("memory_size").(int),
		Handler:      d.Get("handler").(string),
		Runtime:      d.Get("runtime").(string),
		Timeout:      d.Get("time_out").(int),
	}

	if value, ok := d.GetOk("description"); ok {
		result.Description = value.(string)
	}

	if value, ok := d.GetOk("environment"); ok {
		environments := value.(map[string]interface{})
		result.Environment = &api.Environment{
			Variables: make(map[string]string),
		}
		for k, v := range environments {
			result.Environment.Variables[k] = v.(string)
		}
	}

	if value, ok := d.GetOk("vpc_config"); ok && len(value.([]interface{})) > 0 {
		vpcConfig := value.([]interface{})[0].(map[string]interface{})

		result.VpcConfig = &api.VpcConfig{
			SubnetIds:        expandStringSet(vpcConfig["subnet_ids"].(*schema.Set)),
			SecurityGroupIds: expandStringSet(vpcConfig["security_group_ids"].(*schema.Set)),
		}
	}

	if value, ok := d.GetOk("log_type"); ok {
		result.LogType = value.(string)

		if result.LogType == "bos" {
			if value, ok := d.GetOk("log_bos_dir"); ok {
				result.LogBosDir = value.(string)
				if !strings.HasPrefix(result.LogBosDir, "bos://") {
					result.LogBosDir = "bos://" + result.LogBosDir
				}
			}
		}
	}

	result.Code = &api.CodeFile{}
	fileName, fileNameOk := d.GetOk("code_file_name")
	fileDir, fileDirOk := d.GetOk("code_file_dir")
	if fileNameOk {
		zipFile, err := loadFileContent(fileName.(string))
		if err != nil {
			return nil, err
		}

		result.Code.ZipFile = zipFile
	} else if fileDirOk {
		zipFile, err := zipFileDir(fileDir.(string))
		if err != nil {
			return nil, err
		}

		result.Code.ZipFile = zipFile
	} else {
		bucket, bucketOk := d.GetOk("code_bos_bucket")
		object, objectOk := d.GetOk("code_bos_object")

		if !bucketOk || !objectOk {
			return nil, fmt.Errorf("code_bos_bucket and code_bos_object must all be set while using bos code source")
		}

		result.Code.BosBucket = bucket.(string)
		result.Code.BosObject = object.(string)
	}

	return result, nil
}

func buildBaiduCloudUpdateCFCFunctionConfigArgs(d *schema.ResourceData) (update bool, args *api.UpdateFunctionConfigurationArgs) {
	update = false
	args = &api.UpdateFunctionConfigurationArgs{
		FunctionName: d.Id(),
	}

	if d.HasChange("time_out") {
		update = true
		args.Timeout = d.Get("time_out").(int)
	}

	if d.HasChange("description") {
		update = true
		if value, ok := d.GetOk("description"); ok {
			args.Description = value.(string)
		}
	}

	if d.HasChange("handler") {
		update = true
		args.Handler = d.Get("handler").(string)
	}

	if d.HasChange("runtime") {
		update = true
		args.Runtime = d.Get("runtime").(string)
	}

	if d.HasChange("memory_size") {
		update = true
		args.MemorySize = d.Get("memory_size").(int)
	}

	if d.HasChange("environment") {
		update = true
		if value, ok := d.GetOk("environment"); ok {
			environments := value.(map[string]interface{})
			args.Environment = &api.Environment{
				Variables: make(map[string]string),
			}
			for k, v := range environments {
				args.Environment.Variables[k] = v.(string)
			}
		}
	}

	if d.HasChange("vpc_config") {
		update = true
		args.VpcConfig = &api.VpcConfig{}
		if value, ok := d.GetOk("vpc_config"); ok && len(value.([]interface{})) > 0 {
			vpcConfig := value.([]interface{})[0].(map[string]interface{})

			args.VpcConfig.SubnetIds = expandStringSet(vpcConfig["subnet_ids"].(*schema.Set))
			args.VpcConfig.SecurityGroupIds = expandStringSet(vpcConfig["security_group_ids"].(*schema.Set))
		}
	}

	if d.HasChange("log_type") {
		update = true
		if value, ok := d.GetOk("log_type"); ok {
			args.LogType = value.(string)
		}
	}

	if d.HasChange("log_bos_dir") {
		update = true
		if value, ok := d.GetOk("log_bos_dir"); ok {
			args.LogBosDir = value.(string)
			if !strings.HasPrefix(args.LogBosDir, "bos://") {
				args.LogBosDir = "bos://" + args.LogBosDir
			}
		}
	}

	return
}
