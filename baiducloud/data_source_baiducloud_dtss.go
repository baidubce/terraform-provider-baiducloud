/*
Use this data source to query DTS list.

Example Usage

```hcl
data "baiducloud_dtss" "default" {}

output "dtss" {
 value = "${data.baiducloud_dtss.default.dtss}"
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/dts"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudDtss() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudDtssRead,

		Schema: map[string]*schema.Schema{
			"dts_name": {
				Type:        schema.TypeString,
				Description: "Name of the Dts to be queried",
				Optional:    true,
				ForceNew:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "Dtss search result output file",
				Optional:    true,
				ForceNew:    true,
			},
			"type": {
				Type:         schema.TypeString,
				Description:  "type of the task. Available values are migration, sync, subscribe.",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"migration", "sync", "subscribe"}, false),
			},
			"filter": dataSourceFiltersSchema(),
			"dtss": {
				Type:        schema.TypeList,
				Description: "Dts list",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dts_id": {
							Type:        schema.TypeString,
							Description: "Dts task id",
							Optional:    true,
							Computed:    true,
						},
						"task_name": {
							Type:        schema.TypeString,
							Description: "Dts task name",
							Optional:    true,
							Computed:    true,
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
							Computed:    true,
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
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Description: "type",
										Computed:    true,
									},
									"src": {
										Type:        schema.TypeString,
										Description: "src",
										Computed:    true,
									},
									"dst": {
										Type:        schema.TypeString,
										Description: "dst",
										Computed:    true,
									},
									"where": {
										Type:        schema.TypeString,
										Description: "where",
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
							Description: "Dts errmsg",
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
						"product_type": {
							Type:        schema.TypeString,
							Description: "Dts product type",
							Computed:    true,
						},
						"source_instance_type": {
							Type:        schema.TypeString,
							Description: "Dts source instance type",
							Computed:    true,
						},
						"target_instance_type": {
							Type:        schema.TypeString,
							Description: "Dts target instance type",
							Computed:    true,
						},
						"cross_region_tag": {
							Type:        schema.TypeInt,
							Description: "Dts cross region tag",
							Computed:    true,
						},
						"pay_create_time": {
							Type:        schema.TypeInt,
							Description: "Dts pay create time",
							Computed:    true,
						},
						"standard": {
							Type:        schema.TypeString,
							Description: "Dts standard",
							Computed:    true,
						},
						"pay_end_time": {
							Type:        schema.TypeString,
							Description: "Dts pay end time",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudDtssRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	dtsService := DtsService{client}

	listArgs := &dts.ListDtsArgs{}
	if v, ok := d.GetOk("type"); ok && v.(string) != "" {
		listArgs.Type = v.(string)
	}

	action := "List all Dts tasks"
	dtsList, err := dtsService.ListAllDtss(listArgs)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dtss", action, BCESDKGoERROR)
	}
	addDebug(action, dtsList)
	dtsName := ""
	if v, ok := d.GetOk("dts_name"); ok && v.(string) != "" {
		dtsName = v.(string)
	}
	dtsMap := make([]map[string]interface{}, 0, len(dtsList))
	for _, c := range dtsList {
		if dtsName != "" && dtsName != c.TaskName {
			continue
		}

		schema, base := flattenSchemaInfoToMap(c.DynamicInfo)

		dtsMap = append(dtsMap, map[string]interface{}{
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
			"schema":                schema,
			"base":                  base,
			"increment":             flattenIncrementToMap(c.DynamicInfo),
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

	FilterDataSourceResult(d, &dtsMap)

	if err := d.Set("dtss", dtsMap); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dtss", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), dtsMap); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_dtss", action, BCESDKGoERROR)
		}
	}

	return nil
}
