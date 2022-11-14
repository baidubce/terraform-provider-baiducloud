/*
Use this data source to query bls log stores .

Example Usage

```hcl
data "baiducloud_bls_log_stores" "default" {

}

output "log_stores" {
 	value = "${data.baiducloud_bls_log_stores.default.log_stores}"
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/bls/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudBLSLogStores() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudBLSLogStoresRead,

		Schema: map[string]*schema.Schema{
			"name_pattern": {
				Type:        schema.TypeString,
				Description: "Log store namePattern",
				Optional:    true,
				ForceNew:    true,
			},
			"order": {
				Type:        schema.TypeString,
				Description: "search order",
				Optional:    true,
				ForceNew:    true,
			},
			"order_by": {
				Type:        schema.TypeString,
				Description: "order field",
				Optional:    true,
				ForceNew:    true,
			},
			"page_no": {
				Type:        schema.TypeInt,
				Description: "number of page ",
				Optional:    true,
				ForceNew:    true,
			},
			"page_size": {
				Type:        schema.TypeInt,
				Description: "size of page",
				Optional:    true,
				ForceNew:    true,
			},
			"total_count": {
				Type:        schema.TypeInt,
				Description: "Total number of items",
				Optional:    true,
				Computed:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "log stores search result output file",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),

			"log_stores": {
				Type:        schema.TypeList,
				Description: "log store list",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"log_store_name": {
							Type:        schema.TypeString,
							Description: "name of log store",
							Required:    true,
							ForceNew:    true,
						},
						"retention": {
							Type:        schema.TypeInt,
							Description: "retention days of log store",
							Required:    true,
							ForceNew:    true,
						},
						"creation_date_time": {
							Type:        schema.TypeString,
							Description: "log store create date time",
							Computed:    true,
						},
						"last_modified_time": {
							Type:        schema.TypeString,
							Description: "log store last modified time",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudBLSLogStoresRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	blsService := BLSService{client}

	listArgs := &api.QueryConditions{}
	if v, ok := d.GetOk("name_pattern"); ok {
		listArgs.NamePattern = v.(string)
	}
	if v, ok := d.GetOk("order"); ok {
		listArgs.Order = v.(string)
	}

	if v, ok := d.GetOk("order_by"); ok {
		listArgs.OrderBy = v.(string)
	}

	if v, ok := d.GetOk("page_no"); ok {
		listArgs.PageNo = v.(int)
	}

	if v, ok := d.GetOk("page_size"); ok {
		listArgs.PageSize = v.(int)
	}

	action := "Query bls log stores "
	logStoreList, err := blsService.GetBLSLogStoreList(listArgs)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bls_log_stores", action, BCESDKGoERROR)
	}
	addDebug(action, logStoreList)

	blsStoreMap := blsService.FlattenLogStoreModelsToMap(logStoreList.Result)

	FilterDataSourceResult(d, &blsStoreMap)

	if err := d.Set("total_count", logStoreList.TotalCount); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bls_log_stores", action, BCESDKGoERROR)
	}
	if err := d.Set("log_stores", blsStoreMap); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bls_log_stores", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), blsStoreMap); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bls_log_stores", action, BCESDKGoERROR)
		}
	}

	return nil
}
