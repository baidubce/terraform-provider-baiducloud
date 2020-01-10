/*
Use this data source to query CDS list.

Example Usage

```hcl
data "baiducloud_cdss" "default" {}

output "cdss" {
 value = "${data.baiducloud_cdss.default.cdss}"
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudCDSs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudCDSsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Description: "CDS volume bind instance ID",
				Optional:    true,
				ForceNew:    true,
			},
			"zone_name": {
				Type:        schema.TypeString,
				Description: "CDS volume zone name",
				Optional:    true,
				ForceNew:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "CDS volume search result output file",
				Optional:    true,
				ForceNew:    true,
			},

			"cdss": {
				Type:        schema.TypeList,
				Description: "CDS volume list",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cds_id": {
							Type:        schema.TypeString,
							Description: "CDS volume id",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "CDS disk name",
							Computed:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "CDS description",
							Computed:    true,
						},
						"disk_size_in_gb": {
							Type:        schema.TypeInt,
							Description: "CDS disk size, should in [1, 32765], when snapshot_id not set, this parameter is required.",
							Computed:    true,
						},
						"storage_type": {
							Type:        schema.TypeString,
							Description: "CDS dist storage type, support hp1 and std1, default hp1",
							Computed:    true,
						},
						"payment_timing": {
							Type:        schema.TypeString,
							Description: "payment method, support Prepaid or Postpaid",
							Computed:    true,
						},
						"zone_name": {
							Type:        schema.TypeString,
							Description: "Zone name",
							Computed:    true,
						},
						"create_time": {
							Type:        schema.TypeString,
							Description: "CDS disk create time",
							Computed:    true,
						},
						"expire_time": {
							Type:        schema.TypeString,
							Description: "CDS disk expire time",
							Computed:    true,
						},
						"type": {
							Type:        schema.TypeString,
							Description: "CDS disk type",
							Computed:    true,
						},
						"is_system_volume": {
							Type:        schema.TypeBool,
							Description: "CDS disk is system volume or not",
							Computed:    true,
						},
						"source_snapshot_id": {
							Type:        schema.TypeString,
							Description: "CDS disk create source snapshot id",
							Computed:    true,
						},
						"snapshot_num": {
							Type:        schema.TypeString,
							Description: "CDS disk snapshot num",
							Computed:    true,
						},
						"region_id": {
							Type:        schema.TypeString,
							Description: "CDS disk region id",
							Computed:    true,
						},
						"attachments": {
							Type:        schema.TypeList,
							Description: "CDS volume attachments",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"volume_id": {
										Type:        schema.TypeString,
										Description: "CDS attachment volume id",
										Computed:    true,
									},
									"instance_id": {
										Type:        schema.TypeString,
										Description: "CDS attachment instance id",
										Computed:    true,
									},
									"device": {
										Type:        schema.TypeString,
										Description: "CDS attachment device path",
										Computed:    true,
									},
									"serial": {
										Type:        schema.TypeString,
										Description: "CDS attachment serial",
										Computed:    true,
									},
								},
							},
						},
						"auto_snapshot_policy": {
							Type:        schema.TypeList,
							Description: "CDS volume bind auto snapshot policy info",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Description: "Auto Snapshot Policy ID",
										Computed:    true,
									},
									"name": {
										Type:        schema.TypeString,
										Description: "Auto Snapshot Policy name",
										Computed:    true,
									},
									"time_points": {
										Type:        schema.TypeList,
										Description: "Auto Snapshot Policy set snapshot create time points",
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"repeat_weekdays": {
										Type:        schema.TypeList,
										Description: "Auto Snapshot Policy repeat weekdays",
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"status": {
										Type:        schema.TypeString,
										Description: "Auto Snapshot Policy status",
										Computed:    true,
									},
									"retention_days": {
										Type:        schema.TypeInt,
										Description: "Auto Snapshot Policy retention days",
										Computed:    true,
									},
									"created_time": {
										Type:        schema.TypeString,
										Description: "Auto Snapshot Policy created time",
										Computed:    true,
									},
									"updated_time": {
										Type:        schema.TypeString,
										Description: "Auto Snapshot Policy updated time",
										Computed:    true,
									},
									"deleted_time": {
										Type:        schema.TypeString,
										Description: "Auto Snapshot Policy deleted time",
										Computed:    true,
									},
									"last_execute_time": {
										Type:        schema.TypeString,
										Description: "Auto Snapshot Policy last execute time",
										Computed:    true,
									},
									"volume_count": {
										Type:        schema.TypeInt,
										Description: "Auto Snapshot Policy volume count",
										Computed:    true,
									},
								},
							},
						},
						"status": {
							Type:        schema.TypeString,
							Description: "CDS volume status",
							Computed:    true,
						},
						"tags": tagsComputedSchema(),
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudCDSsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bccService := BccService{client}

	listArgs := &api.ListCDSVolumeArgs{}
	if v, ok := d.GetOk("instance_id"); ok && v.(string) != "" {
		listArgs.InstanceId = v.(string)
	}
	if v, ok := d.GetOk("zone_name"); ok && v.(string) != "" {
		listArgs.ZoneName = v.(string)
	}

	action := "Query all CDS volume detail"
	cdsList, err := bccService.ListAllCDSVolumeDetail(listArgs)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cdss", action, BCESDKGoERROR)
	}

	cdsMap := bccService.FlattenCDSVolumeModelToMap(cdsList)
	if err := d.Set("cdss", cdsMap); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cdss", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), cdsMap); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cdss", action, BCESDKGoERROR)
		}
	}

	return nil
}
