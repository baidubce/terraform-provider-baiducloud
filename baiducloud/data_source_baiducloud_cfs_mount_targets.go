/*
Use this data source to get CFS mount target list.

Example Usage

```hcl
data "baiducloud_cfs_mount_targets" "default" {
  fs_id = "cfs-xxxxxxxxxxx"
}

output "cfss" {
 value = "${baiducloud_cfs_mount_targets.default}"
}
```
*/
package baiducloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudCfsMountTargets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudCfsMountTargetsRead,

		Schema: map[string]*schema.Schema{
			"fs_id": {
				Type:        schema.TypeString,
				Description: "CFS ID which you want query",
				Required:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),
			"output_file": {
				Type:        schema.TypeString,
				Description: "CFS search result output file",
				Optional:    true,
				ForceNew:    true,
			},
			"mount_targets": {
				Type:        schema.TypeList,
				Description: "Mount targets info list.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:        schema.TypeString,
							Description: "Subnet ID which mount target belong to.",
							Computed:    true,
						},
						"domain": {
							Type:        schema.TypeString,
							Description: "Domain of mount target.",
							Computed:    true,
						},
						"mount_id": {
							Type:        schema.TypeString,
							Description: "ID of the mount target",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudCfsMountTargetsRead(d *schema.ResourceData, meta interface{}) error {
	action := "Get CFS mount target list"
	client := meta.(*connectivity.BaiduClient)
	cfsService := CfsService{Client: client}
	var fsId string
	if v, ok := d.GetOk("fs_id"); ok && v.(string) != "" {
		fsId = v.(string)
	}
	cfsMountTargetList, err := cfsService.ListCfsMountTarget(fsId)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfs_mount_targets", action, BCESDKGoERROR)
	}
	mountTargetMap := make([]map[string]interface{}, 0)
	for _, model := range cfsMountTargetList {
		mountTargetMap = append(mountTargetMap, map[string]interface{}{
			"subnet_id": model.SubnetID,
			"domain":    model.Domain,
			"mount_id":  model.MountID,
		})
	}

	FilterDataSourceResult(d, &mountTargetMap)
	if err = d.Set("mount_targets", mountTargetMap); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfs_mount_targets", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), mountTargetMap); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfs_mount_targets", action, BCESDKGoERROR)
		}
	}
	return nil
}
