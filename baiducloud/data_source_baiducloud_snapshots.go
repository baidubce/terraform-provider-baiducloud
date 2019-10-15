/*
Use this data source to query Snapshot list.

Example Usage

```hcl
data "baiducloud_snapshots" "default" {}

output "snapshots" {
 value = "${data.baiducloud_snapshots.default.snapshots}"
}
```
*/
package baiducloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudSnapshots() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudSnapshotsRead,

		Schema: map[string]*schema.Schema{
			"volume_id": {
				Type:        schema.TypeString,
				Description: "Volume ID to be attached of snapshots, if volume is system disk, volume ID is instance ID",
				Optional:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "Snapshots search result output file.",
				Optional:    true,
				ForceNew:    true,
			},
			"snapshots": {
				Type:        schema.TypeList,
				Description: "The result of the snapshots list.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "The ID of the snapshot.",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the snapshot.",
							Computed:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "The description of the snapshot.",
							Computed:    true,
						},
						"volume_id": {
							Type:        schema.TypeString,
							Description: "The volume id of the snapshot.",
							Computed:    true,
						},
						"create_time": {
							Type:        schema.TypeString,
							Description: "The creation time of the snapshot.",
							Computed:    true,
						},
						"size_in_gb": {
							Type:        schema.TypeInt,
							Description: "The size of the snapshot in GB.",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "The status of the snapshot.",
							Computed:    true,
						},
						"create_method": {
							Type:        schema.TypeString,
							Description: "The creation method of the snapshot.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudSnapshotsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bccService := BccService{client}

	volumeId := ""
	if v, ok := d.GetOk("volume_id"); ok && v.(string) != "" {
		volumeId = v.(string)
	}

	action := "Data Source Query All Snapshots"
	snapshotList, err := bccService.ListAllSnapshots(volumeId)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_snapshots", action, BCESDKGoERROR)
	}

	snapshotMap := bccService.FlattenSnapshotModelToMap(snapshotList)
	if err := d.Set("snapshots", snapshotMap); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_snapshots", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), snapshotMap); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_snapshots", action, BCESDKGoERROR)
		}
	}

	return nil
}
