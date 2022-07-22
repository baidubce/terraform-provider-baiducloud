/*
Provide a resource to publish a CFC Function Version.

Example Usage

```hcl
resource "baiducloud_cfc_version" "default" {
  function_name       = "terraform-cfc"
  version_description = "terraformVersion"
}
```

```
*/
package baiducloud

import (
	"strings"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/cfc"
	"github.com/baidubce/bce-sdk-go/services/cfc/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudCFCVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudCFCVersionCreate,
		Read:   resourceBaiduCloudCFCVersionRead,
		Update: resourceBaiduCloudCFCVersionUpdate,
		Delete: resourceBaiducloudCFCVersionDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"function_name": {
				Type:        schema.TypeString,
				Description: "CFC function name, length must be between 1 and 64 bytes",
				Required:    true,
				ForceNew:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Function description",
				Computed:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "Function version",
				Computed:    true,
			},
			"version_description": {
				Type:        schema.TypeString,
				Description: "Function version description",
				Optional:    true,
				ForceNew:    true,
			},
			"environment": {
				Type:        schema.TypeMap,
				Description: "CFC Function environment variables",
				Computed:    true,
			},
			"handler": {
				Type:        schema.TypeString,
				Description: "CFC Function execution handler",
				Computed:    true,
			},
			"memory_size": {
				Type:        schema.TypeInt,
				Description: "CFC Function memory size",
				Computed:    true,
			},
			"runtime": {
				Type:        schema.TypeString,
				Description: "CFC Function runtime",
				Computed:    true,
			},
			"time_out": {
				Type:        schema.TypeInt,
				Description: "Function time out",
				Computed:    true,
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
				Optional:    true,
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
				Description: "CFC Function bined VPC Subnet id list",
				Computed:    true,
			},
			"vpc_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_ids": {
							Type:        schema.TypeSet,
							Description: "CFC Function bined VPC Subnet id list",
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
						},
						"security_group_ids": {
							Type:        schema.TypeSet,
							Description: "CFC Function binded Security group list",
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"log_type": {
				Type:         schema.TypeString,
				Description:  "Log save type, support bos/none",
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"none", "bos"}, false),
			},
			"log_bos_dir": {
				Type:        schema.TypeString,
				Description: "Log save dir if log type is bos",
				Computed:    true,
				Optional:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if value, ok := d.GetOk("log_type"); !ok && value.(string) != "bos" {
						return true
					}

					return false
				},
			},
		},
	}
}

func resourceBaiduCloudCFCVersionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	createArgs := &api.PublishVersionArgs{
		FunctionName: d.Get("function_name").(string),
	}

	if value, ok := d.GetOk("version_description"); ok {
		createArgs.Description = value.(string)
	}

	if value, ok := d.GetOk("code_sha256"); ok {
		createArgs.CodeSha256 = value.(string)
	}

	action := "Public CFC Function " + createArgs.FunctionName
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithCFCClient(func(client *cfc.Client) (i interface{}, e error) {
			return client.PublishVersion(createArgs)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		addDebug(action, raw)
		response, _ := raw.(*api.PublishVersionResult)
		d.SetId(response.FunctionName + "-" + response.Version)
		d.Set("version", response.Version)
		d.Set("function_name", response.FunctionName)
		d.Set("function_brn", response.FunctionBrn)

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfc_version", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudCFCVersionUpdate(d, meta)
}

func resourceBaiduCloudCFCVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	cfcService := CFCService{client}

	functionName := d.Get("function_name").(string)
	functionVersion := d.Get("version").(string)
	action := "Query function " + functionName + " with version " + functionVersion

	function, err := cfcService.CFCGetVersionsByFunction(functionName, functionVersion)
	if err != nil {
		d.SetId("")
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfc_alias", action, BCESDKGoERROR)
	}

	d.Set("uid", function.Uid)
	d.Set("description", function.Description)
	d.Set("function_brn", function.FunctionBrn)
	d.Set("function_arn", function.FunctionBrn)
	d.Set("region", function.Region)
	d.Set("time_out", function.Timeout)
	d.Set("version_desc", function.VersionDesc)
	d.Set("update_time", function.UpdatedAt.String())
	d.Set("last_modified", function.LastModified.String())
	d.Set("code_sha256", function.CodeSha256)
	d.Set("code_size", function.CodeSize)
	d.Set("handler", function.Handler)
	d.Set("runtime", function.Runtime)
	d.Set("memory_size", function.MemorySize)
	d.Set("commit_id", function.CommitID)
	d.Set("role", function.Role)
	d.Set("source_type", function.SourceTag)
	d.Set("code_id", function.CodeID)
	d.Set("log_type", function.LogType)
	if function.LogType == "bos" {
		if strings.HasPrefix(function.LogBosDir, "bos://") {
			d.Set("log_bos_dir", function.LogBosDir[6:])
		} else {
			d.Set("log_bos_dir", function.LogBosDir)
		}
	}

	if function.Environment != nil {
		d.Set("environment", function.Environment.Variables)
	}

	if function.VpcConfig != nil {
		d.Set("vpc_config", []interface{}{cfcService.faltternCFCFunctionVpcConfigToMap(function.VpcConfig)})
	}

	return nil
}

func resourceBaiduCloudCFCVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	cfcService := CFCService{client}

	if d.HasChange("log_type") || d.HasChange("log_bos_dir") {
		functionName := d.Get("function_name").(string)
		functionVersion := d.Get("version").(string)
		functionBrn := d.Get("function_brn").(string)

		updateArgs := &api.UpdateFunctionConfigurationArgs{
			FunctionName: functionBrn,
		}

		if value, ok := d.GetOk("log_type"); ok {
			updateArgs.LogType = value.(string)
		}

		if updateArgs.LogType == "bos" {
			if value, ok := d.GetOk("log_bos_dir"); ok {
				updateArgs.LogBosDir = value.(string)
				if !strings.HasPrefix(updateArgs.LogBosDir, "bos://") {
					updateArgs.LogBosDir = "bos://" + updateArgs.LogBosDir
				}
			}
		}

		action := "Update function " + functionName + " with version " + functionVersion
		if err := cfcService.CFCUpdateFunctionConfigure(updateArgs); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfc_function", action, BCESDKGoERROR)
		}
	}

	return resourceBaiduCloudCFCVersionRead(d, meta)
}

func resourceBaiducloudCFCVersionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	deleteArgs := &api.DeleteFunctionArgs{
		FunctionName: d.Get("function_name").(string),
		Qualifier:    d.Get("version").(string),
	}

	action := "Delete CFC Function " + deleteArgs.FunctionName + " version " + deleteArgs.Qualifier
	err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		_, err := client.WithCFCClient(func(client *cfc.Client) (i interface{}, e error) {
			return nil, client.DeleteFunction(deleteArgs)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		addDebug(action, deleteArgs)
		return nil
	})

	if err != nil {
		if IsExceptedErrors(err, ObjectNotFound) {
			return nil
		}

		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfc_version", action, BCESDKGoERROR)
	}

	return nil
}
