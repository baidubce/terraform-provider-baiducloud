/*
Use this data source to query Security Group list.

Example Usage

```hcl
data "baiducloud_security_groups" "default" {}

output "security_groups" {
 value = "${data.baiducloud_security_groups.default.security_groups}"
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

func dataSourceBaiduCloudSecurityGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudSecurityGroupsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Description: "Security Group attached instance ID",
				Optional:    true,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "Security Group attached vpc id",
				Optional:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "Security Group search result output file",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),

			"security_groups": {
				Type:        schema.TypeList,
				Description: "Security Groups search result",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "Security Group ID",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Security Group name",
							Computed:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "Security Group description",
							Computed:    true,
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Description: "Security Group vpc id",
							Computed:    true,
						},
						"tags": tagsComputedSchema(),
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudSecurityGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bccService := BccService{client}

	listSgArgs := &api.ListSecurityGroupArgs{}
	if v, ok := d.GetOk("instance_id"); ok && v.(string) != "" {
		listSgArgs.InstanceId = v.(string)
	}
	if v, ok := d.GetOk("vpc_id"); ok && v.(string) != "" {
		listSgArgs.VpcId = v.(string)
	}

	action := "Data Source Query All Security Groups"
	sgList, err := bccService.ListAllSecurityGroups(listSgArgs)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_security_groups", action, BCESDKGoERROR)
	}

	sgMap := bccService.FlattenSecurityGroupModelToMap(sgList)
	FilterDataSourceResult(d, &sgMap)

	if err := d.Set("security_groups", sgMap); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_security_groups", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), sgMap); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_security_groups", action, BCESDKGoERROR)
		}
	}

	return nil
}
