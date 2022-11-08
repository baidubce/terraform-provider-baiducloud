/*
Use this data source to get CFS list.

Example Usage

```hcl
data baiducloud_cfss "default" {

}

output "cfss" {
 value = "${data.baiducloud_cfss_.default}"
}
```
*/
package baiducloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudCfss() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudCfssRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"output_file": {
				Type:        schema.TypeString,
				Description: "CFS search result output file",
				Optional:    true,
				ForceNew:    true,
			},
			"cfss": {
				Type:        schema.TypeList,
				Description: "cfs list",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fs_id": {
							Type:        schema.TypeString,
							Description: "ID of the CFS.",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the CFS.",
							Computed:    true,
						},
						"type": {
							Type:        schema.TypeString,
							Description: "CFS type, default is cap.",
							Computed:    true,
						},
						"protocol": {
							Type:        schema.TypeString,
							Description: "CFS protocol, available value is nfs and smb, default is nfs.",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "CFS status, available value is available,updating,paused and unavailable.",
							Computed:    true,
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Description: "VPC ID.",
							Computed:    true,
						},
						"mount_target_list": {
							Type:        schema.TypeList,
							Description: "Name of the deployset.",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"subnet_id": {
										Type:        schema.TypeString,
										Description: "ID of subnet which mount target bind.",
										Computed:    true,
									},
									"domain": {
										Type:        schema.TypeString,
										Description: "Domain of the mount target.",
										Computed:    true,
									},
									"mount_id": {
										Type:        schema.TypeString,
										Description: "ID of the mount target.",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudCfssRead(d *schema.ResourceData, meta interface{}) error {
	action := "Get cfs list"
	client := meta.(*connectivity.BaiduClient)
	cfsService := CfsService{Client: client}
	cfsList, err := cfsService.ListCfs()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfs", action, BCESDKGoERROR)
	}
	cfsMap := make([]map[string]interface{}, 0)
	for _, model := range cfsList {
		tempMap := make(map[string]interface{})
		tempMap["fs_id"] = model.FSID
		tempMap["name"] = model.Name
		tempMap["type"] = model.Type
		tempMap["protocol"] = model.Protocol
		tempMap["status"] = model.Status
		tempMap["vpc_id"] = model.VpcID
		cfsMap = append(cfsMap, tempMap)
		mountTargetMap := make([]map[string]interface{}, 0)
		for _, target := range model.MoutTargets {
			mountTargetMap = append(mountTargetMap, map[string]interface{}{
				"subnet_id": target.SubnetID,
				"domain":    target.Domain,
				"mount_id":  target.MountID,
			})
		}
		tempMap["mount_target_list"] = mountTargetMap
	}
	FilterDataSourceResult(d, &cfsMap)
	if err = d.Set("cfss", cfsMap); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfss", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), cfsMap); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfss", action, BCESDKGoERROR)
		}
	}
	return nil
}
