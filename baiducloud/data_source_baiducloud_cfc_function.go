/*
Use this data source to get a function.

Example Usage

```hcl
data "baiducloud_cfc_function" "default" {
   function_name = "terraform-create"
}

output "function" {
 value = "${data.baiducloud_cfc_function.default}"
}
```
*/
package baiducloud

import (
	"strings"

	"github.com/baidubce/bce-sdk-go/services/cfc"
	"github.com/baidubce/bce-sdk-go/services/cfc/api"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudCFCFunction() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudCFCFunctionRead,

		Schema: map[string]*schema.Schema{
			"function_name": {
				Type:         schema.TypeString,
				Description:  "CFC function name, length must be between 1 and 64 bytes",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"qualifier": {
				Type:        schema.TypeString,
				Description: "Function search qualifier",
				Optional:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Function description",
				Computed:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "Function version, should only be $LATEST",
				Computed:    true,
			},
			"environment": {
				Type:        schema.TypeMap,
				Description: "CFC Function environment variables",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"handler": {
				Type:        schema.TypeString,
				Description: "CFC Function execution handler",
				Computed:    true,
			},
			"memory_size": {
				Type:        schema.TypeInt,
				Description: "CFC Function memory size, should be an integer multiple of 128",
				Computed:    true,
			},
			"runtime": {
				Type:        schema.TypeString,
				Description: "CFC Function runtime",
				Computed:    true,
			},
			"time_out": {
				Type:        schema.TypeInt,
				Description: "Function time out, support [1, 300]s",
				Computed:    true,
			},
			"code_zip_file": {
				Type:        schema.TypeString,
				Description: "CFC Function Code base64-encoded data",
				Computed:    true,
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
				Type:        schema.TypeList,
				Description: "Function VPC Config",
				Computed:    true,
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
			"reserved_concurrent_executions": {
				Type:        schema.TypeInt,
				Description: "Function reserved concurrent executions, support [0-90]",
				Computed:    true,
			},
			"log_type": {
				Type:        schema.TypeString,
				Description: "Log save type, support bos/none",
				Computed:    true,
			},
			"log_bos_dir": {
				Type:        schema.TypeString,
				Description: "Log save dir if log type is bos",
				Computed:    true,
			},
			"source_tag": {
				Type:        schema.TypeString,
				Description: "CFC Function source tag",
				Computed:    true,
			},
		},
	}
}

func dataSourceBaiduCloudCFCFunctionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	cfcService := CFCService{client}

	args := &api.GetFunctionArgs{
		FunctionName: d.Get("function_name").(string),
	}

	if value, ok := d.GetOk("qualifier"); ok {
		args.Qualifier = value.(string)
	}

	action := "Get Function " + args.FunctionName + " with qualifier " + args.Qualifier
	raw, err := client.WithCFCClient(func(client *cfc.Client) (i interface{}, e error) {
		return client.GetFunction(args)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfc_function", action, BCESDKGoERROR)
	}
	addDebug(action, raw)

	function := raw.(*api.GetFunctionResult)
	d.Set("uid", function.Configuration.Uid)
	d.Set("description", function.Configuration.Description)
	d.Set("function_brn", function.Configuration.FunctionBrn)
	d.Set("function_arn", function.Configuration.FunctionBrn)
	d.Set("region", function.Configuration.Region)
	d.Set("time_out", function.Configuration.Timeout)
	d.Set("version", function.Configuration.Version)
	d.Set("version_desc", function.Configuration.VersionDesc)
	d.Set("update_time", function.Configuration.UpdatedAt.String())
	d.Set("last_modified", function.Configuration.LastModified.String())
	d.Set("code_sha256", function.Configuration.CodeSha256)
	d.Set("code_size", function.Configuration.CodeSize)
	d.Set("handler", function.Configuration.Handler)
	d.Set("runtime", function.Configuration.Runtime)
	d.Set("memory_size", function.Configuration.MemorySize)
	d.Set("commit_id", function.Configuration.CommitID)
	d.Set("role", function.Configuration.Role)
	d.Set("source_type", function.Configuration.SourceTag)
	d.Set("log_type", function.Configuration.LogType)
	if function.Configuration.LogType == "bos" {
		if strings.HasPrefix(function.Configuration.LogBosDir, "bos://") {
			d.Set("log_bos_dir", function.Configuration.LogBosDir[6:])
		} else {
			d.Set("log_bos_dir", function.Configuration.LogBosDir)
		}
	}

	if function.Configuration.Environment != nil {
		d.Set("environment", function.Configuration.Environment.Variables)
	}

	codeStorage := make(map[string]string)
	codeStorage["location"] = function.Code.Location
	codeStorage["repository_type"] = function.Code.RepositoryType
	d.Set("code_storage", codeStorage)

	if function.Configuration.VpcConfig != nil {
		d.Set("vpc_config", []interface{}{cfcService.faltternCFCFunctionVpcConfigToMap(function.Configuration.VpcConfig)})
	}

	d.SetId(function.Configuration.FunctionName)

	return nil
}
