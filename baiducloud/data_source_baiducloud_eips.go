/*
Use this data source to query EIP list.

Example Usage

```hcl
data "baiducloud_eips" "default" {}

output "eips" {
 value = "${data.baiducloud_eips.default.eips}"
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/eip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudEips() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudEipsRead,

		Schema: map[string]*schema.Schema{
			"eip": {
				Type:         schema.TypeString,
				Description:  "Eip address",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.SingleIP(),
			},
			"status": {
				Type:         schema.TypeString,
				Description:  "Eip status",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"available", "binded", "paused"}, false),
			},
			"instance_type": {
				Type:         schema.TypeString,
				Description:  "Eip bind instance type",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"BCC", "BLB", "VPN", "NAT"}, false),
			},
			"instance_id": {
				Type:        schema.TypeString,
				Description: "Eip bind instance id",
				Optional:    true,
				ForceNew:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "Eips search result output file",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),

			"eips": {
				Type:        schema.TypeList,
				Description: "Eip list",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"eip": {
							Type:        schema.TypeString,
							Description: "Eip address",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Eip name",
							Computed:    true,
						},
						"bandwidth_in_mbps": {
							Type:        schema.TypeInt,
							Description: "Eip bandwidth(Mbps)",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Eip status",
							Computed:    true,
						},
						"eip_instance_type": {
							Type:        schema.TypeString,
							Description: "Eip instance type",
							Computed:    true,
						},
						"share_group_id": {
							Type:        schema.TypeString,
							Description: "Eip share group id",
							Computed:    true,
						},
						"payment_timing": {
							Type:        schema.TypeString,
							Description: "Eip payment timing",
							Computed:    true,
						},
						"billing_method": {
							Type:        schema.TypeString,
							Description: "Eip billing method",
							Computed:    true,
						},
						"create_time": {
							Type:        schema.TypeString,
							Description: "Eip create time",
							Computed:    true,
						},
						"expire_time": {
							Type:        schema.TypeString,
							Description: "Eip expire time",
							Computed:    true,
						},
						"tags": tagsComputedSchema(),
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudEipsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	eipService := EipService{client}

	listArgs := &eip.ListEipArgs{}
	if v, ok := d.GetOk("eip"); ok {
		listArgs.Eip = v.(string)
	}
	if v, ok := d.GetOk("status"); ok {
		listArgs.Status = v.(string)
	}

	if v, ok := d.GetOk("instance_type"); ok {
		listArgs.InstanceType = v.(string)
	}
	if v, ok := d.GetOk("instance_id"); ok && v.(string) != "" {
		listArgs.InstanceId = v.(string)
	}

	action := "Query Eips " + listArgs.Eip
	eipList, err := eipService.ListAllEips(listArgs)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eips", action, BCESDKGoERROR)
	}
	addDebug(action, eipList)

	eipMap := eipService.FlattenEipModelsToMap(eipList)

	FilterDataSourceResult(d, &eipMap)

	if err := d.Set("eips", eipMap); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eips", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), eipMap); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eips", action, BCESDKGoERROR)
		}
	}

	return nil
}
