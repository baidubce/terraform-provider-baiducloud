/*
Use this data source to query BLB list.

Example Usage

```hcl
data "baiducloud_blbs" "default" {
 name = "myLoadBalance"
}

output "blbs" {
 value = "${data.baiducloud_blbs.default.blbs}"
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/blb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudBLBs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudBLBRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the LoadBalance instance to be queried",
				Optional:    true,
			},
			"address": {
				Type:         schema.TypeString,
				Description:  "Address ip of the LoadBalance instance to be queried",
				Optional:     true,
				ValidateFunc: validation.SingleIP(),
			},
			"blb_id": {
				Type:        schema.TypeString,
				Description: "ID of the LoadBalance instance to be queried",
				Optional:    true,
			},
			"bcc_id": {
				Type:        schema.TypeString,
				Description: "ID of the BCC instance bound to the LoadBalance",
				Optional:    true,
			},
			"exactly_match": {
				Type:        schema.TypeBool,
				Description: "Whether the query condition is an exact match or not, default false",
				Optional:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "Query result output file path",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),

			// Attributes used for result
			"blbs": {
				Type:        schema.TypeList,
				Description: "A list of lication LoadBalance Instance",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"blb_id": {
							Type:        schema.TypeString,
							Description: "LoadBalance instance's ID",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "LoadBalance instance's name",
							Computed:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "LoadBalance instance's description",
							Computed:    true,
						},
						"address": {
							Type:        schema.TypeString,
							Description: "LoadBalance instance's service IP, instance can be accessed through this IP",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "LoadBalance instance's status",
							Computed:    true,
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Description: "The VPC short ID to which the LoadBalance instance belongs",
							Computed:    true,
						},
						"vpc_name": {
							Type:        schema.TypeString,
							Description: "The VPC name to which the LoadBalance instance belongs",
							Computed:    true,
						},
						"public_ip": {
							Type:        schema.TypeString,
							Description: "LoadBalance instance's public ip",
							Computed:    true,
						},
						"cidr": {
							Type:        schema.TypeString,
							Description: "Cidr of the network where the LoadBalance instance reside",
							Computed:    true,
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Description: "The subnet ID to which the LoadBalance instance belongs",
							Computed:    true,
						},
						"create_time": {
							Type:        schema.TypeString,
							Description: "LoadBalance instance's create time",
							Computed:    true,
						},
						"listener": {
							Type:        schema.TypeList,
							Description: "List of listeners mounted under the instance",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:        schema.TypeInt,
										Description: "Listening port",
										Computed:    true,
									},
									"type": {
										Type:        schema.TypeString,
										Description: "Listening protocol type",
										Computed:    true,
									},
								},
							},
						},
						"tags": tagsComputedSchema(),
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudBLBRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	blbService := BLBService{client}

	args := &blb.DescribeLoadBalancersArgs{}
	if v, ok := d.GetOk("blb_id"); ok && v.(string) != "" {
		args.BlbId = v.(string)
	}
	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		args.Name = v.(string)
	}
	if v, ok := d.GetOk("bcc_id"); ok && v.(string) != "" {
		args.BccId = v.(string)
	}
	if v, ok := d.GetOk("exactly_match"); ok {
		args.ExactlyMatch = v.(bool)
	}
	if v, ok := d.GetOk("address"); ok && v.(string) != "" {
		args.Address = v.(string)
	}

	action := "Query BLB " + args.BlbId + "_" + args.Name
	blbModels, blbDetails, err := blbService.ListAllBLB(args)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blbs", action, BCESDKGoERROR)
	}

	blbMap := blbService.FlattenBLBDetailsToMap(blbModels, blbDetails)

	FilterDataSourceResult(d, &blbMap)

	if err := d.Set("blbs", blbMap); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blbs", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), blbMap); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blbs", action, BCESDKGoERROR)
		}
	}

	return nil
}
