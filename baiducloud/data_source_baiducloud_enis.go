/*
Use this data source to query ENI list.

Example Usage

```hcl
data "baiducloud_enis" "default" {
  vpc_id      = "vpc-xxxxxx"
}

output "enis" {
 value = "${data.baiducloud_enis.default.enis}"
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/eni"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudEnis() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudEnisRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "Vpc id which ENI belong to",
				Required:    true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Description: "Instance id the ENI bind",
				Optional:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of ENI",
				Optional:    true,
			},
			"private_ip_address": {
				Type:        schema.TypeList,
				Description: "Eni private IP address",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "ENI list result output file",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),
			"enis": {
				Type:        schema.TypeList,
				Description: "ENI list",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"eni_id": {
							Type:        schema.TypeString,
							Description: "ENI ID",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Name of ENI",
							Computed:    true,
						},
						"zone_name": {
							Type:        schema.TypeString,
							Description: "ENI Availability Zone Name",
							Computed:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "Description of ENI",
							Computed:    true,
						},
						"instance_id": {
							Type:        schema.TypeString,
							Description: "Instance id which ENI bind",
							Computed:    true,
						},
						"mac_address": {
							Type:        schema.TypeString,
							Description: "ENI Mac Address",
							Computed:    true,
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Description: "Subnet ID which ENI belong to",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Status of ENI",
							Computed:    true,
						},
						"security_group_ids": {
							Type:        schema.TypeList,
							Description: "ENI security group IDs",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"enterprise_security_group_ids": {
							Type:        schema.TypeList,
							Description: "ENI enterprise security group IDs",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"private_ip_set": {
							Type:        schema.TypeList,
							Description: "ENI private ip set",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"public_ip_address": {
										Type:        schema.TypeString,
										Description: "Public IP address",
										Computed:    true,
									},
									"primary": {
										Type:        schema.TypeBool,
										Description: "True or false, true mean it is primary IP, it's private IP address can not modify, only one primary IP in a ENI",
										Computed:    true,
									},
									"private_ip_address": {
										Type:        schema.TypeString,
										Description: "Private IP address",
										Computed:    true,
									},
								},
							},
						},
						"created_time": {
							Type:        schema.TypeString,
							Description: "ENI create time",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudEnisRead(d *schema.ResourceData, meta interface{}) error {
	action := "Query Eni List"
	client := meta.(*connectivity.BaiduClient)
	eniService := &EniService{
		client: client,
	}

	enis, err := eniService.ListEnis(buildListEniArgs(d))
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eni", action, BCESDKGoERROR)
	}
	enisMap := make([]map[string]interface{}, 0)
	for _, item := range enis {
		enisMap = append(enisMap, eniService.eniToMap(item))
	}
	FilterDataSourceResult(d, &enisMap)
	if err := d.Set("enis", enisMap); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eni", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), enisMap); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_eni", action, BCESDKGoERROR)
		}
	}
	return nil
}

func buildListEniArgs(d *schema.ResourceData) *eni.ListEniArgs {
	res := &eni.ListEniArgs{}
	if v, ok := d.GetOk("name"); ok {
		res.Name = v.(string)
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		res.VpcId = v.(string)
	}
	if v, ok := d.GetOk("instance_id"); ok {
		res.InstanceId = v.(string)
	}
	if v, ok := d.GetOk("private_ip_address"); ok {
		res.PrivateIpAddress = v.([]string)
	}
	return res
}
