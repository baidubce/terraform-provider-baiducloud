/*
Use this data source to query BCC Instance list.

Example Usage

```hcl
data "baiducloud_instances" "default" {}

output "instances" {
 value = "${data.baiducloud_instances.default.instances}"
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudInstancesRead,

		Schema: map[string]*schema.Schema{
			"internal_ip": {
				Type:         schema.TypeString,
				Description:  "Internal ip address of the instance to retrieve.",
				Optional:     true,
				ValidateFunc: validation.SingleIP(),
			},
			"dedicated_host_id": {
				Type:        schema.TypeString,
				Description: "Dedicated host id of the instance to retrieve.",
				Optional:    true,
			},
			"zone_name": {
				Type:        schema.TypeString,
				Description: "Name of the available zone to which the instance belongs.",
				Optional:    true,
			},
			"keypair_id": {
				Type:        schema.TypeString,
				Description: "Keypair ID of the instance.",
				Optional:    true,
			},
			"auto_renew": {
				Type:        schema.TypeBool,
				Description: "Whether to renew automatically.",
				Optional:    true,
			},
			"instance_ids": {
				Type:        schema.TypeString,
				Description: "Multiple instance IDs, separated by commas.",
				Optional:    true,
			},
			"instance_names": {
				Type:        schema.TypeString,
				Description: "Multiple instance names, separated by commas.",
				Optional:    true,
			},
			"cds_ids": {
				Type:        schema.TypeString,
				Description: "Multiple cds disk IDs, separated by commas.",
				Optional:    true,
			},
			"deploy_set_ids": {
				Type:        schema.TypeString,
				Description: "Multiple deployment set IDs, separated by commas.",
				Optional:    true,
			},
			"security_group_ids": {
				Type:        schema.TypeString,
				Description: "Multiple security group IDs, separated by commas.",
				Optional:    true,
			},
			"payment_timing": {
				Type:        schema.TypeString,
				Description: "Payment method. Valid values: `Prepaid`, `Postpaid`.",
				Optional:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Instance status. Valid values: `Recycled`, `Running`, `Stopped`, `Stopping`, `Starting`.",
				Optional:    true,
			},
			"tags": {
				Type:        schema.TypeString,
				Description: "Multiple tags, separated by commas. Format: `tagKey:tagValue` or `tagKey`.",
				Optional:    true,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "Can only be used in combination with the `private_ips` query parameter.",
				Optional:    true,
			},
			"private_ips": {
				Type:        schema.TypeString,
				Description: "Multiple intranet IPs, separated by commas. Must be used in combination with `vpc_id`.",
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
						"instance_type": {
							Type:        schema.TypeString,
							Description: "The type of the instance.",
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
						"gpu_card": {
							Type:        schema.TypeString,
							Description: "The gpu card of the instance.",
							Computed:    true,
						},
						"fpga_card": {
							Type:        schema.TypeString,
							Description: "The fgpa card of the instance.",
							Computed:    true,
						},
						"card_count": {
							Type:        schema.TypeString,
							Description: "The card count of the instance.",
							Computed:    true,
						},
						"memory_capacity_in_gb": {
							Type:        schema.TypeInt,
							Description: "The memory capacity in GB of the instance.",
							Computed:    true,
						},
						"root_disk_size_in_gb": {
							Type:        schema.TypeInt,
							Description: "The system disk size in GB of the instance.",
							Computed:    true,
						},
						"root_disk_storage_type": {
							Type:        schema.TypeString,
							Description: "The system disk storage type of the instance.",
							Computed:    true,
						},
						"ephemeral_disks": {
							Type:        schema.TypeList,
							Description: "The ephemeral disks of the instance.",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"size_in_gb": {
										Type:        schema.TypeInt,
										Description: "The size(GB) of the ephemeral disk.",
										Computed:    true,
									},
									"storage_type": {
										Type:        schema.TypeString,
										Description: "The storage type of the ephemeral disk.",
										Computed:    true,
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
						"placement_policy": {
							Type:        schema.TypeString,
							Description: "The placement policy of the instance.",
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
						"dedicated_host_id": {
							Type:        schema.TypeString,
							Description: "The dedicated host id of the instance.",
							Computed:    true,
						},
						"auto_renew": {
							Type:        schema.TypeBool,
							Description: "Whether to automatically renew.",
							Computed:    true,
						},
						"keypair_id": {
							Type:        schema.TypeString,
							Description: "The key pair id of the instance.",
							Computed:    true,
						},
						"keypair_name": {
							Type:        schema.TypeString,
							Description: "The key pair name of the instance.",
							Computed:    true,
						},
						"instance_spec": {
							Type:        schema.TypeString,
							Description: "spec",
							Computed:    true,
						},
						"deploy_set_ids": {
							Type:        schema.TypeSet,
							Description: "Deploy set ids the instance belong to",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"tags": tagsComputedSchema(),
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bccService := BccService{client}

	listArgs := &api.ListInstanceArgs{}
	if v, ok := d.GetOk("internal_ip"); ok {
		listArgs.InternalIp = v.(string)
	}
	if v, ok := d.GetOk("dedicated_host_id"); ok && v.(string) != "" {
		listArgs.DedicatedHostId = v.(string)
	}
	if v, ok := d.GetOk("zone_name"); ok && v.(string) != "" {
		listArgs.ZoneName = v.(string)
	}

	if v, ok := d.GetOk("keypair_id"); ok && v.(string) != "" {
		listArgs.KeypairId = v.(string)
	}
	if v, ok := d.GetOk("auto_renew"); ok {
		listArgs.AutoRenew = v.(bool)
	}
	if v, ok := d.GetOk("instance_ids"); ok && v.(string) != "" {
		listArgs.InstanceIds = v.(string)
	}
	if v, ok := d.GetOk("instance_names"); ok && v.(string) != "" {
		listArgs.InstanceNames = v.(string)
	}
	if v, ok := d.GetOk("cds_ids"); ok && v.(string) != "" {
		listArgs.CdsIds = v.(string)
	}
	if v, ok := d.GetOk("deploy_set_ids"); ok && v.(string) != "" {
		listArgs.DeploySetIds = v.(string)
	}
	if v, ok := d.GetOk("security_group_ids"); ok && v.(string) != "" {
		listArgs.SecurityGroupIds = v.(string)
	}
	if v, ok := d.GetOk("payment_timing"); ok && v.(string) != "" {
		listArgs.PaymentTiming = v.(string)
	}
	if v, ok := d.GetOk("status"); ok && v.(string) != "" {
		listArgs.Status = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok && v.(string) != "" {
		listArgs.Tags = v.(string)
	}
	if v, ok := d.GetOk("vpc_id"); ok && v.(string) != "" {
		listArgs.VpcId = v.(string)
	}
	if v, ok := d.GetOk("private_ips"); ok && v.(string) != "" {
		listArgs.PrivateIps = v.(string)
	}

	action := "List all Instance "
	instanceList, err := bccService.ListAllInstance(listArgs)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instances", action, BCESDKGoERROR)
	}

	instanceMap, err := bccService.FlattenInstanceModelToMap(instanceList)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instances", action, BCESDKGoERROR)
	}

	FilterDataSourceResult(d, &instanceMap)
	if err = d.Set("instances", instanceMap); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instances", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), instanceMap); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_instances", action, BCESDKGoERROR)
		}
	}

	return nil
}
