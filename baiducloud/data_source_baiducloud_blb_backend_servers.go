/*
Use this data source to query BLB Backend Server list.

Example Usage

```hcl
data "baiducloud_blb_backend_servers" "default" {
 blb_id = "xxxx"
}

output "server_groups" {
 value = "${data.baiducloud_blb_backend_servers.default.backend_server_list}"
}
```
*/
package baiducloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudBLBBackendServer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudBLBBackendServersRead,

		Schema: map[string]*schema.Schema{
			"blb_id": {
				Type:        schema.TypeString,
				Description: "ID of the LoadBalance instance to be queried",
				Required:    true,
				ForceNew:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "Query result output file path",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),

			"backend_server_list": {
				Type:        schema.TypeList,
				Description: "backend server list",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Description: "Backend server instance ID",
							Computed:    true,
						},
						"weight": {
							Type:        schema.TypeInt,
							Description: "Backend server instance weight in this group, range from 0-100",
							Computed:    true,
						},
						"private_ip": {
							Type:        schema.TypeString,
							Description: "Backend server instance bind private ip",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudBLBBackendServersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	blbService := BLBService{client}

	blbId := ""
	if v, ok := d.GetOk("blb_id"); ok && v.(string) != "" {
		blbId = v.(string)
	}

	action := "Query BLB " + blbId + "Backend Server  "
	serverList, err := blbService.BackendServerList(blbId)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb_backend_servers", action, BCESDKGoERROR)
	}
	addDebug(action, serverList)
	FilterDataSourceResult(d, &serverList)

	if err := d.Set("backend_server_list", serverList); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb_backend_servers", action, BCESDKGoERROR)
	}

	d.SetId(resource.UniqueId())
	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), serverList); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb_backend_servers", action, BCESDKGoERROR)
		}
	}

	return nil
}
