/*
Use this data source to query Auto Snapshot Policy list.

Example Usage

```hcl
data "baiducloud_auto_snapshot_policies" "default" {}

output "auto_snapshot_policiess" {
 value = "${data.baiducloud_auto_snapshot_policies.default.auto_snapshot_policies}"
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

func dataSourceBaiduCloudAutoSnapshotPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudAutoSnapshotPoliciesRead,

		Schema: map[string]*schema.Schema{
			"asp_name": {
				Type:        schema.TypeString,
				Description: "Name of the automatic snapshot policy.",
				Optional:    true,
			},
			"volume_name": {
				Type:        schema.TypeString,
				Description: "Name of the volume.",
				Optional:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "Automatic snapshot policies search result output file.",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),

			"auto_snapshot_policies": {
				Type:        schema.TypeList,
				Description: "The automatic snapshot policies search result list.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "The ID of the automatic snapshot policy.",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the automatic snapshot policy.",
							Computed:    true,
						},
						"time_points": {
							Type:        schema.TypeList,
							Description: "The time points of the automatic snapshot policy.",
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"repeat_weekdays": {
							Type:        schema.TypeList,
							Description: "The repeat weekdays of the automatic snapshot policy.",
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"status": {
							Type:        schema.TypeString,
							Description: "The status of the automatic snapshot policy.",
							Computed:    true,
						},
						"retention_days": {
							Type:        schema.TypeInt,
							Description: "The retention days of the automatic snapshot policy.",
							Computed:    true,
						},
						"created_time": {
							Type:        schema.TypeString,
							Description: "The creation time of the automatic snapshot policy.",
							Computed:    true,
						},
						"updated_time": {
							Type:        schema.TypeString,
							Description: "The updation time of the automatic snapshot policy.",
							Computed:    true,
						},
						"deleted_time": {
							Type:        schema.TypeString,
							Description: "The deletion time of the automatic snapshot policy.",
							Computed:    true,
						},
						"last_execute_time": {
							Type:        schema.TypeString,
							Description: "The last execution time of the automatic snapshot policy.",
							Computed:    true,
						},
						"volume_count": {
							Type:        schema.TypeInt,
							Description: "The volume count of the automatic snapshot policy.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudAutoSnapshotPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bccService := BccService{client}

	listAspArgs := &api.ListASPArgs{}
	if v, ok := d.GetOk("asp_name"); ok && v.(string) != "" {
		listAspArgs.AspName = v.(string)
	}
	if v, ok := d.GetOk("volume_name"); ok && v.(string) != "" {
		listAspArgs.VolumeName = v.(string)
	}

	action := "Data Source Query All Auto Snapshot Policies"
	aspList, err := bccService.ListAllAutoSnapshotPolicies(listAspArgs)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_auto_snapshot_policies", action, BCESDKGoERROR)
	}

	aspMap := bccService.FlattenAutoSnapshotPolicyModelToMap(aspList)

	FilterDataSourceResult(d, &aspMap)

	if err := d.Set("auto_snapshot_policies", aspMap); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_auto_snapshot_policies", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), aspMap); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_auto_snapshot_policies", action, BCESDKGoERROR)
		}
	}

	return nil
}
