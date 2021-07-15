/*
Use this data source to query BCC Instance list.

Example Usage

```hcl
data "baiducloud_bbc_instances" "default" {}

output "instances" {
 value = "${data.baiducloud_bbc_instances.default.instances}"
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/bbc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudBbcInstances() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceBaiduCloudBbcInstancesRead,

		Schema: map[string]*schema.Schema{
			"internal_ip": {
				Type:         schema.TypeString,
				Description:  "Internal ip address of the instance to retrieve.",
				Optional:     true,
				ValidateFunc: validation.SingleIP(),
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "id of vpc to which the instance belongs.",
				Optional:    true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Description: "id of bbc instance",
				Optional:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "Output file of the instances search result",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),
			"instances": {
				Type:        schema.TypeList,
				Description: "The result of the instances list.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Description: "The ID of the instance.",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the instance.",
							Computed:    true,
						},
						"flavor_id": {
							Type:        schema.TypeString,
							Description: "flavor id",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "The status of the instance.",
							Computed:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "The description of the instance.",
							Computed:    true,
						},
						"payment_timing": {
							Type:        schema.TypeString,
							Description: "The payment timing of the instance.",
							Computed:    true,
						},
						"create_time": {
							Type:        schema.TypeString,
							Description: "The creation time of the instance.",
							Computed:    true,
						},
						"expire_time": {
							Type:        schema.TypeString,
							Description: "The expire time of the instance.",
							Computed:    true,
						},
						"internal_ip": {
							Type:        schema.TypeString,
							Description: "The internal ip of the instance.",
							Computed:    true,
						},
						"public_ip": {
							Type:        schema.TypeString,
							Description: "The public ip of the instance.",
							Computed:    true,
						},
						"cpu_count": {
							Type:        schema.TypeInt,
							Description: "The cpu count of the instance.",
							Computed:    true,
						},
						"memory_capacity_in_gb": {
							Type:        schema.TypeInt,
							Description: "The memory capacity in GB of the instance.",
							Computed:    true,
						},
						"cds_disks": {
							Type:        schema.TypeList,
							Description: "CDS disks of the instance.",
							Computed:    true,
							MinItems:    1,
							MaxItems:    10,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_size_in_gb": {
										Type:         schema.TypeInt,
										Description:  "The size(GB) of CDS.",
										Optional:     true,
										Default:      0,
										ValidateFunc: validation.IntAtLeast(0),
									},
									"storage_type": {
										Type:         schema.TypeString,
										Description:  "Storage type of the CDS.",
										Optional:     true,
										Default:      bbc.StorageTypeCloudHP1,
										ValidateFunc: validateStorageType(),
									},
									"is_system_volume": {
										Type:        schema.TypeBool,
										Description: "Snapshot ID of CDS.",
										Optional:    true,
									},
								},
							},
						},
						"image_id": {
							Type:        schema.TypeString,
							Description: "The image id of the instance.",
							Computed:    true,
						},
						"network_capacity_in_mbps": {
							Type:        schema.TypeInt,
							Description: "The network capacity in Mbps of the instance.",
							Computed:    true,
						},
						"zone_name": {
							Type:        schema.TypeString,
							Description: "The zone name of the instance.",
							Computed:    true,
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Description: "The subnet ID of the instance.",
							Computed:    true,
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Description: "The VPC ID of the instance.",
							Computed:    true,
						},
						"tags": tagsComputedSchema(),
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudBbcInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bbcService := BbcService{client}
	var instanceList []bbc.InstanceModel = make([]bbc.InstanceModel, 0)
	action := "List all Instance "

	if v, ok := d.GetOk("instance_id"); ok {
		raw, err := client.WithBbcClient(func(bbcClient *bbc.Client) (interface{}, error) {
			return bbcClient.GetInstanceDetail(v.(string))
		})
		action = "Get Instance by id"
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instances", action, BCESDKGoERROR)
		}
		model := raw.(*bbc.InstanceModel)
		instanceList = append(instanceList, *model)
	} else {
		listArgs := &bbc.ListInstancesArgs{}
		if v, ok := d.GetOk("internal_ip"); ok {
			listArgs.InternalIp = v.(string)
		}
		if v, ok := d.GetOk("vpc_id"); ok && v.(string) != "" {
			listArgs.VpcId = v.(string)
		}

		list, err := bbcService.ListAllInstance(listArgs)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instances", action, BCESDKGoERROR)
		}
		instanceList = list
	}

	instanceMap, err := bbcService.FlattenInstanceModelToMap(instanceList)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instances", action, BCESDKGoERROR)
	}

	FilterDataSourceResult(d, &instanceMap)
	if err = d.Set("instances", instanceMap); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instances", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), instanceMap); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instances", action, BCESDKGoERROR)
		}
	}

	return nil
}
