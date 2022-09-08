/*
Use this data source to query BBC Instance list.

Example Usage

```hcl
data "baiducloud_bbc_instances" "data_bbc_instance" {
  internal_ip = "172.16.16.4"
}

output "instances" {
 value = "${data.baiducloud_bbc_instances.data_bbc_instance.instances}"
}

```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/bbc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudBbcInstances() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceBaiduCloudBbcInstancesRead,

		Schema: map[string]*schema.Schema{
			"internal_ip": {
				Type:        schema.TypeString,
				Description: "internal ip.",
				Optional:    true,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "bbc vpc id.",
				Optional:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "Output file of the bbc instances search result",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),
			"instances": {
				Type:        schema.TypeList,
				Description: "The result of the bbc instances list.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Description: "The ID of the BBC instance.",
							Computed:    true,
						},
						"instance_name": {
							Type:        schema.TypeString,
							Description: "The name of the BBC instance.",
							Computed:    true,
						},
						"hostname": {
							Type:        schema.TypeString,
							Description: "The hostname of the BBC instance.",
							Computed:    true,
						},
						"uuid": {
							Type:        schema.TypeString,
							Description: "The UUID of the BBC instance.",
							Computed:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "The description of the BBC instance.",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "The status of the instance.Include starting running stopped deleted",
							Computed:    true,
						},
						"payment_timing": {
							Type:        schema.TypeString,
							Description: "Payment timing of billing, which can be Prepaid or Postpaid. The default is Postpaid.",
							Computed:    true,
						},
						"create_time": {
							Type:        schema.TypeString,
							Description: "BBC create time.",
							Computed:    true,
						},
						"expire_time": {
							Type:        schema.TypeString,
							Description: "BBC expire time.",
							Computed:    true,
						},
						"public_ip": {
							Type:        schema.TypeString,
							Description: "The public IP of the BBC instance.",
							Computed:    true,
						},
						"internal_ip": {
							Type:        schema.TypeString,
							Description: "The internal IP of the BBC instance.",
							Computed:    true,
						},
						"rdma_ip": {
							Type:        schema.TypeString,
							Description: "The rdma IP of the BBC instance.",
							Computed:    true,
						},
						"image_id": {
							Type:        schema.TypeString,
							Description: "The image ID of the BBC instance.",
							Computed:    true,
						},
						"flavor_id": {
							Type:        schema.TypeString,
							Description: "The flavor ID of the BBC instance.",
							Computed:    true,
						},
						"zone": {
							Type:        schema.TypeString,
							Description: "The zone name of the BBC instance.",
							Computed:    true,
						},
						"region": {
							Type:        schema.TypeString,
							Description: "The region of the BBC instance.",
							Computed:    true,
						},
						"has_alive": {
							Type:        schema.TypeInt,
							Description: "BBC instance has alive.",
							Computed:    true,
						},
						"tags": tagsSchema(),
						"switch_id": {
							Type:        schema.TypeString,
							Description: "The switch ID of the BBC instance.",
							Computed:    true,
						},
						"host_id": {
							Type:        schema.TypeString,
							Description: "The host ID of the BBC instance.",
							Computed:    true,
						},
						"deployset_id": {
							Type:        schema.TypeString,
							Description: "The deployset ID of the BBC instance.",
							Computed:    true,
						},
						"network_capacity_in_mbps": {
							Type:        schema.TypeInt,
							Description: "network capacity in mbps.",
							Computed:    true,
						},
						"rack_id": {
							Type:        schema.TypeString,
							Description: "The rack ID of the BBC instance.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
func dataSourceBaiduCloudBbcInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bbcService := BbcService{client}

	listArgs := &bbc.ListInstancesArgs{}
	if v, ok := d.GetOk("internal_ip"); ok {
		listArgs.InternalIp = v.(string)
	}
	if v, ok := d.GetOk("vpc_id"); ok && v.(string) != "" {
		listArgs.VpcId = v.(string)
	}

	action := "List all bbc Instance "
	instanceList, err := bbcService.ListAllBbcInstances(listArgs)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bbc_instances", action, BCESDKGoERROR)
	}
	instanceMap, err := bbcService.FlattenBbcInstanceModelToMap(instanceList)
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
